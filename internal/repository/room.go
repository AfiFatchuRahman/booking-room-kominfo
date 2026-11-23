package repository

import (
	"context"
	"database/sql"
	"errors"
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

func (r *Repository) CreateBookingWithLock(ctx context.Context, roomID, userID int64, title string, start, end time.Time) (*domain.Booking, error) {
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	// Kunci baris room terkait untuk serialisasi request concurrent pada room yang sama
	var isActive bool
	err = tx.QueryRowContext(ctx, `SELECT is_active FROM rooms WHERE id = $1 FOR UPDATE`, roomID).Scan(&isActive)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrRoomNotFound
	}
	if err != nil {
		return nil, err
	}

	if !isActive {
		return nil, domain.ErrRoomInactive
	}

	// Cek irisan jadwal CONFIRMED
	var conflictCount int
	conflictQuery := `
		SELECT COUNT(1)
		FROM bookings
		WHERE room_id = $1
		  AND status = 'CONFIRMED'
		  AND start_time < $3
		  AND end_time > $2
	`
	err = tx.QueryRowContext(ctx, conflictQuery, roomID, start.UTC(), end.UTC()).Scan(&conflictCount)
	if err != nil {
		return nil, err
	}

	if conflictCount > 0 {
		return nil, domain.ErrRoomConflict
	}

	insertQuery := `
		INSERT INTO bookings (room_id, user_id, title, start_time, end_time, status)
		VALUES ($1, $2, $3, $4, $5, 'CONFIRMED')
		RETURNING id, room_id, user_id, title, start_time, end_time, status
	`
	var booking domain.Booking
	err = tx.QueryRowContext(ctx, insertQuery, roomID, userID, title, start.UTC(), end.UTC()).Scan(
		&booking.ID,
		&booking.RoomID,
		&booking.UserID,
		&booking.Title,
		&booking.StartTime,
		&booking.EndTime,
		&booking.Status,
	)
	if err != nil {
		return nil, err
	}

	booking.StartTime = booking.StartTime.UTC()
	booking.EndTime = booking.EndTime.UTC()

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &booking, nil
}

func (r *Repository) CancelBooking(ctx context.Context, bookingID, requesterID int64, requesterRole string) (*domain.CancelBookingResponse, error) {
	var ownerID int64
	err := r.db.QueryRowContext(ctx, `SELECT user_id FROM bookings WHERE id = $1`, bookingID).Scan(&ownerID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrBookingMissing
	}
	if err != nil {
		return nil, err
	}

	if requesterRole != "ADMIN" && ownerID != requesterID {
		return nil, domain.ErrForbidden
	}

	var res domain.CancelBookingResponse
	updateQuery := `
		UPDATE bookings
		SET status = 'CANCELLED', updated_at = NOW()
		WHERE id = $1
		RETURNING id, status
	`
	err = r.db.QueryRowContext(ctx, updateQuery, bookingID).Scan(&res.ID, &res.Status)
	if err != nil {
		return nil, err
	}

	return &res, nil
}
