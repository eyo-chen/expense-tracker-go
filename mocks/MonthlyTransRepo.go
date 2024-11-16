// Code generated by mockery v2.40.1. DO NOT EDIT.

package mocks

import (
	context "context"

	domain "github.com/eyo-chen/expense-tracker-go/internal/domain"

	mock "github.com/stretchr/testify/mock"

	time "time"
)

// MonthlyTransRepo is an autogenerated mock type for the MonthlyTransRepo type
type MonthlyTransRepo struct {
	mock.Mock
}

// Create provides a mock function with given fields: ctx, date, trans
func (_m *MonthlyTransRepo) Create(ctx context.Context, date time.Time, trans []domain.MonthlyAggregatedData) error {
	ret := _m.Called(ctx, date, trans)

	if len(ret) == 0 {
		panic("no return value specified for Create")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, time.Time, []domain.MonthlyAggregatedData) error); ok {
		r0 = rf(ctx, date, trans)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewMonthlyTransRepo creates a new instance of MonthlyTransRepo. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMonthlyTransRepo(t interface {
	mock.TestingT
	Cleanup(func())
}) *MonthlyTransRepo {
	mock := &MonthlyTransRepo{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
