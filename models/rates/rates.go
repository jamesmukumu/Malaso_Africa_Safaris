package rates

import (
	"jamesmukumu/maasaimaratripsportal/models/hotels"
	"jamesmukumu/maasaimaratripsportal/models/users"

	"gorm.io/gorm"
)

type Rates struct {
	gorm.Model
	Title string `gorm:"not null"`
    Description string `gorm:"not null"`
    FilePath string `gorm:"not null"`
    RatesYear string `gorm:"not null"`
	HotelsID uint
	Hotels hotels.Hotels `references:"id"`
	Verified bool `gorm:"not null;default:false"`
    UserID uint
	User users.User `references:"id"`


}