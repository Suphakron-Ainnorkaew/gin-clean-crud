package repository

import (
	"go-clean-api/entity"

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

func (r *postgresProductRepository) EditProduct(product *entity.Product) error {
	return r.db.Save(product).Error
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