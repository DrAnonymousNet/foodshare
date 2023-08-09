package auth

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	UID          uuid.UUID `pg:"type:uuid"`
	FirstName    string    `gorm:"not null"`
	LastName     string    `gorm:"not null"`
	FullName     string
	DOB          time.Time `gorm:"null"`
	Gender       string    //Gender    `sql:"type:gender"`
	Email        string    `gorm:"not null;unique"`
	Password     string    `gorm:"-"`
	PasswordHash string    `gorm:"not null"`
}

func (u *User) isModel() bool {
	return true
}

type JwtToken struct {
	gorm.Model
	UID       uuid.UUID `pg:"type:uuid" gorm:"unique"`
	Token     string
	ExpiresAt time.Time
	UserID    uint
	User      User
}

func (j *JwtToken) isModel() bool {
	return true
}
