package usecase

import (
	"go-clean-api/domain"
	"go-clean-api/entity"

	"github.com/labstack/gommon/log"
	"github.com/pkg/errors"
)

type productUsecase struct {
	productRepo domain.ProductRepository
	shopRepo    domain.ShopRepository
}

func NewProductUsecase(
	productRepo domain.ProductRepository,
	shopRepo domain.ShopRepository,
) domain.ProductUsecase {
	return &productUsecase{
		productRepo: productRepo,
		shopRepo:    shopRepo,
	}
}

func (u *productUsecase) CreateProduct(product *entity.Product) error {
	if product == nil {
		err := errors.New("[Usecase.CreateProduct]: product is nil")
		log.Warn(err)
		return err
	}

	if err := u.productRepo.CreateProduct(product); err != nil {
		err = errors.Wrap(err, "[Usecase.CreateProduct]: failed to create product")
		log.Error(err)
		return err
	}

	log.Info("[Usecase.CreateProduct]: product created successfully")
	return nil
}

func (u *productUsecase) EditProduct(product *entity.Product) error {
	existing, err := u.productRepo.FindProductByID(uint(product.ID))
	if err != nil {
		err = errors.Wrap(err, "[Usecase.EditProduct]: failed to edit product")
		log.Warn(err)
		return err
	}
	if existing == nil {
		err = errors.Wrap(err, "[Usecase.EditProduct]: product not found")
		log.Warn(err)
		return err
	}
	if existing.ShopID != product.ShopID {
		err = errors.Wrap(err, "[Usecase.EditProduct]: not allowed to edit this product")
		log.Warn(err)
		return err
	}

	existing.Product_name = product.Product_name
	existing.Price = product.Price
	existing.Stock = product.Stock

	if err := u.productRepo.EditProduct(existing); err != nil {
		err = errors.Wrap(err, "[Usecase.EditProduct]: failed to edit product")
		log.Error(err)
		return err
	}
	return nil
}

func (u *productUsecase) GetAllProduct() ([]entity.Product, error) {
	product, err := u.productRepo.GetAllProduct()

	if err != nil {
		err = errors.Wrap(err, "[Usecase.GetAllProduct]: failed to get product")
		log.Error(err)
		return nil, err
	}

	return product, err
}

func (u *productUsecase) FindProductByID(id uint) (*entity.Product, error) {
	product, err := u.productRepo.FindProductByID(id)

	if err != nil {
		err = errors.Wrap(err, "[Usecase.FindProductByID]: failed to get product id")
		log.Error(err)
		return product, err
	}
	return product, nil
}

func (u *productUsecase) UpdateProductStock(productID uint, quantity int) error {
	if err := u.productRepo.UpdateProductStock(productID, quantity); err != nil {
		err = errors.Wrap(err, "[Usecase.UpdateProductStock]: failed to update product stock")
		log.Error(err)
		return err
	}

	log.Infof("[Usecase.UpdateProductStock]: stock updated (productID=%d, quantity=%d)", productID, quantity)
	return nil
}

func (u *productUsecase) GetShopByUserID(userID uint) (*entity.Shop, error) {
	if u.shopRepo != nil {
		shop, err := u.shopRepo.FindByUserID(userID)
		if err != nil {
			err = errors.Wrapf(err, "[Usecase.GetShopByUserID]: failed to get shop by userID=%d via shopRepo", userID)
			log.Error(err)
			return nil, err
		}

		if shop == nil {
			log.Warnf("[Usecase.GetShopByUserID]: no shop found for userID=%d (via shopRepo)", userID)
		} else {
			log.Infof("[Usecase.GetShopByUserID]: found shop (ID=%d) via shopRepo", shop.ID)
		}
		return shop, nil
	}

	// fallback หากยังไม่ได้ inject shopRepo
	log.Warn("[Usecase.GetShopByUserID]: shopRepo is nil, fallback to productRepo")

	shop, err := u.productRepo.GetShopByUserID(userID)
	if err != nil {
		err = errors.Wrapf(err, "[Usecase.GetShopByUserID]: failed to get shop by userID=%d via productRepo", userID)
		log.Error(err)
		return nil, err
	}

	if shop == nil {
		log.Warnf("[Usecase.GetShopByUserID]: no shop found for userID=%d (via productRepo)", userID)
	} else {
		log.Infof("[Usecase.GetShopByUserID]: found shop (ID=%d) via productRepo", shop.ID)
	}

	return shop, nil
}

func (u *productUsecase) GetProductsByShopID(shopID uint) ([]entity.Product, error) {
	product, err := u.productRepo.GetProductsByShopID(shopID)
	if err != nil {
		err = errors.Wrap(err, "[Usecase.GetProductsByShopID]: failed to get product by shop id")
		log.Error(err)
		return product, err
	}
	return product, nil
}
