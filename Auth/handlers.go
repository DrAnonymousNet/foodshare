package auth

import (
	"fmt"
	"log"
	"net/http"
	"os"

	core "github.com/DrAnonymousNet/foodshare/Core"
	"github.com/go-chi/jwtauth/v5"
	"github.com/go-chi/render"
)

func CreateUser(w http.ResponseWriter, r *http.Request) {
	// Parse the request body into a CreateUserRequest struct

	data := &CreateUserRequest{}
	if err := render.Bind(r, data); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Invalid request" + fmt.Sprintf("%v", err)})
		return
	}

	
	err := data.Save(r)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{
			"error": "Unable to create a user " + fmt.Sprintf("%v", err),
		})
	}


	// Respond with the created User
	render.Status(r, http.StatusCreated)
	render.JSON(w, r, data)
}


func GenerateJWTTokenHandler(w http.ResponseWriter, r *http.Request) {

	data := &LoginSerializer{}
	if err := render.Bind(r, data); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Invalid request " + fmt.Sprintf("%v", err)})
		return
	}
	//Retrieve the user with the username
	var user *User
	err := core.DB.Model(&User{}).Where("email = ?", data.Email).First(&user).Error
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{
			"error": "Invalid email",
		})
		return
	}
	log.Println(user)

	err = user.ComparePassword(data.Password)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{
			"error": "Invalid password",})
		return
	}
	//Generate token for the user
	tokenAuth := jwtauth.New("HS256", []byte(os.Getenv("SECRET")), nil)
	_, tokenString, err := tokenAuth.Encode(map[string]interface{}{data.Email: data.Password})
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{
			"error": "Unable to generate token",
		})
		return
	}

	render.Status(r, 200)
	render.JSON(w, r, map[string]string{"token": tokenString})

}
