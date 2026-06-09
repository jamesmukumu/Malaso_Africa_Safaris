package email

import (
	"jamesmukumu/maasaimaratripsportal/models/users"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type EmailTemplates struct {
gorm.Model
Title string `gorm:"not null;unique"`
Body string `gorm:"not null"`
UserID uint
User users.User `references:"id"`
Subject string `gorm:"not null"`
}


type Newsletter struct{
gorm.Model
Title string `gorm:"not null;unique"`
Body datatypes.JSON
UserID uint
User users.User `references:"id"`
Subject string
}

type EmailDelivery struct{

	
}