package drafts

import (
	"jamesmukumu/maasaimaratripsportal/models/users"

	"time"

	"gorm.io/gorm"
)

type Draft struct {
gorm.Model
DraftTable string `gorm:"not null"`
DraftClosed bool `gorm:"default:false"`
DraftContent string `gorm:"not null"`
UserID uint
User    users.User  `gorm:"references:id"`
DraftCloses time.Time `gorm:"not null"`
}


func(draft *Draft)DraftExpiry(){
draft.DraftCloses = time.Now().AddDate(0,2,0)
}