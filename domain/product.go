package domain

import (
	"go-clean-api/entity"
)

type ProductUsecase interface {
	CreateProduct(product *entity.Product) error
	EditProduct(product *entity.Product) error
	GetAllProduct() ([]entity.Product, error)
	GetShopByUserID(userID uint) (*entity.Shop, error)
	FindProductByID(id uint) (*entity.Product, error)
	UpdateProductStock(productID uint, quantity int) error
}

type ProductRepository interface {
	CreateProduct(product *entity.Product) error
	EditProduct(product *entity.Product) error
	GetAllProduct() ([]entity.Product, error)
	GetShopByUserID(userID uint) (*entity.Shop, error)
	//GetProductbyID(id uint) (*entity.Product, error)
	FindProductByID(id uint) (*entity.Product, error)
	UpdateProductStock(productID uint, quantity int) error
}
