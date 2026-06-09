package users

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name string `gorm:"unique;not null" json:"name"`
	Email string `gorm:"unique;not null" json:"email"`
	PhoneNumber string `gorm:"unique;not null" json:"phoneNumber"`
    Password string `gorm:"not null" json:"password"`
	Role  string `gorm:"not null" json:"role"`
	SuperUser bool `gorm:"not null" json:"superUser"`
	Org bool `gorm:"default:false"`
	Verified bool `gorm:"false"`
    OrgBanner string 
    OrgPrimaryColors string
	
}