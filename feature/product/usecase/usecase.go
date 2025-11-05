package usecase

import (
	"errors"
	"go-clean-api/config"
	"go-clean-api/domain"
	"go-clean-api/entity"
)

type productUsecase struct {
	productRepo domain.ProductRepository
	shopRepo    domain.ShopRepository
	cfg         config.ToolsConfig
}

func NewProductUsecase(
	productRepo domain.ProductRepository,
	shopRepo domain.ShopRepository,
	cfg config.ToolsConfig,
) domain.ProductUsecase {
	return &productUsecase{
		productRepo: productRepo,
		shopRepo:    shopRepo,
		cfg:         cfg,
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
		u.cfg.Logrus.WithError(err).Error("Failed to edit product")
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
	if u.shopRepo != nil {
		return u.shopRepo.FindByUserID(userID)
	}
	return u.productRepo.GetShopByUserID(userID)
}

func (u *productUsecase) GetProductsByShopID(shopID uint) ([]entity.Product, error) {
	return u.productRepo.GetProductsByShopID(shopID)
}
