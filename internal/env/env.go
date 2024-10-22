package env

import (
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/arafetki/go-tinyurl/internal/config"
)

func GetString(key string, defaultValue string) string {

	value, exist := os.LookupEnv(key)
	if !exist {
		return defaultValue
	}

	return value
}

func GetInt(key string, defaultValue int) int {

	value, exist := os.LookupEnv(key)
	if !exist {
		return defaultValue
	}

	intVal, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}

	return intVal
}

func GetBool(key string, defaultValue bool) bool {

	value, exist := os.LookupEnv(key)
	if !exist {
		return defaultValue
	}

	boolVal, err := strconv.ParseBool(value)
	if err != nil {
		return defaultValue
	}

	return boolVal
}

func GetDuration(key string, defaultValue time.Duration) time.Duration {

	value, exist := os.LookupEnv(key)
	if !exist {
		return defaultValue
	}

	d, err := time.ParseDuration(value)
	if err != nil {
		return defaultValue
	}

	return d
}

func GetAppEnv(key string, defaultValue config.Environment) config.Environment {

	value, exist := os.LookupEnv(key)
	if !exist {
		return defaultValue
	}

	switch strings.ToLower(value) {
	case "development":
		return config.DEVELOPMENT
	case "staging":
		return config.STAGING
	case "production":
		return config.PRODUCTION
	default:
		return defaultValue
	}
}
