package usecase

import (
	"errors"
	"go-clean-api/domain"
	"go-clean-api/entity"
)

type productUsecase struct {
	productRepo domain.ProductRepository
}

func NewProductUsecase(
	productRepo domain.ProductRepository,
) domain.ProductUsecase {
	return &productUsecase{
		productRepo: productRepo,
	}
}

func (u *productUsecase) CreateProduct(product *entity.Product) error {
	return u.productRepo.CreateProduct(product)
}

func (u *productUsecase) EditProduct(product *entity.Product) error {
	existing, err := u.productRepo.FindProductByID(uint(product.ID))
	if err != nil {
		return err
	}
	if existing == nil {
		return errors.New("product not found")
	}
	if existing.ShopID != product.ShopID {
		return errors.New("not allowed to edit this product")
	}

	existing.Product_name = product.Product_name
	existing.Price = product.Price
	existing.Stock = product.Stock

	if err := u.productRepo.EditProduct(existing); err != nil {
		return err
	}
	return nil
}

func (u *productUsecase) GetAllProduct() ([]entity.Product, error) {
	return u.productRepo.GetAllProduct()
}

func (u *productUsecase) FindProductByID(id uint) (*entity.Product, error) {
	return u.productRepo.FindProductByID(id)
}

func (u *productUsecase) UpdateProductStock(productID uint, quantity int) error {
	return u.productRepo.UpdateProductStock(productID, quantity)
}

func (u *productUsecase) GetShopByUserID(userID uint) (*entity.Shop, error) {
	return u.productRepo.GetShopByUserID(userID)
}
