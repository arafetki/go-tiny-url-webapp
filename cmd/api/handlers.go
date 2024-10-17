package main

import (
	"errors"
	"net/http"
	"time"

	"github.com/arafetki/go-tiny-url-webapp/internal/data"
	"github.com/arafetki/go-tiny-url-webapp/internal/db/models"
	"github.com/arafetki/go-tiny-url-webapp/internal/request"
	"github.com/arafetki/go-tiny-url-webapp/internal/response"
	"github.com/go-chi/chi/v5"
	gonanoid "github.com/matoous/go-nanoid/v2"
)

func (app *application) createTinyURLHandler(w http.ResponseWriter, r *http.Request) {

	var input struct {
		LongURL string `json:"long_url"`
	}

	err := request.DecodeJSONStrict(w, r, &input)
	if err != nil {
		app.badRequestResponseHandler(w, r, err)
		return
	}

	short, err := gonanoid.New(7)
	if err != nil {
		app.internalServerErrorResponseHandler(w, r, err)
		return
	}

	tinyurl := &models.TinyURL{
		Short:  short,
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
	tinyurl, err := app.repository.TinyURL.Get(short)
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
		app.errorResponse(w, r, http.StatusGone, "This short URL has expired and is no longer available.", nil)
		return
	}

	http.Redirect(w, r, tinyurl.Long, http.StatusMovedPermanently)
}
