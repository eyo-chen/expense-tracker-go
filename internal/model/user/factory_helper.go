package user

import (
	"database/sql"
	"fmt"
)

func Blueprint(i int, last User) User {
	return User{
		Name:          "test" + fmt.Sprint(i),
		Email:         "test" + fmt.Sprint(i) + "@gmail.com",
		Password_hash: "test" + fmt.Sprint(i),
	}
}

func Inserter(db *sql.DB, u User) (User, error) {
	stmt := `INSERT INTO users (name, email, password_hash) VALUES (?, ?, ?)`
	res, err := db.Exec(stmt, u.Name, u.Email, u.Password_hash)
	if err != nil {
		return User{}, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return User{}, err
	}

	u.ID = id

	return u, nil
}
