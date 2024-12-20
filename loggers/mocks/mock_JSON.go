// Code generated by mockery v2.46.3. DO NOT EDIT.

package loggersmocks

import (
	loggers "github.com/a-novel/golib/loggers"
	mock "github.com/stretchr/testify/mock"
)

// MockJSON is an autogenerated mock type for the JSON type
type MockJSON struct {
	mock.Mock
}

type MockJSON_Expecter struct {
	mock *mock.Mock
}

func (_m *MockJSON) EXPECT() *MockJSON_Expecter {
	return &MockJSON_Expecter{mock: &_m.Mock}
}

// Log provides a mock function with given fields: level, data
func (_m *MockJSON) Log(level loggers.LogLevel, data interface{}) {
	_m.Called(level, data)
}

// MockJSON_Log_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Log'
type MockJSON_Log_Call struct {
	*mock.Call
}

// Log is a helper method to define mock.On call
//   - level loggers.LogLevel
//   - data interface{}
func (_e *MockJSON_Expecter) Log(level interface{}, data interface{}) *MockJSON_Log_Call {
	return &MockJSON_Log_Call{Call: _e.mock.On("Log", level, data)}
}

func (_c *MockJSON_Log_Call) Run(run func(level loggers.LogLevel, data interface{})) *MockJSON_Log_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(loggers.LogLevel), args[1].(interface{}))
	})
	return _c
}

func (_c *MockJSON_Log_Call) Return() *MockJSON_Log_Call {
	_c.Call.Return()
	return _c
}

func (_c *MockJSON_Log_Call) RunAndReturn(run func(loggers.LogLevel, interface{})) *MockJSON_Log_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockJSON creates a new instance of MockJSON. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockJSON(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockJSON {
	mock := &MockJSON{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
