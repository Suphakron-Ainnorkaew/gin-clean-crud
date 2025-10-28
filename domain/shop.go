package domain

import (
	"go-clean-api/entity"
)

type ShopUsecase interface {
	CreateShop(shop *entity.Shop) error
	GetAllShop() ([]entity.Shop, error)
	GetShopByID(id uint) (*entity.Shop, error)
	UpdateShop(shop *entity.Shop) error
	DeleteShop(id uint) error
	CreateProduct(product *entity.Product) error
}

type ShopRepository interface {
	Create(shop *entity.Shop) error
	FindAll() ([]entity.Shop, error)
	FindByID(id uint) (*entity.Shop, error)
	Update(shop *entity.Shop) error
	Delete(id uint) error
	CreateProduct(product *entity.Product) error
}

type ShopCacheRepository interface {
	SetShopCache(shopID uint, shop *entity.Shop) error
	GetShopCache(shopID uint) (*entity.Shop, error)
	DeleteShopCache(shopID uint) error
	SetShopSession(sessionID string, shopID uint) error
	GetShopSession(sessionID string) (uint, error)
}

type ShopMessageRepository interface {
	PublishShopCreated(shop *entity.Shop) error
	PublishShopUpdated(shop *entity.Shop) error
	PublishShopDeleted(shopID uint) error
	SubscribeShopEvents() error
}
