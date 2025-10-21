package shop

import (
	"go-clean-api/internal/shop/domain"
)

type ShopUsecase interface {
	CreateShop(shop *domain.Shop) error
	GetAllShops() ([]domain.Shop, error)
	GetShopByID(id uint) (*domain.Shop, error)
	UpdateShop(shop *domain.Shop) error
	DeleteShop(id uint) error
}

type shopUsecase struct {
	shopRepo ShopRepository
}

func NewShopUsecase(repo ShopRepository) ShopUsecase {
	return &shopUsecase{shopRepo: repo}
}

func (u *shopUsecase) CreateShop(shop *domain.Shop) error {
	return u.shopRepo.Create(shop)
}

func (u *shopUsecase) GetAllShops() ([]domain.Shop, error) {
	return u.shopRepo.FindAll()
}

func (u *shopUsecase) GetShopByID(id uint) (*domain.Shop, error) {
	return u.shopRepo.FindByID(id)
}

func (u *shopUsecase) UpdateShop(shop *domain.Shop) error {
	return u.shopRepo.Update(shop)
}

func (u *shopUsecase) DeleteShop(id uint) error {
	return u.shopRepo.Delete(id)
}
