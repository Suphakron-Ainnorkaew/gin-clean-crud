package delivery

import (
	"go-clean-api/config"
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
	usecase domain.CourierUsecase
	cfg     config.ToolsConfig
}

func NewHandler(e *echo.Group, usecase domain.CourierUsecase, cfg config.ToolsConfig) *Handler {

	handler := &Handler{
		usecase: usecase,
		cfg:     cfg,
	}

	e.POST("/courier", handler.CreateCourier, middleware.RequireRoleFromJWT(entity.UserTypeAdmin))
	e.GET("/courier", handler.GETAllCourier)
	e.GET("/courier/:id", handler.GetCourierByID)
	e.PUT("/courier/:id", handler.UpdateCourier, middleware.RequireRoleFromJWT(entity.UserTypeAdmin))
	e.DELETE("/courier/:id", handler.DeleteCourier, middleware.RequireRoleFromJWT(entity.UserTypeAdmin))
	return handler
}

func (h *Handler) parseID(c echo.Context) (uint, error) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return 0, err
	}
	return uint(id), nil
}

// POST /Courier
func (h *Handler) CreateCourier(c echo.Context) error {

	var delivery entity.Courier

	if err := c.Bind(&delivery); err != nil {
		err = errors.Wrap(err, "[Handler.CreateCourier]: invalid request body")
		log.Warn(err)
		return c.JSON(http.StatusBadRequest, entity.ErrorResponse{
			Message: utils.StandardError(err),
		})
	}
	if err := h.usecase.CreateCourier(&delivery); err != nil {
		err = errors.Wrap(err, "[Handler.CreateCourier]: failed to create courier")
		log.Error(err)
		return c.JSON(http.StatusInternalServerError, entity.ErrorResponse{
			Message: utils.StandardError(err),
		})
	}
	return c.JSON(http.StatusCreated, delivery)
}

// GET /Courier
func (h *Handler) GETAllCourier(c echo.Context) error {

	deliveries, err := h.usecase.GetAllCourier()
	if err != nil {
		err = errors.Wrap(err, "[Handler.GETAllCourier]: failed to get courier")
		log.Error(err)
		return c.JSON(http.StatusInternalServerError, entity.ErrorResponse{
			Message: utils.StandardError(err),
		})
	}
	return c.JSON(http.StatusOK, deliveries)
}

// GET /Courier/:id
func (h *Handler) GetCourierByID(c echo.Context) error {

	id, err := h.parseID(c)
	if err != nil {
		err = errors.Wrap(err, "[Handler.GetCourierByID]: invalid courier ID")
		log.Warn(err)
		return c.JSON(http.StatusBadRequest, entity.ErrorResponse{
			Message: utils.StandardError(err),
		})
	}

	delivery, err := h.usecase.GetCourierByID(id)
	if err != nil {
		err = errors.Wrap(err, "[Handler.GetCourierByID]: failed to get courier id")
		log.Error(err)
		return c.JSON(http.StatusInternalServerError, entity.ErrorResponse{
			Message: utils.StandardError(err),
		})
	}
	return c.JSON(http.StatusOK, delivery)
}

// PUT /Courier/:id
func (h *Handler) UpdateCourier(c echo.Context) error {
	id, err := h.parseID(c)
	if err != nil {
		err = errors.Wrap(err, "[Handler.UpdateCourier]: failed to parse ID")
		log.Warn(err)
		return c.JSON(http.StatusBadRequest, entity.ErrorResponse{
			Message: utils.StandardError(err),
		})
	}

	var courier entity.Courier
	if err := c.Bind(&courier); err != nil {
		err = errors.Wrap(err, "[Handler.UpdateCourier]: invalid request body")
		log.Warn(err)
		return c.JSON(http.StatusBadRequest, entity.ErrorResponse{
			Message: utils.StandardError(err),
		})
	}
	courier.ID = int(id)

	if err := h.usecase.UpdateCourier(&courier); err != nil {
		err = errors.Wrap(err, "[Handler.UpdateCourier]: failed to update courier")
		log.Error(err)
		return c.JSON(http.StatusInternalServerError, entity.ErrorResponse{
			Message: utils.StandardError(err),
		})
	}
	return c.NoContent(http.StatusNoContent)
}

// DELETE /Courier/:id
func (h *Handler) DeleteCourier(c echo.Context) error {
	id, err := h.parseID(c)
	if err != nil {
		err = errors.Wrap(err, "[Handler.DeleteCourier]: invalid courier ID")
		log.Warn(err)
		return c.JSON(http.StatusBadRequest, entity.ErrorResponse{
			Message: utils.StandardError(err),
		})
	}

	if err := h.usecase.DeleteCourier(id); err != nil {
		err = errors.Wrap(err, "[Handler.DeleteCourier]: failed to delete courier")
		log.Error(err)
		return c.JSON(http.StatusInternalServerError, entity.ErrorResponse{
			Message: utils.StandardError(err),
		})
	}
	return c.NoContent(http.StatusNoContent)
}
