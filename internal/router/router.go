package router

import (
	"net/http"

	hd "github.com/OYE0303/expense-tracker-go/internal/handler"
	"github.com/gorilla/mux"
)

// New initializes a new router and returns it
func New(handler *hd.Handler) http.Handler {
	r := mux.NewRouter()

	r.HandleFunc("/v1/user/signup", handler.User.Signup).Methods(http.MethodPost)

	return r
}
