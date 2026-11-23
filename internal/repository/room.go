package repository

import (
	"context"
	"time"

	"github.com/AfiFatchuRahman/booking-room-kominfo/internal/domain"
)

func (r *Repository) CreateRoom(ctx context.Context, name string, capacity int) (*domain.Room, error) {
	query := `
		INSERT INTO rooms (name, capacity, is_active)
		VALUES ($1, $2, TRUE)
		RETURNING id, name, capacity, is_active
	`
	var room domain.Room
	err := r.db.QueryRowContext(ctx, query, name, capacity).Scan(
		&room.ID, &room.Name, &room.Capacity, &room.IsActive,
	)
	if err != nil {
		return nil, err
	}
	return &room, nil
}

func (r *Repository) FindAvailableRooms(ctx context.Context, start, end time.Time) ([]domain.AvailableRoom, error) {
	query := `
		SELECT r.id, r.name, r.capacity
		FROM rooms r
		WHERE r.is_active = TRUE
		  AND NOT EXISTS (
			  SELECT 1
			  FROM bookings b
			  WHERE b.room_id = r.id
			    AND b.status = 'CONFIRMED'
			    AND b.start_time < $2
			    AND b.end_time > $1
		  )
		ORDER BY r.id ASC
	`
	rows, err := r.db.QueryContext(ctx, query, start.UTC(), end.UTC())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	rooms := make([]domain.AvailableRoom, 0)
	for rows.Next() {
		var item domain.AvailableRoom
		if err := rows.Scan(&item.ID, &item.Name, &item.Capacity); err != nil {
			return nil, err
		}
		rooms = append(rooms, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return rooms, nil
}
