package subcateg

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/eyo-chen/expense-tracker-go/internal/domain"
	"github.com/eyo-chen/expense-tracker-go/pkg/errorutil"
	"github.com/eyo-chen/expense-tracker-go/pkg/logger"
)

const (
	uniqueNameUserMainCategory = "sub_categories.unique_name_user_maincategory"
	packageName                = "adapter/repository/subcateg"
)

type Repo struct {
	DB *sql.DB
}

func New(db *sql.DB) *Repo {
	return &Repo{DB: db}
}

type SubCateg struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	UserID      int64  `json:"user_id" gofacto:"foreignKey,struct:User"`
	MainCategID int64  `json:"main_category_id" gofacto:"foreignKey,struct:MainCateg,table:main_categories" mysqlf:"main_category_id"`
}

func (r *Repo) Create(categ *domain.SubCateg, userID int64) error {
	stmt := `INSERT INTO sub_categories (name, user_id, main_category_id) VALUES (?, ?, ?)`

	c := cvtToSubCateg(categ, userID)
	if _, err := r.DB.Exec(stmt, c.Name, c.UserID, c.MainCategID); err != nil {
		if errorutil.ParseError(err, uniqueNameUserMainCategory) {
			return domain.ErrUniqueNameUserMainCateg
		}

		logger.Error("r.DB.Exec failed", "package", packageName, "err", err)
		return err
	}

	return nil
}

func (r *Repo) GetByMainCategID(userID, mainCategID int64) ([]*domain.SubCateg, error) {
	stmt := `SELECT id, name, main_category_id FROM sub_categories WHERE user_id = ? AND main_category_id = ?`

	rows, err := r.DB.Query(stmt, userID, mainCategID)
	if err != nil {
		logger.Error("r.DB.Query failed", "package", packageName, "err", err)
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			logger.Error("Unable to close rows", "package", packageName, "err", err)
		}
	}()

	var categs []*domain.SubCateg
	for rows.Next() {
		var categ SubCateg
		if err := rows.Scan(&categ.ID, &categ.Name, &categ.MainCategID); err != nil {
			logger.Error("rows.Scan failed", "package", packageName, "err", err)
			return nil, err
		}

		categs = append(categs, cvtToDomainSubCateg(&categ))
	}
	defer func() {
		if err := rows.Close(); err != nil {
			logger.Error("Unable to close rows", "package", packageName, "err", err)
		}
	}()

	return categs, nil
}

func (r *Repo) Update(categ *domain.SubCateg) error {
	stmt := `UPDATE sub_categories SET name = ? WHERE id = ?`

	c := cvtToSubCateg(categ, 0)
	if _, err := r.DB.Exec(stmt, c.Name, c.ID); err != nil {
		if errorutil.ParseError(err, uniqueNameUserMainCategory) {
			return domain.ErrUniqueNameUserMainCateg
		}

		logger.Error("r.DB.Exec failed", "package", packageName, "err", err)
		return err
	}

	return nil
}

func (r *Repo) Delete(id int64) error {
	stmt := `DELETE FROM sub_categories WHERE id = ?`

	if _, err := r.DB.Exec(stmt, id); err != nil {
		logger.Error("r.DB.Exec failed", "package", packageName, "err", err)
		return err
	}

	return nil
}

func (r *Repo) GetByID(id, userID int64) (*domain.SubCateg, error) {
	stmt := `SELECT id, name, main_category_id FROM sub_categories WHERE id = ? AND user_id = ?`

	var categ SubCateg
	if err := r.DB.QueryRow(stmt, id, userID).Scan(&categ.ID, &categ.Name, &categ.MainCategID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrSubCategNotFound
		}

		logger.Error("r.DB.QueryRow failed", "package", packageName, "err", err)
		return nil, err
	}

	return cvtToDomainSubCateg(&categ), nil
}

func (r *Repo) BatchCreate(ctx context.Context, categs []domain.SubCateg, userID int64) error {
	var sb strings.Builder
	sb.WriteString(`INSERT INTO sub_categories (name, user_id, main_category_id) VALUES `)
	args := make([]interface{}, 0, len(categs)*3)
	for i, c := range categs {
		sb.WriteString("(?, ?, ?)")
		if i < len(categs)-1 {
			sb.WriteString(", ")
		}

		args = append(args, c.Name, userID, c.MainCategID)
	}

	if _, err := r.DB.ExecContext(ctx, sb.String(), args...); err != nil {
		if errorutil.ParseError(err, uniqueNameUserMainCategory) {
			return domain.ErrUniqueNameUserMainCateg
		}

		logger.Error("r.DB.ExecContext CreateBatch failed", "package", packageName, "err", err)
		return err
	}

	return nil
}
