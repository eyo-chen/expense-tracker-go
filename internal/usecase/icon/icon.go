package icon

import (
	"github.com/OYE0303/expense-tracker-go/internal/domain"
	"github.com/OYE0303/expense-tracker-go/internal/model/interfaces"
)

type IconUC struct {
	Icon interfaces.IconModel
}

func NewIconUC(i interfaces.IconModel) *IconUC {
	a := 0
	return &IconUC{
		Icon: i,
	}
}

func (i *IconUC) List() ([]domain.Icon, error) {
	return i.Icon.List()
}
