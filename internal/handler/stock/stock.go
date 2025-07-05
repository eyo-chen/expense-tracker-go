package stock

import (
	"fmt"
	"net/http"
	"time"

	"github.com/eyo-chen/expense-tracker-go/internal/domain"
	"github.com/eyo-chen/expense-tracker-go/internal/handler/interfaces"
	"github.com/eyo-chen/expense-tracker-go/pkg/ctxutil"
	"github.com/eyo-chen/expense-tracker-go/pkg/errutil"
	"github.com/eyo-chen/expense-tracker-go/pkg/jsonutil"
	"github.com/eyo-chen/expense-tracker-go/pkg/logger"
)

const (
	packageName = "handler/stock"
)

type Hlr struct {
	stockUC interfaces.StockUC
}

func New(stockUC interfaces.StockUC) *Hlr {
	return &Hlr{stockUC: stockUC}
}

type createStockReq struct {
	Symbol     string    `json:"symbol"`
	Price      float64   `json:"price"`
	Quantity   int32     `json:"quantity"`
	ActionType string    `json:"action_type"`
	StockType  string    `json:"stock_type"`
	Date       time.Time `json:"date"`
}

func (h *Hlr) Create(w http.ResponseWriter, r *http.Request) {
	var input createStockReq
	if err := jsonutil.ReadJson(w, r, &input); err != nil {
		logger.Error("jsonutil.ReadJSON failed", "package", packageName, "err", err)
		errutil.BadRequestResponse(w, r, err)
		return
	}

	fmt.Println("input", input)
	fmt.Println("input.Date", input.Date)

	user := ctxutil.GetUser(r)
	stock := domain.CreateStock{
		UserID:     int32(user.ID),
		Symbol:     input.Symbol,
		Quantity:   input.Quantity,
		Price:      input.Price,
		ActionType: input.ActionType,
		StockType:  input.StockType,
		Date:       input.Date,
	}

	ctx := r.Context()
	id, err := h.stockUC.Create(ctx, stock)
	if err != nil {
		errutil.ServerErrorResponse(w, r, err)
		return
	}

	if err := jsonutil.WriteJSON(w, http.StatusCreated, map[string]interface{}{"id": id}, nil); err != nil {
		logger.Error("jsonutil.WriteJSON failed", "package", packageName, "err", err)
		errutil.ServerErrorResponse(w, r, err)
		return
	}
}

func (h *Hlr) GetPortfolioInfo(w http.ResponseWriter, r *http.Request) {
	user := ctxutil.GetUser(r)
	ctx := r.Context()
	portfolio, err := h.stockUC.GetPortfolioInfo(ctx, int32(user.ID))
	if err != nil {
		errutil.ServerErrorResponse(w, r, err)
		return
	}

	respData := map[string]interface{}{
		"total_portfolio_value": portfolio.TotalPortfolioValue,
		"total_gain":            portfolio.TotalGain,
		"roi":                   portfolio.ROI,
	}

	if err := jsonutil.WriteJSON(w, http.StatusOK, respData, nil); err != nil {
		logger.Error("jsonutil.WriteJSON failed", "package", packageName, "err", err)
		errutil.ServerErrorResponse(w, r, err)
		return
	}
}

func (h *Hlr) GetStockInfo(w http.ResponseWriter, r *http.Request) {
	user := ctxutil.GetUser(r)
	ctx := r.Context()
	stockInfo, err := h.stockUC.GetStockInfo(ctx, int32(user.ID))
	if err != nil {
		errutil.ServerErrorResponse(w, r, err)
		return
	}

	respData := map[string]interface{}{
		"stocks": stockInfo.Stocks,
		"etf":    stockInfo.ETF,
		"cash":   stockInfo.Cash,
	}

	if err := jsonutil.WriteJSON(w, http.StatusOK, respData, nil); err != nil {
		logger.Error("jsonutil.WriteJSON failed", "package", packageName, "err", err)
		errutil.ServerErrorResponse(w, r, err)
		return
	}
}
