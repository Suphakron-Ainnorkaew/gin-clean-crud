package courier

import (
	"go-clean-api/internal/courier/domain"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type CourierHandler struct {
	usecase CourierUsecase
}

func NewCourierHandler(e *echo.Echo, usecase CourierUsecase) {
	handler := &CourierHandler{
		usecase: usecase,
	}

	g := e.Group("/courier")
	g.POST("", handler.CreateCourier)
	g.GET("", handler.GetAllCourier)
	g.GET("/:id", handler.GetCourierByID)
	g.PUT("/:id", handler.UpdateCourier)
	g.DELETE("/:id", handler.DeleteCourier)
}

func (h *CourierHandler) parseID(c echo.Context) (uint, error) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return 0, err
	}
	return uint(id), nil
}

// POST /courier
func (h *CourierHandler) CreateCourier(c echo.Context) error {
	var courier domain.Courier

	if err := c.Bind(&courier); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	if err := h.usecase.CreateCourier(&courier); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, courier)
}

// GET /courier
func (h *CourierHandler) GetAllCourier(c echo.Context) error {
	courier, err := h.usecase.GetAllCourier()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, courier)
}

// GET /courier/:id
func (h *CourierHandler) GetCourierByID(c echo.Context) error {
	id, err := h.parseID(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid Courier ID"})
	}

	courier, err := h.usecase.GetCourierByID(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	if courier == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Courier not found"})
	}

	return c.JSON(http.StatusOK, courier)
}

// PUT /courier/:id
func (h *CourierHandler) UpdateCourier(c echo.Context) error {
	id, err := h.parseID(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid courier ID"})
	}

	courier, err := h.usecase.GetCourierByID(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	if courier == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "courier not found"})
	}

	if err := c.Bind(courier); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	if err := h.usecase.UpdateCourier(courier); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, courier)
}

// DELETE /users/:id
func (h *CourierHandler) DeleteCourier(c echo.Context) error {
	id, err := h.parseID(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user ID"})
	}

	if err := h.usecase.DeleteCourier(id); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.NoContent(http.StatusNoContent)
}
