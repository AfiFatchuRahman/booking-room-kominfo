package domain

import (
	"net/http"
	"time"
)

type User struct {
	ID           int64     `json:"id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	Role         string    `json:"role"`
	CreatedAt    time.Time `json:"created_at,omitempty"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type Room struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Capacity int    `json:"capacity"`
	IsActive bool   `json:"is_active"`
}

type CreateRoomRequest struct {
	Name     string `json:"name" binding:"required"`
	Capacity int    `json:"capacity" binding:"required,gt=0"`
}

type AvailableRoom struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Capacity int    `json:"capacity"`
}

type CreateBookingRequest struct {
	RoomID    int64  `json:"room_id" binding:"required"`
	Title     string `json:"title" binding:"required"`
	StartTime string `json:"start_time" binding:"required"`
	EndTime   string `json:"end_time" binding:"required"`
}

type Booking struct {
	ID        int64     `json:"id"`
	RoomID    int64     `json:"room_id"`
	UserID    int64     `json:"user_id"`
	Title     string    `json:"title"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	Status    string    `json:"status"`
}

type ErrorResponse struct {
	ErrorCode string `json:"error_code"`
	Message   string `json:"message"`
}

type AppError struct {
	HTTPStatus int
	Code       string
	Message    string
}

func (e *AppError) Error() string {
	return e.Message
}

func NewAppError(status int, code, msg string) *AppError {
	return &AppError{
		HTTPStatus: status,
		Code:       code,
		Message:    msg,
	}
}

var (
	ErrUnauthorized   = NewAppError(http.StatusUnauthorized, "UNAUTHORIZED", "Email atau password salah")
	ErrForbidden      = NewAppError(http.StatusForbidden, "FORBIDDEN", "Anda tidak memiliki akses ke resource ini")
	ErrInvalidTime    = NewAppError(http.StatusBadRequest, "INVALID_TIME", "Waktu pemesanan tidak valid atau sudah lewat")
	ErrRoomNotFound   = NewAppError(http.StatusNotFound, "ROOM_NOT_FOUND", "Ruang meeting tidak ditemukan")
	ErrRoomInactive   = NewAppError(http.StatusBadRequest, "ROOM_INACTIVE", "Ruang meeting sedang tidak aktif")
	ErrRoomConflict   = NewAppError(http.StatusConflict, "ROOM_CONFLICT", "Ruangan sudah dipesan pada rentang waktu tersebut")
	ErrBookingMissing = NewAppError(http.StatusNotFound, "BOOKING_NOT_FOUND", "Data pemesanan tidak ditemukan")
)
