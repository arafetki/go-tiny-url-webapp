package database

import (
	"context"
	"errors"
	"time"

	"github.com/arafetki/go-tiny-url-webapp/assets"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	*sqlx.DB
}

type Options struct {
	MaxOpenConns int
	MaxIdleConns int
	ConnMaxLife  time.Duration
	ConnMaxIdle  time.Duration
}

func Pool(dsn string, automigrate bool, opts Options) (*DB, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	db, err := sqlx.ConnectContext(ctx, "sqlite3", dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(opts.MaxOpenConns)
	db.SetMaxIdleConns(opts.MaxIdleConns)
	db.SetConnMaxLifetime(opts.ConnMaxLife)
	db.SetConnMaxIdleTime(opts.ConnMaxIdle)

	if automigrate {
		iofsDriver, err := iofs.New(assets.Migrations, "migrations")
		if err != nil {
			return nil, err
		}
		migrator, err := migrate.NewWithSourceInstance("iofs", iofsDriver, "sqlite3://"+dsn)
		if err != nil {
			return nil, err
		}
		err = migrator.Up()
		switch {
		case errors.Is(err, migrate.ErrNoChange):
			break
		default:
			return nil, err
		}
	}

	return &DB{db}, nil
}
