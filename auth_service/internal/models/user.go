package models

import (
	"net/mail"
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        string    `json:"user_id" gorm:"type:uuid;primaryKey"`
	Email     string    `json:"email" gorm:"uniqueIndex;not null"`
	EncPass   string    `json:"enc_pass" gorm:"not null"`
	CreatedAt time.Time `json:"created_at" gorm:"type:timestamptz;default:now();not null;"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if !IsValid(u.Email) {
		return gorm.ErrCheckConstraintViolated
	}
	return nil
}

func IsValid(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}
