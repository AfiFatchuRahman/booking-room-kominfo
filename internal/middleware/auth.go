package middleware

import (
	"net/http"
	"strings"

	"github.com/AfiFatchuRahman/booking-room-kominfo/internal/domain"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	UserID int64  `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func AuthMiddleware(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, domain.ErrorResponse{
				ErrorCode: "UNAUTHORIZED",
				Message:   "Token otorisasi tidak ditemukan atau format salah",
			})
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		claims := &JWTClaims{}

		token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, domain.ErrorResponse{
				ErrorCode: "UNAUTHORIZED",
				Message:   "Token tidak valid atau sudah kedaluwarsa",
			})
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("email", claims.Email)
		c.Set("role", claims.Role)
		c.Next()
	}
}

func RequireRole(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleVal, exists := c.Get("role")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, domain.ErrorResponse{
				ErrorCode: "UNAUTHORIZED",
				Message:   "Informasi role pengguna tidak ditemukan",
			})
			return
		}

		userRole, _ := roleVal.(string)
		for _, allowed := range allowedRoles {
			if userRole == allowed {
				c.Next()
				return
			}
		}

		c.AbortWithStatusJSON(http.StatusForbidden, domain.ErrorResponse{
			ErrorCode: "FORBIDDEN",
			Message:   "Role Anda tidak diizinkan mengakses endpoint ini",
		})
	}
}
