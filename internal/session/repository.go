package session

import "URLShorter/pkg/database"

type Repository struct {
	Database *database.Db
}

func NewRepository(database *database.Db) *Repository {
	return &Repository{
		Database: database,
	}
}

func (r *Repository) Save(Session *Session) (*Session, error) {
	result := r.Database.DB.Create(Session)
	if result.Error != nil {
		return nil, result.Error
	}

	return Session, nil
}
