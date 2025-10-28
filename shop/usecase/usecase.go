package usecase

import (
	"go-clean-api/domain"
	"go-clean-api/entity"
)

type shopUsecase struct {
	shopRepo    domain.ShopRepository
	cacheRepo   domain.ShopCacheRepository
	messageRepo domain.ShopMessageRepository
}

func NewShopUsecase(
	shopRepo domain.ShopRepository,
	cacheRepo domain.ShopCacheRepository,
	messageRepo domain.ShopMessageRepository,
) domain.ShopUsecase {
	return &shopUsecase{
		shopRepo:    shopRepo,
		cacheRepo:   cacheRepo,
		messageRepo: messageRepo,
	}
}

func (u *shopUsecase) CreateShop(shop *entity.Shop) error {
	if err := u.shopRepo.Create(shop); err != nil {
		return err
	}

	go func() {
		u.cacheRepo.SetShopCache(uint(shop.ID), shop)
	}()

	go func() {
		u.messageRepo.PublishShopCreated(shop)
	}()

	return nil
}

func (u *shopUsecase) GetAllShop() ([]entity.Shop, error) {
	return u.shopRepo.FindAll()
}

func (u *shopUsecase) GetShopByID(id uint) (*entity.Shop, error) {
	if shop, err := u.cacheRepo.GetShopCache(id); err == nil && shop != nil {
		return shop, nil
	}

	shop, err := u.shopRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if shop != nil {
		go func() {
			u.cacheRepo.SetShopCache(id, shop)
		}()
	}
	return shop, nil
}

func (u *shopUsecase) UpdateShop(shop *entity.Shop) error {
	if err := u.shopRepo.Update(shop); err != nil {
		return err
	}

	go func() {
		u.cacheRepo.SetShopCache(uint(shop.ID), shop)
	}()

	go func() {
		u.messageRepo.PublishShopUpdated(shop)
	}()

	return nil
}

func (u *shopUsecase) DeleteShop(id uint) error {
	if err := u.shopRepo.Delete(id); err != nil {
		return err
	}

	go func() {
		u.cacheRepo.DeleteShopCache(id)
	}()

	go func() {
		u.messageRepo.PublishShopDeleted(id)
	}()

	return nil
}

func (u *shopUsecase) CreateProduct(product *entity.Product) error {
	return u.shopRepo.CreateProduct(product)
}
