package hotels

import (
	destinations "jamesmukumu/maasaimaratripsportal/models/Destinations"
	"jamesmukumu/maasaimaratripsportal/models/users"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Hotels struct {
gorm.Model
HotelName string `gorm:"not null'" json:"hotelName"`
HotelOverview string `json:"hotelOverview" gorm:"not null"`
HotelDescription string `gorm:"not null" json:"hotelDescription"`
Published bool `gorm:"default:false"`
CancellationPolicy string `gorm:"not null" json:"cancellationPolicy"`
LocationDescription string `gorm:"not null" json:"locationDescription"`
HotelCoordinates datatypes.JSON `json:"hotelCoordinates"`
HotelContactEmail string `gorm:"not null" json:"hotelEmail"`
HotelContactPhonenumber string `json:"hotelContactPhoneNumber" gorm:"not null"`
Slug string `json:"slug"`
Ratings int8 `gorm:"not null" json:"ratings"`
DestinationID uint
Destination destinations.Destination `references:"id"`
UserID uint
User users.User `references:"id"`
ContactPerson string `json:"contactPerson"`
HotelThumbnail string `gorm:"not null"`
HotelPhotos datatypes.JSON `gorm:"not null"`
HotelRates datatypes.JSON `json:"hotelRates" gorm:"not null"`
HotelContacts datatypes.JSON `json:"hotelContacts" gorm:"not null"`
} 