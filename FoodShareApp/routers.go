package foodshare

import "github.com/go-chi/chi"

func DonationRouters(r chi.Router){
	

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