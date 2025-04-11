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

func (r *Repository) FindAllByUserId(userId uint) ([]*Session, error) {
	var sessions []*Session
	result := r.Database.DB.Where("user_id = ?", userId).Find(&sessions)
	if result.Error != nil {
		return nil, result.Error
	}

	return sessions, nil
}

func (r *Repository) DeleteByRefreshTokenHash(refreshTokenHash string) error {
	result := r.Database.DB.Where("refresh_token_hash = ?", refreshTokenHash).Delete(&Session{})
	if result.Error != nil {
		return result.Error
	}

	return nil
}
