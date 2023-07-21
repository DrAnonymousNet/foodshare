package auth

import (
	"net/http"
	"strings"
	"time"

	"errors"

	core "github.com/DrAnonymousNet/foodshare/Core"
)

type CreateUserRequest struct{
	FirstName string 
	LastName string 
	FullName string
	DOB time.Time 
	Gender string 
	Email string 
	Password string 
}

func (c *CreateUserRequest)Bind(r *http.Request) error{
	if c.Email == ""{
		return errors.New("email is required")
	}
	
	var user User
	core.DB.Model(&User{}).Where("Email = ?", c.Email).First(&user)
	if user.ID != 0 {
		return errors.New("email alraedy picked")
	}
	if c.Password == ""{
		return errors.New(("invalid password"))
	}
	//TODO Validate strong password

	c.FirstName = strings.ToTitle(c.FirstName)
	c.LastName = strings.ToTitle(c.LastName)
	c.FullName = c.FirstName + c.LastName
	return nil
	
}
