package data

import (
	database "github.com/arafetki/go-tinyurl/internal/db"
	"github.com/arafetki/go-tinyurl/internal/db/models"
)

type Reposiroty struct {
	TinyURL interface {
		Create(tinyURL *models.TinyURL) error
		Get(short string) (*models.TinyURL, error)
	}
}

func NewRepo(db *database.DB) *Reposiroty {
	return &Reposiroty{
		TinyURL: TinyURLRepo{db},
	}
}
