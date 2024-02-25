package icon

import (
	"database/sql"
	"fmt"
)

func Blueprint(i int, last Icon) Icon {
	return Icon{
		URL: "test" + fmt.Sprint(i),
	}
}

func Inserter(db *sql.DB, i Icon) (Icon, error) {
	stmt := `INSERT INTO icons (url) VALUES (?)`
	res, err := db.Exec(stmt, i.URL)
	if err != nil {
		return Icon{}, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return Icon{}, err
	}

	i.ID = id

	return i, nil
}
