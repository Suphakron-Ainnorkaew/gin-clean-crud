package repository

import (
	"go-clean-api/entity"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type postgresProductRepository struct {
	db *gorm.DB
}

func NewPostgresProductRepository(db *gorm.DB) *postgresProductRepository {
	return &postgresProductRepository{db: db}
}

func (r *postgresProductRepository) CreateProduct(product *entity.Product) error {
	if err := r.db.Create(&product).Error; err != nil {
		return errors.Wrap(err, "[ShopRepository.CreateProduct]: unable to create product")
	}
	return nil
}

func (r *postgresProductRepository) EditProduct(product *entity.Product) error {
	if err := r.db.Model(&entity.Product{}).Where("id = ?", product.ID).Updates(&product).Error; err != nil {
		return errors.Wrap(err, "[ShopRepository.EditProduct]: unable to edit product")
	}
	return nil
}

func (r *postgresProductRepository) FindProductByID(id uint) (*entity.Product, error) {
	var p entity.Product

	if err := r.db.First(&p, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, errors.Wrap(err, "[ShopRepository.FindProductByID]: unable to find product id")
	}
	return &p, nil
}

func (r *postgresProductRepository) GetAllProduct() ([]entity.Product, error) {
	var products []entity.Product
	if err := r.db.Find(&products).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return products, errors.Wrap(err, "[ShopRepository.GetAllProduct]: unable to get product")
	}
	return products, nil
}

func (r *postgresProductRepository) UpdateProductStock(productID uint, quantity int) error {
	err := r.db.Model(&entity.Product{}).
		Where("id = ?", productID).
		UpdateColumn("stock", gorm.Expr("stock + ?", quantity)).Error

	if err != nil {
		return errors.Wrap(err, "[ProductRepository.UpdateProductStock]: unable to update product stock")
	}

	return nil
}

func (r *postgresProductRepository) GetShopByUserID(userID uint) (*entity.Shop, error) {
	var shop entity.Shop

	if err := r.db.Where("user_id = ?", userID).First(&shop).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, errors.Wrap(err, "[ShopRepository.GetShopByUserID]: unable to get product by user id")
	}
	return &shop, nil
}

func (r *postgresProductRepository) GetProductsByShopID(shopID uint) ([]entity.Product, error) {
	var products []entity.Product
	if err := r.db.Where("shop_id = ?", shopID).Find(&products).Error; err != nil {
		return products, errors.Wrap(err, "[ShopRepository.GetProductsByShopID]: unable to get product by shop id")
	}
	return products, nil
}
