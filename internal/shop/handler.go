package shop

import (
	"go-clean-api/internal/shop/domain"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type ShopHandler struct {
	usecase ShopUsecase
}

func NewShopHandler(e *echo.Echo, usecase ShopUsecase) {
	handler := &ShopHandler{
		usecase: usecase,
	}

	g := e.Group("/shops")
	g.POST("", handler.CreateShop)
	g.GET("", handler.GetAllShops)
	g.GET("/:id", handler.GetShopByID)
	g.PUT("/:id", handler.UpdateShop)
	g.DELETE("/:id", handler.DeleteShop)
}

func (h *ShopHandler) parseID(c echo.Context) (uint, error) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return 0, err
	}
	return uint(id), nil
}

// POST /users
func (h *ShopHandler) CreateShop(c echo.Context) error {
	var shop domain.Shop

	if err := c.Bind(&shop); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	if err := h.usecase.CreateShop(&shop); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, shop)
}

// GET /shops
func (h *ShopHandler) GetAllShops(c echo.Context) error {
	shops, err := h.usecase.GetAllShops()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, shops)
}

// GET /shops/:id
func (h *ShopHandler) GetShopByID(c echo.Context) error {
	id, err := h.parseID(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid shop ID"})
	}

	shop, err := h.usecase.GetShopByID(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	if shop == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Shops not found"})
	}

	return c.JSON(http.StatusOK, shop)
}

// PUT /shops/:id
func (h *ShopHandler) UpdateShop(c echo.Context) error {
	id, err := h.parseID(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user ID"})
	}

	shop, err := h.usecase.GetShopByID(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	if shop == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
	}

	if err := c.Bind(shop); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	if err := h.usecase.UpdateShop(shop); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, shop)
}

// DELETE /shops/:id
func (h *ShopHandler) DeleteShop(c echo.Context) error {
	id, err := h.parseID(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user ID"})
	}

	if err := h.usecase.DeleteShop(id); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.NoContent(http.StatusNoContent)
}
