package icon

import (
	"net/http"

	"github.com/eyo-chen/expense-tracker-go/internal/handler/interfaces"
	"github.com/eyo-chen/expense-tracker-go/pkg/errutil"
	"github.com/eyo-chen/expense-tracker-go/pkg/jsonutil"
	"github.com/eyo-chen/expense-tracker-go/pkg/logger"
)

const (
	PackageName = "handler/icon"
)

type Hlr struct {
	Icon interfaces.IconUC
}

func New(i interfaces.IconUC) *Hlr {
	return &Hlr{
		Icon: i,
	}
}

func (h *Hlr) List(w http.ResponseWriter, r *http.Request) {
	icons, err := h.Icon.List()
	if err != nil {
		errutil.ServerErrorResponse(w, r, err)
		return
	}

	respData := map[string]interface{}{
		"icons": icons,
	}
	if err := jsonutil.WriteJSON(w, http.StatusOK, respData, nil); err != nil {
		logger.Error("jsonutil.WriteJSON failed", "package", PackageName, "err", err)
		errutil.ServerErrorResponse(w, r, err)
		return
	}
}
