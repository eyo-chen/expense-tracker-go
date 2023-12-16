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
	categbyUserID, err := s.SubCateg.GetOneByUserID(userID, categ.Name)
	if err != nil && err != domain.ErrDataNotFound {
		logger.Error("s.SubCateg.GetOneByUserID failed", "package", "usecase", "err", err)
		return err
	}
	if categbyUserID != nil {
		return domain.ErrDataAlreadyExists
	}

	if err := s.SubCateg.Create(categ, userID); err != nil {
		logger.Error("s.SubCateg.Create failed", "package", "usecase", "err", err)
		return err
	}

	return nil
}
