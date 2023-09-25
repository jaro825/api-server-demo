// Code generated by mockery. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// MockCache is an autogenerated mock type for the Cache type
type MockCache struct {
	mock.Mock
}

type MockCache_Expecter struct {
	mock *mock.Mock
}

func (_m *MockCache) EXPECT() *MockCache_Expecter {
	return &MockCache_Expecter{mock: &_m.Mock}
}

// Delete provides a mock function with given fields: ctx, key
func (_m *MockCache) Delete(ctx context.Context, key string) error {
	ret := _m.Called(ctx, key)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, key)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockCache_Delete_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Delete'
type MockCache_Delete_Call struct {
	*mock.Call
}

// Delete is a helper method to define mock.On call
//   - ctx context.Context
//   - key string
func (_e *MockCache_Expecter) Delete(ctx interface{}, key interface{}) *MockCache_Delete_Call {
	return &MockCache_Delete_Call{Call: _e.mock.On("Delete", ctx, key)}
}

func (_c *MockCache_Delete_Call) Run(run func(ctx context.Context, key string)) *MockCache_Delete_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *MockCache_Delete_Call) Return(_a0 error) *MockCache_Delete_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockCache_Delete_Call) RunAndReturn(run func(context.Context, string) error) *MockCache_Delete_Call {
	_c.Call.Return(run)
	return _c
}

// Get provides a mock function with given fields: ctx, key, into
func (_m *MockCache) Get(ctx context.Context, key string, into interface{}) error {
	ret := _m.Called(ctx, key, into)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, interface{}) error); ok {
		r0 = rf(ctx, key, into)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockCache_Get_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Get'
type MockCache_Get_Call struct {
	*mock.Call
}

// Get is a helper method to define mock.On call
//   - ctx context.Context
//   - key string
//   - into interface{}
func (_e *MockCache_Expecter) Get(ctx interface{}, key interface{}, into interface{}) *MockCache_Get_Call {
	return &MockCache_Get_Call{Call: _e.mock.On("Get", ctx, key, into)}
}

func (_c *MockCache_Get_Call) Run(run func(ctx context.Context, key string, into interface{})) *MockCache_Get_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(interface{}))
	})
	return _c
}

func (_c *MockCache_Get_Call) Return(_a0 error) *MockCache_Get_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockCache_Get_Call) RunAndReturn(run func(context.Context, string, interface{}) error) *MockCache_Get_Call {
	_c.Call.Return(run)
	return _c
}

// Set provides a mock function with given fields: ctx, key, value
func (_m *MockCache) Set(ctx context.Context, key string, value interface{}) error {
	ret := _m.Called(ctx, key, value)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, interface{}) error); ok {
		r0 = rf(ctx, key, value)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockCache_Set_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Set'
type MockCache_Set_Call struct {
	*mock.Call
}

// Set is a helper method to define mock.On call
//   - ctx context.Context
//   - key string
//   - value interface{}
func (_e *MockCache_Expecter) Set(ctx interface{}, key interface{}, value interface{}) *MockCache_Set_Call {
	return &MockCache_Set_Call{Call: _e.mock.On("Set", ctx, key, value)}
}

func (_c *MockCache_Set_Call) Run(run func(ctx context.Context, key string, value interface{})) *MockCache_Set_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(interface{}))
	})
	return _c
}

func (_c *MockCache_Set_Call) Return(_a0 error) *MockCache_Set_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockCache_Set_Call) RunAndReturn(run func(context.Context, string, interface{}) error) *MockCache_Set_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockCache creates a new instance of MockCache. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockCache(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockCache {
	mock := &MockCache{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
