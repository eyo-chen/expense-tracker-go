package domain

import "errors"

var (
	ErrDataAlreadyExists = errors.New("data already exists")
	ErrDataNotFound      = errors.New("data not found")
	ErrAuthentication    = errors.New("authentication failed")
	ErrAuthToken         = errors.New("invalid auth token")
	ErrServer            = errors.New("internal server error")
)
