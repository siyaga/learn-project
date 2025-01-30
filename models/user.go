package models

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	Name      string         `json:"name"`
	Email     string         `json:"email" gorm:"unique;not null"`
	Password  string         `json:"password"`
	CreatedAt time.Time      `json:"created_at"` // Otomatis diisi saat pertama kali dibuat
	UpdatedAt time.Time      `json:"updated_at"` // Diupdate otomatis oleh GORM
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"` // Soft delete
}

// Hook sebelum user dibuat (generate UUID)
func (user *User) BeforeCreate(tx *gorm.DB) (err error) {
	user.ID = uuid.New()
	return nil
}

// Hash password sebelum disimpan
func (user *User) HashPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return err
	}
	user.Password = string(bytes)
	return nil
}

// Validasi password
func (user *User) CheckPassword(providedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(providedPassword))
}
