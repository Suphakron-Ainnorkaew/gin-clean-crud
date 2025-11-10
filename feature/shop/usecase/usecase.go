package usecase

import (
	"go-clean-api/config"
	"go-clean-api/domain"
	"go-clean-api/entity"

	"github.com/labstack/gommon/log"
	"github.com/pkg/errors"
)

type shopUsecase struct {
	shopRepo domain.ShopRepository
	cfg      config.ToolsConfig
}

func NewShopUsecase(
	shopRepo domain.ShopRepository,
	cfg config.ToolsConfig,
) domain.ShopUsecase {
	return &shopUsecase{
		shopRepo: shopRepo,
		cfg:      cfg,
	}
}

func (u *shopUsecase) CreateShop(shop *entity.Shop) error {

	if shop.UserID == 0 {
		return errors.New("user ID is required")
	}

	existingShop, err := u.shopRepo.FindByUserID(shop.UserID)
	if err != nil {
		err := errors.Wrap(err, "[Usecase.CreateShop]: failed to find shop id")
		log.Warn(err)
		return err
	}
	if existingShop != nil {
		err = errors.Wrap(err, "[Usecase.CreateShop]: shop for this user already exists")
		log.Warn(err)
		return err
	}

	if err := u.shopRepo.Create(shop); err != nil {
		err = errors.Wrap(err, "[Usecase.CreateShop]: failed to create shop")
		log.Error(err)
		return err
	}

	return nil
}

func (u *shopUsecase) GetAllShop() ([]entity.Shop, error) {

	shop, err := u.shopRepo.FindAll()
	if err != nil {
		err := errors.Wrap(err, "[Usecase.GetAllShop]: failed to get shop")
		log.Error(err)
		return nil, err
	}
	return shop, nil
}

func (u *shopUsecase) GetShopByID(id uint) (*entity.Shop, error) {

	shop, err := u.shopRepo.FindByID(id)
	if err != nil {
		err := errors.Wrap(err, "[Usecase.GetShopByID]: failed to get shop id")
		log.Error(err)
		return nil, err
	}
	return shop, nil
}

func (u *shopUsecase) GetShopByUserID(userID uint) (*entity.Shop, error) {

	shop, err := u.shopRepo.FindByUserID(userID)
	if err != nil {
		err := errors.Wrap(err, "[Usecase.GetShopByUserID]: failed to get shop by user id")
		log.Error(err)
		return nil, err
	}
	return shop, nil
}

func (u *shopUsecase) UpdateShop(shop *entity.Shop) error {

	if err := u.shopRepo.Update(shop); err != nil {
		err := errors.Wrap(err, "[Usecase.UpdateShop]: failed to update shop")
		log.Error(err)
		return err
	}
	return nil
}
