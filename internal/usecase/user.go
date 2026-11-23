package usecase

import (
	"context"
	"strings"
	"time"

	"github.com/AfiFatchuRahman/booking-room-kominfo/internal/domain"
	"github.com/AfiFatchuRahman/booking-room-kominfo/internal/middleware"
	"github.com/AfiFatchuRahman/booking-room-kominfo/internal/repository"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type Usecase struct {
	repo      *repository.Repository
	jwtSecret string
}

func NewUsecase(repo *repository.Repository, jwtSecret string) *Usecase {
	return &Usecase{
		repo:      repo,
		jwtSecret: jwtSecret,
	}
}

func (u *Usecase) Login(ctx context.Context, req domain.LoginRequest) (*domain.LoginResponse, error) {
	email := strings.TrimSpace(strings.ToLower(req.Email))
	user, err := u.repo.FindUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, domain.ErrUnauthorized
	}

	now := time.Now().UTC()
	claims := middleware.JWTClaims{
		UserID: user.ID,
		Email:  user.Email,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(24 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(u.jwtSecret))
	if err != nil {
		return nil, err
	}

	return &domain.LoginResponse{Token: signed}, nil
}
