package subcateg

import (
	"database/sql"
	"fmt"
)

func Blueprint(i int, last SubCateg) SubCateg {
	return SubCateg{
		Name: "test" + fmt.Sprint(i),
	}
}

func Inserter(db *sql.DB, s SubCateg) (SubCateg, error) {
	stmt := `INSERT INTO sub_categories (name, user_id, main_category_id) VALUES (?, ?, ?)`
	res, err := db.Exec(stmt, s.Name, s.UserID, s.MainCategID)
	if err != nil {
		return SubCateg{}, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return SubCateg{}, err
	}

	s.ID = id

	return s, nil
}
