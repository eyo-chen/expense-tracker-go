package domain

// DefaultIcon contains default icon information
type DefaultIcon struct {
	ID  int64
	URL string
}

// UserIcon contains user icon information
type UserIcon struct {
	ID        int64
	UserID    int64
	ObjectKey string
}
