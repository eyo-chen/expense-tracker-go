package domain

import "errors"

var (
	ErrDataAlreadyExists    = errors.New("data already exists")
	ErrDataNotFound         = errors.New("data not found")
	ErrAuthentication       = errors.New("authentication failed")
	ErrAuthToken            = errors.New("invalid auth token")
	ErrServer               = errors.New("internal server error")
	ErrInvalidMainCategType = errors.New("invalid main category type")

	// main category
	ErrMainCategNotFound  = errors.New("main category not found")
	ErrUniqueIconUser     = errors.New("icon already used by another main category")
	ErrUniqueNameUserType = errors.New("name already used by another main category with the same type")

	// icon
	ErrIconNotFound = errors.New("icon not found")
)
