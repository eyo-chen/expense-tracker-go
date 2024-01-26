package usecase

import (
	"github.com/OYE0303/expense-tracker-go/internal/domain"
	"github.com/OYE0303/expense-tracker-go/pkg/logger"
)

type mainCategUC struct {
	MainCateg MainCategModel
	Icon      IconModel
}

func newMainCategUC(m MainCategModel, i IconModel) *mainCategUC {
	return &mainCategUC{
		MainCateg: m,
		Icon:      i,
	}
}

func (m *mainCategUC) Create(categ *domain.MainCateg, userID int64) error {
	// check if the icon exists
	_, err := m.Icon.GetByID(categ.Icon.ID)
	if err != nil {
		return err
	}

	if err := m.MainCateg.Create(categ, userID); err != nil {
		return err
	}

	return nil
}

func (m *mainCategUC) GetAll(userID int64) ([]*domain.MainCateg, error) {
	categs, err := m.MainCateg.GetAll(userID)
	if err != nil {
		logger.Error("m.MainCateg.GetAll failed", "package", "usecase", "err", err)
		return nil, err
	}

	return categs, nil
}

func (m *mainCategUC) Update(categ *domain.MainCateg, userID int64) error {
	// check if the main category exists
	if _, err := m.MainCateg.GetByID(categ.ID, userID); err != nil {
		return err
	}

	// check if the icon exists
	if _, err := m.Icon.GetByID(categ.Icon.ID); err != nil {
		return err
	}

	if err := m.MainCateg.Update(categ); err != nil {
		return err
	}

	return nil
}

func (m *mainCategUC) Delete(id int64) error {
	if err := m.MainCateg.Delete(id); err != nil {
		logger.Error("m.MainCateg.Delete failed", "package", "usecase", "err", err)
		return err
	}

	return nil
}
