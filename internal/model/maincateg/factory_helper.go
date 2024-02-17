package maincateg

import (
	"database/sql"
	"fmt"

	"github.com/OYE0303/expense-tracker-go/internal/domain"
	"github.com/OYE0303/expense-tracker-go/internal/model/icon"
	"github.com/OYE0303/expense-tracker-go/internal/model/user"
)

func BluePrint(i int, last MainCateg) MainCateg {
	return MainCateg{
		Name: "test" + fmt.Sprint(i),
		Type: domain.Income.ModelValue(),
	}
}

func Inserter(db *sql.DB, c MainCateg) (MainCateg, error) {
	stmt := `INSERT INTO main_categories (name, type, user_id, icon_id) VALUES (?, ?, ?, ?)`

	res, err := db.Exec(stmt, c.Name, c.Type, c.UserID, c.IconID)
	if err != nil {
		return MainCateg{}, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return MainCateg{}, err
	}

	c.ID = id
	return c, nil
}

func BluePrintIcon(i int, last icon.Icon) icon.Icon {
	return icon.Icon{
		URL: "test" + fmt.Sprint(i),
	}
}

func InsertIcon(db *sql.DB, i icon.Icon) (icon.Icon, error) {
	stmt := `INSERT INTO icons (url) VALUES (?)`
	res, err := db.Exec(stmt, i.URL)
	if err != nil {
		return icon.Icon{}, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return icon.Icon{}, err
	}

	i.ID = id

	return i, nil
}

func BluePrintUser(i int, last user.User) user.User {
	return user.User{
		Name:          "test" + fmt.Sprint(i),
		Email:         "test" + fmt.Sprint(i) + "@gmail.com",
		Password_hash: "test" + fmt.Sprint(i),
	}
}

func InserterUser(db *sql.DB, u user.User) (user.User, error) {
	stmt := `INSERT INTO users (name, email, password_hash) VALUES (?, ?, ?)`
	res, err := db.Exec(stmt, u.Name, u.Email, u.Password_hash)
	if err != nil {
		return user.User{}, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return user.User{}, err
	}

	u.ID = id

	return u, nil
}
