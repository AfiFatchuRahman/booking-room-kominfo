package handler

import (
	"errors"
	"net/http"

	"github.com/AfiFatchuRahman/booking-room-kominfo/internal/domain"
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
