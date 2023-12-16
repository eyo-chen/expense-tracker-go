package usecase

import (
	"github.com/OYE0303/expense-tracker-go/internal/domain"
	"github.com/OYE0303/expense-tracker-go/pkg/logger"
)

type subCategUC struct {
	SubCateg SubCategModel
}

func newSubCategUC(subCateg SubCategModel) *subCategUC {
	return &subCategUC{SubCateg: subCateg}
}

func (s *subCategUC) Create(categ *domain.SubCateg, userID int64) error {
	// check if the sub category name is already taken
	categByUserID, err := s.SubCateg.GetOneByUserID(userID, categ.Name)
	if err != nil && err != domain.ErrDataNotFound {
		logger.Error("s.SubCateg.GetOneByUserID failed", "package", "usecase", "err", err)
		return err
	}
	if categByUserID != nil {
		return domain.ErrDataAlreadyExists
	}

	if err := s.SubCateg.Create(categ, userID); err != nil {
		logger.Error("s.SubCateg.Create failed", "package", "usecase", "err", err)
		return err
	}

	return nil
}

func (s *subCategUC) Update(categ *domain.SubCateg, userID int64) error {
	categByID, err := s.SubCateg.GetByID(categ.ID)
	if err != nil && err != domain.ErrDataNotFound {
		logger.Error("s.SubCateg.GetByID failed", "package", "usecase", "err", err)
		return err
	}
	if categByID == nil {
		return domain.ErrDataNotFound
	}

	// check if the sub category name is already taken
	categByUserID, err := s.SubCateg.GetOneByUserID(userID, categ.Name)
	if err != nil && err != domain.ErrDataNotFound {
		logger.Error("s.SubCateg.GetOneByUserID failed", "package", "usecase", "err", err)
		return err
	}
	if categByUserID != nil {
		return domain.ErrDataAlreadyExists
	}

	if err := s.SubCateg.Update(categ); err != nil {
		logger.Error("s.SubCateg.Update failed", "package", "usecase", "err", err)
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
