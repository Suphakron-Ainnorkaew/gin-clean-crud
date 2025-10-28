package repository

import (
	"go-clean-api/domain"
	"go-clean-api/entity"

	"gorm.io/gorm"
)

type postgresShopRepository struct {
	db *gorm.DB
}

func NewPostgresShopRepository(db *gorm.DB) domain.ShopRepository {
	return &postgresShopRepository{db: db}
}

func (r *postgresShopRepository) Create(user *entity.Shop) error {
	return r.db.Create(user).Error
}

func (r *postgresShopRepository) FindAll() ([]entity.Shop, error) {
	var shops []entity.Shop
	err := r.db.Find(&shops).Error
	return shops, err
}

func (r *postgresShopRepository) FindByID(id uint) (*entity.Shop, error) {
	var shop entity.Shop
	err := r.db.First(&shop, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &shop, nil
}

func (r *postgresShopRepository) Update(shop *entity.Shop) error {
	return r.db.Save(shop).Error
}

func (r *postgresShopRepository) Delete(id uint) error {
	return r.db.Delete(&entity.Shop{}, id).Error
}

func (r *postgresShopRepository) CreateProduct(product *entity.Product) error {
	return r.db.Create(product).Error
}
