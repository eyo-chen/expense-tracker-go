package transaction

import (
	"errors"
	"net/http"
	"time"

	"github.com/OYE0303/expense-tracker-go/internal/domain"
	"github.com/OYE0303/expense-tracker-go/internal/handler/interfaces"
	"github.com/OYE0303/expense-tracker-go/pkg/ctxutil"
	"github.com/OYE0303/expense-tracker-go/pkg/errutil"
	"github.com/OYE0303/expense-tracker-go/pkg/jsonutil"
	"github.com/OYE0303/expense-tracker-go/pkg/logger"
	"github.com/OYE0303/expense-tracker-go/pkg/validator"
)

type TransactionHandler struct {
	transaction interfaces.TransactionUC
}

func NewTransactionHandler(t interfaces.TransactionUC) *TransactionHandler {
	return &TransactionHandler{
		transaction: t,
	}
}

func (t *TransactionHandler) Create(w http.ResponseWriter, r *http.Request) {
	var input struct {
		MainCategID int64      `json:"main_category_id"`
		SubCategID  int64      `json:"sub_category_id"`
		Price       float64    `json:"price"`
		Date        *time.Time `json:"date"`
		Note        string     `json:"note"`
	}
	if err := jsonutil.ReadJson(w, r, &input); err != nil {
		logger.Error("jsonutil.ReadJSON failed", "package", "handler", "err", err)
		errutil.BadRequestResponse(w, r, err)
		return
	}

	user := ctxutil.GetUser(r)
	transaction := domain.Transaction{
		UserID: user.ID,
		MainCateg: &domain.MainCateg{
			ID: input.MainCategID,
		},
		SubCateg: &domain.SubCateg{
			ID: input.SubCategID,
		},
		Price: input.Price,
		Date:  input.Date,
		Note:  input.Note,
	}

	v := validator.New()
	if !v.CreateTransaction(&transaction) {
		errutil.VildateErrorResponse(w, r, v.Error)
		return
	}

	ctx := r.Context()
	if err := t.transaction.Create(ctx, user, &transaction); err != nil {
		if errors.Is(err, domain.ErrDataNotFound) {
			errutil.BadRequestResponse(w, r, err)
			return
		}

		logger.Error("t.transaction.Create failed", "package", "handler", "err", err)
		errutil.ServerErrorResponse(w, r, err)
		return
	}

	if err := jsonutil.WriteJSON(w, http.StatusCreated, nil, nil); err != nil {
		logger.Error("jsonutil.WriteJSON failed", "package", "handler", "err", err)
		errutil.ServerErrorResponse(w, r, err)
		return
	}
}

func (t *TransactionHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	startDate := r.URL.Query().Get("start_date")
	endDate := r.URL.Query().Get("end_date")

	query := &domain.GetQuery{
		StartDate: startDate,
		EndDate:   endDate,
	}

	v := validator.New()
	if !v.GetTransaction(query) {
		errutil.VildateErrorResponse(w, r, v.Error)
		return
	}

	user := ctxutil.GetUser(r)
	ctx := r.Context()
	transactions, err := t.transaction.GetAll(ctx, query, user)
	if err != nil {
		logger.Error("t.transaction.GetAll failed", "package", "handler", "err", err)
		errutil.ServerErrorResponse(w, r, err)
		return
	}

	transResp := cvtToGetTransactionResp(transactions)
	respData := map[string]interface{}{
		"transactions": transResp.Transactions,
	}

	if err := jsonutil.WriteJSON(w, http.StatusOK, respData, nil); err != nil {
		logger.Error("jsonutil.WriteJSON failed", "package", "handler", "err", err)
		errutil.ServerErrorResponse(w, r, err)
		return
	}
}
