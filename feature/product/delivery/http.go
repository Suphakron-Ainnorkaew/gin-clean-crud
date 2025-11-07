package delivery

import (
	"go-clean-api/config"
	"go-clean-api/domain"
	"go-clean-api/entity"
	"go-clean-api/middleware"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	usecase domain.ProductUsecase
	cfg     config.ToolsConfig
}

func NewHandler(e *echo.Group, usecase domain.ProductUsecase, cfg config.ToolsConfig, userFetcher middleware.UserFetcher) *Handler {

	e.Use(middleware.LoggingMiddleware(cfg.Logrus))

	handler := &Handler{
		usecase: usecase,
		cfg:     cfg,
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

	log := c.Get("logger").(*logrus.Entry)

	products, err := h.usecase.GetAllProduct(log)
	if err != nil {
		log.WithError(err).Error("Failed to get all product")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, products)
}

func (h *Handler) CreateProduct(c echo.Context) error {

	log := c.Get("logger").(*logrus.Entry)

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
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	shop, err := h.usecase.GetShopByUserID(log, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}
	if shop == nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "Shop not found for this user. Please create a shop first.",
		})
	}

	product.ShopID = shop.ID

	if err := h.usecase.CreateProduct(log, &product); err != nil {
		log.WithError(err).Error("failed to create product in CreateProduct")
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
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

	log := c.Get("logger").(*logrus.Entry)
	id, err := h.parseID(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid product id"})
	}

	log = log.WithField("product_id", id)

	userIDVal := c.Get("user_id")
	if userIDVal == nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "user_id not found in context"})
	}
	userID, ok := userIDVal.(uint)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid user_id type"})
	}
	log = log.WithField("user_id", userID)

	shop, err := h.usecase.GetShopByUserID(log, userID)
	if err != nil {
		log.WithError(err).Error("failed to get shop by user id")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	if shop == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "shop not found for this user"})
	}

	log = log.WithField("shop_id", shop.ID)

	var payload map[string]interface{}
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	delete(payload, "id")
	delete(payload, "shop_id")
	delete(payload, "created_at")
	delete(payload, "updated_at")
	delete(payload, "deleted_at")

	updatedProduct, err := h.usecase.EditProduct(log, id, shop.ID, payload)
	if err != nil {

		log.WithError(err).Error("failed to update product in EditProduct")
		return c.JSON(http.StatusForbidden, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"message": "product updated", "product": updatedProduct})
}

func (h *Handler) GetProductByID(c echo.Context) error {

	log := c.Get("logger").(*logrus.Entry)

	id, err := h.parseID(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid product id"})
	}

	p, err := h.usecase.FindProductByID(log, id)
	if err != nil {
		log.WithError(err).Error("failed to Get product id in GetProductByID")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	if p == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "product not found"})
	}
	return c.JSON(http.StatusOK, p)
}

func (h *Handler) GetProductsByShopID(c echo.Context) error {
	log := c.Get("logger").(*logrus.Entry)
	shopID, err := h.parseID(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid shop id"})
	}

	products, err := h.usecase.GetProductsByShopID(log, shopID)
	if err != nil {
		log.WithError(err).Error("Failed to Get product by shop id in GetProductsByShopID")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, products)
}
