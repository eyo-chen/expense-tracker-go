package domain

// Icon contains icon information
type Icon struct {
	ID  int64  `json:"id"`
	URL string `json:"url"`
}

// UserIcon contains user icon information
type UserIcon struct {
	ID        int64
	UserID    int64
	ObjectKey string
}
