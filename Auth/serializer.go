package auth

import (
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"errors"

	core "github.com/DrAnonymousNet/foodshare/Core"
	"github.com/go-chi/jwtauth/v5"
	"github.com/google/uuid"
)

type CreateUserSerializer struct {
	FirstName string    `json:"first_name" validate:"required"`
	LastName  string    `json:"last_name" validate:"required"`
	FullName  string    `json:"-"`
	DOB       time.Time `json:"date_of_birth"`
	Gender    string    `json:"gender" validate:"required"`
	Email     string    `json:"email" validate:"required"`
	Password  string
	UID       uuid.UUID `json:"uid"`
}

func (c *CreateUserSerializer) Bind(r *http.Request) error {
	if c.Email == "" {
		return errors.New("email is required")
	}

	var user User
	core.DB.Model(&User{}).Where("Email = ?", c.Email).First(&user)
	if user.ID != 0 {
		return errors.New("email alraedy picked")
	}
	if c.Password == "" {
		return errors.New(("invalid password"))
	}
	//TODO Validate strong password

	c.FirstName = strings.ToTitle(c.FirstName)
	c.LastName = strings.ToTitle(c.LastName)
	c.FullName = c.FirstName + c.LastName
	c.UID = uuid.New()
	log.Println(string(c.UID.String()))
	return nil

}

func (d *CreateUserSerializer) Save(r *http.Request) error {
	//var user auth.User
	//core.DB.Model(&auth.User{}).Where("ID = ?", d.DonorID)

	userRequest := User{
		UID:       d.UID,
		FirstName: d.FirstName,
		LastName:  d.LastName,
		FullName:  d.FullName,
		DOB:       d.DOB,
		Email:     d.Email,
		Gender:    d.Gender,
		Password:  d.Password,
	}
	err := userRequest.SetPassword(userRequest.Password)
	if err != nil {
		return err
	}
	err = core.DB.Model(&User{}).Create(&userRequest).Error
	if err != nil {
		return err
	}
	d.Password = ""
	return nil
}

type LoginSerializer struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
	user *User `json:"-"`
}

func (l *LoginSerializer) Bind(r *http.Request) error {
	if l.Email == "" {
		return errors.New("email is required")
	}
	if l.Password == "" {
		return errors.New("password is required")
	}
	err := l.Validate(r)
	if err != nil {
		return err
	}
	return nil
}

func (l *LoginSerializer) Validate(r *http.Request) error {
	//Retrieve the user with the username
	var user *User
	err := core.DB.Model(&User{}).Where("email = ?", l.Email).First(&user).Error
	if err != nil {
		return err
	}

	err = user.ComparePassword(l.Password)
	if err != nil {
		return err
	}
	l.user = user
	return nil
}

func (l *LoginSerializer) Save(r *http.Request) (string, error) {

	//Generate the JWT token
	tokenAuth := jwtauth.New("HS256", []byte(os.Getenv("SECRET")), nil)
	token, tokenString, err := tokenAuth.Encode(map[string]interface{}{l.Email: l.Password})
	if err != nil {
		return "", err
	}
	jwtToken := JwtToken{
		UID:       uuid.New(),
		Token:     tokenString,
		ExpiresAt: token.Expiration(),
		UserID:    l.user.ID,   
	}
	err = core.DB.Model(&JwtToken{}).Create(&jwtToken).Error
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
