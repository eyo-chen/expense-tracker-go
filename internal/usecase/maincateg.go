package usecase

import (
	"errors"

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
	// check if the main category name is already taken
	categbyUserID, err := m.MainCateg.GetOne(categ, userID)
	if err != nil && !errors.Is(err, domain.ErrDataNotFound) {
		logger.Error("m.MainCateg.GetOne failed", "package", "usecase", "err", err)
		return err
	}
	if categbyUserID != nil {
		return domain.ErrDataAlreadyExists
	}

	// check if the icon exists
	icon, err := m.Icon.GetByID(categ.Icon.ID)
	if err != nil && !errors.Is(err, domain.ErrDataNotFound) {
		logger.Error("m.Icon.GetByID failed", "package", "usecase", "err", err)
		return err
	}
	if icon == nil {
		return domain.ErrDataNotFound
	}

	if err := m.MainCateg.Create(categ, userID); err != nil {
		logger.Error("m.MainCateg.Create failed", "package", "usecase", "err", err)
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
	_, err := m.MainCateg.GetByID(categ.ID, userID)
	if errors.Is(err, domain.ErrDataNotFound) {
		return domain.ErrDataNotFound
	}
	if err != nil {
		logger.Error("m.MainCateg.GetByID failed", "package", "usecase", "err", err)
		return err
	}

	// check if the main category name is already taken
	categbyUserID, err := m.MainCateg.GetOne(categ, userID)
	if err != nil && !errors.Is(err, domain.ErrDataNotFound) {
		logger.Error("m.MainCateg.GetOne failed", "package", "usecase", "err", err)
		return err
	}
	if categbyUserID != nil {
		return domain.ErrDataAlreadyExists
	}

	// check if the icon exists
	icon, err := m.Icon.GetByID(categ.Icon.ID)
	if err != nil && !errors.Is(err, domain.ErrDataNotFound) {
		logger.Error("m.Icon.GetByID failed", "package", "usecase", "err", err)
		return err
	}
	if icon == nil {
		return domain.ErrDataNotFound
	}

	if err := m.MainCateg.Update(categ); err != nil {
		logger.Error("m.MainCateg.Update failed", "package", "usecase", "err", err)
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
