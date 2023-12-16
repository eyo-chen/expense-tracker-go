package router

import (
	"net/http"

	hd "github.com/OYE0303/expense-tracker-go/internal/handler"
	"github.com/OYE0303/expense-tracker-go/internal/middleware"
	"github.com/gorilla/mux"
)

// New initializes a new router and returns it
func New(handler *hd.Handler) http.Handler {
	r := mux.NewRouter()

	// user
	r.HandleFunc("/v1/user/signup", handler.User.Signup).Methods(http.MethodPost)
	r.HandleFunc("/v1/user/login", handler.User.Login).Methods(http.MethodPost)

	// main category
	r.HandleFunc("/v1/main-category", handler.MainCateg.CreateMainCateg).Methods(http.MethodPost)
	r.HandleFunc("/v1/main-category/{id}", handler.MainCateg.UpdateMainCateg).Methods(http.MethodPatch)
	r.HandleFunc("/v1/main-category/{id}", handler.MainCateg.DeleteMainCateg).Methods(http.MethodDelete)

	// sub category
	r.HandleFunc("/v1/sub-category", handler.SubCateg.CreateSubCateg).Methods(http.MethodPost)
	r.HandleFunc("/v1/sub-category/{id}", handler.SubCateg.UpdateSubCateg).Methods(http.MethodPatch)

	return middleware.LogRequest(middleware.Authenticate(r))
}
