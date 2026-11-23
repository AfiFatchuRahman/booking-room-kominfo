package repository

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/AfiFatchuRahman/booking-room-kominfo/internal/domain"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) FindUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `SELECT id, name, email, password_hash, role FROM users WHERE email = $1`
	var u domain.User
	err := r.db.QueryRowContext(ctx, query, email).Scan(&u.ID, &u.Name, &u.Email, &u.PasswordHash, &u.Role)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrUnauthorized
	}
	if err != nil {
		log.Print(err)
		return nil, err
	}
	return &u, nil
}
