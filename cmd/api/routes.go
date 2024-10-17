package main

import (
	"net/http"

	"github.com/arafetki/go-tiny-url-webapp/assets"
	"github.com/arafetki/go-tiny-url-webapp/internal/response"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (app *application) Routes() *chi.Mux {

	mux := chi.NewMux()

	// Middlewares
	mux.Use(middleware.Logger)
	mux.Use(middleware.Recoverer)

	// Health endpoint for load balancers
	mux.Get("/health", app.healthCheckHandler)

	// Static files
	fs := http.FileServer(http.FS(assets.Static))
	mux.Handle("/static/*", app.secureFS(fs))

	// HomePage
	mux.Get("/", func(w http.ResponseWriter, r *http.Request) {
		app.tmpl.ExecuteTemplate(w, "home.html", map[string]any{
			"name": "Arafet",
		})
	})

	// API version 1 routes
	v1 := mux.Group(func(r chi.Router) {
		r.Get("/hello", func(w http.ResponseWriter, r *http.Request) {
			response.JSON(w, http.StatusOK, map[string]string{"message": "Hello World"})
		})
	})

	mux.Mount("/v1", v1)

	return mux
}
