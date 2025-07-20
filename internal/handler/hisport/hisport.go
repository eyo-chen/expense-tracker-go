package hisport

import (
	"fmt"
	"net/http"
	"time"

	"github.com/eyo-chen/expense-tracker-go/internal/handler/interfaces"
	"github.com/eyo-chen/expense-tracker-go/pkg/ctxutil"
	"github.com/eyo-chen/expense-tracker-go/pkg/errutil"
	"github.com/eyo-chen/expense-tracker-go/pkg/jsonutil"
	"github.com/eyo-chen/expense-tracker-go/pkg/logger"
)

const (
	packageName = "handler/hisport"
)

type Hlr struct {
	historicalPortfolioUC interfaces.HistoricalPortfolioUC
}

func New(historicalPortfolioUC interfaces.HistoricalPortfolioUC) *Hlr {
	return &Hlr{historicalPortfolioUC: historicalPortfolioUC}
}

type createHistoricalPortfolioReq struct {
	Date string `json:"date"`
}

func (h *Hlr) Create(w http.ResponseWriter, r *http.Request) {
	var input createHistoricalPortfolioReq
	if err := jsonutil.ReadJson(w, r, &input); err != nil {
		logger.Error("jsonutil.ReadJSON failed", "package", packageName, "err", err)
		errutil.BadRequestResponse(w, r, err)
		return
	}

	date, err := time.Parse("2006-01-02", input.Date)
	if err != nil {
		logger.Error("Invalid date format", "package", packageName, "err", err)
		errutil.BadRequestResponse(w, r, fmt.Errorf("invalid date format, expected YYYY-MM-DD"))
		return
	}

	user := ctxutil.GetUser(r)
	ctx := r.Context()
	if err := h.historicalPortfolioUC.Create(ctx, int32(user.ID), date); err != nil {
		logger.Error("Create failed", "package", packageName, "err", err)
		errutil.ServerErrorResponse(w, r, err)
		return
	}

	if err := jsonutil.WriteJSON(w, http.StatusCreated, map[string]interface{}{"message": "Historical portfolio created successfully"}, nil); err != nil {
		logger.Error("jsonutil.WriteJSON failed", "package", packageName, "err", err)
		errutil.ServerErrorResponse(w, r, err)
		return
	}
}

func (h *Hlr) GetPortfolioValue(w http.ResponseWriter, r *http.Request) {
	dateOption := r.URL.Query().Get("date_option")
	if dateOption == "" {
		logger.Error("date_option query parameter is required", "package", packageName)
		errutil.BadRequestResponse(w, r, fmt.Errorf("date_option query parameter is required"))
		return
	}

	user := ctxutil.GetUser(r)
	ctx := r.Context()
	dates, values, err := h.historicalPortfolioUC.GetPortfolioValue(ctx, int32(user.ID), dateOption)
	if err != nil {
		logger.Error("GetPortfolioValue failed", "package", packageName, "err", err)
		errutil.ServerErrorResponse(w, r, err)
		return
	}

	respData := map[string]interface{}{
		"dates":  dates,
		"values": values,
	}

	if err := jsonutil.WriteJSON(w, http.StatusOK, respData, nil); err != nil {
		logger.Error("jsonutil.WriteJSON failed", "package", packageName, "err", err)
		errutil.ServerErrorResponse(w, r, err)
		return
	}
}