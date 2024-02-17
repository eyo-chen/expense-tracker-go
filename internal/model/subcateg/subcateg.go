package subcateg

import (
	"database/sql"

	"github.com/OYE0303/expense-tracker-go/internal/domain"
	"github.com/OYE0303/expense-tracker-go/pkg/errorutil"
	"github.com/OYE0303/expense-tracker-go/pkg/logger"
)

const (
	UniqueNameUserMainCategory = "sub_categories.unique_name_user_maincategory"
)

type SubCategModel struct {
	DB *sql.DB
}

func NewSubCategModel(db *sql.DB) *SubCategModel {
	return &SubCategModel{DB: db}
}

type SubCateg struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	UserID      int64  `json:"user_id" factory:"User,users"`
	MainCategID int64  `json:"main_category_id" factory:"MainCateg,main_categories"`
}

func (m *SubCategModel) Create(categ *domain.SubCateg, userID int64) error {
	stmt := `INSERT INTO sub_categories (name, user_id, main_category_id) VALUES (?, ?, ?)`

	c := cvtToSubCateg(categ, userID)
	if _, err := m.DB.Exec(stmt, c.Name, c.UserID, c.MainCategID); err != nil {
		if errorutil.ParseError(err, UniqueNameUserMainCategory) {
			return domain.ErrUniqueNameUserMainCateg
		}

		logger.Error("m.DB.Exec failed", "package", "model", "err", err)
		return err
	}

	return nil
}

func (m *SubCategModel) GetAll(userID int64) ([]*domain.SubCateg, error) {
	stmt := `SELECT id, name, main_category_id FROM sub_categories WHERE user_id = ?`

	rows, err := m.DB.Query(stmt, userID)
	if err != nil {
		logger.Error("m.DB.Query failed", "package", "model", "err", err)
		return nil, err
	}
	defer rows.Close()

	var categs []*domain.SubCateg
	for rows.Next() {
		var categ SubCateg
		if err := rows.Scan(&categ.ID, &categ.Name, &categ.MainCategID); err != nil {
			logger.Error("rows.Scan failed", "package", "model", "err", err)
			return nil, err
		}

		categs = append(categs, cvtToDomainSubCateg(&categ))
	}

	return categs, nil
}

func (m *SubCategModel) GetByMainCategID(userID, mainCategID int64) ([]*domain.SubCateg, error) {
	stmt := `SELECT id, name, main_category_id FROM sub_categories WHERE user_id = ? AND main_category_id = ?`

	rows, err := m.DB.Query(stmt, userID, mainCategID)
	if err != nil {
		logger.Error("m.DB.Query failed", "package", "model", "err", err)
		return nil, err
	}
	defer rows.Close()

	var categs []*domain.SubCateg
	for rows.Next() {
		var categ SubCateg
		if err := rows.Scan(&categ.ID, &categ.Name, &categ.MainCategID); err != nil {
			logger.Error("rows.Scan failed", "package", "model", "err", err)
			return nil, err
		}

		categs = append(categs, cvtToDomainSubCateg(&categ))
	}

	return categs, nil
}

func (m *SubCategModel) Update(categ *domain.SubCateg) error {
	stmt := `UPDATE sub_categories SET name = ? WHERE id = ?`

	c := cvtToSubCateg(categ, 0)
	if _, err := m.DB.Exec(stmt, c.Name, c.ID); err != nil {
		if errorutil.ParseError(err, UniqueNameUserMainCategory) {
			return domain.ErrUniqueNameUserMainCateg
		}

		logger.Error("m.DB.Exec failed", "package", "model", "err", err)
		return err
	}

	return nil
}

func (m *SubCategModel) Delete(id int64) error {
	stmt := `DELETE FROM sub_categories WHERE id = ?`

	if _, err := m.DB.Exec(stmt, id); err != nil {
		logger.Error("m.DB.Exec failed", "package", "model", "err", err)
		return err
	}

	return nil
}

func (m *SubCategModel) GetByID(id, userID int64) (*domain.SubCateg, error) {
	stmt := `SELECT id, name, main_category_id FROM sub_categories WHERE id = ? AND user_id = ?`

	var categ SubCateg
	if err := m.DB.QueryRow(stmt, id, userID).Scan(&categ.ID, &categ.Name, &categ.MainCategID); err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrSubCategNotFound
		}

		logger.Error("m.DB.QueryRow failed", "package", "model", "err", err)
		return nil, err
	}

	return cvtToDomainSubCateg(&categ), nil
}
