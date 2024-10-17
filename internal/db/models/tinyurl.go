package models

import "time"

type TinyURL struct {
	Short   string    `json:"short"`
	Long    string    `json:"long_url"`
	Expiry  time.Time `json:"expiry"`
	Created time.Time `json:"-"`
}
