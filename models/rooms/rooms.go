package rooms

import (
	"jamesmukumu/maasaimaratripsportal/models/hotels"
	"jamesmukumu/maasaimaratripsportal/models/users"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Rooms struct {
	gorm.Model
	RoomName string `json:"roomName" gorm:"not null"`
	RoomOverview string `gorm:"not null"`
	RoomDescription string `gorm:"not null"`
    RoomPhoto string
    RoomPhotos datatypes.JSON `gorm:"not null"`
	RoomRates datatypes.JSON
	RoomMeals datatypes.JSON
	HotelsID uint
	Hotels hotels.Hotels `references:"id"`
	UserID uint
	User users.User `references:"id"`
}