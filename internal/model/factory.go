package model

import (
	"database/sql"
	"reflect"

	"github.com/OYE0303/expense-tracker-go/internal/domain"
	"github.com/OYE0303/expense-tracker-go/internal/model/icon"
	"github.com/OYE0303/expense-tracker-go/internal/model/maincateg"
	"github.com/OYE0303/expense-tracker-go/internal/model/subcateg"
	"github.com/OYE0303/expense-tracker-go/internal/model/user"
)

type Factory struct {
	db *sql.DB
}

func NewFactory(db *sql.DB) *Factory {
	return &Factory{db: db}
}

// Field: ID, Name, Email, Password_hash
func (f *Factory) NewUser(overwrites ...map[string]any) (*user.User, error) {
	stmt := `INSERT INTO users (name, email, password_hash) VALUES (?, ?, ?)`
	user := &user.User{
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
func (f *Factory) NewIcon(overwrites ...map[string]any) (*icon.Icon, error) {
	stmt := `INSERT INTO icons (url) VALUES (?)`
	icon := &icon.Icon{
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
func (f *Factory) NewMainCateg(user *user.User, overwrites ...map[string]any) (*maincateg.MainCateg, error) {
	if user == nil {
		var err error
		user, err = f.NewUser()
		if err != nil {
			return nil, err
		}
	}

	icon, err := f.NewIcon()
	if err != nil {
		return nil, err
	}

	categ := &maincateg.MainCateg{
		Name:   "test",
		Type:   domain.Expense.ToModelValue(),
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

func (f *Factory) NewSubCateg(user *user.User, mainCateg *maincateg.MainCateg, overwrites ...map[string]any) (*subcateg.SubCateg, error) {
	if user == nil {
		var err error
		user, err = f.NewUser()
		if err != nil {
			return nil, err
		}
	}

	if mainCateg == nil {
		var err error
		mainCateg, err = f.NewMainCateg(user)
		if err != nil {
			return nil, err
		}
	}

	categ := &subcateg.SubCateg{
		Name:        "test",
		MainCategID: mainCateg.ID,
	}

	for _, o := range overwrites {
		merge(categ, o)
	}

	stmt := `INSERT INTO sub_categories (name, user_id, main_category_id) VALUES (?, ?, ?)`

	res, err := f.db.Exec(stmt, categ.Name, user.ID, mainCateg.ID)
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
