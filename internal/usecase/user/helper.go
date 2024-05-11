package user

import (
	"os"
	"time"

	"github.com/OYE0303/expense-tracker-go/internal/domain"
	"github.com/OYE0303/expense-tracker-go/pkg/logger"
	"github.com/golang-jwt/jwt/v5"
)

func genJWTToken(user domain.User) (string, error) {
	key := []byte(os.Getenv("JWT_SECRET_KEY"))
	claims := Claims{
		UserID:    user.ID,
		UserName:  user.Name,
		UserEmail: user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	s, err := token.SignedString(key)
	if err != nil {
		logger.Error("token.SignedString failed", "err", err)
		return "", err
	}

	return s, nil
}
