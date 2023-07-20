package auth

import (
	"net/http"
	"strings"
	"time"

	"errors"

	auth "github.com/DrAnonymousNet/foodshare/Auth"
	core "github.com/DrAnonymousNet/foodshare/Core"
)

type CreateUserRequest struct{
	FirstName string 
	LastName string 
	DOB time.Time 
	Gender string 
	Email string 
	Password string 
}

func (c *CreateUserRequest)Bind(r *http.Request) error{
	if c.Email == ""{
		return errors.New("email is required")
	}
	
	var user auth.User
	core.DB.Model(&auth.User{}).Where("Email = ?", c.Email).First(&user)
	if user.ID != 0 {
		return errors.New("email alraedy picked")
	}
	if c.Password == ""{
		return errors.New(("invalid password"))
	}
	//TODO Validate strong password

	c.FirstName = strings.Title(c.FirstName)
	c.LastName = strings.Title(c.LastName)
	return nil
	
}
