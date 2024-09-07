package maincateg

import (
	"context"
	"database/sql"

	"github.com/eyo-chen/expense-tracker-go/internal/adapter/repository/icon"
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
	ID     int64  `json:"id"`
	Name   string `json:"name"`
	Type   string `json:"type"`
	IconID int64  `json:"icon_id" gofacto:"foreignKey,struct:Icon"`
	UserID int64  `json:"user_id" gofacto:"foreignKey,struct:User"`
}

func (r *Repo) Create(categ *domain.MainCateg, userID int64) error {
	stmt := `INSERT INTO main_categories (name, type, user_id, icon_id) VALUES (?, ?, ?, ?)`

	c := cvtToMainCateg(categ, userID)
	if _, err := r.DB.Exec(stmt, c.Name, c.Type, c.UserID, c.IconID); err != nil {
		if errorutil.ParseError(err, uniqueNameUserType) {
			return domain.ErrUniqueNameUserType
		}

		logger.Error("r.DB.Exec failed", "package", packageName, "err", err)
		return err
	}

	return nil
}

func (r *Repo) GetAll(ctx context.Context, userID int64, transType domain.TransactionType) ([]domain.MainCateg, error) {
	stmt := `SELECT mc.id, mc.name, mc.type, i.id, i.url
					 FROM main_categories AS mc
					 LEFT JOIN icons AS i 
					 ON mc.icon_id = i.id
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
		var icon icon.Icon
		if err := rows.Scan(&categ.ID, &categ.Name, &categ.Type, &icon.ID, &icon.URL); err != nil {
			logger.Error("rows.Scan failed", "package", packageName, "err", err)
			return nil, err
		}

		categs = append(categs, cvtToDomainMainCateg(categ, icon))
	}
	defer rows.Close()

	return categs, nil
}

func (r *Repo) Update(categ *domain.MainCateg) error {
	stmt := `UPDATE main_categories SET name = ?, type = ?, icon_id = ? WHERE id = ?`

	c := cvtToMainCateg(categ, 0)
	if _, err := r.DB.Exec(stmt, c.Name, c.Type, c.IconID, c.ID); err != nil {
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
	stmt := `SELECT id, name, type FROM main_categories WHERE id = ? AND user_id = ?`

	var categ MainCateg
	if err := r.DB.QueryRow(stmt, id, userID).Scan(&categ.ID, &categ.Name, &categ.Type); err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrMainCategNotFound
		}

		logger.Error("r.DB.QueryRow failed", "package", packageName, "err", err)
		return nil, err
	}

	domainCateg := cvtToDomainMainCateg(categ, icon.Icon{})
	return &domainCateg, nil
}

func (r *Repo) BatchCreate(ctx context.Context, categs []domain.MainCateg, userID int64) error {
	stmt := `INSERT INTO main_categories (name, type, user_id, icon_id) VALUES `
	args := make([]interface{}, 0, len(categs)*4)
	for i, c := range categs {
		stmt += "(?, ?, ?, ?)"
		if i < len(categs)-1 {
			stmt += ", "
		}

		args = append(args, c.Name, c.Type.ToModelValue(), userID, c.Icon.ID)
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
