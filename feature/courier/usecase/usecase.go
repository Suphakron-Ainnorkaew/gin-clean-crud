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

func NewCourierUsecase(courierRepo domain.CourierRepository, cfg config.ToolsConfig) domain.CourierUsecase {
	return &courierUsecase{
		courierRepo: courierRepo,
		cfg:         cfg,
	}
}

func (u *courierUsecase) CreateCourier(courier *entity.Courier) error {
	if err := u.courierRepo.Create(courier); err != nil {
		u.cfg.Logrus.WithError(err).Error("failed to create courier")
		return err
	}
	u.cfg.Logrus.Info("courier created")
	return nil
}

func (u *courierUsecase) GetCourierByID(id uint) (*entity.Courier, error) {
	courier, err := u.courierRepo.GetByID(id)
	if err != nil {
		u.cfg.Logrus.WithError(err).WithField("courierID", id).Error("failed to get courier by id")
		return nil, err
	}
	if courier == nil {
		u.cfg.Logrus.WithField("courierID", id).Warn("courier not found")
		return nil, nil
	}
	u.cfg.Logrus.WithField("courierID", id).Info("get courier success")
	return courier, nil
}

func (u *courierUsecase) GetAllCourier() ([]entity.Courier, error) {
	courier, err := u.courierRepo.GetAll()
	if err != nil {
		u.cfg.Logrus.WithError(err).Error("failed to get all courier")
		return nil, err
	}
	if courier == nil {
		u.cfg.Logrus.WithError(err).Error("courier all not found")
		return nil, nil
	}
	u.cfg.Logrus.WithError(err).Info("get all courier success")
	return courier, nil
}

func (u *courierUsecase) UpdateCourier(courier *entity.Courier) error {

	if err := u.courierRepo.Update(courier); err != nil {
		u.cfg.Logrus.WithError(err).WithField("courierID", courier.ID).Error("failed to update courier")
		return err
	}
	u.cfg.Logrus.WithField("courierID", courier.ID).Info("update courier success")
	return nil
}

func (u *courierUsecase) DeleteCourier(id uint) error {
	if err := u.courierRepo.Delete(id); err != nil {
		u.cfg.Logrus.WithError(err).WithField("courierID", id).Error("failed to delete courier")
		return err
	}
	u.cfg.Logrus.WithField("courierID", id).Info("delete courier success")
	return nil
}
