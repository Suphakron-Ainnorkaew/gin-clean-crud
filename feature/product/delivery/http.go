package delivery

import (
	"go-clean-api/domain"
	"go-clean-api/entity"
	"go-clean-api/middleware"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	usecase domain.ProductUsecase
}

func NewHandler(e *echo.Group, usecase domain.ProductUsecase, jwtSecret string, userFetcher middleware.UserFetcher) *Handler {
	handler := &Handler{
		usecase: usecase,
	}

	e.Use(middleware.NewJWTAuth(jwtSecret))
	e.GET("/products", handler.GetAllProduct)
	e.POST("/products", handler.CreateProduct, middleware.RequireRole(userFetcher, entity.UserTypeShop))
	e.PUT("/products/:id", handler.EditProduct, middleware.RequireRole(userFetcher, entity.UserTypeShop))

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
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
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
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	shop, err := h.usecase.GetShopByUserID(userID)
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

	if err := h.usecase.CreateProduct(&product); err != nil {
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
	id, err := h.parseID(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid product id"})
	}

	userIDVal := c.Get("user_id")
	if userIDVal == nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "user_id not found in context"})
	}
	userID, ok := userIDVal.(uint)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid user_id type"})
	}

	var payload entity.Product
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	shop, err := h.usecase.GetShopByUserID(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	if shop == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "shop not found for this user"})
	}

	payload.ID = int(id)
	payload.ShopID = shop.ID

	if err := h.usecase.EditProduct(&payload); err != nil {
		return c.JSON(http.StatusForbidden, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"message": "product updated", "product": payload})
}
