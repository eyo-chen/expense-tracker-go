package maincateg

import (
	"context"
	"database/sql"
	"errors"

	"github.com/eyo-chen/expense-tracker-go/internal/domain"
	"github.com/eyo-chen/expense-tracker-go/pkg/errorutil"
	"github.com/eyo-chen/expense-tracker-go/pkg/logger"
)

const (
	uniqueNameUserType = "main_categories.unique_name_user_type"
	packageName        = "adapter/repository/maincateg"
)

type Repo struct {
	DB *sql.DB
}

func New(db *sql.DB) *Repo {
	return &Repo{DB: db}
}

type MainCateg struct {
	ID       int64
	Name     string
	Type     string
	UserID   int64 `gofacto:"foreignKey,struct:User"`
	IconType string
	IconData string
}

func (r *Repo) Create(ctx context.Context, categ domain.MainCateg, userID int64) error {
	stmt := `INSERT INTO main_categories (name, type, user_id, icon_type, icon_data) VALUES (?, ?, ?, ?, ?)`

	c := cvtToMainCateg(categ, userID)
	if _, err := r.DB.ExecContext(ctx, stmt, c.Name, c.Type, c.UserID, c.IconType, c.IconData); err != nil {
		if errorutil.ParseError(err, uniqueNameUserType) {
			return domain.ErrUniqueNameUserType
		}

		logger.Error("r.DB.Exec failed", "package", packageName, "err", err)
		return err
	}

	return nil
}

func (r *Repo) GetAll(ctx context.Context, userID int64, transType domain.TransactionType) ([]domain.MainCateg, error) {
	stmt := `SELECT id, name, type, icon_type, icon_data
					 FROM main_categories
					 WHERE user_id = ?`

	if transType.IsValid() {
		stmt += ` AND type = ` + transType.ToModelValue()
	}

	rows, err := r.DB.QueryContext(ctx, stmt, userID)
	if err != nil {
		logger.Error("r.DB.Query failed", "package", packageName, "err", err)
		return nil, err
	}
	defer rows.Close()

	var categs []domain.MainCateg
	for rows.Next() {
		var categ MainCateg
		if err := rows.Scan(&categ.ID, &categ.Name, &categ.Type, &categ.IconType, &categ.IconData); err != nil {
			logger.Error("rows.Scan failed", "package", packageName, "err", err)
			return nil, err
		}

		categs = append(categs, cvtToDomainMainCateg(categ))
	}
	defer rows.Close()

	return categs, nil
}

func (r *Repo) Update(ctx context.Context, categ domain.MainCateg) error {
	stmt := `UPDATE main_categories SET name = ?, type = ?, icon_type = ?, icon_data = ? WHERE id = ?`

	c := cvtToMainCateg(categ, 0)
	if _, err := r.DB.ExecContext(ctx, stmt, c.Name, c.Type, c.IconType, c.IconData, c.ID); err != nil {
		if errorutil.ParseError(err, uniqueNameUserType) {
			return domain.ErrUniqueNameUserType
		}

		logger.Error("r.DB.Exec failed", "package", packageName, "err", err)
		return err
	}

	return nil
}

func (r *Repo) Delete(id int64) error {
	stmt := `DELETE FROM main_categories WHERE id = ?`

	if _, err := r.DB.Exec(stmt, id); err != nil {
		logger.Error("r.DB.Exec failed", "package", packageName, "err", err)
		return err
	}

	return nil
}

func (r *Repo) GetByID(id, userID int64) (*domain.MainCateg, error) {
	stmt := `SELECT id, name, type, icon_type, icon_data FROM main_categories WHERE id = ? AND user_id = ?`

	var categ MainCateg
	if err := r.DB.QueryRow(stmt, id, userID).Scan(&categ.ID, &categ.Name, &categ.Type, &categ.IconType, &categ.IconData); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrMainCategNotFound
		}

		logger.Error("r.DB.QueryRow failed", "package", packageName, "err", err)
		return nil, err
	}

	domainCateg := cvtToDomainMainCateg(categ)
	return &domainCateg, nil
}

func (r *Repo) BatchCreate(ctx context.Context, categs []domain.MainCateg, userID int64) error {
	stmt := `INSERT INTO main_categories (name, type, user_id, icon_type, icon_data) VALUES `
	args := make([]interface{}, 0, len(categs)*6)
	for i, c := range categs {
		stmt += "(?, ?, ?, ?, ?)"
		if i < len(categs)-1 {
			stmt += ", "
		}

		args = append(args, c.Name, c.Type.ToModelValue(), userID, c.IconType.ToModelValue(), c.IconData)
	}

	if _, err := r.DB.ExecContext(ctx, stmt, args...); err != nil {
		if errorutil.ParseError(err, uniqueNameUserType) {
			return domain.ErrUniqueNameUserType
		}

		logger.Error("r.DB.Exec failed", "package", packageName, "err", err)
		return err
	}

	return nil
}
