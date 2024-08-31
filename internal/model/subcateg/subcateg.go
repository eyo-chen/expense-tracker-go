package subcateg

import (
	"context"
	"database/sql"

	"github.com/eyo-chen/expense-tracker-go/internal/domain"
	"github.com/eyo-chen/expense-tracker-go/pkg/errorutil"
	"github.com/eyo-chen/expense-tracker-go/pkg/logger"
)

const (
	UniqueNameUserMainCategory = "sub_categories.unique_name_user_maincategory"
	packageName                = "model/subcateg"
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
	UserID      int64  `json:"user_id" gofacto:"foreignKey,struct:User"`
	MainCategID int64  `json:"main_category_id" gofacto:"foreignKey,struct:MainCateg,table:main_categories" mysqlf:"main_category_id"`
}

func (s *SubCategModel) Create(categ *domain.SubCateg, userID int64) error {
	stmt := `INSERT INTO sub_categories (name, user_id, main_category_id) VALUES (?, ?, ?)`

	c := cvtToSubCateg(categ, userID)
	if _, err := s.DB.Exec(stmt, c.Name, c.UserID, c.MainCategID); err != nil {
		if errorutil.ParseError(err, UniqueNameUserMainCategory) {
			return domain.ErrUniqueNameUserMainCateg
		}

		logger.Error("m.DB.Exec failed", "package", packageName, "err", err)
		return err
	}

	return nil
}

func (s *SubCategModel) GetByMainCategID(userID, mainCategID int64) ([]*domain.SubCateg, error) {
	stmt := `SELECT id, name, main_category_id FROM sub_categories WHERE user_id = ? AND main_category_id = ?`

	rows, err := s.DB.Query(stmt, userID, mainCategID)
	if err != nil {
		logger.Error("m.DB.Query failed", "package", packageName, "err", err)
		return nil, err
	}
	defer rows.Close()

	var categs []*domain.SubCateg
	for rows.Next() {
		var categ SubCateg
		if err := rows.Scan(&categ.ID, &categ.Name, &categ.MainCategID); err != nil {
			logger.Error("rows.Scan failed", "package", packageName, "err", err)
			return nil, err
		}

		categs = append(categs, cvtToDomainSubCateg(&categ))
	}
	defer rows.Close()

	return categs, nil
}

func (s *SubCategModel) Update(categ *domain.SubCateg) error {
	stmt := `UPDATE sub_categories SET name = ? WHERE id = ?`

	c := cvtToSubCateg(categ, 0)
	if _, err := s.DB.Exec(stmt, c.Name, c.ID); err != nil {
		if errorutil.ParseError(err, UniqueNameUserMainCategory) {
			return domain.ErrUniqueNameUserMainCateg
		}

		logger.Error("m.DB.Exec failed", "package", packageName, "err", err)
		return err
	}

	return nil
}

func (s *SubCategModel) Delete(id int64) error {
	stmt := `DELETE FROM sub_categories WHERE id = ?`

	if _, err := s.DB.Exec(stmt, id); err != nil {
		logger.Error("m.DB.Exec failed", "package", packageName, "err", err)
		return err
	}

	return nil
}

func (s *SubCategModel) GetByID(id, userID int64) (*domain.SubCateg, error) {
	stmt := `SELECT id, name, main_category_id FROM sub_categories WHERE id = ? AND user_id = ?`

	var categ SubCateg
	if err := s.DB.QueryRow(stmt, id, userID).Scan(&categ.ID, &categ.Name, &categ.MainCategID); err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrSubCategNotFound
		}

		logger.Error("m.DB.QueryRow failed", "package", packageName, "err", err)
		return nil, err
	}

	return cvtToDomainSubCateg(&categ), nil
}

func (s *SubCategModel) BatchCreate(ctx context.Context, categs []domain.SubCateg, userID int64) error {
	stmt := `INSERT INTO sub_categories (name, user_id, main_category_id) VALUES `
	args := make([]interface{}, 0, len(categs)*3)
	for i, c := range categs {
		stmt += "(?, ?, ?)"
		if i < len(categs)-1 {
			stmt += ", "
		}

		args = append(args, c.Name, userID, c.MainCategID)
	}

	if _, err := s.DB.ExecContext(ctx, stmt, args...); err != nil {
		if errorutil.ParseError(err, UniqueNameUserMainCategory) {
			return domain.ErrUniqueNameUserMainCateg
		}

		logger.Error("m.DB.ExecContext CreateBatch failed", "package", packageName, "err", err)
		return err
	}

	return nil
}
