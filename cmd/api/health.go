package main

import (
	"net/http"
	"time"

	"github.com/arafetki/go-tiny-url-webapp/internal/response"
)

func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]any{
		"timestamp": time.Now().Nanosecond(),
	}

	err := response.JSON(w, http.StatusOK, data)
	if err != nil {
		app.logger.Error(err.Error())
	}
}
