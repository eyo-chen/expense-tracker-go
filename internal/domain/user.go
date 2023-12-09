package domain

// User contains user information
type User struct {
	ID            int64
	Name          string
	Email         string
	CountryID     int
	Password      string
	Password_hash string
}
