package auth

import "github.com/go-chi/chi/v5"

func AuthRoutes() chi.Router {
	r := chi.NewRouter()
	return r.Route("/auth", func(r chi.Router) {
		r.Post("/signup", CreateUser)
		r.Post("/login", GenerateJWTTokenHandler)
	})
}
