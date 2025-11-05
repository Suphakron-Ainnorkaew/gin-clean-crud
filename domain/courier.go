package domain

import (
	"go-clean-api/entity"
)

type CourierUsecase interface {
	CreateCourier(courier *entity.Courier) error
	GetCourierByID(id uint) (*entity.Courier, error)
	GetAllCourier() ([]entity.Courier, error)
	UpdateCourier(courier *entity.Courier) error
	DeleteCourier(id uint) error
}

type CourierRepository interface {
	Create(courier *entity.Courier) error
	GetByID(id uint) (*entity.Courier, error)
	GetAll() ([]entity.Courier, error)
	Update(courier *entity.Courier) error
	Delete(id uint) error
}
