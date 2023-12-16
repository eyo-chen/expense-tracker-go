package usecase

import (
	"github.com/OYE0303/expense-tracker-go/internal/domain"
	"github.com/OYE0303/expense-tracker-go/pkg/logger"
)

type mainCategUC struct {
	MainCategModel MainCategModel
	IconModel      IconModel
}

func newMainCategUC(m MainCategModel, i IconModel) *mainCategUC {
	return &mainCategUC{
		MainCategModel: m,
		IconModel:      i,
	}
}

func (m *mainCategUC) Add(categ *domain.MainCateg, userID int64, iconID int64) error {
	categbyUserID, err := m.MainCategModel.GetOneByUserID(userID, categ.Name)
	if err != nil && err != domain.ErrDataNotFound {
		logger.Error("m.MainCategModel.GetOneByUserID failed", "package", "usecase", "err", err)
		return err
	}

	if categbyUserID != nil {
		return domain.ErrDataAlreadyExists
	}

	icon, err := m.IconModel.GetByID(iconID)
	if err != nil && err != domain.ErrDataNotFound {
		logger.Error("m.IconModel.GetByID failed", "package", "usecase", "err", err)
		return err
	}

	if icon == nil {
		return domain.ErrDataNotFound
	}

	if err := m.MainCategModel.Create(categ, userID, iconID); err != nil {
		logger.Error("m.MainCategModel.Create failed", "package", "usecase", "err", err)
		return err
	}

	return nil
}
