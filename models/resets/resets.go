package resets

import (
	"jamesmukumu/maasaimaratripsportal/models/users"

	"gorm.io/gorm"
)

type Resets struct {
	gorm.Model
	ResetID string `gorm:"unique;not null"`
	Email string `gorm:"not null"`
	UserID uint `gorm:"not null"`
    User users.User `references:"id"`
}