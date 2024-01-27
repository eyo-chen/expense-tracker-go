package subcateg

import (
	"github.com/OYE0303/expense-tracker-go/internal/domain"
	"github.com/OYE0303/expense-tracker-go/internal/usecase/interfaces"
	"github.com/OYE0303/expense-tracker-go/pkg/logger"
)

type SubCategUC struct {
	SubCateg  interfaces.SubCategModel
	MainCateg interfaces.MainCategModel
}

func NewSubCategUC(s interfaces.SubCategModel, m interfaces.MainCategModel) *SubCategUC {
	return &SubCategUC{
		SubCateg:  s,
		MainCateg: m,
	}
}

func (s *SubCategUC) Create(categ *domain.SubCateg, userID int64) error {
	// check if the main category exists
	if _, err := s.MainCateg.GetByID(categ.MainCategID, userID); err != nil {
		return err
	}

	if err := s.SubCateg.Create(categ, userID); err != nil {
		return err
	}

	return nil
}

func (s *SubCategUC) GetAll(userID int64) ([]*domain.SubCateg, error) {
	categs, err := s.SubCateg.GetAll(userID)
	if err != nil {
		logger.Error("s.SubCateg.GetAll failed", "package", "usecase", "err", err)
		return nil, err
	}

	return categs, nil
}

func (s *SubCategUC) GetByMainCategID(userID, mainCategID int64) ([]*domain.SubCateg, error) {
	categs, err := s.SubCateg.GetByMainCategID(userID, mainCategID)
	if err != nil {
		logger.Error("s.SubCateg.GetByMainCategID failed", "package", "usecase", "err", err)
		return nil, err
	}

	return categs, nil
}

func (s *SubCategUC) Update(categ *domain.SubCateg, userID int64) error {
	// check if the sub category exists
	subCategByID, err := s.SubCateg.GetByID(categ.ID, userID)
	if err != nil {
		return err
	}

	// check if the input main category ID and the database main category ID are the same
	if subCategByID.MainCategID != categ.MainCategID {
		return domain.ErrMainCategNotFound
	}

	if err := s.SubCateg.Update(categ); err != nil {
		return err
	}

	return nil
}

func (s *SubCategUC) Delete(id int64) error {
	if err := s.SubCateg.Delete(id); err != nil {
		logger.Error("s.SubCateg.Delete failed", "package", "usecase", "err", err)
		return err
	}

	return nil
}
