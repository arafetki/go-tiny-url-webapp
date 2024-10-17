package main

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/arafetki/go-tiny-url-webapp/internal/response"
)

func (app *application) logError(r *http.Request, err error) {

	var (
		message = err.Error()
		method  = r.Method
		url     = r.URL.String()
	)
	requestAttrs := slog.Group("request", "method", method, "url", url)
	app.logger.Error(message, requestAttrs)
}

func (app *application) errorResponse(w http.ResponseWriter, r *http.Request, status int, message any, headers http.Header) {
	err := response.JSONWithHeaders(w, status, envelope{"error": message}, headers)
	if err != nil {
		app.logError(r, err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (app *application) notFoundResponseHandler(w http.ResponseWriter, r *http.Request) {
	message := "the requested resource could not be found"
	app.errorResponse(w, r, http.StatusNotFound, message, nil)
}

func (app *application) methodNotAllowedResponseHandler(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf("the %s method is not supported for this resource", r.Method)
	app.errorResponse(w, r, http.StatusMethodNotAllowed, message, nil)
}

func (app *application) internalServerErrorResponseHandler(w http.ResponseWriter, r *http.Request, err error) {
	message := "the server encountered a problem and could not process your request"
	app.logError(r, err)
	app.errorResponse(w, r, http.StatusInternalServerError, message, nil)

}

func (app *application) badRequestResponseHandler(w http.ResponseWriter, r *http.Request, err error) {
	message := "The request could not be understood by the server due to malformed syntax or incorrect parameter type"
	app.logError(r, err)
	app.errorResponse(w, r, http.StatusBadRequest, message, nil)

}
