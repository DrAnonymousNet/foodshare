package auth

import (
	"errors"
	"log"
	"strings"

	core "github.com/DrAnonymousNet/foodshare/Core"
	"golang.org/x/crypto/bcrypt"
)

func (u *User) SetPassword(password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.PasswordHash = string(hash)
	return nil
}

func (u *User) ComparePassword(password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	log.Println(u.PasswordHash, password)
	return err //returns nil if error is nil
}


func ParseAuthorizationHeader(header string) (string, string, error) {
	realm, token := strings.Split(header, " ")[0], strings.Split(header, " ")[1]
	if realm != "Bearer" {
		return "", "", errors.New("invalid authorization header")
	}
	return realm, token, nil
}

func GetUserFromToken(tokenString string) (*User, error){
	token := &JwtToken{}
	err := core.DB.Model(&JwtToken{}).Where("token = ?", tokenString).First(&token).Error
	if err != nil {
		return nil, err
	}
	user := &User{}
	err = core.DB.Model(&User{}).Where("id = ?", token.UserID).First(&user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}