package domain

import "fmt"

// DefaultIcon contains default icon information
type DefaultIcon struct {
	ID  int64  `json:"id"`
	URL string `json:"url"`
}

// UserIcon contains user icon information
type UserIcon struct {
	ID        int64  `json:"id"`
	UserID    int64  `json:"user_id"`
	ObjectKey string `json:"object_key"`
}

// Icon contains icon information
type Icon struct {
	ID   int64
	Type IconType
	URL  string
}

// GenUserIconCacheKey generates a cache key for user icon
func GenUserIconCacheKey(objectKey string) string {
	return fmt.Sprintf("user_icon-%s", objectKey)
}
