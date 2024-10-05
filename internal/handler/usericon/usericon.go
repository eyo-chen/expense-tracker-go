package usericon

import (
	"net/http"

	"github.com/eyo-chen/expense-tracker-go/internal/handler/interfaces"
	"github.com/eyo-chen/expense-tracker-go/pkg/ctxutil"
	"github.com/eyo-chen/expense-tracker-go/pkg/errutil"
	"github.com/eyo-chen/expense-tracker-go/pkg/jsonutil"
	"github.com/eyo-chen/expense-tracker-go/pkg/logger"
	"github.com/eyo-chen/expense-tracker-go/pkg/validator"
)

type Hlr struct {
	UserIcon interfaces.UserIconUC
}

func New(userIcon interfaces.UserIconUC) *Hlr {
	return &Hlr{UserIcon: userIcon}
}

type getPutObjectURLInput struct {
	FileName string `json:"file_name"`
}

func (h *Hlr) GetPutObjectURL(w http.ResponseWriter, r *http.Request) {
	var input getPutObjectURLInput

	if err := jsonutil.ReadJson(w, r, &input); err != nil {
		logger.Error("jsonutil.ReadJson failed", "package", "handler", "err", err)
		errutil.BadRequestResponse(w, r, err)
		return
	}

	v := validator.New()
	if !v.GetPutObjectURL(input.FileName) {
		errutil.VildateErrorResponse(w, r, v.Error)
		return
	}

	user := ctxutil.GetUser(r)
	url, err := h.UserIcon.GetPutObjectURL(r.Context(), input.FileName, user.ID)
	if err != nil {
		errutil.ServerErrorResponse(w, r, err)
		return
	}

	resp := map[string]interface{}{
		"url": url,
	}
	if err := jsonutil.WriteJSON(w, http.StatusOK, resp, nil); err != nil {
		logger.Error("jsonutil.WriteJSON failed", "package", "handler", "err", err)
		errutil.ServerErrorResponse(w, r, err)
		return
	}
}
