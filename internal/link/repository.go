package link

import "URLShorter/pkg/database"

type Repository struct {
	Database *database.Db
}

func NewRepository(database *database.Db) *Repository {
	return &Repository{
		Database: database,
	}
}

func (r *Repository) Create(link *Link) (*Link, error) {
	result := r.Database.DB.Create(link)
	if result.Error != nil {
		return nil, result.Error
	}

	return link, nil
}

func (r *Repository) GetByHash(hash string) (*Link, error) {
	link := &Link{}
	result := r.Database.DB.Where("hash = ?", hash).First(link)
	if result.Error != nil {
		return nil, result.Error
	}

	return link, nil
}
