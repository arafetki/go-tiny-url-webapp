package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func (app *application) start() error {

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", app.cfg.Server.Port),
		Handler:      app.Routes(),
		ErrorLog:     slog.NewLogLogger(app.logger.Handler(), slog.LevelWarn),
		ReadTimeout:  app.cfg.Server.ReadTimeout,
		WriteTimeout: app.cfg.Server.WriteTimeout,
		IdleTimeout:  app.cfg.Server.IdleTimeout,
	}

	shutdownErrChan := make(chan error)

	go func() {

		quitChan := make(chan os.Signal, 1)
		signal.Notify(quitChan, syscall.SIGINT, syscall.SIGTERM)
		<-quitChan

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		shutdownErrChan <- srv.Shutdown(ctx)

	}()

	app.logger.Info(fmt.Sprintf("started %s server on %s", app.cfg.App.Env, srv.Addr))
	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutdownErrChan
	if err != nil {
		return err
	}

	app.logger.Warn("server stopped")
	app.wg.Wait()

	return nil
}
