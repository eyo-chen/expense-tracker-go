package model

import (
	"database/sql"
	"reflect"

	"github.com/OYE0303/expense-tracker-go/internal/domain"
)

type factory struct {
	db *sql.DB
}

func newFactory(db *sql.DB) *factory {
	return &factory{db: db}
}

// Field: ID, Name, Email, Password_hash
func (f *factory) newUser(overwrites ...map[string]any) (*User, error) {
	stmt := `INSERT INTO users (name, email, password_hash) VALUES (?, ?, ?)`
	user := &User{
		Name:          "test",
		Email:         "test@gmail.com",
		Password_hash: "test",
	}

	for _, o := range overwrites {
		merge(user, o)
	}

	res, err := f.db.Exec(stmt, user.Name, user.Email, user.Password_hash)
	if err != nil {
		return nil, err
	}

	user.ID, err = res.LastInsertId()
	if err != nil {
		return nil, err
	}

	return user, nil
}

// Field: ID, URL
func (f *factory) newIcon(overwrites ...map[string]any) (*Icon, error) {
	stmt := `INSERT INTO icons (url) VALUES (?)`
	icon := &Icon{
		URL: "https://test.com",
	}

	for _, o := range overwrites {
		merge(icon, o)
	}

	res, err := f.db.Exec(stmt, icon.URL)
	if err != nil {
		return nil, err
	}

	icon.ID, err = res.LastInsertId()
	if err != nil {
		return nil, err
	}

	return icon, nil
}

// Field: ID, Name, Type, IconID
func (f *factory) newMainCateg(user *User, overwrites ...map[string]any) (*MainCateg, error) {
	if user == nil {
		var err error
		user, err = f.newUser()
		if err != nil {
			return nil, err
		}
	}

	icon, err := f.newIcon()
	if err != nil {
		return nil, err
	}

	categ := &MainCateg{
		Name:   "test",
		Type:   domain.Expense.ModelValue(),
		IconID: icon.ID,
	}

	for _, o := range overwrites {
		merge(categ, o)
	}

	stmt := `INSERT INTO main_categories (name, type, user_id, icon_id) VALUES (?, ?, ?, ?)`

	res, err := f.db.Exec(stmt, categ.Name, categ.Type, user.ID, categ.IconID)
	if err != nil {
		return nil, err
	}

	categ.ID, err = res.LastInsertId()
	if err != nil {
		return nil, err
	}

	return categ, nil
}

func merge(obj any, values map[string]any) {
	st := reflect.ValueOf(obj).Elem()

	for k, v := range values {
		f := st.FieldByName(k)
		v := reflect.ValueOf(v)
		f.Set(v)
	}
}
