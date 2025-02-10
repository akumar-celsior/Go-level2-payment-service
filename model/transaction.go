package model

import "time"

type AuditLog struct {
	AuditLogId    string    `gorm:"primaryKey;column=audit_log_id;size:36"`
	TransactionID string    `json:"transaction_id"`
	Action        string    `json:"action"`
	CreatedAt     time.Time `gorm:"column=created_at;autoCreateTime"`
	Details       string    `json:"details"`
}

type Payee struct {
	PayeeID   string    `gorm:"primaryKey;column=payee_id;size:36"` // Primary key
	Name      string    `gorm:"column=name;size:36"`                // Payee's name
	Email     string    `gorm:"column=email;size:36"`               // Payee's email
	Address   string    `gorm:"column=address;size:50"`             // Payee's address (optional)
	Balance   float64   `gorm:"column=balance;not null"`            // Cannot be null
	Status    string    `gorm:"column=status;size:50"`              // Status (active/inactive/suspended)
	CreatedAt time.Time `gorm:"column=created_at;autoCreateTime"`   // Auto-set at creation
	UpdatedAt time.Time `gorm:"column=updated_at;autoUpdateTime"`   // Auto-set on update
}

type PaymentMethod struct {
	PaymentMethodID string    `gorm:"primaryKey;column=payment_method_id;size:36"` // Primary key
	PayerID         string    `gorm:"column=payer_id;size:36;not null"`            // Foreign key to Payer table
	MethodType      string    `gorm:"column=method_type;size:36"`                  // Payment method type (e.g., card, wallet)
	Details         string    `gorm:"column=details;size:50"`                      // Masked or tokenized payment details
	ExpiryDate      string    `gorm:"column=expiry_date;size:50"`                  // Expiry date of the payment method
	CVV             string    `gorm:"column=cvv;size:50"`                          // CVV (if required)
	Status          string    `gorm:"column=status;size:50"`                       // Status of the payment method
	CreatedAt       time.Time `gorm:"column=created_at;autoCreateTime"`            // Automatically set at creation
	UpdatedAt       time.Time `gorm:"column=updated_at;autoUpdateTime"`            // Automatically set on update
}

type Payer struct {
	PayerID         string    `gorm:"primaryKey;column=payer_id;size:36"` // Primary key
	Name            string    `gorm:"column=name;size:36"`
	Email           string    `gorm:"column=email;size:36"`
	PhoneNumber     string    `gorm:"column=phone_number;size:36"`      // Optional field
	Address         string    `gorm:"column=address;size:50"`           // Optional field
	PaymentMethodID string    `gorm:"column=payment_method_id;size:50"` // Optional field
	Balance         float64   `gorm:"column=balance;not null"`          // Cannot be null
	Status          string    `gorm:"column=status;size:36"`            // Status (e.g., active/inactive)
	CreatedAt       time.Time `gorm:"column=created_at;autoCreateTime"` // Automatically set at creation
	UpdatedAt       time.Time `gorm:"column=updated_at;autoUpdateTime"` // Automatically set on update
}

type Request struct {
	TransactionID   string            `json:"transaction_id"`
	PayerID         string            `json:"payer_id"`
	PayeeID         string            `json:"payee_id"`
	Amount          float64           `json:"amount"`
	TransactionType string            `json:"transaction_type"`
	PaymentMethodID string            `json:"payment_method_id"`
	Details         map[string]string `json:"details"`
}

type Transaction struct {
	TransactionID   string    `gorm:"primaryKey;size:36"` // Unique identifier for the transaction
	PayerID         string    `gorm:"size:36"`            // Payer's ID (foreign key to Payer table)
	PayeeID         string    `gorm:"size:36"`            // Payee's ID (foreign key to Payee table)
	Amount          float64   `gorm:"not null"`           // Amount to be transacted
	TransactionType string    `gorm:"size:50"`            // Type of transaction (Debit, Credit, Refund)
	Status          string    `gorm:"size:50"`            // Status of the transaction (Pending, Completed, Failed, Reserved)
	CreatedAt       time.Time `gorm:"autoCreateTime"`     // Timestamp for when the transaction is created
	UpdatedAt       time.Time `gorm:"autoUpdateTime"`     // Timestamp for when the transaction is updated
	ReservedAmount  float64   // Reserved balance for the transaction
}
