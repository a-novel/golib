// Code generated by mockery v2.46.0. DO NOT EDIT.

package loggersmocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// MockGRPCCallback is an autogenerated mock type for the GRPCCallback type
type MockGRPCCallback[In interface{}, Out interface{}] struct {
	mock.Mock
}

type MockGRPCCallback_Expecter[In interface{}, Out interface{}] struct {
	mock *mock.Mock
}

func (_m *MockGRPCCallback[In, Out]) EXPECT() *MockGRPCCallback_Expecter[In, Out] {
	return &MockGRPCCallback_Expecter[In, Out]{mock: &_m.Mock}
}

// Execute provides a mock function with given fields: ctx, in
func (_m *MockGRPCCallback[In, Out]) Execute(ctx context.Context, in In) (Out, error) {
	ret := _m.Called(ctx, in)

	if len(ret) == 0 {
		panic("no return value specified for Execute")
	}

	var r0 Out
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, In) (Out, error)); ok {
		return rf(ctx, in)
	}
	if rf, ok := ret.Get(0).(func(context.Context, In) Out); ok {
		r0 = rf(ctx, in)
	} else {
		r0 = ret.Get(0).(Out)
	}

	if rf, ok := ret.Get(1).(func(context.Context, In) error); ok {
		r1 = rf(ctx, in)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockGRPCCallback_Execute_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Execute'
type MockGRPCCallback_Execute_Call[In interface{}, Out interface{}] struct {
	*mock.Call
}

// Execute is a helper method to define mock.On call
//   - ctx context.Context
//   - in In
func (_e *MockGRPCCallback_Expecter[In, Out]) Execute(ctx interface{}, in interface{}) *MockGRPCCallback_Execute_Call[In, Out] {
	return &MockGRPCCallback_Execute_Call[In, Out]{Call: _e.mock.On("Execute", ctx, in)}
}

func (_c *MockGRPCCallback_Execute_Call[In, Out]) Run(run func(ctx context.Context, in In)) *MockGRPCCallback_Execute_Call[In, Out] {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(In))
	})
	return _c
}

func (_c *MockGRPCCallback_Execute_Call[In, Out]) Return(_a0 Out, _a1 error) *MockGRPCCallback_Execute_Call[In, Out] {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockGRPCCallback_Execute_Call[In, Out]) RunAndReturn(run func(context.Context, In) (Out, error)) *MockGRPCCallback_Execute_Call[In, Out] {
	_c.Call.Return(run)
	return _c
}

// NewMockGRPCCallback creates a new instance of MockGRPCCallback. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockGRPCCallback[In interface{}, Out interface{}](t interface {
	mock.TestingT
	Cleanup(func())
}) *MockGRPCCallback[In, Out] {
	mock := &MockGRPCCallback[In, Out]{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
