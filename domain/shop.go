package domain

import (
	"go-clean-api/entity"

	"github.com/sirupsen/logrus"
)

type ShopUsecase interface {
	CreateShop(log *logrus.Entry, shop *entity.Shop) error
	GetAllShop(log *logrus.Entry) ([]entity.Shop, error)
	GetShopByID(log *logrus.Entry, id uint) (*entity.Shop, error)
	GetShopByUserID(log *logrus.Entry, userID uint) (*entity.Shop, error)
	UpdateShop(log *logrus.Entry, shop *entity.Shop) error
}

type ShopRepository interface {
	Create(shop *entity.Shop) error
	FindAll() ([]entity.Shop, error)
	FindByID(id uint) (*entity.Shop, error)
	FindByUserID(userID uint) (*entity.Shop, error)
	Update(shop *entity.Shop) error
}

// type ShopCacheRepository interface {
// 	SetShopCache(shopID uint, shop *entity.Shop) error
// 	GetShopCache(shopID uint) (*entity.Shop, error)
// 	DeleteShopCache(shopID uint) error
// 	SetShopSession(sessionID string, shopID uint) error
// 	GetShopSession(sessionID string) (uint, error)
// }

// type ShopMessageRepository interface {
// 	PublishShopCreated(shop *entity.Shop) error
// 	PublishShopUpdated(shop *entity.Shop) error
// 	PublishShopDeleted(shopID uint) error
// 	SubscribeShopEvents() error
// }
