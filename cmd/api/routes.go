package main

import (
	"net/http"

	"github.com/arafetki/go-tiny-url-webapp/assets"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (app *application) Routes() *chi.Mux {

	mux := chi.NewMux()

	// Middlewares
	mux.Use(middleware.Logger)
	mux.Use(middleware.Recoverer)

	// notFound and methodNotAllowed handlers
	mux.NotFound(app.notFoundResponseHandler)
	mux.MethodNotAllowed(app.methodNotAllowedResponseHandler)

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
		r.Post("/shorten", app.createTinyURLHandler)
		r.Get("/resolve/{short}", app.resolveTinyURLHandler)
	})

	mux.Mount("/v1", v1)

	return mux
}
