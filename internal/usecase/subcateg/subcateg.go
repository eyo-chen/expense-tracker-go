package subcateg

import (
	"github.com/eyo-chen/expense-tracker-go/internal/domain"
	"github.com/eyo-chen/expense-tracker-go/internal/usecase/interfaces"
)

type UC struct {
	SubCateg  interfaces.SubCategRepo
	MainCateg interfaces.MainCategRepo
}

func New(s interfaces.SubCategRepo, m interfaces.MainCategRepo) *UC {
	return &UC{
		SubCateg:  s,
		MainCateg: m,
	}
}

func (u *UC) Create(categ *domain.SubCateg, userID int64) error {
	// check if the main category exists
	if _, err := u.MainCateg.GetByID(categ.MainCategID, userID); err != nil {
		return err
	}

	return u.SubCateg.Create(categ, userID)
}

func (u *UC) GetByMainCategID(userID, mainCategID int64) ([]*domain.SubCateg, error) {
	return u.SubCateg.GetByMainCategID(userID, mainCategID)
}

func (u *UC) Update(categ *domain.SubCateg, userID int64) error {
	// check if the sub category exists
	subCategByID, err := u.SubCateg.GetByID(categ.ID, userID)
	if err != nil {
		return err
	}

	// check if the input main category ID and the database main category ID are the same
	if subCategByID.MainCategID != categ.MainCategID {
		return domain.ErrMainCategNotFound
	}

	if err := u.SubCateg.Update(categ); err != nil {
		return err
	}

	return nil
}

func (u *UC) Delete(id int64) error {
	return u.SubCateg.Delete(id)
}
