package adventures

import (
	destinations "jamesmukumu/maasaimaratripsportal/models/Destinations"

	"jamesmukumu/maasaimaratripsportal/models/users"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type AdventuresCategory struct {
gorm.Model
AdventureCategoryName string `gorm:"not null;unique"`
AventureCategoryDescription string `gorm:"not null"`
UserID uint
User users.User `references:"id"`
AdventureCategoryImage string `gorm:"not null"`
AdventureSlug string
}

type Adventure struct{
gorm.Model
Name string `gorm:"not null;unique"`
Slug string `gorm:"not null;unique"`
Published bool `gorm:"default:false"`
DestinationID uint
Destination destinations.Destination `references:"id"`
UserID uint
AdventuresCategoryID uint
AdventuresCategory AdventuresCategory `references:"id"`
User users.User `references:"id"`
AdventureOverview string `gorm:"not null"`
AdventureDescription string `gorm:"not null"`
AdventurePhoto string
AdventurePhotos datatypes.JSON `gorm:"not null"`
AdventureInclusions string `gorm:"not null"`
AdventureExclusions string `gorm:"not null"`
AdventureCharges int `gorm:"not null"`
AdventureChargeCurrency string
SpecialNotes string
ModeOfTransport string
StartValidity string
EndValidty string
MaximumCapacity int
Opened bool `gorm:"default:true"`
}