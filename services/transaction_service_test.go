package services

import (
	"encoding/base64"
	"errors"
	"fmt"
	"goPocDemo/initializer"
	"goPocDemo/model"
	"log"
	"testing"
	"time"

	spannergorm "github.com/googleapis/go-gorm-spanner"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB

func dbConnection() {
	// Build the Spanner connection string
	dsn := fmt.Sprintf("projects/%s/instances/%s/databases/%s", "iconic-star-447805-v6", "go-instance-poc", "go_spanner_db")

	// Connect to Spanner using the GORM driver for Spanner
	db, _ = gorm.Open(spannergorm.New(spannergorm.Config{
		DriverName: "spanner", // Spanner driver name
		DSN:        dsn,       // The connection string (Data Source Name)
	}), &gorm.Config{
		PrepareStmt:                      true,
		IgnoreRelationshipsWhenMigrating: true,
		Logger:                           logger.Default.LogMode(logger.Error),
	})
	//db, err := initializer.ConnectSpannerDB()
	//assert.NoError(t, err)

	// Initialize UserService with GORM DB
	//userService = &UserService{DB: db}
}

func TestValidateRequest(t *testing.T) {
	tests := []struct {
		name    string
		req     model.Request
		expects error
	}{
		{
			name: "Valid Request",
			req: model.Request{
				TransactionID:   "txn_12345",
				PayerID:         "payer_001",
				PayeeID:         "payee_002",
				Amount:          100.50,
				TransactionType: "Debit",
				PaymentMethodID: "pm_123",
				Details: map[string]string{
					"card_number": "1234567812345678",
					"expiry_date": "12/24",
					"cvv":         "123",
				},
			},
			expects: nil,
		},
		{
			name: "Missing TransactionID",
			req: model.Request{
				PayerID:         "payer_001",
				PayeeID:         "payee_002",
				Amount:          100.50,
				TransactionType: "Debit",
				PaymentMethodID: "pm_123",
				Details: map[string]string{
					"card_number": "1234567812345678",
					"expiry_date": "12/24",
					"cvv":         "123",
				},
			},
			expects: errors.New("TransactionID is required"),
		},
		{
			name: "Invalid Transaction Type",
			req: model.Request{
				TransactionID:   "txn_12345",
				PayerID:         "payer_001",
				PayeeID:         "payee_002",
				Amount:          100.50,
				TransactionType: "InvalidType",
				PaymentMethodID: "pm_123",
				Details: map[string]string{
					"card_number": "1234567812345678",
					"expiry_date": "12/24",
					"cvv":         "123",
				},
			},
			expects: errors.New("TransactionType must be either 'Debit' or 'Credit' or 'Refund'"),
		},
		{
			name: "Missing Card Details",
			req: model.Request{
				TransactionID:   "txn_12345",
				PayerID:         "payer_001",
				PayeeID:         "payee_002",
				Amount:          100.50,
				TransactionType: "Debit",
				PaymentMethodID: "pm_123",
				Details:         map[string]string{}, // Empty details
			},
			expects: errors.New("details is required"), // Update this message
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateRequest(&tt.req)
			if (err == nil && tt.expects != nil) || (err != nil && tt.expects == nil) || (err != nil && tt.expects != nil && err.Error() != tt.expects.Error()) {
				t.Errorf("Test %s failed: expected %v, got %v", tt.name, tt.expects, err)
			}
		})
	}
}

// TestCheckDuplicateTransaction tests the CheckDuplicateTransaction function
func TestCheckDuplicateTransaction(t *testing.T) {
	//t.Skip()
	dbConnection()
	initializer.DbInstance = db

	transactionID := "txn_12345"

	// /  **Cleanup before each test case**
	var txnCount int64
	db.Model(&model.Transaction{}).Where("transaction_id = ?", "txn_12345").Count(&txnCount)
	log.Printf("🔍 Debug: Found %d records for transaction_id = txn_12345", txnCount)
	if txnCount > 0 {
		db.Where("transaction_id = ?", transactionID).Delete(&model.Transaction{})
	}

	// **Test Case 1: No Duplicate Transaction**
	t.Run("No Duplicate Transaction", func(t *testing.T) {
		err := CheckDuplicateTransaction(transactionID)
		assert.NoError(t, err) // Should return nil (no duplicate)
	})

	// **Test Case 2: Duplicate Transaction Exists**
	t.Run("Duplicate Transaction Exists", func(t *testing.T) {
		// Insert transaction
		newTransaction := &model.Transaction{
			TransactionID:   "txn_12345",
			PayerID:         "234",
			PayeeID:         "34",
			Amount:          100.0,
			TransactionType: "Debit",
			Status:          "Pending",
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
			ReservedAmount:  100.0,
		}

		insertErr := db.Create(newTransaction).Error
		assert.NoError(t, insertErr, "Failed to insert test transaction")

		err := CheckDuplicateTransaction(transactionID)
		assert.Error(t, err)
		assert.Equal(t, "duplicate transaction id", err.Error(), "Expected duplicate error")
	})

	// **Test Case 3: Database Error**
	// t.Run("Database Error", func(t *testing.T) {
	// 	// Close DB connection to simulate an error
	// 	dbConn, _ := db.DB()
	// 	dbConn.Close() // Simulate DB failure

	// 	err := CheckDuplicateTransaction(userService, transactionID)
	// 	assert.Error(t, err)
	// 	assert.Contains(t, err.Error(), "database is closed", "Unexpected error message")
	// })
}

func TestValidatePaymentMethod(t *testing.T) {
	//t.Skip()
	dbConnection()

	initializer.DbInstance = db

	// /  **Cleanup before each test case**
	var pmCount int64
	db.Model(&model.PaymentMethod{}).Where("payment_method_id = ?", "pm_12345").Count(&pmCount)
	log.Printf("🔍 Debug: Found %d records for payment_method_id = pm_12345", pmCount)
	if pmCount > 0 {
		db.Where("payment_method_id = ?", "pm_12345").Delete(&model.PaymentMethod{})
	}

	//Insert a valid active payment method
	validPayment := model.PaymentMethod{
		PaymentMethodID: "pm_12345",
		PayerID:         "234",
		MethodType:      "card",
		Details:         "4111111111111111",
		ExpiryDate:      time.Now().AddDate(1, 0, 0).Format("01/06"), // Valid expiry date
		CVV:             "123",
		Status:          "active",
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
	db.Create(&validPayment)

	// **Test Case 1: Valid Payment Method**
	t.Run("Valid Payment Method", func(t *testing.T) {
		paymentDetails := map[string]string{
			"card_number": "4111111111111111",
			"expiry_date": validPayment.ExpiryDate,
			"cvv":         validPayment.CVV,
		}
		err := ValidatePaymentMethod("pm_12345", paymentDetails)
		assert.NoError(t, err, "Expected valid payment method")
	})

	//**Test Case 2: Payment Method Not Found**
	t.Run("Payment Method Not Found", func(t *testing.T) {
		err := ValidatePaymentMethod("pm_99999", map[string]string{})
		assert.Error(t, err)
		assert.Equal(t, "payment method not found", err.Error())
	})

	// **Test Case 3: Payment Method is Inactive**
	t.Run("Payment Method is Inactive", func(t *testing.T) {
		inactivePayment := validPayment
		inactivePayment.PaymentMethodID = "pm_inactive"
		inactivePayment.Status = "inactive"
		// /  **Cleanup before each test case**
		var pmInActiveCount int64
		db.Model(&model.PaymentMethod{}).Where("payment_method_id = ?", "pm_inactive").Count(&pmInActiveCount)
		log.Printf("🔍 Debug: Found %d records for payment_method_id = pm_inactive", pmInActiveCount)
		if pmInActiveCount > 0 {
			db.Where("payment_method_id = ?", "pm_inactive").Delete(&model.PaymentMethod{})
		}
		db.Create(&inactivePayment)

		err := ValidatePaymentMethod("pm_inactive", map[string]string{})
		assert.Error(t, err)
		assert.Equal(t, "payment method is not active", err.Error())
	})

	// **Test Case 4: Missing Required Card Details**
	t.Run("Missing Required Card Details", func(t *testing.T) {
		err := ValidatePaymentMethod("pm_12345", map[string]string{"card_number": "4111111111111111"})
		assert.Error(t, err)
		assert.Equal(t, "missing required card details", err.Error())
	})

	// **Test Case 5: Invalid Card Number**
	t.Run("Invalid Card Number", func(t *testing.T) {
		paymentDetails := map[string]string{
			"card_number": "1234567890123456", // Invalid card number
			"expiry_date": validPayment.ExpiryDate,
			"cvv":         validPayment.CVV,
		}
		err := ValidatePaymentMethod("pm_12345", paymentDetails)
		assert.Error(t, err)
		assert.Equal(t, "invalid card number", err.Error())
	})

	// **Test Case 6: Expired Card**
	t.Run("Expired Card", func(t *testing.T) {
		expiredPayment := validPayment
		expiredPayment.PaymentMethodID = "pm_expired"
		expiredPayment.ExpiryDate = time.Now().AddDate(-1, 0, 0).Format("01/06") // Expired card
		var pmExpiredDateCount int64
		db.Model(&model.PaymentMethod{}).Where("payment_method_id = ?", "pm_expired").Count(&pmExpiredDateCount)
		log.Printf("🔍 Debug: Found %d records for payment_method_id = pm_expired", pmExpiredDateCount)
		if pmExpiredDateCount > 0 {
			db.Where("payment_method_id = ?", "pm_expired").Delete(&model.PaymentMethod{})
		}
		db.Create(&expiredPayment)

		paymentDetails := map[string]string{
			"card_number": "4111111111111111",
			"expiry_date": expiredPayment.ExpiryDate,
			"cvv":         expiredPayment.CVV,
		}
		err := ValidatePaymentMethod("pm_expired", paymentDetails)
		assert.Error(t, err)
		assert.Equal(t, "card has expired or invalid expiry date", err.Error())
	})

	// **Test Case 7: Invalid CVV**
	t.Run("Invalid CVV", func(t *testing.T) {
		paymentDetails := map[string]string{
			"card_number": "4111111111111111",
			"expiry_date": validPayment.ExpiryDate,
			"cvv":         "999", // Invalid CVV
		}
		err := ValidatePaymentMethod("pm_12345", paymentDetails)
		assert.Error(t, err)
		assert.Equal(t, "invalid CVV", err.Error())
	})
}

func TestValidateCardNumber(t *testing.T) {
	result := validateCardNumber("1234567812345678", "1234567812345678")
	if !result {
		t.Errorf("Expected true, got %v", result)
	}
}

// validateExpiryDate is assumed to be imported from the actual package
func TestValidateExpiryDate(t *testing.T) {
	tests := []struct {
		name                    string
		expiryDate              string
		paymentMethodExpiryDate string
		expectedResult          bool
	}{
		{
			name:                    "Valid expiry date",
			expiryDate:              "12/26",
			paymentMethodExpiryDate: "12/26",
			expectedResult:          true,
		},
		{
			name:                    "Mismatched expiry date",
			expiryDate:              "12/26",
			paymentMethodExpiryDate: "11/26",
			expectedResult:          false,
		},
		{
			name:                    "Expired card",
			expiryDate:              "01/22",
			paymentMethodExpiryDate: "01/22",
			expectedResult:          false,
		},
		{
			name:                    "Invalid format",
			expiryDate:              "13/25",
			paymentMethodExpiryDate: "13/25",
			expectedResult:          false,
		},
		{
			name:                    "Current month and year",
			expiryDate:              time.Now().Format("01/06"),
			paymentMethodExpiryDate: time.Now().Format("01/06"),
			expectedResult:          true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := validateExpiryDate(tc.expiryDate, tc.paymentMethodExpiryDate)
			assert.Equal(t, tc.expectedResult, result)
		})
	}
}

func TestValidateCVV(t *testing.T) {
	tests := []struct {
		name      string
		cvv       string
		storedCvv string
		expected  bool
	}{
		{
			name:      "Valid 3-digit CVV",
			cvv:       "123",
			storedCvv: "123",
			expected:  true,
		},
		{
			name:      "Valid 4-digit CVV",
			cvv:       "1234",
			storedCvv: "1234",
			expected:  true,
		},
		{
			name:      "Mismatched CVV",
			cvv:       "123",
			storedCvv: "321",
			expected:  false,
		},
		{
			name:      "Invalid CVV format (letters included)",
			cvv:       "12a",
			storedCvv: "12a",
			expected:  false,
		},
		{
			name:      "Invalid CVV length (too short)",
			cvv:       "12",
			storedCvv: "12",
			expected:  false,
		},
		{
			name:      "Invalid CVV length (too long)",
			cvv:       "12345",
			storedCvv: "12345",
			expected:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := validateCVV(tt.cvv, tt.storedCvv)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCreateAuditLog(t *testing.T) {
	dbConnection()
	initializer.DbInstance = db

	// Cleanup before each test case
	var auditLogCount int64
	db.Model(&model.AuditLog{}).Where("audit_log_id = ?", "au_1234").Count(&auditLogCount)
	log.Printf("🔍 Debug: Found %d records for audit_log_id = au_1234", auditLogCount)
	if auditLogCount > 0 {
		db.Where("audit_log_id = ?", "au_1234").Delete(&model.AuditLog{})
	}

	//Insert a valid active payment method
	auditLogData := model.AuditLog{
		AuditLogId:    "au_1234",
		TransactionID: "897",
		Action:        "card",
		CreatedAt:     time.Now(),
		Details:       "test",
	}
	db.Create(&auditLogData)

	tests := []struct {
		name          string
		transactionID string
		action        string
		details       string
		expectError   bool
	}{
		{
			name:          "Valid Audit Log",
			transactionID: "897",
			action:        "CREATE",
			details:       "Transaction created successfully",
			expectError:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			auditLog, err := CreateAuditLog(tt.transactionID, tt.action, tt.details)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, auditLog)
				assert.Equal(t, tt.transactionID, auditLog.TransactionID)
				assert.Equal(t, tt.action, auditLog.Action)
				assert.Equal(t, tt.details, auditLog.Details)
			}
		})
	}
}

func TestCreateTransaction(t *testing.T) {
	dbConnection()
	initializer.DbInstance = db

	// Cleanup before each test case
	var transactionCount int64
	db.Model(&model.Transaction{}).Where("transaction_id = ?", "txn_12347").Count(&transactionCount)
	log.Printf("🔍 Debug: Found %d records for transaction_id = txn_12347", transactionCount)
	if transactionCount > 0 {
		db.Where("transaction_id = ?", "txn_12347").Delete(&model.Transaction{})
	}

	tests := []struct {
		name    string
		req     model.Request
		status  string
		expects error
	}{
		{
			name: "Valid Transaction",
			req: model.Request{
				TransactionID:   "txn_12347",
				PayerID:         "payer_001",
				PayeeID:         "payee_002",
				Amount:          100.50,
				TransactionType: "Debit",
				PaymentMethodID: "pm_123",
				Details: map[string]string{
					"info": "newDetails",
				},
			},
			status:  "Pending",
			expects: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			transaction, err := CreateTransaction(tt.req, tt.status)
			if (err == nil && tt.expects != nil) || (err != nil && tt.expects == nil) || (err != nil && tt.expects != nil && err.Error() != tt.expects.Error()) {
				t.Errorf("Test %s failed: expected %v, got %v", tt.name, tt.expects, err)
			}
			if err == nil {
				assert.NotNil(t, transaction)
				assert.Equal(t, tt.req.TransactionID, transaction.TransactionID)
				assert.Equal(t, tt.req.PayerID, transaction.PayerID)
				assert.Equal(t, tt.req.PayeeID, transaction.PayeeID)
				assert.Equal(t, tt.req.Amount, transaction.Amount)
				assert.Equal(t, tt.req.TransactionType, transaction.TransactionType)
				assert.Equal(t, tt.status, transaction.Status)
			}
		})
	}
}

func TestUpdateTransactionStatus(t *testing.T) {
	dbConnection()
	initializer.DbInstance = db

	// Cleanup before each test case
	var transactionCount int64
	db.Model(&model.Transaction{}).Where("transaction_id = ?", "txn_12348").Count(&transactionCount)
	log.Printf("🔍 Debug: Found %d records for transaction_id = txn_12348", transactionCount)
	if transactionCount > 0 {
		db.Where("transaction_id = ?", "txn_12348").Delete(&model.Transaction{})
	}

	// Insert a transaction to update
	newTransaction := &model.Transaction{
		TransactionID:   "txn_12348",
		PayerID:         "payer_001",
		PayeeID:         "payee_002",
		Amount:          100.0,
		TransactionType: "Debit",
		Status:          "Pending",
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
		ReservedAmount:  100.0,
	}
	db.Create(newTransaction)

	tests := []struct {
		name          string
		transactionID string
		status        string
		expectError   bool
	}{
		{
			name:          "Valid Status Update",
			transactionID: "txn_12348",
			status:        "Completed",
			expectError:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := UpdateTransactionStatus(tt.transactionID, tt.status)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)

				// Verify the status update
				var updatedTransaction model.Transaction
				db.Where("transaction_id = ?", tt.transactionID).First(&updatedTransaction)
				assert.Equal(t, tt.status, updatedTransaction.Status)
				assert.Equal(t, 0.0, updatedTransaction.ReservedAmount)
			}
		})
	}
}

func TestDecodeCardDetails(t *testing.T) {
	tests := []struct {
		name            string
		paymentDetails  map[string]string
		expectedResult  []byte
		expectedSuccess bool
	}{
		{
			name: "Invalid Base64 Encoding",
			paymentDetails: map[string]string{
				"card_number": "invalid_base64",
			},
			expectedResult:  nil,
			expectedSuccess: false,
		},
		{
			name: "Decryption Failure",
			paymentDetails: map[string]string{
				"card_number": base64.StdEncoding.EncodeToString([]byte("invalid_encrypted_data")),
			},
			expectedResult:  nil,
			expectedSuccess: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, success := DecodeCardDetails(tt.paymentDetails)
			assert.Equal(t, tt.expectedSuccess, success)
			assert.Equal(t, tt.expectedResult, result)
		})
	}
}

func TestValidatePayerBalance(t *testing.T) {
	dbConnection()
	initializer.DbInstance = db

	// Cleanup before each test case
	var payerCount int64
	db.Model(&model.Payer{}).Where("payer_id = ?", "payer_001").Count(&payerCount)
	log.Printf("🔍 Debug: Found %d records for payer_id = payer_001", payerCount)
	if payerCount > 0 {
		db.Where("payer_id = ?", "payer_001").Delete(&model.Payer{})
	}

	var transactionCount int64
	db.Model(&model.Transaction{}).Where("payer_id = ?", "payer_001").Count(&transactionCount)
	log.Printf("🔍 Debug: Found %d records for payer_id = payer_001", transactionCount)
	if transactionCount > 0 {
		db.Where("payer_id = ?", "payer_001").Delete(&model.Transaction{})
	}

	// Insert a payer
	payer := &model.Payer{
		PayerID: "payer_001",
		Balance: 200.0,
	}
	db.Create(payer)

	tests := []struct {
		name    string
		req     model.Request
		setup   func()
		expects error
	}{
		{
			name: "Payer Not Found",
			req: model.Request{
				PayerID: "non_existent_payer",
				Amount:  100.0,
			},
			setup:   func() {},
			expects: errors.New("payer not found"),
		},
		{
			name: "Insufficient Funds",
			req: model.Request{
				TransactionID: "txn_12349",
				PayerID:       "payer_001",
				Amount:        300.0,
			},
			setup: func() {
				// Insert a pending transaction
				transaction := &model.Transaction{
					TransactionID:   "txn_12349",
					PayerID:         "payer_001",
					PayeeID:         "payee_002",
					Amount:          100.0,
					TransactionType: "Debit",
					Status:          "Pending",
					CreatedAt:       time.Now(),
					UpdatedAt:       time.Now(),
					ReservedAmount:  100.0,
				}
				db.Create(transaction)
			},
			expects: errors.New("insufficient funds"),
		},
		{
			name: "Sufficient Funds",
			req: model.Request{
				TransactionID: "txn_12350",
				PayerID:       "payer_001",
				Amount:        50.0,
			},
			setup: func() {
				// Insert a pending transaction
				transaction := &model.Transaction{
					TransactionID:   "txn_12350",
					PayerID:         "payer_001",
					PayeeID:         "payee_002",
					Amount:          50.0,
					TransactionType: "Debit",
					Status:          "Pending",
					CreatedAt:       time.Now(),
					UpdatedAt:       time.Now(),
					ReservedAmount:  50.0,
				}
				db.Create(transaction)
			},
			expects: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			err := ValidatePayerBalance(tt.req)
			if (err == nil && tt.expects != nil) || (err != nil && tt.expects == nil) || (err != nil && tt.expects != nil && err.Error() != tt.expects.Error()) {
				t.Errorf("Test %s failed: expected %v, got %v", tt.name, tt.expects, err)
			}
		})
	}
}

func BenchmarkValidateCardNumber(b *testing.B) {
	cardNumber := "1234567812345678"
	paymentMethodCardDetails := "1234567812345678"

	for i := 0; i < b.N; i++ {
		_ = validateCardNumber(cardNumber, paymentMethodCardDetails)
	}
}
