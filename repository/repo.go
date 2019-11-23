package repository

import (
	"github.com/Dimitriy14/notifyme/models"
	"github.com/Dimitriy14/notifyme/postgres"
)

type Repository interface {
	GetFilters() ([]models.ProductFiler, error)
	SaveFilter(filter []models.ProductFiler) ([]models.ProductFiler, error)
	DeleteFilter(filter models.ProductFiler) error
}

type repoImpl struct {
	db postgres.PGClient
}

func NewRepo(db postgres.PGClient) Repository {
	return &repoImpl{
		db: db,
	}
}

func (r *repoImpl) GetFilters() ([]models.ProductFiler, error) {
	var filters []models.ProductFiler
	if err := r.db.Session.Find(&filters).Error; err != nil {
		return nil, err
	}
	return filters, nil
}

func (r *repoImpl) SaveFilter(filters []models.ProductFiler) ([]models.ProductFiler, error) {
	for _, filter := range filters {
		if err := r.db.Session.Save(&filter).Error; err != nil {
			return nil, err
		}
	}

	return filters, nil
}
func (r *repoImpl) DeleteFilter(filter models.ProductFiler) error {
	if err := r.db.Session.Delete(&filter).Error; err != nil {
		return err
	}
	return nil
}
