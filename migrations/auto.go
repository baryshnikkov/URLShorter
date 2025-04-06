package main

import (
	"URLShorter/internal/session"
	"URLShorter/internal/user"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}

	db, err := gorm.Open(postgres.Open(os.Getenv("DB_DSN")), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(&user.User{}, &session.Session{})
	if err != nil {
		panic(err)
	}
}
