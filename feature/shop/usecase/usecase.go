package usecase

import (
	"errors"
	"go-clean-api/config"
	"go-clean-api/domain"
	"go-clean-api/entity"
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
		return err
	}
	if existingShop != nil {
		return errors.New("shop for this user already exists")
	}

	if err := u.shopRepo.Create(shop); err != nil {
		return err
	}

	return nil
}

func (u *shopUsecase) GetAllShop() ([]entity.Shop, error) {

	shop, err := u.shopRepo.FindAll()
	if err != nil {
		return nil, err
	}
	return shop, nil
}

func (u *shopUsecase) GetShopByID(id uint) (*entity.Shop, error) {

	shop, err := u.shopRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	return shop, nil
}

func (u *shopUsecase) GetShopByUserID(userID uint) (*entity.Shop, error) {

	shop, err := u.shopRepo.FindByUserID(userID)
	if err != nil {
		return nil, err
	}
	return shop, nil
}

func (u *shopUsecase) UpdateShop(shop *entity.Shop) error {

	if err := u.shopRepo.Update(shop); err != nil {
		return err
	}
	return nil
}
