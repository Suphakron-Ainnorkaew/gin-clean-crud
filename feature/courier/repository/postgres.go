package repository

import (
	"go-clean-api/domain"
	"go-clean-api/entity"

	"gorm.io/gorm"
)

type postgresCourierRepository struct {
	db *gorm.DB
}

func NewPostgresCourierRepository(db *gorm.DB) domain.CourierRepository {
	return &postgresCourierRepository{db: db}
}

func (r *postgresCourierRepository) Create(courier *entity.Courier) error {
	return r.db.Create(courier).Error
}

func (r *postgresCourierRepository) GetByID(id uint) (*entity.Courier, error) {
	var courier entity.Courier
	err := r.db.First(&courier, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &courier, nil
}

func (r *postgresCourierRepository) GetAll() ([]entity.Courier, error) {
	var couriers []entity.Courier
	err := r.db.Find(&couriers).Error
	return couriers, err
}

func (r *postgresCourierRepository) Update(courier *entity.Courier) error {
	return r.db.Save(courier).Error
}

func (r *postgresCourierRepository) Delete(id uint) error {
	return r.db.Delete(&entity.Courier{}, id).Error
}
