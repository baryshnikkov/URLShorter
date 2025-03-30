package configs

import "os"

type dbConfig struct {
	Dsn string
}

type AppConfig struct {
	Db dbConfig
}

func LoadAppConfig() *AppConfig {
	return &AppConfig{
		Db: dbConfig{
			Dsn: os.Getenv("DB_DSN"),
		},
	}
}
