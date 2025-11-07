package usecase

import (
	"errors"
	"go-clean-api/config"
	"go-clean-api/domain"
	"go-clean-api/entity"

	"github.com/sirupsen/logrus"
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

func (u *shopUsecase) CreateShop(log *logrus.Entry, shop *entity.Shop) error {

	if shop.UserID == 0 {
		return errors.New("user ID is required")
	}

	existingShop, err := u.shopRepo.FindByUserID(shop.UserID)
	if err != nil {
		return err
	}
	if existingShop != nil {
		return errors.New("shop for this user already exists")
	}

	if err := u.shopRepo.Create(shop); err != nil {
		log.WithError(err).Error("Failed to create shop in repo")
		return err
	}
	return nil
}

func (u *shopUsecase) GetAllShop(log *logrus.Entry) ([]entity.Shop, error) {

	shop, err := u.shopRepo.FindAll()
	if err != nil {
		log.WithError(err).Error("Failed to get-all shop from repo")
		return nil, err
	}

	return shop, nil
}

func (u *shopUsecase) GetShopByID(log *logrus.Entry, id uint) (*entity.Shop, error) {

	log = log.WithField("shop_id", id)

	shop, err := u.shopRepo.FindByID(id)
	if err != nil {
		log.WithError(err).Error("Failed to get shop id from repo")
		return nil, err
	}

	return shop, nil
}

func (u *shopUsecase) GetShopByUserID(log *logrus.Entry, userID uint) (*entity.Shop, error) {

	log = log.WithField("user_id", userID)

	shop, err := u.shopRepo.FindByUserID(userID)
	if err != nil {
		log.WithError(err).Error("Failed to get shop by user id from repo")
		return nil, err
	}
	return shop, nil
}

func (u *shopUsecase) UpdateShop(log *logrus.Entry, shop *entity.Shop) error {

	log = log.WithField("shop_id", shop.ID)

	if err := u.shopRepo.Update(shop); err != nil {
		log.WithError(err).Error("Failed to update shop in repo")
		return err
	}

	return nil
}
