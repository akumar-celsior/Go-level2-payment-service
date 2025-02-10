package model

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAuditLog(t *testing.T) {
	auditLog := AuditLog{
		AuditLogId:    "1",
		TransactionID: "123",
		Action:        "CREATE",
		CreatedAt:     time.Now(),
		Details:       "Transaction created",
	}

	assert.Equal(t, "1", auditLog.AuditLogId)
	assert.Equal(t, "123", auditLog.TransactionID)
	assert.Equal(t, "CREATE", auditLog.Action)
	assert.NotNil(t, auditLog.CreatedAt)
	assert.Equal(t, "Transaction created", auditLog.Details)
}

func TestPayee(t *testing.T) {
	payee := Payee{
		PayeeID:   "1",
		Name:      "John Doe",
		Email:     "john.doe@example.com",
		Address:   "123 Main St",
		Balance:   100.0,
		Status:    "active",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	assert.Equal(t, "1", payee.PayeeID)
	assert.Equal(t, "John Doe", payee.Name)
	assert.Equal(t, "john.doe@example.com", payee.Email)
	assert.Equal(t, "123 Main St", payee.Address)
	assert.Equal(t, 100.0, payee.Balance)
	assert.Equal(t, "active", payee.Status)
	assert.NotNil(t, payee.CreatedAt)
	assert.NotNil(t, payee.UpdatedAt)
}

func TestPaymentMethod(t *testing.T) {
	paymentMethod := PaymentMethod{
		PaymentMethodID: "1",
		PayerID:         "123",
		MethodType:      "card",
		Details:         "**** **** **** 1234",
		ExpiryDate:      "12/23",
		CVV:             "123",
		Status:          "active",
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	assert.Equal(t, "1", paymentMethod.PaymentMethodID)
	assert.Equal(t, "123", paymentMethod.PayerID)
	assert.Equal(t, "card", paymentMethod.MethodType)
	assert.Equal(t, "**** **** **** 1234", paymentMethod.Details)
	assert.Equal(t, "12/23", paymentMethod.ExpiryDate)
	assert.Equal(t, "123", paymentMethod.CVV)
	assert.Equal(t, "active", paymentMethod.Status)
	assert.NotNil(t, paymentMethod.CreatedAt)
	assert.NotNil(t, paymentMethod.UpdatedAt)
}

func TestPayer(t *testing.T) {
	payer := Payer{
		PayerID:         "1",
		Name:            "Jane Doe",
		Email:           "jane.doe@example.com",
		PhoneNumber:     "123-456-7890",
		Address:         "456 Main St",
		PaymentMethodID: "1",
		Balance:         200.0,
		Status:          "active",
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	assert.Equal(t, "1", payer.PayerID)
	assert.Equal(t, "Jane Doe", payer.Name)
	assert.Equal(t, "jane.doe@example.com", payer.Email)
	assert.Equal(t, "123-456-7890", payer.PhoneNumber)
	assert.Equal(t, "456 Main St", payer.Address)
	assert.Equal(t, "1", payer.PaymentMethodID)
	assert.Equal(t, 200.0, payer.Balance)
	assert.Equal(t, "active", payer.Status)
	assert.NotNil(t, payer.CreatedAt)
	assert.NotNil(t, payer.UpdatedAt)
}

func TestRequest(t *testing.T) {
	details := map[string]string{"key": "value"}
	request := Request{
		TransactionID:   "1",
		PayerID:         "123",
		PayeeID:         "456",
		Amount:          50.0,
		TransactionType: "Debit",
		PaymentMethodID: "1",
		Details:         details,
	}

	assert.Equal(t, "1", request.TransactionID)
	assert.Equal(t, "123", request.PayerID)
	assert.Equal(t, "456", request.PayeeID)
	assert.Equal(t, 50.0, request.Amount)
	assert.Equal(t, "Debit", request.TransactionType)
	assert.Equal(t, "1", request.PaymentMethodID)
	assert.Equal(t, details, request.Details)
}

func TestTransaction(t *testing.T) {
	transaction := Transaction{
		TransactionID:   "1",
		PayerID:         "123",
		PayeeID:         "456",
		Amount:          75.0,
		TransactionType: "Credit",
		Status:          "Completed",
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
		ReservedAmount:  25.0,
	}

	assert.Equal(t, "1", transaction.TransactionID)
	assert.Equal(t, "123", transaction.PayerID)
	assert.Equal(t, "456", transaction.PayeeID)
	assert.Equal(t, 75.0, transaction.Amount)
	assert.Equal(t, "Credit", transaction.TransactionType)
	assert.Equal(t, "Completed", transaction.Status)
	assert.NotNil(t, transaction.CreatedAt)
	assert.NotNil(t, transaction.UpdatedAt)
	assert.Equal(t, 25.0, transaction.ReservedAmount)
}
