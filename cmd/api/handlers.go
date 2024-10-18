package main

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/arafetki/go-tiny-url-webapp/internal/data"
	"github.com/arafetki/go-tiny-url-webapp/internal/db/models"
	"github.com/arafetki/go-tiny-url-webapp/internal/nanoid"
	"github.com/arafetki/go-tiny-url-webapp/internal/request"
	"github.com/arafetki/go-tiny-url-webapp/internal/response"
	"github.com/go-chi/chi/v5"
)

func (app *application) createTinyURLHandler(w http.ResponseWriter, r *http.Request) {

	var input struct {
		LongURL string `json:"long_url" validate:"http_url"`
	}

	err := request.DecodeJSONStrict(w, r, &input)
	if err != nil {
		app.badRequestResponseHandler(w, r, err)
		return
	}

	err = app.validate.Struct(input)
	if err != nil {
		app.badRequestResponseHandler(w, r, err)
		return
	}

	short, err := nanoid.Generate(7)
	if err != nil {
		app.internalServerErrorResponseHandler(w, r, err)
		return
	}

	tinyurl := &models.TinyURL{
		Short:  strings.ToLower(short),
		Long:   input.LongURL,
		Expiry: time.Now().Add(24 * time.Hour),
	}

	err = app.repository.TinyURL.Create(tinyurl)
	if err != nil {
		app.internalServerErrorResponseHandler(w, r, err)
		return
	}

	err = response.JSON(w, http.StatusCreated, envelope{
		"data": tinyurl,
	})
	if err != nil {
		app.internalServerErrorResponseHandler(w, r, err)
	}
}

func (app *application) resolveTinyURLHandler(w http.ResponseWriter, r *http.Request) {

	short := chi.URLParam(r, "short")
	err := app.validate.Var(short, "len=7,nanoid_charset")
	if err != nil {
		app.notFoundResponseHandler(w, r)
		return
	}

	tinyurl, err := app.repository.TinyURL.Get(strings.ToLower(short))
	if err != nil {
		switch {
		case errors.Is(err, data.ErrNotFound):
			app.notFoundResponseHandler(w, r)
		default:
			app.internalServerErrorResponseHandler(w, r, err)
		}
		return
	}

	if time.Now().After(tinyurl.Expiry) {
		app.errorResponse(w, r, http.StatusGone, "the requested resource has expired and is no longer available", nil)
		return
	}

	http.Redirect(w, r, tinyurl.Long, http.StatusMovedPermanently)
}
