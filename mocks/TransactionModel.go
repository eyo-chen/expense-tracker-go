// Code generated by mockery v2.40.1. DO NOT EDIT.

package mocks

import (
	context "context"

	domain "github.com/OYE0303/expense-tracker-go/internal/domain"

	mock "github.com/stretchr/testify/mock"
)

// TransactionModel is an autogenerated mock type for the TransactionModel type
type TransactionModel struct {
	mock.Mock
}

// Create provides a mock function with given fields: ctx, trans
func (_m *TransactionModel) Create(ctx context.Context, trans domain.CreateTransactionInput) error {
	ret := _m.Called(ctx, trans)

	if len(ret) == 0 {
		panic("no return value specified for Create")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, domain.CreateTransactionInput) error); ok {
		r0 = rf(ctx, trans)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Delete provides a mock function with given fields: ctx, id
func (_m *TransactionModel) Delete(ctx context.Context, id int64) error {
	ret := _m.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for Delete")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int64) error); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetAccInfo provides a mock function with given fields: ctx, query, userID
func (_m *TransactionModel) GetAccInfo(ctx context.Context, query domain.GetAccInfoQuery, userID int64) (domain.AccInfo, error) {
	ret := _m.Called(ctx, query, userID)

	if len(ret) == 0 {
		panic("no return value specified for GetAccInfo")
	}

	var r0 domain.AccInfo
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, domain.GetAccInfoQuery, int64) (domain.AccInfo, error)); ok {
		return rf(ctx, query, userID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, domain.GetAccInfoQuery, int64) domain.AccInfo); ok {
		r0 = rf(ctx, query, userID)
	} else {
		r0 = ret.Get(0).(domain.AccInfo)
	}

	if rf, ok := ret.Get(1).(func(context.Context, domain.GetAccInfoQuery, int64) error); ok {
		r1 = rf(ctx, query, userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAll provides a mock function with given fields: ctx, query, userID
func (_m *TransactionModel) GetAll(ctx context.Context, query domain.GetQuery, userID int64) ([]domain.Transaction, error) {
	ret := _m.Called(ctx, query, userID)

	if len(ret) == 0 {
		panic("no return value specified for GetAll")
	}

	var r0 []domain.Transaction
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, domain.GetQuery, int64) ([]domain.Transaction, error)); ok {
		return rf(ctx, query, userID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, domain.GetQuery, int64) []domain.Transaction); ok {
		r0 = rf(ctx, query, userID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]domain.Transaction)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, domain.GetQuery, int64) error); ok {
		r1 = rf(ctx, query, userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetByIDAndUserID provides a mock function with given fields: ctx, id, userID
func (_m *TransactionModel) GetByIDAndUserID(ctx context.Context, id int64, userID int64) (domain.Transaction, error) {
	ret := _m.Called(ctx, id, userID)

	if len(ret) == 0 {
		panic("no return value specified for GetByIDAndUserID")
	}

	var r0 domain.Transaction
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int64, int64) (domain.Transaction, error)); ok {
		return rf(ctx, id, userID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int64, int64) domain.Transaction); ok {
		r0 = rf(ctx, id, userID)
	} else {
		r0 = ret.Get(0).(domain.Transaction)
	}

	if rf, ok := ret.Get(1).(func(context.Context, int64, int64) error); ok {
		r1 = rf(ctx, id, userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetDailyBarChartData provides a mock function with given fields: ctx, dateRange, transactionType, mainCategIDs, userID
func (_m *TransactionModel) GetDailyBarChartData(ctx context.Context, dateRange domain.ChartDateRange, transactionType domain.TransactionType, mainCategIDs *[]int64, userID int64) (domain.DateToChartData, error) {
	ret := _m.Called(ctx, dateRange, transactionType, mainCategIDs, userID)

	if len(ret) == 0 {
		panic("no return value specified for GetDailyBarChartData")
	}

	var r0 domain.DateToChartData
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, domain.ChartDateRange, domain.TransactionType, *[]int64, int64) (domain.DateToChartData, error)); ok {
		return rf(ctx, dateRange, transactionType, mainCategIDs, userID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, domain.ChartDateRange, domain.TransactionType, *[]int64, int64) domain.DateToChartData); ok {
		r0 = rf(ctx, dateRange, transactionType, mainCategIDs, userID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(domain.DateToChartData)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, domain.ChartDateRange, domain.TransactionType, *[]int64, int64) error); ok {
		r1 = rf(ctx, dateRange, transactionType, mainCategIDs, userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetMonthlyBarChartData provides a mock function with given fields: ctx, dateRange, transactionType, userID
func (_m *TransactionModel) GetMonthlyBarChartData(ctx context.Context, dateRange domain.ChartDateRange, transactionType domain.TransactionType, userID int64) (domain.DateToChartData, error) {
	ret := _m.Called(ctx, dateRange, transactionType, userID)

	if len(ret) == 0 {
		panic("no return value specified for GetMonthlyBarChartData")
	}

	var r0 domain.DateToChartData
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, domain.ChartDateRange, domain.TransactionType, int64) (domain.DateToChartData, error)); ok {
		return rf(ctx, dateRange, transactionType, userID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, domain.ChartDateRange, domain.TransactionType, int64) domain.DateToChartData); ok {
		r0 = rf(ctx, dateRange, transactionType, userID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(domain.DateToChartData)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, domain.ChartDateRange, domain.TransactionType, int64) error); ok {
		r1 = rf(ctx, dateRange, transactionType, userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetMonthlyData provides a mock function with given fields: ctx, dateRange, userID
func (_m *TransactionModel) GetMonthlyData(ctx context.Context, dateRange domain.GetMonthlyDateRange, userID int64) (domain.MonthDayToTransactionType, error) {
	ret := _m.Called(ctx, dateRange, userID)

	if len(ret) == 0 {
		panic("no return value specified for GetMonthlyData")
	}

	var r0 domain.MonthDayToTransactionType
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, domain.GetMonthlyDateRange, int64) (domain.MonthDayToTransactionType, error)); ok {
		return rf(ctx, dateRange, userID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, domain.GetMonthlyDateRange, int64) domain.MonthDayToTransactionType); ok {
		r0 = rf(ctx, dateRange, userID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(domain.MonthDayToTransactionType)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, domain.GetMonthlyDateRange, int64) error); ok {
		r1 = rf(ctx, dateRange, userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetPieChartData provides a mock function with given fields: ctx, dataRange, transactionType, userID
func (_m *TransactionModel) GetPieChartData(ctx context.Context, dataRange domain.ChartDateRange, transactionType domain.TransactionType, userID int64) (domain.ChartData, error) {
	ret := _m.Called(ctx, dataRange, transactionType, userID)

	if len(ret) == 0 {
		panic("no return value specified for GetPieChartData")
	}

	var r0 domain.ChartData
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, domain.ChartDateRange, domain.TransactionType, int64) (domain.ChartData, error)); ok {
		return rf(ctx, dataRange, transactionType, userID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, domain.ChartDateRange, domain.TransactionType, int64) domain.ChartData); ok {
		r0 = rf(ctx, dataRange, transactionType, userID)
	} else {
		r0 = ret.Get(0).(domain.ChartData)
	}

	if rf, ok := ret.Get(1).(func(context.Context, domain.ChartDateRange, domain.TransactionType, int64) error); ok {
		r1 = rf(ctx, dataRange, transactionType, userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Update provides a mock function with given fields: ctx, trans
func (_m *TransactionModel) Update(ctx context.Context, trans domain.UpdateTransactionInput) error {
	ret := _m.Called(ctx, trans)

	if len(ret) == 0 {
		panic("no return value specified for Update")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, domain.UpdateTransactionInput) error); ok {
		r0 = rf(ctx, trans)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewTransactionModel creates a new instance of TransactionModel. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewTransactionModel(t interface {
	mock.TestingT
	Cleanup(func())
}) *TransactionModel {
	mock := &TransactionModel{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
