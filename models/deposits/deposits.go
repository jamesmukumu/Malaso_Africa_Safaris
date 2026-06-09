package deposits

import (
	"jamesmukumu/maasaimaratripsportal/models/users"

	"gorm.io/gorm"
)

type Deposits struct {
	gorm.Model
	Invoice_Id string `gorm:"not null;unique" json:"invoice_id"`
    Amount string `json:"amount" gorm:"not null"`
	CustomersName string `json:"name" gorm:"not null"`
	CustomersPhoneNumber string `json:"phone" gorm:"not null"`
	UserID uint 
	User users.User 

}