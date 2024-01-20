package model

import (
	"database/sql"

	"github.com/OYE0303/expense-tracker-go/internal/domain"
	"github.com/OYE0303/expense-tracker-go/pkg/errorutil"
	"github.com/OYE0303/expense-tracker-go/pkg/logger"
)

type MainCategModel struct {
	DB *sql.DB
}

func newMainCategModel(db *sql.DB) *MainCategModel {
	return &MainCategModel{DB: db}
}

const (
	UniqueIconUser     = "main_categories.unique_icon_user"
	UniqueNameUserType = "main_categories.unique_name_user_type"
)

type MainCateg struct {
	ID     int64  `json:"id"`
	Name   string `json:"name"`
	Type   string `json:"type"`
	IconID int64  `json:"icon_id"`
}

func (m *MainCategModel) Create(categ *domain.MainCateg, userID int64) error {
	stmt := `INSERT INTO main_categories (name, type, user_id, icon_id) VALUES (?, ?, ?, ?)`

	if _, err := m.DB.Exec(stmt, categ.Name, categ.Type.ModelValue(), userID, categ.Icon.ID); err != nil {
		if errorutil.ParseError(err, UniqueNameUserType) {
			return domain.ErrUniqueNameUserType
		}

		if errorutil.ParseError(err, UniqueIconUser) {
			return domain.ErrUniqueIconUser
		}

		logger.Error("m.DB.Exec failed", "package", "model", "err", err)
		return err
	}

	return nil
}

func (m *MainCategModel) GetAll(userID int64) ([]*domain.MainCateg, error) {
	stmt := `SELECT id, name, type, icon_id FROM main_categories WHERE user_id = ?`

	rows, err := m.DB.Query(stmt, userID)
	if err != nil {
		logger.Error("m.DB.Query failed", "package", "model", "err", err)
		return nil, err
	}
	defer rows.Close()

	var categs []*domain.MainCateg
	for rows.Next() {
		var categ MainCateg
		if err := rows.Scan(&categ.ID, &categ.Name, &categ.Type, &categ.IconID); err != nil {
			logger.Error("rows.Scan failed", "package", "model", "err", err)
			return nil, err
		}

		categs = append(categs, cvtToDomainMainCateg(&categ, nil))
	}

	return categs, nil
}

func (m *MainCategModel) Update(categ *domain.MainCateg) error {
	stmt := `UPDATE main_categories SET name = ?, type = ?, icon_id = ? WHERE id = ?`

	if _, err := m.DB.Exec(stmt, categ.Name, categ.Type.ModelValue(), categ.Icon.ID, categ.ID); err != nil {
		logger.Error("m.DB.Exec failed", "package", "model", "err", err)
		return err
	}

	return nil
}

func (m *MainCategModel) Delete(id int64) error {
	stmt := `DELETE FROM main_categories WHERE id = ?`

	if _, err := m.DB.Exec(stmt, id); err != nil {
		logger.Error("m.DB.Exec failed", "package", "model", "err", err)
		return err
	}

	return nil
}

func (m *MainCategModel) GetByID(id, userID int64) (*domain.MainCateg, error) {
	stmt := `SELECT id, name, type FROM main_categories WHERE id = ? AND user_id = ?`

	var categ MainCateg
	if err := m.DB.QueryRow(stmt, id, userID).Scan(&categ.ID, &categ.Name, &categ.Type); err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrDataNotFound
		}

		logger.Error("m.DB.QueryRow failed", "package", "model", "err", err)
		return nil, err
	}

	return cvtToDomainMainCateg(&categ, nil), nil
}

func (m *MainCategModel) GetOne(inputCateg *domain.MainCateg, userID int64) (*domain.MainCateg, error) {
	stmt := `SELECT id, name, type FROM main_categories WHERE user_id = ? AND name = ? AND type = ?`

	var categ MainCateg
	if err := m.DB.QueryRow(stmt, userID, inputCateg.Name, inputCateg.Type.ModelValue()).Scan(&categ.ID, &categ.Name, &categ.Type); err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrDataNotFound
		}

		logger.Error("m.DB.QueryRow failed", "package", "model", "err", err)
		return nil, err
	}

	return cvtToDomainMainCateg(&categ, nil), nil
}

func (m *MainCategModel) GetFullInfoByID(id, userID int64) (*domain.MainCateg, error) {
	return nil, nil
}
