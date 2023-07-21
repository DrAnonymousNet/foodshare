package auth

import "github.com/go-chi/chi/v5"

func AuthRoutes() chi.Router {
	r := chi.NewRouter()
	return r.Route("/accounts", func(r chi.Router) {
		r.Post("/", CreateUser)
		r.Post("/login", GenerateJWTTokenHandler)
	})
}
