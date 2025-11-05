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
	usecase        domain.OrderUsecase
}

func NewHandler(e *echo.Group, usecase domain.OrderUsecase, jwtSecret string, userFetcher middleware.UserFetcher) *Handler {
	handler := &Handler{
		usecase: usecase,
	}

	e.POST("/orders", handler.CreateOrder)
	e.GET("/orders", handler.GetMyOrders)
	e.GET("/orders/:id", handler.GetOrderByID)

	e.PATCH("/orders/:id/status", handler.UpdateOrderStatus, middleware.RequireRole(userFetcher, entity.UserTypeShop))

	e.PATCH("/orders/:id/payment", handler.UpdatePaymentStatus, middleware.RequireRole(userFetcher, entity.UserTypeGeneral))

	// ร้านค้า ดูรายละเอียดคำสั่งซื้อของร้านค้าตัวเอง
	e.GET("/shops/orders", handler.GetShopOrders, middleware.RequireRole(userFetcher, entity.UserTypeShop))
	e.PATCH("/shops/orders/:id/status", handler.UpdateShopOrderStatus, middleware.RequireRole(userFetcher, entity.UserTypeShop))
	e.PATCH("/shops/orders/:id/cancel", handler.CancelShopOrder, middleware.RequireRole(userFetcher, entity.UserTypeShop))

	return handler

}

type CreateOrderRequest struct {
	ShopID    int                `json:"shop_id" validate:"required"`
	CourierID int                `json:"courier_id" validate:"required"`
	Items     []OrderItemRequest `json:"items" validate:"required,min=1"`
}

type OrderItemRequest struct {
	ProductID int `json:"product_id" validate:"required"`
	Quantity  int `json:"quantity" validate:"required,min=1"`
}

func (h *Handler) CreateOrder(c echo.Context) error {
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

	typeUserVal := c.Get("type_user")
	if typeUserVal == nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "type_user not found in context",
		})
	}
	typeUser, ok := typeUserVal.(string)
	if !ok || typeUser != "general" {
		return c.JSON(http.StatusForbidden, map[string]string{
			"error": "only general can create orders",
		})
	}

	var req CreateOrderRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	if len(req.Items) == 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Order must have at least one item",
		})
	}

	orderItems := make([]entity.OrderItem, len(req.Items))
	for i, item := range req.Items {
		orderItems[i] = entity.OrderItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
		}
	}

	order := &entity.Order{
		UserID:        int(userID),
		ShopID:        req.ShopID,
		CourierID:     req.CourierID,
		PaymentStatus: entity.PaymentStatusPending,
		Status:        entity.OrderStatusPending,
	}

	if err := h.usecase.CreateOrder(order, orderItems); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message":  "Order created successfully",
		"order_id": order.ID,
		"total":    order.Total,
	})
}

func (h *Handler) GetMyOrders(c echo.Context) error {
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

	orders, err := h.usecase.GetOrdersByUserID(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"orders": orders,
		"count":  len(orders),
	})
}

// GET /orders/:id
func (h *Handler) GetOrderByID(c echo.Context) error {
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

	typeUserVal := c.Get("type_user")
	if typeUserVal == nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "type_user not found in context",
		})
	}
	typeUser, ok := typeUserVal.(string)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "invalid type_user",
		})
	}

	orderID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid order ID",
		})
	}

	order, err := h.usecase.GetOrderByID(uint(orderID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}
	if order == nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "Order not found",
		})
	}

	// ตรวจสอบ authorization ตาม user type
	userTypeEnum := entity.UserType(typeUser)
	hasPermission := false

	switch userTypeEnum {
	case entity.UserTypeGeneral:
		// General user ต้องเป็นผู้สั่งซื้อ
		hasPermission = order.UserID == int(userID)
	case entity.UserTypeShop:
		// Shop ต้องเป็นเจ้าของร้าน (ใช้ Shop ที่ preload ไว้แล้ว)
		hasPermission = order.Shop.UserID == userID
	default:
		hasPermission = false
	}

	if !hasPermission {
		return c.JSON(http.StatusForbidden, map[string]string{
			"error": "You don't have permission to view this order",
		})
	}

	return c.JSON(http.StatusOK, order)
}

// PATCH /orders/:id/status - อัพเดทสถานะคำสั่งซื้อ
func (h *Handler) UpdateOrderStatus(c echo.Context) error {
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

	orderID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid order ID",
		})
	}

	var req struct {
		Status string `json:"status" validate:"required"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	status := entity.OrderStatus(req.Status)
	if status != entity.OrderStatusPending &&
		status != entity.OrderStatusShipped &&
		status != entity.OrderStatusDelivered &&
		status != entity.OrderStatusCancelled {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid status. Must be: pending, shipped, delivered, or cancelled",
		})
	}

	if err := h.usecase.UpdateOrderStatus(uint(orderID), status, userID); err != nil {
		if err.Error() == "order not found" {
			return c.JSON(http.StatusNotFound, map[string]string{
				"error": err.Error(),
			})
		}
		if err.Error() == "you don't have permission to update this order status" {
			return c.JSON(http.StatusForbidden, map[string]string{
				"error": err.Error(),
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Order status updated successfully",
		"status":  status,
	})
}

// PATCH /orders/:id/payment - อัพเดทสถานะการชำระเงิน
func (h *Handler) UpdatePaymentStatus(c echo.Context) error {
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

	orderID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid order ID",
		})
	}

	var req struct {
		PaymentStatus string `json:"payment_status" validate:"required"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	paymentStatus := entity.PaymentStatus(req.PaymentStatus)
	if paymentStatus != entity.PaymentStatusPending &&
		paymentStatus != entity.PaymentStatusComplete {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid payment status. Must be: pending or complete",
		})
	}

	if err := h.usecase.UpdatePaymentStatus(uint(orderID), paymentStatus, userID); err != nil {
		if err.Error() == "order not found" {
			return c.JSON(http.StatusNotFound, map[string]string{
				"error": err.Error(),
			})
		}
		if err.Error() == "you don't have permission to update this payment status" {
			return c.JSON(http.StatusForbidden, map[string]string{
				"error": err.Error(),
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":        "Payment status updated successfully",
		"payment_status": paymentStatus,
	})
}

// GET /shops/orders
func (h *Handler) GetShopOrders(c echo.Context) error {
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

	orders, err := h.usecase.GetShopOrders(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}
	return c.JSON(http.StatusOK, orders)
}

// PATCH /shops/orders/:id/status - ร้านค้าอัพเดทสถานะคำสั่งซื้อ
func (h *Handler) UpdateShopOrderStatus(c echo.Context) error {
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

	orderID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid order ID",
		})
	}

	var req struct {
		Status string `json:"status" validate:"required"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	status := entity.OrderStatus(req.Status)
	if status != entity.OrderStatusPending &&
		status != entity.OrderStatusShipped &&
		status != entity.OrderStatusDelivered &&
		status != entity.OrderStatusCancelled {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid status. Must be: pending, shipped, delivered, or cancelled",
		})
	}

	if err := h.usecase.UpdateShopOrderStatus(uint(orderID), status, userID); err != nil {
		if err.Error() == "order not found" {
			return c.JSON(http.StatusNotFound, map[string]string{
				"error": err.Error(),
			})
		}
		if err.Error() == "you don't have permission to update this order status" {
			return c.JSON(http.StatusForbidden, map[string]string{
				"error": err.Error(),
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Order status updated successfully",
		"status":  status,
	})
}

// PATCH /shops/orders/:id/cancel - ร้านค้ายกเลิกคำสั่งซื้อ
func (h *Handler) CancelShopOrder(c echo.Context) error {
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

	orderID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid order ID",
		})
	}

	if err := h.usecase.CancelShopOrder(uint(orderID), userID); err != nil {
		if err.Error() == "order not found" {
			return c.JSON(http.StatusNotFound, map[string]string{
				"error": err.Error(),
			})
		}
		if err.Error() == "you don't have permission to update this order status" {
			return c.JSON(http.StatusForbidden, map[string]string{
				"error": err.Error(),
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Order cancelled successfully",
		"status":  entity.OrderStatusCancelled,
	})
}
