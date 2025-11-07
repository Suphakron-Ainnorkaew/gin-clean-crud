package domain

import (
	"go-clean-api/entity"

	"github.com/sirupsen/logrus"
)

type ProductUsecase interface {
	CreateProduct(log *logrus.Entry, product *entity.Product) error
	EditProduct(log *logrus.Entry, productID uint, shopID int, changes map[string]interface{}) (*entity.Product, error)
	GetAllProduct(log *logrus.Entry) ([]entity.Product, error)
	GetShopByUserID(log *logrus.Entry, userID uint) (*entity.Shop, error)
	FindProductByID(log *logrus.Entry, id uint) (*entity.Product, error)
	UpdateProductStock(log *logrus.Entry, productID uint, quantity int) error
	GetProductsByShopID(log *logrus.Entry, shopID uint) ([]entity.Product, error)
}

type ProductRepository interface {
	CreateProduct(product *entity.Product) error
	UpdateProduct(log *logrus.Entry, productID uint, shopID int, changes map[string]interface{}) (*entity.Product, error)
	GetAllProduct() ([]entity.Product, error)
	GetShopByUserID(userID uint) (*entity.Shop, error)
	//GetProductbyID(id uint) (*entity.Product, error)
	FindProductByID(id uint) (*entity.Product, error)
	UpdateProductStock(productID uint, quantity int) error
	GetProductsByShopID(shopID uint) ([]entity.Product, error)
}
