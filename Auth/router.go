package auth

import "github.com/go-chi/chi/v5"

func AuthRoutes() chi.Router {
	r := chi.NewRouter()
	r.Post("/signup", CreateUser)
	r.Post("/login", GenerateJWTTokenHandler)
	return r
}

