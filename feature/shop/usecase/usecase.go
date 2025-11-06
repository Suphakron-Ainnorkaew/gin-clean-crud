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

	log = log.WithField("user_id", shop.UserID)
	log.Info("CreateShop usecase started")

	if shop.UserID == 0 {
		log.Warn("CreateShop validation failed: UserID is required")
		return errors.New("user ID is required")
	}

	existingShop, err := u.shopRepo.FindByUserID(shop.UserID)
	if err != nil {
		log.WithError(err).Error("Failed to FindByUserID in repo")
		return err
	}
	if existingShop != nil {
		log.Warn("CreateShop failed: Shop for this user already exists")
		return errors.New("shop for this user already exists")
	}

	if err := u.shopRepo.Create(shop); err != nil {
		log.WithError(err).Error("Failed to create shop in repo")
		return err
	}

	log.WithField("shop_id", shop.ID).Info("Shop created successfully in usecase")
	return nil
}

func (u *shopUsecase) GetAllShop(log *logrus.Entry) ([]entity.Shop, error) {
	log.Info("GetAllShop usecase started")

	shop, err := u.shopRepo.FindAll()
	if err != nil {
		log.WithError(err).Error("Failed to get-all shop from repo")
		return nil, err
	}

	log.WithField("count", len(shop)).Info("GetAllShop usecase successful")
	return shop, nil
}

func (u *shopUsecase) GetShopByID(log *logrus.Entry, id uint) (*entity.Shop, error) {

	log = log.WithField("shop_id", id)
	log.Info("GetShopByID usecase started")

	shop, err := u.shopRepo.FindByID(id)
	if err != nil {
		log.WithError(err).Error("Failed to get shop id from repo")
		return nil, err
	}

	log.Info("GetShopByID usecase successful")
	return shop, nil
}

func (u *shopUsecase) GetShopByUserID(log *logrus.Entry, userID uint) (*entity.Shop, error) {

	log = log.WithField("user_id", userID)
	log.Info("GetShopByUserID usecase started")

	shop, err := u.shopRepo.FindByUserID(userID)
	if err != nil {
		log.WithError(err).Error("Failed to get shop by user id from repo")
		return nil, err
	}

	log.Info("GetShopByUserID usecase successful")
	return shop, nil
}

func (u *shopUsecase) UpdateShop(log *logrus.Entry, shop *entity.Shop) error {

	log = log.WithField("shop_id", shop.ID)
	log.Info("UpdateShop usecase started")

	if err := u.shopRepo.Update(shop); err != nil {
		log.WithError(err).Error("Failed to update shop in repo")
		return err
	}

	log.Info("UpdateShop usecase successful")
	return nil
}
