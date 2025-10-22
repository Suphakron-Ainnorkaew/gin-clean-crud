package courier

import (
	"go-clean-api/internal/courier/domain"

	"gorm.io/gorm"
)

type CourierRepository interface {
	Create(courier *domain.Courier) error
	FindAll() ([]domain.Courier, error)
	FindByID(id uint) (*domain.Courier, error)
	Update(courier *domain.Courier) error
	Delete(id uint) error
}

type gormCourierRepository struct {
	db *gorm.DB
}

func NewGormCourierRepository(db *gorm.DB) CourierRepository {
	return &gormCourierRepository{db: db}
}

func (r *gormCourierRepository) Create(courier *domain.Courier) error {
	return r.db.Create(courier).Error
}

func (r *gormCourierRepository) FindAll() ([]domain.Courier, error) {
	var courier []domain.Courier
	err := r.db.Find(&courier).Error
	return courier, err
}

func (r *gormCourierRepository) FindByID(id uint) (*domain.Courier, error) {
	var courier domain.Courier
	err := r.db.First(&courier, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &courier, nil
}

func (r *gormCourierRepository) Update(courier *domain.Courier) error {
	return r.db.Save(courier).Error
}

func (r *gormCourierRepository) Delete(id uint) error {
	return r.db.Delete(&domain.Courier{}, id).Error
}
