package customizedpackages

import (
	"jamesmukumu/maasaimaratripsportal/models/users"
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type CustomizedPackage struct {
	gorm.Model
	PackageTitle          string `json:"packageTitle" gorm:"packageTitle"`
	PackagePhoto          string `gorm:"not null"`
	PackageCategoryID     uint
	PackageItenerary      datatypes.JSON
	PackagePhotos         datatypes.JSON `gorm:"not null"`
	MeansOfTransport      string         `gorm:"not null"`
	PackageOverview       string         `gorm:"not null"`
	Inclusions            datatypes.JSON
	Exclusions            datatypes.JSON
	SpecialNotes          string `gorm:"not null"`
	UserID                uint
	User                  users.User `references:"id"`
	Published             bool       `gorm:"default:true"`
	PackageSlug           string     `gorm:"not null;unique"`
	PackageCharge         int        `gorm:"not null"`
	PackageChargeCurrency string
    ClosesOn time.Time `gorm:"not null"`
    ClientsDetail datatypes.JSON 
}


type Client struct {
FirstName string `json:"firstName"`
LastName string `json:"lastName"`
Email string `json:"email"`
PhoneNumber string `json:"phoneNumber"`
}

func(pack *CustomizedPackage)ExpirySetter(){
pack.ClosesOn=  time.Now().AddDate(0,0,7)
}