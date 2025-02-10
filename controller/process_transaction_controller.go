package controller

import (
	"errors"
	"goPocDemo/initializer"
	"goPocDemo/model"
	"goPocDemo/services"
	"log"

	"github.com/kataras/iris/v12"
	"gorm.io/gorm"
)

// StartTransactionHandler handles the transaction request
func StartTransactionHandler(ctx iris.Context) {
	var req model.Request

	// Parse and validate the request
	if err := ctx.ReadJSON(&req); err != nil {
		sendErrorResponse(ctx, iris.StatusBadRequest, "Failed to parse request")
		return
	}

	services.CreateAuditLog(req.TransactionID, req.TransactionType, "Transaction Initiated")

	if err := services.ValidateRequest(&req); err != nil {
		logAndRespond(ctx, req.TransactionID, req.TransactionType, "Invalid request body", iris.StatusBadRequest, err.Error())
		return
	}

	// Check for duplicate transactions
	if req.TransactionType != "Refund" {
		if err := services.CheckDuplicateTransaction(req.TransactionID); err != nil {
			logAndRespond(ctx, req.TransactionID, req.TransactionType, "Duplicate transaction", iris.StatusConflict, err.Error())
			return
		}
	}

	// Create a new transaction
	if req.TransactionType != "Refund" {
		_, err := services.CreateTransaction(req, "Pending")
		if err != nil {
			logAndRespond(ctx, req.TransactionID, req.TransactionType, "Failed to create transaction", iris.StatusInternalServerError, err.Error())
			return
		}
		services.CreateAuditLog(req.TransactionID, req.TransactionType, "Transaction is pending")
	}

	// Validate payment method
	if err := services.ValidatePaymentMethod(req.PaymentMethodID, req.Details); err != nil {
		logAndRespond(ctx, req.TransactionID, req.TransactionType, "Invalid payment details", iris.StatusConflict, err.Error())
		services.UpdateTransactionStatus(req.TransactionID, "Failed")
		return
	}

	// Encrypt sensitive data
	if err := services.EncryptData(&req); err != nil {
		sendErrorResponse(ctx, iris.StatusInternalServerError, err.Error())
		return
	}

	// Handle specific transaction types
	if err := handleTransactionTypes(ctx, req); err != nil {
		logAndRespond(ctx, req.TransactionID, req.TransactionType, "Invalid request body", iris.StatusBadRequest, err.Error())
		return
	}
	logSuceessDetails(ctx, req)
}

// Credit amount to payee
func creditAmountToPayee(payeeID string, requestAmount float64) error {
	dbInstance := initializer.GetDB()
	var payee model.Payee
	if err := dbInstance.Where("payee_id = ?", payeeID).First(&payee).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("payee not found")
		}
		return err
	}

	// Update balance
	result := dbInstance.Model(&model.Payee{}).Where("payee_id = ?", payeeID).Update("balance", payee.Balance+requestAmount)
	return result.Error
}

// Deduct amount from payer
func deductAmountFromPayer(payerID string, requestAmount float64) error {
	dbInstance := initializer.GetDB()
	var payer model.Payer
	if err := dbInstance.Where("payer_id = ?", payerID).First(&payer).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("payer not found")
		}
		return err
	}

	// Deduct balance
	result := dbInstance.Model(&model.Payer{}).Where("payer_id = ?", payerID).Update("balance", payer.Balance-requestAmount)
	return result.Error
}

func handleCreditFunctionality(ctx iris.Context, req model.Request) error {
	// Reserve the transaction amount
	if err := reserveTransactionAmount(req.TransactionID, req.Amount); err != nil {
		logAndRespond(ctx, req.TransactionID, req.TransactionType, "Failed to reserve transaction amount", iris.StatusInternalServerError, err.Error())
		return err
	}

	// Credit amount to payee
	if err := creditAmountToPayee(req.PayeeID, req.Amount); err != nil {
		logAndRespond(ctx, req.TransactionID, req.TransactionType, "Failed to credit amount to payee", iris.StatusInternalServerError, err.Error())
		return err
	}

	// Log success
	services.CreateAuditLog(req.TransactionID, req.TransactionType, "Credit transaction completed successfully")
	return nil
}

func handleDebitFunctionality(ctx iris.Context, req model.Request) error {
	// Validate payer balance
	if err := services.ValidatePayerBalance(req); err != nil {
		logAndRespond(ctx, req.TransactionID, req.TransactionType, "Insufficient balance", iris.StatusInternalServerError, err.Error())
		return err
	}

	// Deduct amount from payer
	if err := deductAmountFromPayer(req.PayerID, req.Amount); err != nil {
		logAndRespond(ctx, req.TransactionID, req.TransactionType, "Failed to deduct amount from payer", iris.StatusInternalServerError, err.Error())
		return err
	}

	// Reserve the transaction amount
	if err := reserveTransactionAmount(req.TransactionID, req.Amount); err != nil {
		logAndRespond(ctx, req.TransactionID, req.TransactionType, "Failed to reserve transaction amount", iris.StatusInternalServerError, err.Error())
		return err
	}

	// Credit amount to payee
	if err := creditAmountToPayee(req.PayeeID, req.Amount); err != nil {
		logAndRespond(ctx, req.TransactionID, req.TransactionType, "Failed to credit amount to payee", iris.StatusInternalServerError, err.Error())
		return err
	}

	// Log success
	services.CreateAuditLog(req.TransactionID, req.TransactionType, "Debit transaction completed successfully")
	return nil
}

func handleRefundFunctionality(req model.Request) error {
	// Fetch the original transaction
	var transaction model.Transaction
	dbInstance := initializer.GetDB()
	if err := dbInstance.Where("transaction_id = ?", req.TransactionID).First(&transaction).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("transaction not found")
		}
		return err
	}

	// check if the original transaction is already refunded TransactionType is Refund and Status is Completed
	if transaction.TransactionType == "Refund" && transaction.Status == "Completed" {
		return errors.New("original transaction is already refunded")
	}
	// Fetch payer and payee details
	var payer model.Payer
	if err := dbInstance.Where("payer_id = ?", req.PayerID).First(&payer).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("payer not found")
		}
		return err
	}

	var payee model.Payee
	if err := dbInstance.Where("payee_id = ?", req.PayeeID).First(&payee).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("payee not found")
		}
		return err
	}

	// Adjust balances based on the original transaction type
	switch transaction.TransactionType {
	case "Debit":
		dbInstance.Model(&model.Payer{}).Where("payer_id = ?", req.PayerID).Update("balance", payer.Balance+req.Amount)
		dbInstance.Model(&model.Payee{}).Where("payee_id = ?", req.PayeeID).Update("balance", payee.Balance-req.Amount)
	case "Credit":
		dbInstance.Model(&model.Payee{}).Where("payee_id = ?", req.PayeeID).Update("balance", payee.Balance-req.Amount)
	}

	// Update the transaction type to "Refund"
	if err := dbInstance.Model(&model.Transaction{}).Where("transaction_id = ?", req.TransactionID).Update("transaction_type", "Refund").Error; err != nil {
		return err
	}

	// Log success
	services.CreateAuditLog(req.TransactionID, req.TransactionType, "Refund transaction completed successfully")
	return nil
}

func handleTransactionTypes(ctx iris.Context, req model.Request) error {
	var err error
	switch req.TransactionType {
	case "Debit":
		err = handleDebitFunctionality(ctx, req)
	case "Credit":
		err = handleCreditFunctionality(ctx, req)
	case "Refund":
		err = handleRefundFunctionality(req)
	}
	if err != nil {
		sendErrorResponse(ctx, iris.StatusInternalServerError, err.Error())
		return err
	}

	return nil
}

// Helper function to create audit logs and send responses
func logAndRespond(ctx iris.Context, transactionID, transactionType, message string, statusCode int, errMsg string) {
	services.CreateAuditLog(transactionID, transactionType, message)
	sendErrorResponse(ctx, statusCode, errMsg)
}

func logSuceessDetails(ctx iris.Context, req model.Request) {
	// Update transaction status to completed
	services.UpdateTransactionStatus(req.TransactionID, "Completed")

	// Respond with success
	ctx.StatusCode(iris.StatusCreated)
	ctx.JSON(map[string]string{
		"status":         "Completed",
		"transaction_id": req.TransactionID,
		"message":        "Transaction created successfully",
	})
	log.Printf("Transaction %s is successful", req.TransactionID)
}

// Reserve transaction amount
func reserveTransactionAmount(transactionID string, amount float64) error {
	dbInstance := initializer.GetDB()
	result := dbInstance.Model(&model.Transaction{}).Where("transaction_id = ?", transactionID).Updates(map[string]interface{}{
		"reserved_amount": amount,
		"status":          "Reserved",
	})
	return result.Error
}

// Helper function to send error responses
func sendErrorResponse(ctx iris.Context, statusCode int, errMsg string) {
	ctx.StatusCode(statusCode)
	ctx.JSON(map[string]string{"error": errMsg})
	log.Println("error", errMsg)
}
