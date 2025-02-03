// models/product.go
package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Product struct {
    ID          uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
    Name        string    `json:"name" gorm:"not null"`
    Description string    `json:"description"`
    Price       float64   `json:"price" gorm:"not null"`
		CreatedAt time.Time      `json:"created_at"` // Otomatis diisi saat pertama kali dibuat
		UpdatedAt time.Time      `json:"updated_at"` // Diupdate otomatis oleh GORM
		DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"` // Soft delete
}

func (product *Product) BeforeCreate(tx *gorm.DB) (err error) {
    product.ID = uuid.New() // Generate a new UUID for the product
    return
}