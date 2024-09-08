package user

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"time"

	"github.com/eyo-chen/expense-tracker-go/internal/domain"
	"github.com/eyo-chen/expense-tracker-go/pkg/logger"
	"github.com/golang-jwt/jwt/v5"
)

type claims struct {
	UserID    int64  `json:"user_id"`
	UserName  string `json:"user_name"`
	UserEmail string `json:"user_email"`
	jwt.RegisteredClaims
}

func hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

func genJWTToken(user domain.User) (string, error) {
	key := []byte(os.Getenv("JWT_SECRET_KEY"))
	claims := claims{
		UserID:    user.ID,
		UserName:  user.Name,
		UserEmail: user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
		},
	}

	return genToken(key, claims)
}

func genRefreshToken() (string, error) {
	key := []byte(os.Getenv("JWT_SECRET_KEY"))
	claims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(30 * 24 * time.Hour)),
	}

	return genToken(key, claims)
}

func genToken(key []byte, claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	s, err := token.SignedString(key)
	if err != nil {
		logger.Error("token.SignedString failed", "err", err)
		return "", err
	}

	return s, nil

}
