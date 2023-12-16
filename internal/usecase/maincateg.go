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

func (m *mainCategUC) Add(categ *domain.MainCateg, userID int64) error {
	// check if the main category name is already taken
	categbyUserID, err := m.MainCategModel.GetOneByUserID(userID, categ.Name)
	if err != nil && err != domain.ErrDataNotFound {
		logger.Error("m.MainCategModel.GetOneByUserID failed", "package", "usecase", "err", err)
		return err
	}
	if categbyUserID != nil {
		return domain.ErrDataAlreadyExists
	}

	// check if the icon exists
	icon, err := m.IconModel.GetByID(categ.IconID)
	if err != nil && err != domain.ErrDataNotFound {
		logger.Error("m.IconModel.GetByID failed", "package", "usecase", "err", err)
		return err
	}
	if icon == nil {
		return domain.ErrDataNotFound
	}

	if err := m.MainCategModel.Create(categ, userID); err != nil {
		logger.Error("m.MainCategModel.Create failed", "package", "usecase", "err", err)
		return err
	}

	return nil
}

func (m *mainCategUC) Update(categ *domain.MainCateg, userID int64) error {
	categbyID, err := m.MainCategModel.GetByID(categ.ID)
	if err != nil && err != domain.ErrDataNotFound {
		logger.Error("m.MainCategModel.GetByID failed", "package", "usecase", "err", err)
		return err
	}
	if categbyID == nil {
		return domain.ErrDataNotFound
	}

	// check if the main category name is already taken
	categbyUserID, err := m.MainCategModel.GetOneByUserID(userID, categ.Name)
	if err != nil && err != domain.ErrDataNotFound {
		logger.Error("m.MainCategModel.GetOneByUserID failed", "package", "usecase", "err", err)
		return err
	}
	if categbyUserID != nil {
		return domain.ErrDataAlreadyExists
	}

	// check if the icon exists
	icon, err := m.IconModel.GetByID(categ.IconID)
	if err != nil && err != domain.ErrDataNotFound {
		logger.Error("m.IconModel.GetByID failed", "package", "usecase", "err", err)
		return err
	}
	if icon == nil {
		return domain.ErrDataNotFound
	}

	if err := m.MainCategModel.Update(categ); err != nil {
		logger.Error("m.MainCategModel.Update failed", "package", "usecase", "err", err)
		return err
	}

	return nil
}

func (m *mainCategUC) Delete(id int64) error {
	if err := m.MainCategModel.Delete(id); err != nil {
		logger.Error("m.MainCategModel.Delete failed", "package", "usecase", "err", err)
		return err
	}

	return nil
}
