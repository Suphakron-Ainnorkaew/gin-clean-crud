package usecase

import (
	"go-clean-api/config"
	"go-clean-api/domain"
	"go-clean-api/entity"
)

type courierUsecase struct {
	courierRepo domain.CourierRepository
	cfg         config.ToolsConfig
}

func NewCourierUsecase(
	courierRepo domain.CourierRepository,
	cfg config.ToolsConfig,
) domain.CourierUsecase {
	return &courierUsecase{
		courierRepo: courierRepo,
		cfg:         cfg,
	}
}

func (u *courierUsecase) CreateCourier(courier *entity.Courier) error {
	if err := u.courierRepo.Create(courier); err != nil {
		return err
	}
	return nil
}

func (u *courierUsecase) GetCourierByID(id uint) (*entity.Courier, error) {

	courier, err := u.courierRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	return courier, nil
}

func (u *courierUsecase) GetAllCourier() ([]entity.Courier, error) {

	courier, err := u.courierRepo.GetAll()
	if err != nil {
	}

	return courier, nil
}

func (u *courierUsecase) UpdateCourier(courier *entity.Courier) error {

	if err := u.courierRepo.Update(courier); err != nil {
		return err
	}

	return nil
}

func (u *courierUsecase) DeleteCourier(id uint) error {

	if err := u.courierRepo.Delete(id); err != nil {
	}
	return u.courierRepo.Delete(id)
}
