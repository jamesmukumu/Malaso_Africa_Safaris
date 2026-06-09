package attractions



import (
	destinations "jamesmukumu/maasaimaratripsportal/models/Destinations"
	"jamesmukumu/maasaimaratripsportal/models/users"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Attraction struct {
	gorm.Model
	AttractionName        string `gorm:"unique"`
	AttractionSlug        string `gorm:"unique"`
	AttractionOverview    string `gorm:"not null"`
	AttractionDescription string `gorm:"not null"`
	AttractionPhoto       string
	AttractionPhotos      datatypes.JSON
    DestinationID uint
    Destination destinations.Destination `references:"id"`
	UserID uint
    User users.User `references:"id"`
	Published bool `gorm:"false"`
}