package foodshare

import "github.com/go-chi/chi"

func DonationRoutes() chi.Router {
	d := SetViewSet()
	r := chi.NewRouter()
	return r.Route("/donations", func(r chi.Router) {
		r.Get("/", d.ListDonations)
		r.Post("/create", d.CreateDonation)
		//r.Get("/search", SearchDonations)

		r.Route("/{uid}", func(r chi.Router) {
			r.Get("/", d.GetDonation)
			r.Patch("/", d.UpdateDonation)
			r.Delete("/", d.DeleteDonation)
		})
	})

}

//r.Route("/articles", func(r chi.Router) {
//	r.With(paginate).Get("/", ListArticles)
//	r.Post("/", CreateArticle)       // POST /articles
//	r.Get("/search", SearchArticles) // GET /articles/search

//	r.Route("/{articleID}", func(r chi.Router) {
//		r.Use(ArticleCtx)            // Load the *Article on the request context
//		r.Get("/", GetArticle)       // GET /articles/123
//		r.Put("/", UpdateArticle)    // PUT /articles/123
//		r.Delete("/", DeleteArticle) // DELETE /articles/123
//	})

// GET /articles/whats-up
//	r.With(ArticleCtx).Get("/{articleSlug:[a-z-]+}", GetArticle)
//})
