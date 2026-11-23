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
