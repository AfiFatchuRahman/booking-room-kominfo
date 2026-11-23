package usecase

import (
	"context"
	"strings"
	"time"

	"github.com/AfiFatchuRahman/booking-room-kominfo/internal/domain"
)

func (u *Usecase) CreateRoom(ctx context.Context, req domain.CreateRoomRequest) (*domain.Room, error) {
	name := strings.TrimSpace(req.Name)
	if name == "" || req.Capacity <= 0 {
		return nil, domain.NewAppError(400, "INVALID_INPUT", "Nama ruangan dan kapasitas harus valid")
	}
	return u.repo.CreateRoom(ctx, name, req.Capacity)
}

func (u *Usecase) GetAvailableRooms(ctx context.Context, startStr, endStr string) ([]domain.AvailableRoom, error) {
	start, err := time.Parse(time.RFC3339, startStr)
	if err != nil {
		return nil, domain.ErrInvalidTime
	}
	end, err := time.Parse(time.RFC3339, endStr)
	if err != nil {
		return nil, domain.ErrInvalidTime
	}

	if !end.After(start) {
		return nil, domain.ErrInvalidTime
	}

	return u.repo.FindAvailableRooms(ctx, start.UTC(), end.UTC())
}

func (u *Usecase) CreateBooking(ctx context.Context, userID int64, req domain.CreateBookingRequest) (*domain.Booking, error) {
	start, err := time.Parse(time.RFC3339, req.StartTime)
	if err != nil {
		return nil, domain.ErrInvalidTime
	}
	end, err := time.Parse(time.RFC3339, req.EndTime)
	if err != nil {
		return nil, domain.ErrInvalidTime
	}

	now := time.Now().UTC()
	if !end.After(start) || start.UTC().Before(now) {
		return nil, domain.ErrInvalidTime
	}

	title := strings.TrimSpace(req.Title)
	if title == "" {
		return nil, domain.NewAppError(400, "INVALID_INPUT", "Judul rapat wajib diisi")
	}

	return u.repo.CreateBookingWithLock(ctx, req.RoomID, userID, title, start.UTC(), end.UTC())
}
