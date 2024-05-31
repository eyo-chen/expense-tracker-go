package domain

import "errors"

var (
	// email not found error
	ErrEmailNotFound = errors.New("email not found")

	// user id not found error
	ErrUserIDNotFound = errors.New("user id not found")

	// email already exists error
	ErrEmailAlreadyExists = errors.New("email already exists")

	// data already exists error
	ErrDataAlreadyExists = errors.New("data already exists")

	// data not found error
	ErrDataNotFound = errors.New("data not found")

	// authentication error
	ErrAuthentication = errors.New("authentication failed")

	// authorization error
	ErrAuthToken = errors.New("invalid auth token")

	// internal server error
	ErrServer = errors.New("internal server error")

	// main category not found error
	ErrMainCategNotFound = errors.New("main category not found")

	// main category unique icon error
	ErrUniqueIconUser = errors.New("icon already used by another main category")

	// main category unique name error
	ErrUniqueNameUserType = errors.New("name already used by another main category with the same type")

	// sub category not found error
	ErrSubCategNotFound = errors.New("sub category not found")

	// sub category unique name error
	ErrUniqueNameUserMainCateg = errors.New("name already used by another sub category with the same main category")

	// icon not found error
	ErrIconNotFound = errors.New("icon not found")

	// main category in sub category is not consistent with the main category
	ErrMainCategNotConsistent = errors.New("main category in sub category is not consistent with the main category")

	// type in main category is not consistent with the tranmsaction type
	ErrTypeNotConsistent = errors.New("type in main category is not consistent with the tranmsaction type")

	// transaction data not found error
	ErrTransactionDataNotFound = errors.New("transaction data not found")

	// sort by type not valid error
	ErrSortByTypeNotValid = errors.New("sort by type not valid")

	// sort direction type not valid error
	ErrSortDirTypeNotValid = errors.New("sort direction type not valid")
)
