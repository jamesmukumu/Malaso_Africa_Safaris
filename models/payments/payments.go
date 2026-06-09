package payments

import (
	enquiries "jamesmukumu/maasaimaratripsportal/models/Enquiries"

	"gorm.io/gorm"
)

type InitializedPayments struct {
gorm.Model
AuthorizationUrl string `gorm:"not null" json:"authorization_url"`
AccessCode string `gorm:"not null" json:"access_code"`
Reference string `gorm:"not null;" json:"reference"`
Email string `gorm:"not null" json:"email"`


EnquiriesID uint
Enquiries   enquiries.Enquiries `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
 

type CompletedPayment struct{
gorm.Model
TransactionID uint `gorm:"not null"`
Reference string `gorm:"not null"`
Amount int
GateWayResponse string
Currency string
InitializedPaymentsID uint 
InitializedPayments InitializedPayments `gorm:"references:id"` 
}