//go:generate pnpm build

package main

import (
	"html/template"
	"log/slog"
	"os"
	"runtime/debug"
	"sync"
	"time"

	"github.com/arafetki/go-tiny-url-webapp/assets"
	"github.com/arafetki/go-tiny-url-webapp/internal/config"
	"github.com/arafetki/go-tiny-url-webapp/internal/data"
	database "github.com/arafetki/go-tiny-url-webapp/internal/db"
	"github.com/arafetki/go-tiny-url-webapp/internal/env"
	"github.com/arafetki/go-tiny-url-webapp/internal/logging"
	"github.com/arafetki/go-tiny-url-webapp/internal/version"
)

type application struct {
	cfg        config.Config
	logger     *slog.Logger
	tmpl       *template.Template
	repository *data.Reposiroty
	wg         sync.WaitGroup
}

func main() {

	var cfg config.Config

	cfg.App.Version = version.Get()
	cfg.App.Env = env.GetAppEnv("APP_ENV", config.DEVELOPMENT)
	cfg.App.Debug = env.GetBool("APP_DEBUG", true)

	cfg.Server.Port = env.GetInt("SERVER_PORT", 8080)
	cfg.Server.ReadTimeout = env.GetDuration("SERVER_READ_TIMEOUT", 5*time.Second)
	cfg.Server.WriteTimeout = env.GetDuration("SERVER_WRITE_TIMEOUT", 5*time.Second)
	cfg.Server.IdleTimeout = env.GetDuration("SERVER_IDLE_TIMEOUT", time.Minute)

	cfg.Database.Dsn = env.GetString("DATABASE_DSN", "sqlite.db")
	cfg.Database.MaxOpenConns = env.GetInt("DATABASE_MAX_OPEN_CONNS", 25)
	cfg.Database.MaxIdleConns = env.GetInt("DATABASE_MAX_IDLE_CONNS", 25)
	cfg.Database.ConnMaxLife = env.GetDuration("DATABASE_CONN_MAX_LIFE", 5*time.Minute)
	cfg.Database.ConnMaxIdle = env.GetDuration("DATABASE_CONN_MAX_IDLE", time.Hour)
	cfg.Database.AutoMigrate = env.GetBool("DATABASE_AUTO_MIGRATE", true)

	logger := logging.NewLogger(logging.Options{Writer: os.Stdout, Debug: cfg.App.Debug})

	db, err := database.Pool(cfg.Database.Dsn, cfg.Database.AutoMigrate, database.Options{
		MaxOpenConns: cfg.Database.MaxOpenConns,
		MaxIdleConns: cfg.Database.MaxIdleConns,
		ConnMaxLife:  cfg.Database.ConnMaxLife,
		ConnMaxIdle:  cfg.Database.ConnMaxIdle,
	})

	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	defer db.Close()

	logger.Info("database connection has been established successfully")

	tmpl, err := template.ParseFS(assets.Templates, "templates/*.html")
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	app := &application{
		cfg:        cfg,
		logger:     logger,
		tmpl:       tmpl,
		repository: data.NewRepo(db),
	}

	err = app.start()
	if err != nil {
		trace := string(debug.Stack())
		app.logger.Error(err.Error(), "trace", trace)
		os.Exit(1)
	}
}
