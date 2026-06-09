package packages

import (
	destinations "jamesmukumu/maasaimaratripsportal/models/Destinations"
	packagescategory "jamesmukumu/maasaimaratripsportal/models/PackagesCategory"
	"jamesmukumu/maasaimaratripsportal/models/users"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Package struct {
gorm.Model
PackageTitle string `json:"packageTitle" gorm:"packageTitle"`
PackagePhoto string `gorm:"not null"`
PackageCategoryID uint
PackageCategory packagescategory.PackageCategory `gorm:"references:id"`
PackageItenerary datatypes.JSON
PackagePhotos datatypes.JSON `gorm:"not null"`
MeansOfTransport string `gorm:"not null"`
PackageOverview string `gorm:"not null"`
Inclusions datatypes.JSON
Exclusions datatypes.JSON
SpecialNotes string `gorm:"not null"`
DestinationID uint
Destination destinations.Destination `references:"id"`
UserID uint
User users.User `references:"id"`
Published bool `gorm:"default:false"`
PackageSlug string `gorm:"not null;unique"`
PackageCharge int `gorm:"not null"`
PackageChargeCurrency string
StartValidity string `gorm:"not null"`
EndValidity string `gorm:"not null"`
}