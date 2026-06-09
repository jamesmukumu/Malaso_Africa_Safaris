package wallets

import (
	"jamesmukumu/maasaimaratripsportal/models/users"
	

	"gorm.io/gorm"
)

type Wallet struct {
gorm.Model
WalletLabel string `gorm:"not null;" json:"walletLabel"`
WalletBalance int8 `gorm:"not null"`
UserID uint 
User users.User 
}