package auth

import (
	"net/http"
	"os"

	core "github.com/DrAnonymousNet/foodshare/Core"
	"github.com/go-chi/jwtauth/v5"
	"github.com/go-chi/render"
)

func CreateUser(w http.ResponseWriter, r *http.Request) {
	// Parse the request body into a CreateUserRequest struct

	var user User
	data := &CreateUserRequest{}
	if err := render.Bind(r, data); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Invalid request"})
		return
	}

	// Create a new User with the provided data

	if err := core.DB.Create(user).Error; err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{
			"error": "Unable to create a user",
		})
	}

	err := user.SetPassword(user.Password)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": err.Error()})
		return
	}
	user.Password = "" //Set the password to nil before sending the request.

	// Respond with the created User
	render.Status(r, http.StatusCreated)
	render.JSON(w, r, user)
}

func GenerateJWTTokenHandler(w http.ResponseWriter, r *http.Request) {

	type user_auth_data struct {
		username string
		password string
	}
	//Decode the users post data
	auth_data := &user_auth_data{}
	if err := render.DecodeJSON(r.Body, auth_data); err != nil {
		core.ErrInvalidRequest(err)
	}

	//Retrieve the user with the username
	var user User
	core.DB.Model(&User{}).Where("username = ?", auth_data.username).First(&user)

	err := user.ComparePassword(auth_data.password)
	if err != nil {
		core.ErrInvalidRequest(err)
	}
	//Generate token for the user
	tokenAuth := jwtauth.New("HS256", []byte(os.Getenv("SECRET")), nil)
	_, tokenString, _ := tokenAuth.Encode(map[string]interface{}{auth_data.username: auth_data.password})

	render.Status(r, 200)
	render.JSON(w, r, map[string]string{"token": tokenString})

}
