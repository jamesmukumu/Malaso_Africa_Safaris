package destinations

import (

	"jamesmukumu/maasaimaratripsportal/models/users"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Destination struct {
gorm.Model
DestinationTitle string `gorm:"not null"`
DestinationOverview string `gorm:"not null"`
DestinationDescription string `gorm:"not null"`
DestinationPhoto string `gorm:"not null"`
DestinationPhotos datatypes.JSON
UserID uint
Published bool `gorm:"not null;default:false"`
User    users.User  `gorm:"references:id"`
DestinationSlug string `gorm:"unique"`
}