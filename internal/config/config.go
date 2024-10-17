package config

import "time"

type Config struct {
	App struct {
		Version string
		Env     Environment
		Debug   bool
	}
	Server struct {
		Port           int
		ReadTimeout    time.Duration
		WriteTimeout   time.Duration
		IdleTimeout    time.Duration
		ShutdownPeriod time.Duration
	}
	Database struct {
		Dsn          string
		AutoMigrate  bool
		MaxOpenConns int
		MaxIdleConns int
		ConnMaxLife  time.Duration
		ConnMaxIdle  time.Duration
	}
}
