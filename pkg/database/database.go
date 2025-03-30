package database

import (
	"URLShorter/configs"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Db struct {
	*gorm.DB
}

func New(config *configs.AppConfig) *Db {
	db, err := gorm.Open(postgres.Open(config.Db.Dsn), &gorm.Config{})
	if err != nil {
		zap.L().Fatal("failed to connect database", zap.Error(err))
	}

	return &Db{db}
}
