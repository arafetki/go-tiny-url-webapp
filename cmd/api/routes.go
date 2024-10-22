package main

import (
	"net/http"

	"github.com/arafetki/go-tinyurl/assets"
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
	fs := http.FileServer(http.FS(assets.StaticFiles))
	mux.Handle("/static/*", app.secureFS(fs))

	// HomePage
	mux.Get("/", func(w http.ResponseWriter, r *http.Request) {
		err := app.tmpl.ExecuteTemplate(w, "home.html", map[string]any{
			"name": "Arafet",
		})
		if err != nil {
			app.internalServerErrorResponseHandler(w, r, err)
		}
	})

	// API version 1 routes
	mux.Route("/v1", func(r chi.Router) {
		r.Post("/tinyurl", app.createTinyURLHandler)
		r.Get("/tinyurl/{short}", app.resolveTinyURLHandler)
	})

	return mux
}
