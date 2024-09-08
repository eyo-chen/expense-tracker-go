// Code generated by mockery v2.40.1. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// S3Service is an autogenerated mock type for the S3Service type
type S3Service struct {
	mock.Mock
}

// DeleteObject provides a mock function with given fields: ctx, objectKey
func (_m *S3Service) DeleteObject(ctx context.Context, objectKey string) error {
	ret := _m.Called(ctx, objectKey)

	if len(ret) == 0 {
		panic("no return value specified for DeleteObject")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, objectKey)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetObjectUrl provides a mock function with given fields: ctx, objectKey, lifetimeSecs
func (_m *S3Service) GetObjectUrl(ctx context.Context, objectKey string, lifetimeSecs int64) (string, error) {
	ret := _m.Called(ctx, objectKey, lifetimeSecs)

	if len(ret) == 0 {
		panic("no return value specified for GetObjectUrl")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, int64) (string, error)); ok {
		return rf(ctx, objectKey, lifetimeSecs)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, int64) string); ok {
		r0 = rf(ctx, objectKey, lifetimeSecs)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, int64) error); ok {
		r1 = rf(ctx, objectKey, lifetimeSecs)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PutObjectUrl provides a mock function with given fields: ctx, objectKey, lifetimeSecs
func (_m *S3Service) PutObjectUrl(ctx context.Context, objectKey string, lifetimeSecs int64) (string, error) {
	ret := _m.Called(ctx, objectKey, lifetimeSecs)

	if len(ret) == 0 {
		panic("no return value specified for PutObjectUrl")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, int64) (string, error)); ok {
		return rf(ctx, objectKey, lifetimeSecs)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, int64) string); ok {
		r0 = rf(ctx, objectKey, lifetimeSecs)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, int64) error); ok {
		r1 = rf(ctx, objectKey, lifetimeSecs)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewS3Service creates a new instance of S3Service. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewS3Service(t interface {
	mock.TestingT
	Cleanup(func())
}) *S3Service {
	mock := &S3Service{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}