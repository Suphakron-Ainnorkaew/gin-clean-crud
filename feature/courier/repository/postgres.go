package repository

import (
	"go-clean-api/domain"
	"go-clean-api/entity"

	"github.com/pkg/errors"

	"gorm.io/gorm"
)

type postgresCourierRepository struct {
	db *gorm.DB
}

func NewPostgresCourierRepository(db *gorm.DB) domain.CourierRepository {
	return &postgresCourierRepository{db: db}
}

func (r *postgresCourierRepository) Create(courier *entity.Courier) error {
	if err := r.db.Create(&courier).Error; err != nil {
		return errors.Wrap(err, "[CourierRepository.Create]: unable to create courier")
	}
	return nil
}

func (r *postgresCourierRepository) GetByID(id uint) (*entity.Courier, error) {
	var courier entity.Courier
	if err := r.db.First(&courier, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, errors.Wrap(err, "[CourierRepository.GetByID]: unable to get courier id")
	}

	return &courier, nil
}

func (r *postgresCourierRepository) GetAll() ([]entity.Courier, error) {
	var couriers []entity.Courier
	if err := r.db.Find(&couriers).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return couriers, errors.Wrap(err, "[CourierRepository.GetAll]: unable to get courier")
	}
	return couriers, nil
}

func (r *postgresCourierRepository) Update(courier *entity.Courier) error {
	if err := r.db.Model(&entity.Courier{}).Where("id = ?", courier.ID).Updates(&courier).Error; err != nil {
		return errors.Wrap(err, "[CourierRepository.Update]: unable to update courier")
	}
	return nil
}

func (r *postgresCourierRepository) Delete(id uint) error {
	if err := r.db.Delete(&entity.Courier{}, id).Error; err != nil {
		return errors.Wrap(err, "[CourierRepository.Delete]: unable to Delete courier")
	}
	return nil
}
