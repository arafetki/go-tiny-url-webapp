package nanoid

import (
	"strings"

	"github.com/go-playground/validator/v10"
	gonanoid "github.com/matoous/go-nanoid/v2"
)

const Charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

func Generate(size int) (string, error) {

	id, err := gonanoid.Generate(Charset, size)
	if err != nil {
		return "", err
	}

	return id, nil
}

func CharsetValidate(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	for _, char := range value {
		if !strings.ContainsRune(Charset, char) {
			return false
		}
	}
	return true
}
