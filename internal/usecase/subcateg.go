package usecase

import (
	"github.com/OYE0303/expense-tracker-go/internal/domain"
	"github.com/OYE0303/expense-tracker-go/pkg/logger"
)

type subCategUC struct {
	SubCateg  SubCategModel
	MainCateg MainCategModel
}

func newSubCategUC(s SubCategModel, m MainCategModel) *subCategUC {
	return &subCategUC{
		SubCateg:  s,
		MainCateg: m,
	}
}

func (s *subCategUC) Create(categ *domain.SubCateg, userID int64) error {
	// check if the main category exists
	if _, err := s.MainCateg.GetByID(categ.MainCategID, userID); err != nil {
		return err
	}

	if err := s.SubCateg.Create(categ, userID); err != nil {
		return err
	}

	return nil
}

func (s *subCategUC) GetAll(userID int64) ([]*domain.SubCateg, error) {
	categs, err := s.SubCateg.GetAll(userID)
	if err != nil {
		logger.Error("s.SubCateg.GetAll failed", "package", "usecase", "err", err)
		return nil, err
	}

	return categs, nil
}

func (s *subCategUC) GetByMainCategID(userID, mainCategID int64) ([]*domain.SubCateg, error) {
	categs, err := s.SubCateg.GetByMainCategID(userID, mainCategID)
	if err != nil {
		logger.Error("s.SubCateg.GetByMainCategID failed", "package", "usecase", "err", err)
		return nil, err
	}

	return categs, nil
}

func (s *subCategUC) Update(categ *domain.SubCateg, userID int64) error {
	// check if the sub category exists
	if _, err := s.SubCateg.GetByID(categ.ID, userID); err != nil {
		return err
	}

	if err := s.SubCateg.Update(categ); err != nil {
		return err
	}

	return nil
}

func (s *subCategUC) Delete(id int64) error {
	if err := s.SubCateg.Delete(id); err != nil {
		logger.Error("s.SubCateg.Delete failed", "package", "usecase", "err", err)
		return err
	}

	return nil
}
