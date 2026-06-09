package packagescategory

import (
	"jamesmukumu/maasaimaratripsportal/models/users"

	"gorm.io/gorm"
)

type PackageCategory struct {
	gorm.Model
	PackageCategoryName string `json:"name"`
	UserID uint
	User    users.User  `gorm:"references:id"`
	PackageDescription string `json:"descp"`
	PackageCategoryImage string 
}