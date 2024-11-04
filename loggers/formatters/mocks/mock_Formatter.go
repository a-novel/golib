// Code generated by mockery v2.46.3. DO NOT EDIT.

package formattersmocks

import (
	loggers "github.com/a-novel/golib/loggers"
	formatters "github.com/a-novel/golib/loggers/formatters"

	mock "github.com/stretchr/testify/mock"
)

// MockFormatter is an autogenerated mock type for the Formatter type
type MockFormatter struct {
	mock.Mock
}

type MockFormatter_Expecter struct {
	mock *mock.Mock
}

func (_m *MockFormatter) EXPECT() *MockFormatter_Expecter {
	return &MockFormatter_Expecter{mock: &_m.Mock}
}

// Log provides a mock function with given fields: content, level
func (_m *MockFormatter) Log(content formatters.LogContent, level loggers.LogLevel) {
	_m.Called(content, level)
}

// MockFormatter_Log_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Log'
type MockFormatter_Log_Call struct {
	*mock.Call
}

// Log is a helper method to define mock.On call
//   - content formatters.LogContent
//   - level loggers.LogLevel
func (_e *MockFormatter_Expecter) Log(content interface{}, level interface{}) *MockFormatter_Log_Call {
	return &MockFormatter_Log_Call{Call: _e.mock.On("Log", content, level)}
}

func (_c *MockFormatter_Log_Call) Run(run func(content formatters.LogContent, level loggers.LogLevel)) *MockFormatter_Log_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(formatters.LogContent), args[1].(loggers.LogLevel))
	})
	return _c
}

func (_c *MockFormatter_Log_Call) Return() *MockFormatter_Log_Call {
	_c.Call.Return()
	return _c
}

func (_c *MockFormatter_Log_Call) RunAndReturn(run func(formatters.LogContent, loggers.LogLevel)) *MockFormatter_Log_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockFormatter creates a new instance of MockFormatter. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockFormatter(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockFormatter {
	mock := &MockFormatter{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
