package configs

import "os"

type dbConfig struct {
	Dsn string
}

type AuthConfig struct {
	SecretKey string
}

type AppConfig struct {
	Db   dbConfig
	Auth AuthConfig
}

func LoadAppConfig() *AppConfig {
	return &AppConfig{
		Db: dbConfig{
			Dsn: os.Getenv("DB_DSN"),
		},
		Auth: AuthConfig{
			SecretKey: os.Getenv("JWT_SECRET"),
		},
	}
}
