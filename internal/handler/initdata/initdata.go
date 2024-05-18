package initdata

import (
	"net/http"

	"github.com/OYE0303/expense-tracker-go/internal/usecase/interfaces"
	"github.com/OYE0303/expense-tracker-go/pkg/errutil"
	"github.com/OYE0303/expense-tracker-go/pkg/jsonutil"
	"github.com/OYE0303/expense-tracker-go/pkg/logger"
)

const (
	packageName = "handler/initdata"
)

type InitDataHlr struct {
	InitData interfaces.InitDataUC
}

func New(i interfaces.InitDataUC) *InitDataHlr {
	return &InitDataHlr{
		InitData: i,
	}
}

func (i *InitDataHlr) List(w http.ResponseWriter, r *http.Request) {
	initData, err := i.InitData.List()
	if err != nil {
		errutil.ServerErrorResponse(w, r, err)
		return
	}

	respData := map[string]interface{}{
		"init_data": initData,
	}
	if err := jsonutil.WriteJSON(w, http.StatusOK, respData, nil); err != nil {
		logger.Error("jsonutil.WriteJSON failed", "package", packageName, "err", err)
		errutil.ServerErrorResponse(w, r, err)
		return
	}
}
