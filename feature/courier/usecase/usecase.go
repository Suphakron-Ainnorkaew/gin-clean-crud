package usecase

import (
	"go-clean-api/domain"
	"go-clean-api/entity"
)

type courierUsecase struct {
	courierRepo domain.CourierRepository
}

func NewCourierUsecase(courierRepo domain.CourierRepository) domain.CourierUsecase {
	return &courierUsecase{
		courierRepo: courierRepo,
	}
}

func (u *courierUsecase) CreateCourier(courier *entity.Courier) error {
	return u.courierRepo.Create(courier)
}

func (u *courierUsecase) GetCourierByID(id uint) (*entity.Courier, error) {
	return u.courierRepo.GetByID(id)
}

func (u *courierUsecase) GetAllCourier() ([]entity.Courier, error) {
	return u.courierRepo.GetAll()
}

func (u *courierUsecase) UpdateCourier(courier *entity.Courier) error {
	return u.courierRepo.Update(courier)
}

func (u *courierUsecase) DeleteCourier(id uint) error {
	return u.courierRepo.Delete(id)
}
