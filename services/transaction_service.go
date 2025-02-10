package services

import (
	"encoding/base64"
	"errors"
	"fmt"
	"goPocDemo/initializer"
	"goPocDemo/model"
	"goPocDemo/utils"
	"log"
	"regexp"
	"time"

	"gorm.io/gorm"
)

// Check for duplicate transaction
func CheckDuplicateTransaction(transactionID string) error {
	var existingTransaction model.Transaction
	dbInstance := initializer.GetDB()
	// Perform lookup in the database
	if err := dbInstance.Where("transaction_id = ?", transactionID).First(&existingTransaction).Error; err != nil {
		// If no record is found, return nil (not a duplicate)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		// If another error occurs, return it
		return err
	}

	// If a record is found, it's a duplicate
	return errors.New("duplicate transaction id")
}

func CreateAuditLog(transaction_id, action string, details string) (*model.AuditLog, error) {
	auditLog := &model.AuditLog{
		AuditLogId:    utils.GenerateUniqueID(),
		TransactionID: transaction_id,
		Action:        action,
		CreatedAt:     time.Now(),
		Details:       details,
	}

	// Insert user into the database using GORM
	if err := initializer.GetDB().Create(auditLog).Error; err != nil {
		return nil, err
	}

	// Return user and nil error
	return auditLog, nil
}

func CreateTransaction(req model.Request, status string) (*model.Transaction, error) {
	transaction := &model.Transaction{
		TransactionID:   req.TransactionID,
		PayerID:         req.PayerID,
		PayeeID:         req.PayeeID,
		Amount:          req.Amount,
		TransactionType: req.TransactionType,
		Status:          status,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
		ReservedAmount:  req.Amount,
	}

	// Insert user into the database using GORM
	if err := initializer.GetDB().Create(transaction).Error; err != nil {
		return nil, err
	}

	// Return user and nil error
	return transaction, nil
}

func DecodeCardDetails(paymentDetails map[string]string) ([]byte, bool) {
	decoded, err := base64.StdEncoding.DecodeString(paymentDetails["card_number"])
	if err != nil {
		return nil, false
	}

	decrypted, err := Decrypt(decoded, secretKey)
	if err != nil {
		return nil, false
	}

	fmt.Println(paymentDetails, paymentDetails["card_number"], decrypted, secretKey)

	return decrypted, true
}

func UpdateTransactionStatus(transactionID, status string) error {
	// Update the status field in the transaction table
	result := initializer.GetDB().Model(&model.Transaction{}).Where("transaction_id = ?", transactionID).Updates(map[string]interface{}{
		"reserved_amount": 0.0,
		"status":          status,
	})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// Do we have to implement Luhn's algorithm  to Validate a card number
func validateCardNumber(cardNumber string, paymentMethodCardDetails string) bool {
	cardNumberLength := len(cardNumber)
	return cardNumberLength == 16 && paymentMethodCardDetails == cardNumber
}

// Validate CVV (typically 3 or 4 digits)
func validateCVV(cvv string, paymentMethodCvv string) bool {
	match, _ := regexp.MatchString(`^\d{3,4}$`, cvv)
	return match && cvv == paymentMethodCvv
}

func validateExpiryDate(expiryDate string, paymentMethodExpiryDate string) bool {
	// Check format using regex
	match, _ := regexp.MatchString(`^(0[1-9]|1[0-2])/(\d{2})$`, expiryDate)
	if !match {
		return false
	}

	// Parse the date
	now := time.Now()
	currentYear := now.Year() % 100 // Get the last two digits of the current year
	currentMonth := int(now.Month())
	month := int(expiryDate[0]-'0')*10 + int(expiryDate[1]-'0')
	year := int(expiryDate[3]-'0')*10 + int(expiryDate[4]-'0')

	// Validate the card expiry
	return (year > currentYear || (year == currentYear && month >= currentMonth)) && expiryDate == paymentMethodExpiryDate
}

func ValidatePayerBalance(req model.Request) error {
	dbInstance := initializer.GetDB()
	var payer model.Payer
	var transactions []model.Transaction
	var reserveAmount float64
	if err := dbInstance.Where("payer_id = ?", req.PayerID).First(&payer).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("payer not found")
		}
		return err
	}

	if err := dbInstance.Where("status = ? AND payer_id = ?", "Pending", req.PayerID).Find(&transactions).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("transaction not found")
		}
		return err
	}

	for _, transaction := range transactions {
		reserveAmount += transaction.ReservedAmount
	}

	if payer.Balance < (req.Amount + reserveAmount) {
		dbInstance.Model(&model.Transaction{}).Where("transaction_id = ?", req.TransactionID).Update("status", "Failed")
		return errors.New("insufficient funds")
	}

	return nil
}

// Validate the payment method
func ValidatePaymentMethod(paymentMethodID string, paymentDetails map[string]string) error {
	var paymentMethod model.PaymentMethod
	dbInstance := initializer.GetDB()

	// Fetch the payment method by ID
	if err := dbInstance.Where("payment_method_id = ?", paymentMethodID).First(&paymentMethod).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("payment method not found")
		}
		return err
	}

	// Check if the payment method is active
	if paymentMethod.Status != "active" {
		return errors.New("payment method is not active")
	}

	if paymentMethod.MethodType == "card" {
		cardNumber, ok1 := paymentDetails["card_number"]
		//cardNumber, ok1 := DecodeCardDetails(paymentDetails)
		expiryDate, ok2 := paymentDetails["expiry_date"]
		cvv, ok3 := paymentDetails["cvv"]

		// Check if all required fields are present
		if !ok1 || !ok2 || !ok3 {
			return errors.New("missing required card details")
		}

		// Validate card details
		if !validateCardNumber(cardNumber, paymentMethod.Details) {
			return errors.New("invalid card number")
		}

		// encrypt the card number details again

		if !validateExpiryDate(expiryDate, paymentMethod.ExpiryDate) {
			return errors.New("card has expired or invalid expiry date")
		}

		if !validateCVV(cvv, paymentMethod.CVV) {
			return errors.New("invalid CVV")
		}
	}
	// All validations passed
	return nil
}

func ValidateRequest(req *model.Request) error {
	log.Println("Validating requests")
	if req.TransactionID == "" {
		return errors.New("TransactionID is required")
	}

	if req.PayerID == "" {
		return errors.New("PayerID is required")
	}

	if req.PayeeID == "" {
		return errors.New("PayeeID is required")
	}

	if req.PayerID == req.PayeeID {
		return errors.New("PayerID and PayeeID cannot be the same")
	}

	if req.Amount == 0 {
		return errors.New("Amount is required")
	}

	if req.Amount <= 0 {
		return errors.New("amount must be greater than zero")
	}

	if req.TransactionType == "" {
		return errors.New("TransactionType is required")
	}

	validTransactionTypes := map[string]bool{"Debit": true, "Credit": true, "Refund": true}
	if !validTransactionTypes[req.TransactionType] {
		return errors.New("TransactionType must be either 'Debit' or 'Credit' or 'Refund'")
	}

	if req.PaymentMethodID == "" {
		return errors.New("PaymentMethodID is required")
	}

	if len(req.Details) == 0 {
		return errors.New("details is required")
	}

	_, ok1 := req.Details["card_number"]
	_, ok2 := req.Details["expiry_date"]
	_, ok3 := req.Details["cvv"]

	// Check if all required fields are present
	if !ok1 || !ok2 || !ok3 {
		return errors.New("missing required card details")
	}

	return nil
}
