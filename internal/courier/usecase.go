package courier

import (
	"go-clean-api/internal/courier/domain"
)

type CourierUsecase interface {
	CreateCourier(courier *domain.Courier) error
	GetAllCourier() ([]domain.Courier, error)
	GetCourierByID(id uint) (*domain.Courier, error)
	UpdateCourier(courier *domain.Courier) error
	DeleteCourier(id uint) error
}

type courierUsecase struct {
	courierRepo CourierRepository
}

func NewCourierUsecase(repo CourierRepository) CourierUsecase {
	return &courierUsecase{courierRepo: repo}
}

func (u *courierUsecase) CreateCourier(courier *domain.Courier) error {
	return u.courierRepo.Create(courier)
}

func (u *courierUsecase) GetAllCourier() ([]domain.Courier, error) {
	return u.courierRepo.FindAll()
}

func (u *courierUsecase) GetCourierByID(id uint) (*domain.Courier, error) {
	return u.courierRepo.FindByID(id)
}

func (u *courierUsecase) UpdateCourier(user *domain.Courier) error {
	return u.courierRepo.Update(user)
}

func (u *courierUsecase) DeleteCourier(id uint) error {
	return u.courierRepo.Delete(id)
}
