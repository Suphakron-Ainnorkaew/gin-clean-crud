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
	usecase domain.ShopUsecase
	cfg     config.ToolsConfig
}

func NewHandler(e *echo.Group, usecase domain.ShopUsecase, cfg config.ToolsConfig) *Handler {

	e.Use(middleware.LoggingMiddleware(cfg.Logrus))

	handler := &Handler{
		usecase: usecase,
		cfg:     cfg,
	}

	e.POST("/shops", handler.CreateShop)
	e.GET("/shops", handler.GetAllShop)
	e.PUT("/shops/:id", handler.UpdateShop)

	return handler
}

func (h *Handler) parseID(c echo.Context) (uint, error) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return 0, err
	}
	return uint(id), nil
}

func (h *Handler) CreateShop(c echo.Context) error {
	var shop entity.Shop

	log := c.Get("logger").(*logrus.Entry)

	userID, ok := c.Get("user_id").(uint)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid user id")
	}

	if err := c.Bind(&shop); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
	}

	shop.UserID = userID

	if err := h.usecase.CreateShop(log, &shop); err != nil {
		log.WithError(err).Error("Failed to create shop")
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, shop)
}

// GET /shops
func (h *Handler) GetAllShop(c echo.Context) error {

	log := c.Get("logger").(*logrus.Entry)

	shops, err := h.usecase.GetAllShop(log)
	if err != nil {
		log.WithError(err).Error("Failed to get all shops")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, shops)
}

func (h *Handler) UpdateShop(c echo.Context) error {

	log := c.Get("logger").(*logrus.Entry)

	id, err := h.parseID(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid shop id"})
	}

	var payload entity.Shop
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	payload.ID = int(id)

	if err := h.usecase.UpdateShop(log, &payload); err != nil {
		log.WithError(err).Error("Failed to update shop in UpdateShop")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"message": "shop updated", "shop": payload})
}
