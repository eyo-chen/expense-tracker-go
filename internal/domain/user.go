package domain

// User contains user information
type User struct {
	ID                int64
	Name              string
	Email             string
	IsSetInitCategory bool
	Password          string
	Password_hash     string
}

// UpdateUserOpt contains option to update user
type UpdateUserOpt struct {
	IsSetInitCategory *bool
}

// Token contains access token and refresh token
type Token struct {
	Access  string
	Refresh string
}
