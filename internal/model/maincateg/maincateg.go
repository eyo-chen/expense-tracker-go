package maincateg

import (
	"database/sql"

	"github.com/OYE0303/expense-tracker-go/internal/domain"
	"github.com/OYE0303/expense-tracker-go/internal/model/icon"
	"github.com/OYE0303/expense-tracker-go/pkg/errorutil"
	"github.com/OYE0303/expense-tracker-go/pkg/logger"
)

const (
	UniqueIconUser     = "main_categories.unique_icon_user"
	UniqueNameUserType = "main_categories.unique_name_user_type"
)

type MainCategModel struct {
	DB *sql.DB
}

func NewMainCategModel(db *sql.DB) *MainCategModel {
	return &MainCategModel{DB: db}
}

type MainCateg struct {
	ID     int64  `json:"id"`
	Name   string `json:"name"`
	Type   string `json:"type"`
	IconID int64  `json:"icon_id" factory:"Icon,icons"`
	UserID int64  `json:"user_id" factory:"User,users"`
}

func (m *MainCategModel) Create(categ *domain.MainCateg, userID int64) error {
	stmt := `INSERT INTO main_categories (name, type, user_id, icon_id) VALUES (?, ?, ?, ?)`

	c := cvtToMainCateg(categ, userID)
	if _, err := m.DB.Exec(stmt, c.Name, c.Type, c.UserID, c.IconID); err != nil {
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

func (m *MainCategModel) GetAll(userID int64, categType domain.MainCategType) ([]domain.MainCateg, error) {
	stmt := `SELECT mc.id, mc.name, mc.type, i.id, i.url
					 FROM main_categories AS mc
					 LEFT JOIN icons AS i 
					 ON mc.icon_id = i.id
					 WHERE user_id = ?`

	if categType != domain.UnSpecified {
		stmt += ` AND type = ` + categType.ModelValue()
	}

	rows, err := m.DB.Query(stmt, userID)
	if err != nil {
		logger.Error("m.DB.Query failed", "package", "model", "err", err)
		return nil, err
	}
	defer rows.Close()

	var categs []domain.MainCateg
	for rows.Next() {
		var categ MainCateg
		var icon icon.Icon
		if err := rows.Scan(&categ.ID, &categ.Name, &categ.Type, &icon.ID, &icon.URL); err != nil {
			logger.Error("rows.Scan failed", "package", "model", "err", err)
			return nil, err
		}

		categs = append(categs, getAllCvtToDomainMainCateg(categ, icon))
	}

	return categs, nil
}

func (m *MainCategModel) Update(categ *domain.MainCateg) error {
	stmt := `UPDATE main_categories SET name = ?, type = ?, icon_id = ? WHERE id = ?`

	c := cvtToMainCateg(categ, 0)
	if _, err := m.DB.Exec(stmt, c.Name, c.Type, c.IconID, c.ID); err != nil {
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
			return nil, domain.ErrMainCategNotFound
		}

		logger.Error("m.DB.QueryRow failed", "package", "model", "err", err)
		return nil, err
	}

	return cvtToDomainMainCateg(&categ, nil), nil
}
