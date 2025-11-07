package usecase

import (
	"go-clean-api/config"
	"go-clean-api/domain"
	"go-clean-api/entity"

	"github.com/sirupsen/logrus"
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

func (u *productUsecase) CreateProduct(log *logrus.Entry, product *entity.Product) error {
	if err := u.productRepo.CreateProduct(product); err != nil {
		log.WithError(err).Error("failed to create product in repo")
		return err
	}
	return nil
}

func (u *productUsecase) EditProduct(log *logrus.Entry, productID uint, shopID int, changes map[string]interface{}) (*entity.Product, error) {

	log = log.WithFields(logrus.Fields{
		"product_id": productID,
		"shop_id":    shopID,
	})

	updatedProduct, err := u.productRepo.UpdateProduct(log, productID, shopID, changes)
	if err != nil {
		log.WithError(err).Error("failed to update product in repo")
		return nil, err
	}
	return updatedProduct, nil
}

func (u *productUsecase) GetAllProduct(log *logrus.Entry) ([]entity.Product, error) {
	product, err := u.productRepo.GetAllProduct()
	if err != nil {
		log.WithError(err).Error("Failed to get-all product from repo")
		return nil, err
	}
	return product, nil
}

func (u *productUsecase) FindProductByID(log *logrus.Entry, id uint) (*entity.Product, error) {
	log = log.WithField("product_id", id)
	product, err := u.productRepo.FindProductByID(id)
	if err != nil {
		log.WithError(err).Error("Failed to get product id from repo")
		return nil, err
	}
	return product, nil
}

func (u *productUsecase) UpdateProductStock(log *logrus.Entry, productID uint, quantity int) error {
	log = log.WithField("product_id", productID)

	if err := u.productRepo.UpdateProductStock(productID, quantity); err != nil {
		log.WithError(err).Error("Failed to update stock product from repo")
		return err
	}

	return nil
}

func (u *productUsecase) GetShopByUserID(log *logrus.Entry, userID uint) (*entity.Shop, error) {

	log = log.WithField("user_id", userID)
	product, err := u.productRepo.GetShopByUserID(userID)

	if err != nil {
		log.WithError(err).Error("Failed to get shop by user id from repo")
		return nil, err
	}

	return product, nil
}

func (u *productUsecase) GetProductsByShopID(log *logrus.Entry, shopID uint) ([]entity.Product, error) {

	log = log.WithField("shop_id", shopID)
	product, err := u.productRepo.GetProductsByShopID(shopID)

	if err != nil {
		log.WithError(err).Error("failed to get product by shop id from repo")
		return nil, err
	}

	return product, nil
}
