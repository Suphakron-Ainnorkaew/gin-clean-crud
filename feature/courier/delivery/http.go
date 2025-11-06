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
		h.cfg.Logrus.WithError(err).Warn("invalid request body for CreateCourier")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}
	h.cfg.Logrus.WithFields(logrus.Fields{"brand": delivery.Brand}).Info("CreateCourier request")
	if err := h.usecase.CreateCourier(&delivery); err != nil {
		h.cfg.Logrus.WithError(err).Error("CreateCourier failed")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	h.cfg.Logrus.WithFields(logrus.Fields{"courierID": delivery.ID}).Info("CreateCourier success")
	return c.JSON(http.StatusCreated, delivery)
}

// GET /Courier
func (h *Handler) GETAllCourier(c echo.Context) error {
	h.cfg.Logrus.Info("GETAllCourier request")
	deliveries, err := h.usecase.GetAllCourier()
	if err != nil {
		h.cfg.Logrus.WithError(err).Error("GETAllCourier failed")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	h.cfg.Logrus.WithField("count", len(deliveries)).Info("GETAllCourier success")
	return c.JSON(http.StatusOK, deliveries)
}

// GET /Courier/:id
func (h *Handler) GetCourierByID(c echo.Context) error {
	id, err := h.parseID(c)
	if err != nil {
		h.cfg.Logrus.WithError(err).Warn("invalid courier id in GetCourierByID")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid delivery ID"})
	}

	delivery, err := h.usecase.GetCourierByID(id)
	if err != nil {
		h.cfg.Logrus.WithError(err).Error("GetCourierByID failed")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	if delivery == nil {
		h.cfg.Logrus.WithField("courierID", id).Warn("courier not found")
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Courier not found"})
	}
	h.cfg.Logrus.WithField("courierID", id).Info("GetCourierByID success")
	return c.JSON(http.StatusOK, delivery)
}

// PUT /Courier/:id
func (h *Handler) UpdateCourier(c echo.Context) error {
	id, err := h.parseID(c)
	if err != nil {
		h.cfg.Logrus.WithError(err).Warn("invalid courier id in UpdateCourier")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid delivery ID"})
	}

	var courier entity.Courier
	if err := c.Bind(&courier); err != nil {
		h.cfg.Logrus.WithError(err).Warn("invalid request body for UpdateCourier")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}
	courier.ID = int(id)

	if err := h.usecase.UpdateCourier(&courier); err != nil {
		h.cfg.Logrus.WithError(err).Error("UpdateCourier failed")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	h.cfg.Logrus.WithField("courierID", id).Info("UpdateCourier success")
	return c.NoContent(http.StatusNoContent)
}

// DELETE /Courier/:id
func (h *Handler) DeleteCourier(c echo.Context) error {
	id, err := h.parseID(c)
	if err != nil {
		h.cfg.Logrus.WithError(err).Warn("invalid courier id in DeleteCourier")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid delivery ID"})
	}

	if err := h.usecase.DeleteCourier(id); err != nil {
		h.cfg.Logrus.WithError(err).Error("DeleteCourier failed")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	h.cfg.Logrus.WithField("courierID", id).Info("DeleteCourier success")
	return c.NoContent(http.StatusNoContent)
}
