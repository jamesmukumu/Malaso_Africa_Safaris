package bulkmails

import "gorm.io/gorm"

type BulkMails struct {
gorm.Model
UserName string `gorm:"null" json:"userName"`
Email string `gorm:"not null;unique" json:"email"`
PhoneNumber string `gorm:"unique;not null" json:"phoneNumber"`
}