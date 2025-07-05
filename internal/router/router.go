package router

import (
	"net/http"

	hd "github.com/eyo-chen/expense-tracker-go/internal/handler"
	"github.com/eyo-chen/expense-tracker-go/internal/middleware"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
)

// New initializes a new router and returns it
func New(handler *hd.Handler) http.Handler {
	r := mux.NewRouter()

	// user
	r.HandleFunc("/v1/user/signup", handler.User.Signup).Methods(http.MethodPost)
	r.HandleFunc("/v1/user/login", handler.User.Login).Methods(http.MethodPost)
	r.HandleFunc("/v1/user/token", handler.User.Token).Methods(http.MethodGet)

	// init data
	r.Handle("/v1/init-data", http.HandlerFunc(handler.InitData.List)).Methods(http.MethodGet)

	// icon
	r.Handle("/v1/icon", http.HandlerFunc(handler.Icon.List)).Methods(http.MethodGet)

	auth := alice.New(middleware.Authenticate)

	// user with auth
	r.Handle("/v1/user", auth.ThenFunc(handler.User.GetInfo)).Methods(http.MethodGet)

	// user icon
	r.Handle("/v1/user-icon", auth.ThenFunc(handler.Icon.ListByUserID)).Methods(http.MethodGet)
	r.Handle("/v1/user-icon/url", auth.ThenFunc(handler.UserIcon.GetPutObjectURL)).Methods(http.MethodPost)
	r.Handle("/v1/user-icon", auth.ThenFunc(handler.UserIcon.Create)).Methods(http.MethodPost)

	// init data with auth
	r.Handle("/v1/init-data", auth.ThenFunc(handler.InitData.Create)).Methods(http.MethodPost)

	// main category
	r.Handle("/v1/main-category", auth.ThenFunc(handler.MainCateg.Create)).Methods(http.MethodPost)
	r.Handle("/v1/main-category", auth.ThenFunc(handler.MainCateg.GetAll)).Methods(http.MethodGet)
	r.Handle("/v1/main-category/{id}", auth.ThenFunc(handler.MainCateg.Update)).Methods(http.MethodPatch)
	r.Handle("/v1/main-category/{id}", auth.ThenFunc(handler.MainCateg.Delete)).Methods(http.MethodDelete)

	// sub category
	r.Handle("/v1/sub-category", auth.ThenFunc(handler.SubCateg.CreateSubCateg)).Methods(http.MethodPost)
	r.Handle("/v1/main-category/{id}/sub-category", auth.ThenFunc(handler.SubCateg.GetByMainCategID)).Methods(http.MethodGet)
	r.Handle("/v1/sub-category/{id}", auth.ThenFunc(handler.SubCateg.UpdateSubCateg)).Methods(http.MethodPatch)
	r.Handle("/v1/sub-category/{id}", auth.ThenFunc(handler.SubCateg.DeleteSubCateg)).Methods(http.MethodDelete)

	// transaction
	r.Handle("/v1/transaction", auth.ThenFunc(handler.Transaction.Create)).Methods(http.MethodPost)
	r.Handle("/v1/transaction", auth.ThenFunc(handler.Transaction.GetAll)).Methods(http.MethodGet)
	r.Handle("/v1/transaction/{id}", auth.ThenFunc(handler.Transaction.Update)).Methods(http.MethodPut)
	r.Handle("/v1/transaction/{id}", auth.ThenFunc(handler.Transaction.Delete)).Methods(http.MethodDelete)
	r.Handle("/v1/transaction/info", auth.ThenFunc(handler.Transaction.GetAccInfo)).Methods(http.MethodGet)
	r.Handle("/v1/transaction/bar-chart", auth.ThenFunc(handler.Transaction.GetBarChartData)).Methods(http.MethodGet)
	r.Handle("/v1/transaction/pie-chart", auth.ThenFunc(handler.Transaction.GetPieChartData)).Methods(http.MethodGet)
	r.Handle("/v1/transaction/line-chart", auth.ThenFunc(handler.Transaction.GetLineChartData)).Methods(http.MethodGet)
	r.Handle("/v1/transaction/monthly-data", auth.ThenFunc(handler.Transaction.GetMonthlyData)).Methods(http.MethodGet)

	// stock
	r.Handle("/v1/stock", auth.ThenFunc(handler.Stock.Create)).Methods(http.MethodPost)
	r.Handle("/v1/stock/portfolio", auth.ThenFunc(handler.Stock.GetPortfolioInfo)).Methods(http.MethodGet)
	r.Handle("/v1/stock/info", auth.ThenFunc(handler.Stock.GetStockInfo)).Methods(http.MethodGet)

	regular := alice.New(middleware.LogRequest, middleware.EnableCORS)

	return regular.Then(r)
}
