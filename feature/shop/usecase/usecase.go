package usecase

import (
	"errors"
	"go-clean-api/domain"
	"go-clean-api/entity"
)

type shopUsecase struct {
	shopRepo domain.ShopRepository
}

func NewShopUsecase(
	shopRepo domain.ShopRepository,
) domain.ShopUsecase {
	return &shopUsecase{
		shopRepo: shopRepo,
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
	return u.shopRepo.FindAll()
}

func (u *shopUsecase) GetShopByID(id uint) (*entity.Shop, error) {

	shop, err := u.shopRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	return shop, nil
}

func (u *shopUsecase) GetShopByUserID(userID uint) (*entity.Shop, error) {
	return u.shopRepo.FindByUserID(userID)
}

func (u *shopUsecase) UpdateShop(shop *entity.Shop) error {
	if err := u.shopRepo.Update(shop); err != nil {
		return err
	}

	return nil
}

func (u *shopUsecase) DeleteShop(id uint) error {
	if err := u.shopRepo.Delete(id); err != nil {
		return err
	}

	return nil
}
