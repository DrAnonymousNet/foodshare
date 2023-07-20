package auth

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Role string

const (
	Male Role = "Male"
	Female Role = "Female"
)

type User struct{
	gorm.Model
	UID uuid.UUID `gorm:"default:generate_uuid_v4"`
	FirstName string `gorm:"not null"`
	LastName string `gorm:"not null"`
	FullName  string `gorm:"->;type:GENERATED ALWAYS AS (concat(firstname,' ',lastname));default:(-);"`
	DOB time.Time `gorm:"null"`
	Gender string `gorm:":not null:type:ENUM('Male', 'Female')"`
	Email string `gorm:"not null;unique"`
	Password string `gorm:"-"`
	PasswordHash string `gorm:"not null"`
	

}

type JwtToken struct{
	gorm.Model
	UID uuid.UUID `gorm:"default:generate_uuid_v4"`
	Token string
	ExpiresAt time.Time
	UserID uint8
	User User
}
