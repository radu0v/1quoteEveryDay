package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/radu0v/1quoteEveryDay/internal/handlers"
)

func routes() http.Handler {
	mux := chi.NewRouter()
	mux.Use(middleware.Recoverer)
	mux.Use(NoSurf)
	mux.Get("/", handlers.Repo.Home)
	mux.Post("/", handlers.Repo.PostHome)
	mux.Get("/unsubscribe", handlers.Repo.Unsubscribe)
	mux.Post("/unsubscribe", handlers.Repo.UnsubscribePost)
	mux.Get("/feedback", handlers.Repo.Feedback)
	mux.Post("/feedback", handlers.Repo.FeedbackPost)
	mux.Get("/privacy-policy", handlers.Repo.PrivacyPolicy)

	mux.Mount("/admin", adminRouter())
	//enabling static files
	fileserver := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileserver))

	return mux
}

// admin router function that defines the admin routing
func adminRouter() http.Handler {
	r := chi.NewRouter()
	r.Use(NoSurf)
	r.Use(middleware.Recoverer)
	r.Use(SessionLoad)
	r.Get("/", handlers.Repo.Admin)
	r.Get("/login", handlers.Repo.Login)
	r.Post("/login", handlers.Repo.PostLogin)
	r.Get("/quotes", handlers.Repo.AdminQuotes)
	r.Post("/quotes", handlers.Repo.DeleteQuote)
	r.Get("/subscribers", handlers.Repo.Subscribers)
	r.Get("/quotes/add", handlers.Repo.AdminAddQuote)
	r.Post("/quotes/add", handlers.Repo.AdminAddQuotePost)
	return r
}
