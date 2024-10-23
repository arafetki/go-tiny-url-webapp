package main

import (
	"errors"
	"net/http"
	"time"

	"github.com/arafetki/go-tinyurl/internal/data"
	"github.com/arafetki/go-tinyurl/internal/db/models"
	"github.com/arafetki/go-tinyurl/internal/nanoid"
	"github.com/arafetki/go-tinyurl/internal/request"
	"github.com/arafetki/go-tinyurl/internal/response"
	"github.com/arafetki/go-tinyurl/internal/utils"
	"github.com/go-chi/chi/v5"
)

func (app *application) createTinyURLHandler(w http.ResponseWriter, r *http.Request) {

	var input struct {
		LongURL string `json:"long_url" validate:"http_url"`
	}

	if err := request.DecodeJSONStrict(w, r, &input); err != nil {
		app.badRequestResponseHandler(w, r, err)
		return
	}

	if err := app.validate.Struct(input); err != nil {
		app.badRequestResponseHandler(w, r, err)
		return
	}

	short, err := nanoid.Generate(9)
	if err != nil {
		app.internalServerErrorResponseHandler(w, r, err)
		return
	}

	tinyurl := &models.TinyURL{
		Short:  short,
		Long:   input.LongURL,
		Expiry: time.Now().Add(720 * time.Hour),
	}

	if err := app.repository.TinyURL.Create(tinyurl); err != nil {
		app.internalServerErrorResponseHandler(w, r, err)
		return
	}

	if err := response.JSON(w, http.StatusCreated, envelope{
		"data": tinyurl,
	}); err != nil {
		app.internalServerErrorResponseHandler(w, r, err)
	}
}

func (app *application) resolveTinyURLHandler(w http.ResponseWriter, r *http.Request) {

	short := chi.URLParam(r, "short")
	if err := app.validate.Var(short, "len=9,nanoid_charset"); err != nil {
		app.notFoundResponseHandler(w, r)
		return
	}

	if cachedTinyURL, found := app.cache.Get(short); found {

		http.Redirect(w, r, cachedTinyURL.Long, http.StatusMovedPermanently)
		return
	}

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

	if utils.IsExpired(tinyurl.Expiry) {
		app.errorResponse(w, r, http.StatusGone, "the requested resource has expired and is no longer available", nil)
		return
	}

	app.cache.Set(tinyurl.Short, tinyurl, 24*time.Hour)

	http.Redirect(w, r, tinyurl.Long, http.StatusMovedPermanently)
}
