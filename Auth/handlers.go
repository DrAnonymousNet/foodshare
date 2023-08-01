package auth

import (
	"fmt"
	"net/http"
	"github.com/go-chi/render"
)

func CreateUser(w http.ResponseWriter, r *http.Request) {
	// Parse the request body into a CreateUserRequest struct

	data := &CreateUserSerializer{}
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

	tokenString, err := data.Save(r)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{
			"error": "Unable to generate token " + fmt.Sprintf("%v", err),
		})
		return 
	}


	render.Status(r, 200)
	render.JSON(w, r, map[string]string{"token": tokenString})

}
