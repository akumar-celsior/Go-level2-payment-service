package controller

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"testing"
	"time"

	"goPocDemo/initializer"
	"goPocDemo/model"

	spannergorm "github.com/googleapis/go-gorm-spanner"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/httptest"
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
}

func TestDeductAmountFromPayer(t *testing.T) {
	dbConnection()
	initializer.DbInstance = db

	// Create a test payer
	payer := &model.Payer{
		PayerID:         "payer1",
		Name:            "test1",
		Email:           "ak@123",
		PhoneNumber:     "123",
		Address:         "newadd",
		PaymentMethodID: "23",
		Balance:         100.0,
		Status:          "active",
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
	// /  **Cleanup before each test case**
	var count int64
	db.Model(&model.Payer{}).Where("payer_id = ?", "payer1").Count(&count)
	log.Printf("🔍 Debug: Found %d records for payer_id = payer1", count)
	if count > 0 {
		db.Where("payer_id = ?", "payer1").Delete(&model.Payer{})
	}
	db.Create(payer)

	tests := []struct {
		name            string
		payerID         string
		requestAmount   float64
		expectedError   error
		expectedBalance float64
	}{
		{
			name:            "Successful deduction",
			payerID:         "payer1",
			requestAmount:   50.0,
			expectedError:   nil,
			expectedBalance: 50.0,
		},
		{
			name:            "Insufficient balance",
			payerID:         "payer1",
			requestAmount:   150.0,
			expectedError:   nil,
			expectedBalance: -100.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := deductAmountFromPayer(tt.payerID, tt.requestAmount)
			if tt.expectedError != nil {
				assert.EqualError(t, err, tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}

			// Verify the balance
			var updatedPayer model.Payer
			db.Where("payer_id = ?", tt.payerID).First(&updatedPayer)
			assert.Equal(t, tt.expectedBalance, updatedPayer.Balance)
		})
	}
}
func TestCreditAmountToPayee(t *testing.T) {
	dbConnection()
	initializer.DbInstance = db

	// Create a test payee
	payee := &model.Payee{
		PayeeID:   "payee1",
		Name:      "test1",
		Email:     "ak@123",
		Address:   "newadd",
		Balance:   100.0,
		Status:    "active",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	// **Cleanup before each test case**
	var count int64
	db.Model(&model.Payee{}).Where("payee_id = ?", "payee1").Count(&count)
	log.Printf("🔍 Debug: Found %d records for payee_id = payee1", count)
	if count > 0 {
		db.Where("payee_id = ?", "payee1").Delete(&model.Payee{})
	}
	db.Create(payee)

	tests := []struct {
		name            string
		payeeID         string
		requestAmount   float64
		expectedError   error
		expectedBalance float64
	}{
		{
			name:            "Successful credit",
			payeeID:         "payee1",
			requestAmount:   50.0,
			expectedError:   nil,
			expectedBalance: 150.0,
		},
		{
			name:            "Payee not found",
			payeeID:         "nonexistent",
			requestAmount:   50.0,
			expectedError:   errors.New("payee not found"),
			expectedBalance: 100.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := creditAmountToPayee(tt.payeeID, tt.requestAmount)
			if tt.expectedError != nil {
				assert.EqualError(t, err, tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}

			// Verify the balance
			var updatedPayee model.Payee
			db.Where("payee_id = ?", tt.payeeID).First(&updatedPayee)
			if tt.expectedError == nil {
				assert.Equal(t, tt.expectedBalance, updatedPayee.Balance)
			}
		})
	}
}
func TestReserveTransactionAmount(t *testing.T) {
	dbConnection()
	initializer.DbInstance = db

	// Create a test transaction
	transaction := &model.Transaction{
		TransactionID:   "txn1",
		PayerID:         "234",
		PayeeID:         "34",
		Amount:          100.0,
		TransactionType: "Debit",
		Status:          "Pending",
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
		ReservedAmount:  100.0,
	}
	// **Cleanup before each test case**
	var count int64
	db.Model(&model.Transaction{}).Where("transaction_id = ?", "txn1").Count(&count)
	log.Printf("🔍 Debug: Found %d records for transaction_id = txn1", count)
	if count > 0 {
		db.Where("transaction_id = ?", "txn1").Delete(&model.Transaction{})
	}
	db.Create(transaction)

	tests := []struct {
		name            string
		transactionID   string
		amount          float64
		expectedError   error
		expectedStatus  string
		expectedReserve float64
	}{
		{
			name:            "Successful reservation",
			transactionID:   "txn1",
			amount:          50.0,
			expectedError:   nil,
			expectedStatus:  "Reserved",
			expectedReserve: 50.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := reserveTransactionAmount(tt.transactionID, tt.amount)
			if tt.expectedError != nil {
				assert.EqualError(t, err, tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}

			// Verify the transaction status and reserved amount
			var updatedTransaction model.Transaction
			db.Where("transaction_id = ?", tt.transactionID).First(&updatedTransaction)
			if tt.expectedError == nil {
				assert.Equal(t, tt.expectedStatus, updatedTransaction.Status)
				assert.Equal(t, tt.expectedReserve, updatedTransaction.ReservedAmount)
			}
		})
	}
}
func TestHandleDebitFunctionality(t *testing.T) {
	dbConnection()
	initializer.DbInstance = db

	// Mock the context
	//app := iris.New()
	ctx := new(iris.Context)
	// ctx.Reset(httptest.NewRecorder(), nil)
	// defer app.ReleaseContext(ctx)

	// Create a test payer
	payer := &model.Payer{
		PayerID:         "payer1",
		Name:            "test1",
		Email:           "ak@123",
		PhoneNumber:     "123",
		Address:         "newadd",
		PaymentMethodID: "23",
		Balance:         100.0,
		Status:          "active",
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
	// Create a test payee
	payee := &model.Payee{
		PayeeID:   "payee1",
		Name:      "test1",
		Email:     "ak@123",
		Address:   "newadd",
		Balance:   100.0,
		Status:    "active",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	// Create a test transaction
	transaction := &model.Transaction{
		TransactionID:   "txn1",
		PayerID:         "payer1",
		PayeeID:         "payee1",
		Amount:          50.0,
		TransactionType: "Debit",
		Status:          "Pending",
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	// Cleanup before each test case
	db.Where("payer_id = ?", "payer1").Delete(&model.Payer{})
	db.Where("payee_id = ?", "payee1").Delete(&model.Payee{})
	db.Where("transaction_id = ?", "txn1").Delete(&model.Transaction{})

	db.Create(payer)
	db.Create(payee)
	db.Create(transaction)

	tests := []struct {
		name          string
		req           model.Request
		expectedError error
	}{
		{
			name: "Successful debit transaction",
			req: model.Request{
				TransactionID:   "txn1",
				PayerID:         "payer1",
				PayeeID:         "payee1",
				Amount:          50.0,
				TransactionType: "Debit",
			},
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := handleDebitFunctionality(*ctx, tt.req)
			if tt.expectedError != nil {
				assert.EqualError(t, err, tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}

			// Verify the payer balance
			var updatedPayer model.Payer
			db.Where("payer_id = ?", tt.req.PayerID).First(&updatedPayer)
			if tt.expectedError == nil {
				assert.Equal(t, payer.Balance-tt.req.Amount, updatedPayer.Balance)
			}

			// Verify the payee balance
			var updatedPayee model.Payee
			db.Where("payee_id = ?", tt.req.PayeeID).First(&updatedPayee)
			if tt.expectedError == nil {
				assert.Equal(t, payee.Balance+tt.req.Amount, updatedPayee.Balance)
			}

			// Verify the transaction status
			var updatedTransaction model.Transaction
			db.Where("transaction_id = ?", tt.req.TransactionID).First(&updatedTransaction)
			if tt.expectedError == nil {
				assert.Equal(t, "Reserved", updatedTransaction.Status)
				assert.Equal(t, tt.req.Amount, updatedTransaction.ReservedAmount)
			}
		})
	}
}

func TestSendErrorResponse(t *testing.T) {
	// Initialize Iris App
	app := iris.New()
	app.Get("/error", func(ctx iris.Context) {
		sendErrorResponse(ctx, http.StatusBadRequest, "Invalid request")
	})

	// Create test instance
	e := httptest.New(t, app)

	// Perform GET request to trigger the error response
	resp := e.GET("/error").Expect().Status(http.StatusBadRequest).Body().Raw()

	// Parse JSON response
	var result map[string]string
	err := json.Unmarshal([]byte(resp), &result)
	assert.NoError(t, err)

	// Validate response content
	assert.Contains(t, result, "error")
	assert.Equal(t, "Invalid request", result["error"])
}
func TestHandleCreditFunctionality(t *testing.T) {
	dbConnection()
	initializer.DbInstance = db

	// Mock the context
	app := iris.New()
	ctx := app.ContextPool.Acquire(httptest.NewRecorder(), nil)
	defer app.ContextPool.Release(ctx)

	// Create a test payee
	payee := &model.Payee{
		PayeeID:   "payee1",
		Name:      "test1",
		Email:     "ak@123",
		Address:   "newadd",
		Balance:   100.0,
		Status:    "active",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	// Create a test transaction
	transaction := &model.Transaction{
		TransactionID:   "txn1",
		PayerID:         "payer1",
		PayeeID:         "payee1",
		Amount:          50.0,
		TransactionType: "Credit",
		Status:          "Pending",
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
		ReservedAmount:  50.0,
	}

	// Cleanup before each test case
	db.Where("payee_id = ?", "payee1").Delete(&model.Payee{})
	db.Where("transaction_id = ?", "txn1").Delete(&model.Transaction{})

	db.Create(payee)
	db.Create(transaction)

	tests := []struct {
		name          string
		req           model.Request
		expectedError error
	}{
		{
			name: "Successful credit transaction",
			req: model.Request{
				TransactionID:   "txn1",
				PayerID:         "payer1",
				PayeeID:         "payee1",
				Amount:          50.0,
				TransactionType: "Credit",
				PaymentMethodID: "23",
				Details:         map[string]string{"card_number": "1234567890123456"},
			},
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := handleCreditFunctionality(ctx, tt.req)
			if tt.expectedError != nil {
				assert.EqualError(t, err, tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}

			// Verify the payee balance
			var updatedPayee model.Payee
			db.Where("payee_id = ?", tt.req.PayeeID).First(&updatedPayee)
			if tt.expectedError == nil {
				assert.Equal(t, payee.Balance+tt.req.Amount, updatedPayee.Balance)
			}

			// Verify the transaction status
			var updatedTransaction model.Transaction
			db.Where("transaction_id = ?", tt.req.TransactionID).First(&updatedTransaction)
			if tt.expectedError == nil {
				assert.Equal(t, "Reserved", updatedTransaction.Status)
				assert.Equal(t, tt.req.Amount, updatedTransaction.ReservedAmount)
			}
		})
	}
}
func TestHandleRefundFunctionality(t *testing.T) {
	dbConnection()
	initializer.DbInstance = db

	// Create a test payer
	payer := &model.Payer{
		PayerID:         "payer1",
		Name:            "test1",
		Email:           "ak@123",
		PhoneNumber:     "123",
		Address:         "newadd",
		PaymentMethodID: "23",
		Balance:         100.0,
		Status:          "active",
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
	// Create a test payee
	payee := &model.Payee{
		PayeeID:   "payee1",
		Name:      "test1",
		Email:     "ak@123",
		Address:   "newadd",
		Balance:   100.0,
		Status:    "active",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	// Create a test transaction
	transaction := &model.Transaction{
		TransactionID:   "txn1",
		PayerID:         "payer1",
		PayeeID:         "payee1",
		Amount:          50.0,
		TransactionType: "Debit",
		Status:          "Completed",
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
		ReservedAmount:  50.0,
	}

	// Cleanup before each test case
	db.Where("payer_id = ?", "payer1").Delete(&model.Payer{})
	db.Where("payee_id = ?", "payee1").Delete(&model.Payee{})
	db.Where("transaction_id = ?", "txn1").Delete(&model.Transaction{})

	db.Create(payer)
	db.Create(payee)
	db.Create(transaction)

	tests := []struct {
		name          string
		req           model.Request
		expectedError error
	}{
		{
			name: "Successful refund transaction",
			req: model.Request{
				TransactionID:   "txn1",
				PayerID:         "payer1",
				PayeeID:         "payee1",
				Amount:          50.0,
				TransactionType: "Refund",
				PaymentMethodID: "23",
				Details:         map[string]string{"card_number": "1234567890123456"},
			},
			expectedError: nil,
		},
		{
			name: "Transaction not found",
			req: model.Request{
				TransactionID:   "nonexistent",
				PayerID:         "payer1",
				PayeeID:         "payee1",
				Amount:          50.0,
				TransactionType: "Refund",
				PaymentMethodID: "23",
				Details:         map[string]string{"card_number": "1234567890123456"},
			},
			expectedError: errors.New("transaction not found"),
		},
		{
			name: "Original transaction already refunded",
			req: model.Request{
				TransactionID:   "txn1",
				PayerID:         "payer1",
				PayeeID:         "payee1",
				Amount:          50.0,
				TransactionType: "Refund",
				PaymentMethodID: "23",
				Details:         map[string]string{"card_number": "1234567890123456"},
			},
			expectedError: errors.New("original transaction is already refunded"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := handleRefundFunctionality(tt.req)
			if tt.expectedError != nil {
				assert.EqualError(t, err, tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}

			if tt.expectedError == nil {
				// Verify the payer balance
				var updatedPayer model.Payer
				db.Where("payer_id = ?", tt.req.PayerID).First(&updatedPayer)
				assert.Equal(t, payer.Balance+tt.req.Amount, updatedPayer.Balance)

				// Verify the payee balance
				var updatedPayee model.Payee
				db.Where("payee_id = ?", tt.req.PayeeID).First(&updatedPayee)
				assert.Equal(t, payee.Balance-tt.req.Amount, updatedPayee.Balance)

				// Verify the transaction type
				var updatedTransaction model.Transaction
				db.Where("transaction_id = ?", tt.req.TransactionID).First(&updatedTransaction)
				assert.Equal(t, "Refund", updatedTransaction.TransactionType)
			}
		})
	}
}
func TestHandleTransactionTypes(t *testing.T) {
	dbConnection()
	initializer.DbInstance = db

	// Mock the context
	app := iris.New()
	ctx := app.ContextPool.Acquire(httptest.NewRecorder(), nil)
	defer app.ContextPool.Release(ctx)

	// Create a test payer
	payer := &model.Payer{
		PayerID:         "payer1",
		Name:            "test1",
		Email:           "ak@123",
		PhoneNumber:     "123",
		Address:         "newadd",
		PaymentMethodID: "23",
		Balance:         100.0,
		Status:          "active",
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
	// Create a test payee
	payee := &model.Payee{
		PayeeID:   "payee1",
		Name:      "test1",
		Email:     "ak@123",
		Address:   "newadd",
		Balance:   100.0,
		Status:    "active",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	// Create a test transaction
	transaction := &model.Transaction{
		TransactionID:   "txn1",
		PayerID:         "payer1",
		PayeeID:         "payee1",
		Amount:          50.0,
		TransactionType: "Debit",
		Status:          "Pending",
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
		ReservedAmount:  50.0,
	}

	// Cleanup before each test case
	db.Where("payer_id = ?", "payer1").Delete(&model.Payer{})
	db.Where("payee_id = ?", "payee1").Delete(&model.Payee{})
	db.Where("transaction_id = ?", "txn1").Delete(&model.Transaction{})

	db.Create(payer)
	db.Create(payee)
	db.Create(transaction)

	tests := []struct {
		name           string
		req            model.Request
		expectedError  error
		expectedStatus string
	}{
		{
			name: "Successful debit transaction",
			req: model.Request{
				TransactionID:   "txn1",
				PayerID:         "payer1",
				PayeeID:         "payee1",
				Amount:          50.0,
				TransactionType: "Debit",
				PaymentMethodID: "23",
				Details:         map[string]string{"card_number": "1234567890123456"},
			},
			expectedError:  nil,
			expectedStatus: "Reserved",
		},
		{
			name: "Successful credit transaction",
			req: model.Request{
				TransactionID:   "txn1",
				PayerID:         "payer1",
				PayeeID:         "payee1",
				Amount:          50.0,
				TransactionType: "Credit",
				PaymentMethodID: "23",
				Details:         map[string]string{"card_number": "1234567890123456"},
			},
			expectedError:  nil,
			expectedStatus: "Reserved",
		},
		{
			name: "Successful refund transaction",
			req: model.Request{
				TransactionID:   "txn1",
				PayerID:         "payer1",
				PayeeID:         "payee1",
				Amount:          50.0,
				TransactionType: "Refund",
				PaymentMethodID: "23",
				Details:         map[string]string{"card_number": "1234567890123456"},
			},
			expectedError:  nil,
			expectedStatus: "Reserved",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := handleTransactionTypes(ctx, tt.req)
			if tt.expectedError != nil {
				assert.EqualError(t, err, tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}

			// Verify the transaction status
			var updatedTransaction model.Transaction
			db.Where("transaction_id = ?", tt.req.TransactionID).First(&updatedTransaction)
			if tt.expectedError == nil {
				assert.Equal(t, tt.expectedStatus, updatedTransaction.Status)
			}
		})
	}
}
func TestLogAndRespond(t *testing.T) {
	// Initialize Iris App
	app := iris.New()
	app.Get("/logAndRespond", func(ctx iris.Context) {
		logAndRespond(ctx, "txn1", "Debit", "Test message", http.StatusBadRequest, "Invalid request")
	})

	// Create test instance
	e := httptest.New(t, app)

	// Perform GET request to trigger the logAndRespond function
	resp := e.GET("/logAndRespond").Expect().Status(http.StatusBadRequest).Body().Raw()

	// Parse JSON response
	var result map[string]string
	err := json.Unmarshal([]byte(resp), &result)
	assert.NoError(t, err)

	// Validate response content
	assert.Contains(t, result, "error")
	assert.Equal(t, "Invalid request", result["error"])

	// Here you would typically also verify that the audit log was created
	// Since this is a unit test, you might need to mock the services.CreateAuditLog function
}
func TestLogSuccessDetails(t *testing.T) {
	// Initialize Iris App
	app := iris.New()
	app.Post("/logSuccessDetails", func(ctx iris.Context) {
		req := model.Request{
			TransactionID:   "txn1",
			TransactionType: "Debit",
		}
		logSuceessDetails(ctx, req)
	})

	// Create test instance
	e := httptest.New(t, app)

	// Perform POST request to trigger the logSuceessDetails function
	resp := e.POST("/logSuccessDetails").Expect().Status(iris.StatusCreated).Body().Raw()

	// Parse JSON response
	var result map[string]string
	err := json.Unmarshal([]byte(resp), &result)
	assert.NoError(t, err)

	// Validate response content
	assert.Equal(t, "Completed", result["status"])
	assert.Equal(t, "txn1", result["transaction_id"])
	assert.Equal(t, "Transaction created successfully", result["message"])
}
func TestStartTransactionHandler(t *testing.T) {
	t.Skip("Skip the test")
	dbConnection()
	initializer.DbInstance = db

	// Mock the context
	app := iris.New()
	ctx := app.ContextPool.Acquire(httptest.NewRecorder(), httptest.NewRequest("POST", "/start-transaction", nil))
	defer app.ContextPool.Release(ctx)

	tests := []struct {
		name           string
		requestBody    model.Request
		expectedStatus int
		expectedError  string
	}{
		{
			name: "Invalid request body",
			requestBody: model.Request{
				TransactionID:   "",
				PayerID:         "",
				PayeeID:         "",
				Amount:          0.0,
				TransactionType: "",
				PaymentMethodID: "",
				Details:         nil,
			},
			expectedStatus: iris.StatusBadRequest,
			expectedError:  "",
		},
		{
			name: "Invalid payment details",
			requestBody: model.Request{
				TransactionID:   "txn2",
				PayerID:         "payer2",
				PayeeID:         "payee2",
				Amount:          50.0,
				TransactionType: "Debit",
				PaymentMethodID: "invalid",
				Details:         map[string]string{"card_number": "invalid"},
			},
			expectedStatus: iris.StatusBadRequest,
			expectedError:  "",
		},
		{
			name: "Failed to create transaction",
			requestBody: model.Request{
				TransactionID:   "txn3",
				PayerID:         "payer3",
				PayeeID:         "payee3",
				Amount:          50.0,
				TransactionType: "Debit",
				PaymentMethodID: "23",
				Details:         map[string]string{"card_number": "1234567890123456"},
			},
			expectedStatus: iris.StatusBadRequest,
			expectedError:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Prepare the request body
			body, _ := json.Marshal(tt.requestBody)
			ctx.Request().Body = ioutil.NopCloser(bytes.NewReader(body))

			// Call the handler
			StartTransactionHandler(ctx)

			// Verify the response status
			assert.Equal(t, tt.expectedStatus, ctx.ResponseWriter().StatusCode())

			// Verify the response body
			if tt.expectedError != "" {
				var result map[string]string
				err := json.Unmarshal(ctx.Recorder().Body(), &result)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedError, result["error"])
			}
		})
	}
}
