package usecase

import (
	"go-clean-api/config"
	"go-clean-api/domain"
	"go-clean-api/entity"

	"github.com/labstack/gommon/log"
	"github.com/pkg/errors"
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
		err = errors.Wrap(err, "[Usecase.CreateCourier]: failed to create courier")
		log.Error(err)
		return err
	}
	return nil
}

func (u *courierUsecase) GetCourierByID(id uint) (*entity.Courier, error) {

	courier, err := u.courierRepo.GetByID(id)
	if err != nil {
		err = errors.Wrap(err, "[Usecase.GetCourierByID]: failed to get courier id")
		log.Error(err)
		return nil, err
	}

	return courier, nil
}

func (u *courierUsecase) GetAllCourier() ([]entity.Courier, error) {

	courier, err := u.courierRepo.GetAll()
	if err != nil {
		err = errors.Wrap(err, "[Usecase.GetAllCourier]: failed to get courier")
		log.Error(err)
		return courier, err
	}

	return courier, nil
}

func (u *courierUsecase) UpdateCourier(courier *entity.Courier) error {

	if err := u.courierRepo.Update(courier); err != nil {
		err = errors.Wrap(err, "[Usecase.UpdateCourier]: failed to update courier")
		log.Error(err)
		return err
	}

	return nil
}

func (u *courierUsecase) DeleteCourier(id uint) error {

	if err := u.courierRepo.Delete(id); err != nil {
		err = errors.Wrap(err, "[Usecase.DeleteCourier]: failed to delete courier")
		log.Error(err)
		return err
	}
	return nil
}
