package icon

import (
	"github.com/eyo-chen/expense-tracker-go/internal/domain"
	"github.com/eyo-chen/expense-tracker-go/internal/model/interfaces"
)

type IconUC struct {
	Icon interfaces.IconModel
}

func NewIconUC(i interfaces.IconModel) *IconUC {
	return &IconUC{
		Icon: i,
	}
}

func (i *IconUC) List() ([]domain.Icon, error) {
	return i.Icon.List()
}
