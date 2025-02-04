package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Bank struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	UserID    uuid.UUID      `gorm:"type:uuid;not null" json:"user_id"`
	BankName  string         `gorm:"not null" json:"bank_name"`
	AccountNo string         `gorm:"unique;not null" json:"account_no"`
	Nominal   float64        `gorm:"default:0" json:"nominal"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

// Hook before creating a bank account (generate UUID)
func (bank *Bank) BeforeCreate(tx *gorm.DB) (err error) {
	bank.ID = uuid.New()
	return nil
}
