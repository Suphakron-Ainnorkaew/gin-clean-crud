package repository

import (
	"go-clean-api/domain"
	"go-clean-api/entity"

	"github.com/pkg/errors"

	"gorm.io/gorm"
)

type postgresShopRepository struct {
	db *gorm.DB
}

func NewPostgresShopRepository(db *gorm.DB) domain.ShopRepository {
	return &postgresShopRepository{db: db}
}

func (r *postgresShopRepository) Create(user *entity.Shop) error {
	if err := r.db.Create(&user).Error; err != nil {
		return errors.Wrap(err, "[ShopRepository.Create]: unable to create shop")
	}
	return nil
}

func (r *postgresShopRepository) FindAll() ([]entity.Shop, error) {
	var shops []entity.Shop
	if err := r.db.Find(&shops).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return shops, errors.Wrap(err, "[ShopRepository.FindAll]: unable to find shop")
	}
	return shops, nil
}

func (r *postgresShopRepository) FindByID(id uint) (*entity.Shop, error) {
	var shop entity.Shop
	if err := r.db.First(&shop, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, errors.Wrap(err, "[ShopRepository.FindByID]: unable to find shop id")
	}
	return &shop, nil
}

func (r *postgresShopRepository) FindByUserID(userID uint) (*entity.Shop, error) {
	var shop entity.Shop
	if err := r.db.Where("user_id = ?", userID).First(&shop).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, errors.Wrap(err, "[ShopRepository.FindByUserID]: unable to find shop by user id")
	}
	return &shop, nil
}

func (r *postgresShopRepository) Update(shop *entity.Shop) error {
	if err := r.db.Model(&entity.Shop{}).Where("id = ?", shop.ID).Updates(&shop).Error; err != nil {
		return errors.Wrap(err, "[ShopRepository.Update]: unable to update shop")
	}
	return nil
}

func (r *postgresShopRepository) Delete(id uint) error {
	if err := r.db.Delete(&entity.Shop{}, id).Error; err != nil {
		return errors.Wrap(err, "[ShopRepository.Delete]: unable to Delete shop")
	}
	return nil
}
