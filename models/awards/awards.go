package awards;
import (
	"jamesmukumu/maasaimaratripsportal/models/users"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Award struct {
	gorm.Model
	Title       string         `gorm:"unique;not null"`
	AwardThumbnailPhoto string `gorm:"not null;type:text"`
	Slug        string         `gorm:"unique;not null"`
	Inclusions  datatypes.JSON `gorm:"type:jsonb;not null"`
	Exclusions  datatypes.JSON `gorm:"type:jsonb;not null"`
	UserID      uint
	User        *users.User  `gorm:"references:id" json:"User,omitempty"`
	PointsNeeded int `gorm:"not null"`
	Published bool `gorm:"not null"`
	AwardDescription string `gorm:"not null"`  
}
