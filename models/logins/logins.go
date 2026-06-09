package logins

import (
	"jamesmukumu/maasaimaratripsportal/models/users"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Logins struct {
gorm.Model
TokenString string `gorm:"unique;not null"`
LoginID string `gorm:"not null;unique"`
UserID uint
User users.User `references:"id"`
ExpiryTime time.Time `gorm:"not null"`
}

func (login *Logins)Preset(){
login.ExpiryTime = time.Now().UTC().Add(3 * time.Hour)
login.LoginID = uuid.NewString()
}

