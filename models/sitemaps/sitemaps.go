package sitemaps

import (
	"time"

	"gorm.io/gorm"
)

type SiteMap struct {
gorm.Model
Filename string `gorm:"not null"`
FilePath string `gorm:"not null"`
IsActive bool `gorm:"not null;default:true"`
ValidTill time.Time `gorm:"not null"`
}


func (site *SiteMap)Validity(){
site.ValidTill = time.Now().AddDate(0,3,0)
}