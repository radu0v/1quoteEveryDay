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

	mux.Route("/admin", func(r chi.Router) {
		//mux.Use(Auth)
		mux.Get("/login", handlers.Repo.Login)
		mux.Get("/dashboard", handlers.Repo.Admin)
		mux.Get("/quotes", handlers.Repo.AdminQuotes)
		mux.Get("/subscribers", handlers.Repo.Subscribers)
	})

	//enabling static files
	fileserver := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileserver))

	return mux
}
