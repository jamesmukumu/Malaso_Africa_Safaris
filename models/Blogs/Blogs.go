package blogs

import (
	"jamesmukumu/maasaimaratripsportal/models/users"

	"gorm.io/gorm"
)

type BlogsCategory struct {
	gorm.Model
	Title string `gorm:"not null"`
	Description string `gorm:"not null"`
	UserID uint
	User users.User `references:"id"`
}


type Blogs struct{
gorm.Model
BlogTitle string `gorm:"not null"`
BlogDescription string `gorm:"not null"`
BlogsCategoryID uint
BlogsCategory BlogsCategory `references:"id"`
BlogOverview string `gorm:"not null"`
UserID uint
User users.User `references:"id"`
Slug string `gorm:"unique"`
Published bool `gorm:"default:false"`
BlogPhoto string `gorm:"not null"`
}