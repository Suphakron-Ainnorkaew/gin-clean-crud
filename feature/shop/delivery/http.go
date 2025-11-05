package delivery

import (
	"errors"
	"go-clean-api/config"
	"go-clean-api/domain"
	"go-clean-api/entity"
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

// POST /shops
func (h *Handler) CreateShop(c echo.Context) error {

	var shop entity.Shop

	log := h.cfg.Logrus.WithFields(logrus.Fields{
		"endpoint": "POST /shops",
		"method":   c.Request().Method,
		"path":     c.Path(),
	})

	log.Info("Request started")

	userID, ok := c.Get("user_id").(uint)
	if !ok {
		log.WithError(errors.New("invalid user id")).Error("Failed to create shop")
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid user id")
	}

	if err := c.Bind(&shop); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
	}

	shop.UserID = userID

	if err := h.usecase.CreateShop(&shop); err != nil {
		log.WithError(err).Error("Failed to create shop")
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	log.WithField("shop", shop).Info("Shop created successfully")
	return c.JSON(http.StatusCreated, shop)
}

// GET /shops
func (h *Handler) GetAllShop(c echo.Context) error {
	shops, err := h.usecase.GetAllShop()
	if err != nil {
		h.cfg.Logrus.WithError(err).Error("Failed to get all shops")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, shops)
}

func (h *Handler) UpdateShop(c echo.Context) error {
	id, err := h.parseID(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid shop id"})
	}

	var payload entity.Shop
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	payload.ID = int(id)

	if err := h.usecase.UpdateShop(&payload); err != nil {
		h.cfg.Logrus.WithError(err).Error("Failed to update shop")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"message": "shop updated", "shop": payload})
}
