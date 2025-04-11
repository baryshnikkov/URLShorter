package user

import (
	"URLShorter/pkg/database"
)

type Repository struct {
	Database *database.Db
}

func NewRepository(database *database.Db) *Repository {
	return &Repository{
		Database: database,
	}
}

func (r *Repository) Create(User *User) (*User, error) {
	result := r.Database.DB.Create(User)
	if result.Error != nil {
		return nil, result.Error
	}

	return User, nil
}

func (r *Repository) GetByLogin(login string) (*User, error) {
	user := &User{}

	result := r.Database.DB.First(user, "login = ?", login)
	if result.Error != nil {
		return nil, result.Error
	}

	return user, nil
}

func (r *Repository) GetByEmail(email string) (*User, error) {
	user := &User{}

	result := r.Database.DB.First(user, "email = ?", email)
	if result.Error != nil {
		return nil, result.Error
	}

	return user, nil
}

func (r *Repository) GetEmailById(id uint) (email string, err error) {
	user := &User{}
	result := r.Database.DB.First(user, "id = ?", id)
	if result.Error != nil {
		return "", result.Error
	}

	return user.Email, nil
}
