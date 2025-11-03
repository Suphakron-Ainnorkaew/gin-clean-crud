package delivery

import (
	"go-clean-api/domain"
	"go-clean-api/entity"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	usecase domain.CourierUsecase
}

func NewHandler(e *echo.Group, usecase domain.CourierUsecase) *Handler {
	handler := &Handler{
		usecase: usecase,
	}

	e.POST("/courier", handler.CreateCourier)
	e.GET("/courier", handler.GETAllCourier)
	e.GET("/courier/:id", handler.GetCourierByID)
	e.PUT("/courier/:id", handler.UpdateCourier)
	e.DELETE("/courier/:id", handler.DeleteCourier)
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
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}
	if err := h.usecase.CreateCourier(&delivery); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, delivery)
}

// GET /Courier
func (h *Handler) GETAllCourier(c echo.Context) error {
	deliveries, err := h.usecase.GETAllCourier()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, deliveries)
}

// GET /Courier/:id
func (h *Handler) GetCourierByID(c echo.Context) error {
	id, err := h.parseID(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid delivery ID"})
	}

	delivery, err := h.usecase.GetCourierByID(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, delivery)
}

// PUT /Courier/:id
func (h *Handler) UpdateCourier(c echo.Context) error {
	id, err := h.parseID(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid delivery ID"})
	}

	var courier entity.Courier
	if err := c.Bind(&courier); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}
	courier.ID = int(id)

	if err := h.usecase.UpdateCourier(&courier); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.NoContent(http.StatusNoContent)
}

// DELETE /Courier/:id
func (h *Handler) DeleteCourier(c echo.Context) error {
	id, err := h.parseID(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid delivery ID"})
	}

	if err := h.usecase.DeleteCourier(id); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.NoContent(http.StatusNoContent)
}
