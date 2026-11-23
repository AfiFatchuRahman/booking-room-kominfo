package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/AfiFatchuRahman/booking-room-kominfo/internal/domain"
	"github.com/AfiFatchuRahman/booking-room-kominfo/internal/middleware"
	"github.com/AfiFatchuRahman/booking-room-kominfo/internal/usecase"

	"github.com/gin-gonic/gin"
)

type HTTPHandler struct {
	uc *usecase.Usecase
}

func NewHTTPHandler(uc *usecase.Usecase) *HTTPHandler {
	return &HTTPHandler{uc: uc}
}

func (h *HTTPHandler) RegisterRoutes(r *gin.Engine, jwtSecret string) {
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, domain.ErrorResponse{
			ErrorCode: "NOT_FOUND",
			Message:   "Endpoint tidak ditemukan",
		})
	})

	v1 := r.Group("/api/v1")
	{
		v1.POST("/auth/login", h.Login)

		protected := v1.Group("")
		protected.Use(middleware.AuthMiddleware(jwtSecret))
		{
			protected.POST("/rooms", middleware.RequireRole("ADMIN"), h.CreateRoom)
			protected.GET("/rooms/available", h.GetAvailableRooms)
			protected.POST("/bookings", h.CreateBooking)
			protected.PATCH("/bookings/:id/cancel", h.CancelBooking)

		}
	}
}

func (h *HTTPHandler) Login(c *gin.Context) {
	var req domain.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{
			ErrorCode: "INVALID_REQUEST",
			Message:   "Payload email dan password tidak valid",
		})
		return
	}

	res, err := h.uc.Login(c.Request.Context(), req)
	if err != nil {
		h.writeError(c, err)
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *HTTPHandler) CreateRoom(c *gin.Context) {
	var req domain.CreateRoomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{
			ErrorCode: "INVALID_REQUEST",
			Message:   "Payload nama ruangan dan kapasitas tidak valid",
		})
		return
	}

	room, err := h.uc.CreateRoom(c.Request.Context(), req)
	if err != nil {
		h.writeError(c, err)
		return
	}

	c.JSON(http.StatusCreated, room)
}

func (h *HTTPHandler) GetAvailableRooms(c *gin.Context) {
	start := c.Query("start")
	end := c.Query("end")
	if start == "" || end == "" {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{
			ErrorCode: "INVALID_TIME",
			Message:   "Parameter query start dan end wajib diisi",
		})
		return
	}

	rooms, err := h.uc.GetAvailableRooms(c.Request.Context(), start, end)
	if err != nil {
		h.writeError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": rooms,
	})
}

func (h *HTTPHandler) CreateBooking(c *gin.Context) {
	var req domain.CreateBookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{
			ErrorCode: "INVALID_REQUEST",
			Message:   "Payload pemesanan ruangan tidak valid",
		})
		return
	}

	userID := c.GetInt64("user_id")
	booking, err := h.uc.CreateBooking(c.Request.Context(), userID, req)
	if err != nil {
		h.writeError(c, err)
		return
	}

	c.JSON(http.StatusCreated, booking)
}

func (h *HTTPHandler) CancelBooking(c *gin.Context) {
	idParam := c.Param("id")
	bookingID, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil || bookingID <= 0 {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{
			ErrorCode: "INVALID_ID",
			Message:   "ID booking tidak valid",
		})
		return
	}

	userID := c.GetInt64("user_id")
	role := c.GetString("role")

	res, err := h.uc.CancelBooking(c.Request.Context(), bookingID, userID, role)
	if err != nil {
		h.writeError(c, err)
		return
	}

	c.JSON(http.StatusOK, res)
}
func (h *HTTPHandler) writeError(c *gin.Context, err error) {
	var appErr *domain.AppError
	if errors.As(err, &appErr) {
		c.JSON(appErr.HTTPStatus, domain.ErrorResponse{
			ErrorCode: appErr.Code,
			Message:   appErr.Message,
		})
		return
	}

	c.JSON(http.StatusInternalServerError, domain.ErrorResponse{
		ErrorCode: "INTERNAL_SERVER_ERROR",
		Message:   "Terjadi kesalahan pada server",
	})
}
