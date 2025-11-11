package delivery

import (
	"go-clean-api/domain"
	"go-clean-api/entity"
	"go-clean-api/middleware"
	"go-clean-api/utils"
	"net/http"
	"strconv"

	"github.com/pkg/errors"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

type Handler struct {
	usecase domain.ProductUsecase
}

func NewHandler(e *echo.Group, usecase domain.ProductUsecase, jwtSecret string, userFetcher middleware.UserFetcher) *Handler {
	handler := &Handler{
		usecase: usecase,
	}

	e.GET("/products", handler.GetAllProduct)
	e.POST("/products", handler.CreateProduct, middleware.RequireRole(userFetcher, entity.UserTypeShop))
	e.PUT("/products/:id", handler.EditProduct, middleware.RequireRole(userFetcher, entity.UserTypeShop))
	e.GET("/shops/:id/products", handler.GetProductsByShopID)

	return handler
}

func (h *Handler) parseID(c echo.Context) (uint, error) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return 0, err
	}
	return uint(id), nil
}

func (h *Handler) GetAllProduct(c echo.Context) error {
	products, err := h.usecase.GetAllProduct()
	if err != nil {
		err = errors.Wrap(err, "[Handler.GetAllProduct]: failed to get products")
		log.Error(err)
		return c.JSON(http.StatusInternalServerError, entity.ErrorResponse{
			Message: utils.StandardError(err),
		})
	}
	return c.JSON(http.StatusOK, products)
}

func (h *Handler) CreateProduct(c echo.Context) error {
	userIDVal := c.Get("user_id")
	if userIDVal == nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "user_id not found in context",
		})
	}
	userID, ok := userIDVal.(uint)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "invalid user_id type",
		})
	}

	var product entity.Product
	if err := c.Bind(&product); err != nil {
		err = errors.Wrap(err, "[Handler.CreateProduct]: invalid request body")
		return c.JSON(http.StatusBadRequest, entity.ErrorResponse{
			Message: utils.StandardError(err),
		})
	}

	shop, err := h.usecase.GetShopByUserID(userID)
	if err != nil {
		err = errors.Wrap(err, "[Handler.CreateProduct]: failed to get product by shop id")
		log.Warn(err)
		return c.JSON(http.StatusInternalServerError, entity.ErrorResponse{
			Message: utils.StandardError(err),
		})
	}
	if shop == nil {
		err = errors.Wrap(err, "[Handler.CreateProduct]: Shop not found for this user Please create a shop first")
		log.Warn(err)
		return c.JSON(http.StatusNotFound, entity.ErrorResponse{
			Message: utils.StandardError(err),
		})
	}

	product.ShopID = shop.ID

	if err := h.usecase.CreateProduct(&product); err != nil {
		err = errors.Wrap(err, "[Handler.CreateProduct]: failed to create product")
		log.Error(err)
		return c.JSON(http.StatusInternalServerError, entity.ErrorResponse{
			Message: utils.StandardError(err),
		})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message": "Product created successfully",
		"product": product,
		"shop": map[string]interface{}{
			"id":   shop.ID,
			"name": shop.Name,
		},
	})
}

func (h *Handler) EditProduct(c echo.Context) error {
	id, err := h.parseID(c)
	if err != nil {
		err = errors.Wrap(err, "[Handler.EditProduct]: invalid product id")
		log.Error(err)
		return c.JSON(http.StatusBadRequest, entity.ErrorResponse{
			Message: utils.StandardError(err),
		})
	}

	userIDVal := c.Get("user_id")
	if userIDVal == nil {
		err = errors.Wrap(err, "[Handler.EditProduct]: user_id not found in context")
		log.Warn(err)
		return c.JSON(http.StatusUnauthorized, entity.ErrorResponse{
			Message: utils.StandardError(err),
		})
	}
	userID, ok := userIDVal.(uint)
	if !ok {
		err = errors.Wrap(err, "[Handler.EditProduct]: invalid user_id type")
		log.Warn(err)
		return c.JSON(http.StatusUnauthorized, entity.ErrorResponse{
			Message: utils.StandardError(err),
		})
	}

	var payload entity.Product
	if err := c.Bind(&payload); err != nil {
		err = errors.Wrap(err, "[Handler.EditProduct]: invalid request body")
		log.Warn(err)
		return c.JSON(http.StatusBadRequest, entity.ErrorResponse{
			Message: utils.StandardError(err),
		})
	}

	shop, err := h.usecase.GetShopByUserID(userID)
	if err != nil {
		err = errors.Wrap(err, "[Handler.EditProduct]: failed to get shop by user id")
		log.Warn(err)
		return c.JSON(http.StatusInternalServerError, entity.ErrorResponse{
			Message: utils.StandardError(err),
		})
	}
	if shop == nil {
		err = errors.Wrap(err, "[Handler.EditProduct]: shop not found for this user")
		log.Warn(err)
		return c.JSON(http.StatusNotFound, entity.ErrorResponse{
			Message: utils.StandardError(err),
		})
	}

	payload.ID = int(id)
	payload.ShopID = shop.ID

	if err := h.usecase.EditProduct(&payload); err != nil {
		err = errors.Wrap(err, "[Handler.EditProduct]: failed to edit product")
		log.Error(err)
		return c.JSON(http.StatusForbidden, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"message": "product updated", "product": payload})
}

func (h *Handler) GetProductByID(c echo.Context) error {
	id, err := h.parseID(c)
	if err != nil {
		err = errors.Wrap(err, "[Handler.GetProductByID]: invalid product id")
		log.Warn(err)
		return c.JSON(http.StatusBadRequest, entity.ErrorResponse{
			Message: utils.StandardError(err),
		})
	}

	p, err := h.usecase.FindProductByID(id)
	if err != nil {
		err = errors.Wrap(err, "[Handler.GetProductByID]: failed to edit product")
		return c.JSON(http.StatusInternalServerError, entity.ErrorResponse{
			Message: utils.StandardError(err),
		})
	}
	if p == nil {
		err = errors.Wrap(err, "[Handler.GetProductByID]: product not found")
		return c.JSON(http.StatusNotFound, entity.ErrorResponse{
			Message: utils.StandardError(err),
		})
	}
	return c.JSON(http.StatusOK, p)
}

func (h *Handler) GetProductsByShopID(c echo.Context) error {
	shopID, err := h.parseID(c)
	if err != nil {
		err = errors.Wrap(err, "[Handler.GetProductsByShopID]: invalid get product by shop id")
		log.Warn(err)
		return c.JSON(http.StatusBadRequest, entity.ErrorResponse{
			Message: utils.StandardError(err),
		})
	}

	products, err := h.usecase.GetProductsByShopID(shopID)
	if err != nil {
		err = errors.Wrap(err, "[Handler.GetProductsByShopID]: failed to get product by shop id")
		log.Error(err)
		return c.JSON(http.StatusInternalServerError, entity.ErrorResponse{
			Message: utils.StandardError(err),
		})
	}
	return c.JSON(http.StatusOK, products)
}
