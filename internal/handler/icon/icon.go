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

type IconHandler struct {
	Icon interfaces.IconUC
}

func NewIconHandler(i interfaces.IconUC) *IconHandler {
	return &IconHandler{
		Icon: i,
	}
}

func (i *IconHandler) List(w http.ResponseWriter, r *http.Request) {
	icons, err := i.Icon.List()
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
