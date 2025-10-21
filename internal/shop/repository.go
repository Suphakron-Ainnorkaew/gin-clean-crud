package shop

import (
	"go-clean-api/internal/shop/domain"

	"gorm.io/gorm"
)

type ShopRepository interface {
	Create(shop *domain.Shop) error
	FindAll() ([]domain.Shop, error)
	FindByID(id uint) (*domain.Shop, error)
	Update(shop *domain.Shop) error
	Delete(id uint) error
}

type gormShopRepository struct {
	db *gorm.DB
}

func NewGormShopRepository(db *gorm.DB) ShopRepository {
	return &gormShopRepository{db: db}
}

func (r *gormShopRepository) Create(shop *domain.Shop) error {
	return r.db.Create(shop).Error
}

func (r *gormShopRepository) FindAll() ([]domain.Shop, error) {
	var shops []domain.Shop
	err := r.db.Debug().Find(&shops).Error
	return shops, err
}

func (r *gormShopRepository) FindByID(id uint) (*domain.Shop, error) {
	var shop domain.Shop
	err := r.db.First(&shop, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &shop, nil
}

func (r *gormShopRepository) Update(shop *domain.Shop) error {
	return r.db.Save(shop).Error
}

func (r *gormShopRepository) Delete(id uint) error {
	return r.db.Delete(&domain.Shop{}, id).Error
}
