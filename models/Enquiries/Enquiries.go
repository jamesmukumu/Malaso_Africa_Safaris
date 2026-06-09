package enquiries

import (

	"gorm.io/gorm"
)

type Enquiries struct {
gorm.Model
FirstName string `json:"firstName" gorm:"not null"`
LastName string `json:"lastName" gorm:"not null"`
Email string `json:"email" gorm:"not null"`
ContactPreference string `gorm:"default:email;not null"`
PhoneNumber string `gorm:"not null"`
AdultsCount int8
ChildrenCount int8
EnquiryDescription string 
StartStayDate string
EndStayDate string
Addressed bool `gorm:"default:false"`
RoomsCount int8
KidsAges string
Nationality string  `gorm:"not null"`
}