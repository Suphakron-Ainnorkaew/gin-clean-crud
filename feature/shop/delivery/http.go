package delivery

import (
	"go-clean-api/config"
	"go-clean-api/domain"
	"go-clean-api/entity"
	"go-clean-api/utils"
	"net/http"
	"strconv"

	"github.com/labstack/gommon/log"
	"github.com/pkg/errors"

	"github.com/labstack/echo/v4"
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
	e.GET("/shops/:id", handler.GetShopByID)
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

	userID, ok := c.Get("user_id").(uint)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid user id")
	}

	if err := c.Bind(&shop); err != nil {
		err = errors.Wrap(err, "[Handler.CreateShop]: invalid request body")
		log.Warn(err)
		return c.JSON(http.StatusBadRequest, entity.ErrorResponse{
			Message: utils.StandardError(err),
		})
	}

	shop.UserID = userID

	if err := h.usecase.CreateShop(&shop); err != nil {
		err = errors.Wrap(err, "[Handler.CreateShop]: failed to create shop")
		log.Warn(err)
		return c.JSON(http.StatusInternalServerError, entity.ErrorResponse{
			Message: utils.StandardError(err),
		})
	}
	return c.JSON(http.StatusCreated, shop)
}

// GET /shops
func (h *Handler) GetAllShop(c echo.Context) error {

	shops, err := h.usecase.GetAllShop()
	if err != nil {
		err = errors.Wrap(err, "[Handler.GetAllShop]: failed to get shop")
		log.Warn(err)
		return c.JSON(http.StatusInternalServerError, entity.ErrorResponse{
			Message: utils.StandardError(err),
		})
	}

	return c.JSON(http.StatusOK, shops)
}

func (h *Handler) UpdateShop(c echo.Context) error {

	id, err := h.parseID(c)
	if err != nil {
		err = errors.Wrap(err, "[Handler.UpdateShop]: invalid shop id")
		return c.JSON(http.StatusBadRequest, entity.ErrorResponse{
			Message: utils.StandardError(err),
		})
	}

	var payload entity.Shop
	if err := c.Bind(&payload); err != nil {
		err = errors.Wrap(err, "[Handler.UpdateShop]: invalid request body")
		return c.JSON(http.StatusBadRequest, entity.ErrorResponse{
			Message: utils.StandardError(err),
		})
	}

	payload.ID = int(id)

	if err := h.usecase.UpdateShop(&payload); err != nil {
		err = errors.Wrap(err, "[Handler.UpdateShop]: failed to update shop")
		return c.JSON(http.StatusInternalServerError, entity.ErrorResponse{
			Message: utils.StandardError(err),
		})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{"message": "shop updated", "shop": payload})
}

func (h *Handler) GetShopByID(c echo.Context) error {

	id, err := h.parseID(c)
	if err != nil {
		err = errors.Wrap(err, "[Handler.GetShopByID]: invalid shop id")
		log.Warn(err)
		return c.JSON(http.StatusBadRequest, entity.ErrorResponse{
			Message: utils.StandardError(err),
		})
	}

	shop, err := h.usecase.GetShopByID(id)

	if err != nil {
		err = errors.Wrap(err, "[Handler.GetShopByID]: failed to get shop id")
		log.Warn(err)
		return c.JSON(http.StatusInternalServerError, entity.ErrorResponse{
			Message: utils.StandardError(err),
		})
	}

	if shop == nil {
		err = errors.Wrap(err, "[Handler.GetShopByID]: shop not found")
		log.Warn(err)
		return c.JSON(http.StatusNotFound, entity.ErrorResponse{
			Message: utils.StandardError(err),
		})
	}

	return c.JSON(http.StatusOK, shop)
}
