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
	usecase domain.ShopUsecase
}

// NewHandler now accepts jwtSecret and a userFetcher so we can protect
// product-creation routes to "shop" users.
func NewHandler(e *echo.Group, usecase domain.ShopUsecase, jwtSecret string, userFetcher middleware.UserFetcher) *Handler {
	handler := &Handler{
		usecase: usecase,
	}

	// public shop routes
	e.POST("/shops", handler.CreateShop)
	e.GET("/shops", handler.GetAllShop)

	// protected routes under /shops that require authentication
	g := e.Group("/shops")
	g.Use(middleware.NewJWTAuth(jwtSecret))

	// example protected route: create product (only shop users)
	// Option 1: Check from database (more secure, slower)
	// g.POST("/products", handler.CreateProduct, middleware.RequireRole(userFetcher, entity.UserTypeShop))
	
	// Option 2: Check from JWT token (faster, recommended for most cases)
	g.POST("/products", handler.CreateProduct, middleware.RequireRoleFromJWT(entity.UserTypeShop))

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

	if err := c.Bind(&shop); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}
	if err := h.usecase.CreateShop(&shop); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, shop)
}

// GET /shops
func (h *Handler) GetAllShop(c echo.Context) error {
	shops, err := h.usecase.GetAllShop()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, shops)
}

// CreateProduct is a protected endpoint — only shop users can create products.
// The middleware RequireRoleFromJWT ensures that only users with type_user="shop" can access this.
func (h *Handler) CreateProduct(c echo.Context) error {
	// Get user_id from context (set by JWT middleware)
	userID := c.Get("user_id")
	typeUser := c.Get("type_user")
	
	// Parse product data
	var product entity.Product
	if err := c.Bind(&product); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}
	
	// Save product to database
	if err := h.usecase.CreateProduct(&product); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}
	
	// Return success with product info and user context
	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message": "Product created successfully",
		"product": product,
		"created_by_user_id": userID,
		"user_type": typeUser,
	})
}
