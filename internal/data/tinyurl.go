package data

import (
	"context"
	"database/sql"
	"errors"
	"time"

	database "github.com/arafetki/go-tiny-url-webapp/internal/db"
	"github.com/arafetki/go-tiny-url-webapp/internal/db/models"
)

type TinyURLRepo struct {
	db *database.DB
}

var (
	ErrNotFound = errors.New("record not found")
)

func (turl TinyURLRepo) Create(tinyurl *models.TinyURL) error {

	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	query := `INSERT INTO tinyurls (short,long,expiry) VALUES ($1,$2,$3) RETURNING created;`

	args := []any{tinyurl.Short, tinyurl.Long, tinyurl.Expiry}

	return turl.db.QueryRowxContext(ctx, query, args...).Scan(&tinyurl.Created)
}

func (turl TinyURLRepo) Get(short string) (*models.TinyURL, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	query := `SELECT * FROM tinyurls WHERE short=$1;`

	var tinyurl models.TinyURL

	err := turl.db.QueryRowxContext(ctx, query, short).Scan(
		&tinyurl.Short,
		&tinyurl.Long,
		&tinyurl.Expiry,
		&tinyurl.Created,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	return &tinyurl, nil
}
