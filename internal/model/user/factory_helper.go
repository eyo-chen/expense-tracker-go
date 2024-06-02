package user

import "fmt"

func userBluePrint(i int, _ User) User {
	return User{
		ID:                int64(i),
		Name:              fmt.Sprintf("name%d", i),
		Email:             fmt.Sprintf("email%d@gmail.com", i),
		IsSetInitCategory: false,
		Password_hash:     fmt.Sprintf("password_hash%d", i),
	}
}
