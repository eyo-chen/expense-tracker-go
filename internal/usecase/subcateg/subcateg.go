package subcateg

import (
	"github.com/eyo-chen/expense-tracker-go/internal/domain"
	"github.com/eyo-chen/expense-tracker-go/internal/usecase/interfaces"
)

type SubCategUC struct {
	SubCateg  interfaces.SubCategRepo
	MainCateg interfaces.MainCategRepo
}

func NewSubCategUC(s interfaces.SubCategRepo, m interfaces.MainCategRepo) *SubCategUC {
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

	return s.SubCateg.Create(categ, userID)
}

func (s *SubCategUC) GetByMainCategID(userID, mainCategID int64) ([]*domain.SubCateg, error) {
	return s.SubCateg.GetByMainCategID(userID, mainCategID)
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
	return s.SubCateg.Delete(id)
}
