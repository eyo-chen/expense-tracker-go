package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/OYE0303/expense-tracker-go/internal/domain"
	"github.com/OYE0303/expense-tracker-go/pkg/ctxutil"
	"github.com/OYE0303/expense-tracker-go/pkg/errutil"
	"github.com/OYE0303/expense-tracker-go/pkg/logger"
	"github.com/golang-jwt/jwt/v5"
)

func LogRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Info("request started", "method", r.Method, "url", r.URL)
		next.ServeHTTP(w, r)
		logger.Info("request completed", "method", r.Method, "url", r.URL)
	})
}

func Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorizationHeader := r.Header.Get("Authorization")
		if authorizationHeader == "" {
			errutil.AuthenticationErrorResponse(w, r, domain.ErrAuthentication)
			return
		}

		auth := strings.Split(authorizationHeader, " ")
		if len(auth) != 2 || auth[0] != "Bearer" {
			errutil.AuthenticationErrorResponse(w, r, domain.ErrAuthToken)
			return
		}

		token, err := jwt.Parse(auth[1], func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				logger.Error("Unexpected signing method", "package", "middleware")
				errutil.ServerErrorResponse(w, r, domain.ErrServer)
			}

			return []byte(os.Getenv("JWT_SECRET_KEY")), nil
		})
		if err != nil {
			logger.Error("jwt.Parse failed", "package", "middleware", "err", err)
			errutil.AuthenticationErrorResponse(w, r, err)
			return
		}

		user := domain.User{
			ID:    int64(token.Claims.(jwt.MapClaims)["user_id"].(float64)),
			Email: token.Claims.(jwt.MapClaims)["user_email"].(string),
			Name:  token.Claims.(jwt.MapClaims)["user_name"].(string),
		}

		next.ServeHTTP(w, ctxutil.SetUser(r, &user))
	})
}

func EnableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Authorization, Content-Type")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
