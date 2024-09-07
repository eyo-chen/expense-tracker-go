// Code generated by mockery v2.40.1. DO NOT EDIT.

package mocks

import (
	context "context"

	domain "github.com/eyo-chen/expense-tracker-go/internal/domain"

	mock "github.com/stretchr/testify/mock"
)

// UserRepo is an autogenerated mock type for the UserRepo type
type UserRepo struct {
	mock.Mock
}

// Create provides a mock function with given fields: name, email, passwordHash
func (_m *UserRepo) Create(name string, email string, passwordHash string) error {
	ret := _m.Called(name, email, passwordHash)

	if len(ret) == 0 {
		panic("no return value specified for Create")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string, string) error); ok {
		r0 = rf(name, email, passwordHash)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// FindByEmail provides a mock function with given fields: email
func (_m *UserRepo) FindByEmail(email string) (domain.User, error) {
	ret := _m.Called(email)

	if len(ret) == 0 {
		panic("no return value specified for FindByEmail")
	}

	var r0 domain.User
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (domain.User, error)); ok {
		return rf(email)
	}
	if rf, ok := ret.Get(0).(func(string) domain.User); ok {
		r0 = rf(email)
	} else {
		r0 = ret.Get(0).(domain.User)
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(email)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetInfo provides a mock function with given fields: userID
func (_m *UserRepo) GetInfo(userID int64) (domain.User, error) {
	ret := _m.Called(userID)

	if len(ret) == 0 {
		panic("no return value specified for GetInfo")
	}

	var r0 domain.User
	var r1 error
	if rf, ok := ret.Get(0).(func(int64) (domain.User, error)); ok {
		return rf(userID)
	}
	if rf, ok := ret.Get(0).(func(int64) domain.User); ok {
		r0 = rf(userID)
	} else {
		r0 = ret.Get(0).(domain.User)
	}

	if rf, ok := ret.Get(1).(func(int64) error); ok {
		r1 = rf(userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Update provides a mock function with given fields: ctx, userID, opt
func (_m *UserRepo) Update(ctx context.Context, userID int64, opt domain.UpdateUserOpt) error {
	ret := _m.Called(ctx, userID, opt)

	if len(ret) == 0 {
		panic("no return value specified for Update")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int64, domain.UpdateUserOpt) error); ok {
		r0 = rf(ctx, userID, opt)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewUserRepo creates a new instance of UserRepo. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewUserRepo(t interface {
	mock.TestingT
	Cleanup(func())
}) *UserRepo {
	mock := &UserRepo{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}