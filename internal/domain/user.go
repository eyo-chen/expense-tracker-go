package domain

// User contains user information
type User struct {
	ID            int64
	Name          string
	Email         string
	Password      string
	Password_hash string
}
