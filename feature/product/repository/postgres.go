package repository

import (
	"errors"
	"go-clean-api/entity"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type postgresProductRepository struct {
	db *gorm.DB
}

func NewPostgresProductRepository(db *gorm.DB) *postgresProductRepository {
	return &postgresProductRepository{db: db}
}

func (r *postgresProductRepository) CreateProduct(product *entity.Product) error {
	return r.db.Create(product).Error
}

func (r *postgresProductRepository) UpdateProduct(log *logrus.Entry, productID uint, shopID int, changes map[string]interface{}) (*entity.Product, error) {

	var product entity.Product
	product.ID = int(productID)

	tx := r.db.Model(&product).Where("shop_id = ?", shopID).Updates(changes)
	if tx.Error != nil {
		log.WithError(tx.Error).Error("GORM Updates failed")
		return nil, tx.Error
	}

	if tx.RowsAffected == 0 {

		var checkProduct entity.Product
		if err := r.db.First(&checkProduct, productID).Error; err != nil {
			return nil, errors.New("product not found")
		}

		if checkProduct.ShopID != shopID {
			return nil, errors.New("you do not have permission to edit this product")
		}

		return nil, errors.New("product not found or no changes detected")
	}

	if err := r.db.First(&product, productID).Error; err != nil {
		log.WithError(err).Error("Failed to fetch updated product after update")
		return nil, err
	}

	return &product, nil
}

func (r *postgresProductRepository) FindProductByID(id uint) (*entity.Product, error) {
	var p entity.Product
	err := r.db.First(&p, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &p, nil
}

func (r *postgresProductRepository) GetAllProduct() ([]entity.Product, error) {
	var products []entity.Product
	err := r.db.Find(&products).Error
	return products, err
}

func (r *postgresProductRepository) UpdateProductStock(productID uint, quantity int) error {
	return r.db.Model(&entity.Product{}).
		Where("id = ?", productID).
		UpdateColumn("stock", gorm.Expr("stock + ?", quantity)).Error
}

func (r *postgresProductRepository) GetShopByUserID(userID uint) (*entity.Shop, error) {
	var shop entity.Shop
	err := r.db.Where("user_id = ?", userID).First(&shop).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &shop, nil
}

func (r *postgresProductRepository) GetProductsByShopID(shopID uint) ([]entity.Product, error) {
	var products []entity.Product
	err := r.db.Where("shop_id = ?", shopID).Find(&products).Error
	return products, err
}
