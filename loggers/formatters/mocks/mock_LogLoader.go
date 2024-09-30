// Code generated by mockery v2.46.0. DO NOT EDIT.

package formattersmocks

import (
	formatters "github.com/a-novel/golib/loggers/formatters"
	mock "github.com/stretchr/testify/mock"
)

// MockLogLoader is an autogenerated mock type for the LogLoader type
type MockLogLoader struct {
	mock.Mock
}

type MockLogLoader_Expecter struct {
	mock *mock.Mock
}

func (_m *MockLogLoader) EXPECT() *MockLogLoader_Expecter {
	return &MockLogLoader_Expecter{mock: &_m.Mock}
}

// RenderConsole provides a mock function with given fields:
func (_m *MockLogLoader) RenderConsole() string {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for RenderConsole")
	}

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// MockLogLoader_RenderConsole_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'RenderConsole'
type MockLogLoader_RenderConsole_Call struct {
	*mock.Call
}

// RenderConsole is a helper method to define mock.On call
func (_e *MockLogLoader_Expecter) RenderConsole() *MockLogLoader_RenderConsole_Call {
	return &MockLogLoader_RenderConsole_Call{Call: _e.mock.On("RenderConsole")}
}

func (_c *MockLogLoader_RenderConsole_Call) Run(run func()) *MockLogLoader_RenderConsole_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockLogLoader_RenderConsole_Call) Return(_a0 string) *MockLogLoader_RenderConsole_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockLogLoader_RenderConsole_Call) RunAndReturn(run func() string) *MockLogLoader_RenderConsole_Call {
	_c.Call.Return(run)
	return _c
}

// RenderConsoleDynamic provides a mock function with given fields: renderer
func (_m *MockLogLoader) RenderConsoleDynamic(renderer func(string)) formatters.LogDynamicContent {
	ret := _m.Called(renderer)

	if len(ret) == 0 {
		panic("no return value specified for RenderConsoleDynamic")
	}

	var r0 formatters.LogDynamicContent
	if rf, ok := ret.Get(0).(func(func(string)) formatters.LogDynamicContent); ok {
		r0 = rf(renderer)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(formatters.LogDynamicContent)
		}
	}

	return r0
}

// MockLogLoader_RenderConsoleDynamic_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'RenderConsoleDynamic'
type MockLogLoader_RenderConsoleDynamic_Call struct {
	*mock.Call
}

// RenderConsoleDynamic is a helper method to define mock.On call
//   - renderer func(string)
func (_e *MockLogLoader_Expecter) RenderConsoleDynamic(renderer interface{}) *MockLogLoader_RenderConsoleDynamic_Call {
	return &MockLogLoader_RenderConsoleDynamic_Call{Call: _e.mock.On("RenderConsoleDynamic", renderer)}
}

func (_c *MockLogLoader_RenderConsoleDynamic_Call) Run(run func(renderer func(string))) *MockLogLoader_RenderConsoleDynamic_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(func(string)))
	})
	return _c
}

func (_c *MockLogLoader_RenderConsoleDynamic_Call) Return(_a0 formatters.LogDynamicContent) *MockLogLoader_RenderConsoleDynamic_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockLogLoader_RenderConsoleDynamic_Call) RunAndReturn(run func(func(string)) formatters.LogDynamicContent) *MockLogLoader_RenderConsoleDynamic_Call {
	_c.Call.Return(run)
	return _c
}

// RenderJSON provides a mock function with given fields:
func (_m *MockLogLoader) RenderJSON() interface{} {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for RenderJSON")
	}

	var r0 interface{}
	if rf, ok := ret.Get(0).(func() interface{}); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interface{})
		}
	}

	return r0
}

// MockLogLoader_RenderJSON_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'RenderJSON'
type MockLogLoader_RenderJSON_Call struct {
	*mock.Call
}

// RenderJSON is a helper method to define mock.On call
func (_e *MockLogLoader_Expecter) RenderJSON() *MockLogLoader_RenderJSON_Call {
	return &MockLogLoader_RenderJSON_Call{Call: _e.mock.On("RenderJSON")}
}

func (_c *MockLogLoader_RenderJSON_Call) Run(run func()) *MockLogLoader_RenderJSON_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockLogLoader_RenderJSON_Call) Return(_a0 interface{}) *MockLogLoader_RenderJSON_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockLogLoader_RenderJSON_Call) RunAndReturn(run func() interface{}) *MockLogLoader_RenderJSON_Call {
	_c.Call.Return(run)
	return _c
}

// SetChild provides a mock function with given fields: child
func (_m *MockLogLoader) SetChild(child formatters.LogContent) formatters.LogLoader {
	ret := _m.Called(child)

	if len(ret) == 0 {
		panic("no return value specified for SetChild")
	}

	var r0 formatters.LogLoader
	if rf, ok := ret.Get(0).(func(formatters.LogContent) formatters.LogLoader); ok {
		r0 = rf(child)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(formatters.LogLoader)
		}
	}

	return r0
}

// MockLogLoader_SetChild_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SetChild'
type MockLogLoader_SetChild_Call struct {
	*mock.Call
}

// SetChild is a helper method to define mock.On call
//   - child formatters.LogContent
func (_e *MockLogLoader_Expecter) SetChild(child interface{}) *MockLogLoader_SetChild_Call {
	return &MockLogLoader_SetChild_Call{Call: _e.mock.On("SetChild", child)}
}

func (_c *MockLogLoader_SetChild_Call) Run(run func(child formatters.LogContent)) *MockLogLoader_SetChild_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(formatters.LogContent))
	})
	return _c
}

func (_c *MockLogLoader_SetChild_Call) Return(_a0 formatters.LogLoader) *MockLogLoader_SetChild_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockLogLoader_SetChild_Call) RunAndReturn(run func(formatters.LogContent) formatters.LogLoader) *MockLogLoader_SetChild_Call {
	_c.Call.Return(run)
	return _c
}

// SetCompleted provides a mock function with given fields:
func (_m *MockLogLoader) SetCompleted() formatters.LogLoader {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for SetCompleted")
	}

	var r0 formatters.LogLoader
	if rf, ok := ret.Get(0).(func() formatters.LogLoader); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(formatters.LogLoader)
		}
	}

	return r0
}

// MockLogLoader_SetCompleted_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SetCompleted'
type MockLogLoader_SetCompleted_Call struct {
	*mock.Call
}

// SetCompleted is a helper method to define mock.On call
func (_e *MockLogLoader_Expecter) SetCompleted() *MockLogLoader_SetCompleted_Call {
	return &MockLogLoader_SetCompleted_Call{Call: _e.mock.On("SetCompleted")}
}

func (_c *MockLogLoader_SetCompleted_Call) Run(run func()) *MockLogLoader_SetCompleted_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockLogLoader_SetCompleted_Call) Return(_a0 formatters.LogLoader) *MockLogLoader_SetCompleted_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockLogLoader_SetCompleted_Call) RunAndReturn(run func() formatters.LogLoader) *MockLogLoader_SetCompleted_Call {
	_c.Call.Return(run)
	return _c
}

// SetDescription provides a mock function with given fields: description
func (_m *MockLogLoader) SetDescription(description string) formatters.LogLoader {
	ret := _m.Called(description)

	if len(ret) == 0 {
		panic("no return value specified for SetDescription")
	}

	var r0 formatters.LogLoader
	if rf, ok := ret.Get(0).(func(string) formatters.LogLoader); ok {
		r0 = rf(description)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(formatters.LogLoader)
		}
	}

	return r0
}

// MockLogLoader_SetDescription_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SetDescription'
type MockLogLoader_SetDescription_Call struct {
	*mock.Call
}

// SetDescription is a helper method to define mock.On call
//   - description string
func (_e *MockLogLoader_Expecter) SetDescription(description interface{}) *MockLogLoader_SetDescription_Call {
	return &MockLogLoader_SetDescription_Call{Call: _e.mock.On("SetDescription", description)}
}

func (_c *MockLogLoader_SetDescription_Call) Run(run func(description string)) *MockLogLoader_SetDescription_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *MockLogLoader_SetDescription_Call) Return(_a0 formatters.LogLoader) *MockLogLoader_SetDescription_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockLogLoader_SetDescription_Call) RunAndReturn(run func(string) formatters.LogLoader) *MockLogLoader_SetDescription_Call {
	_c.Call.Return(run)
	return _c
}

// SetError provides a mock function with given fields:
func (_m *MockLogLoader) SetError() formatters.LogLoader {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for SetError")
	}

	var r0 formatters.LogLoader
	if rf, ok := ret.Get(0).(func() formatters.LogLoader); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(formatters.LogLoader)
		}
	}

	return r0
}

// MockLogLoader_SetError_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SetError'
type MockLogLoader_SetError_Call struct {
	*mock.Call
}

// SetError is a helper method to define mock.On call
func (_e *MockLogLoader_Expecter) SetError() *MockLogLoader_SetError_Call {
	return &MockLogLoader_SetError_Call{Call: _e.mock.On("SetError")}
}

func (_c *MockLogLoader_SetError_Call) Run(run func()) *MockLogLoader_SetError_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockLogLoader_SetError_Call) Return(_a0 formatters.LogLoader) *MockLogLoader_SetError_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockLogLoader_SetError_Call) RunAndReturn(run func() formatters.LogLoader) *MockLogLoader_SetError_Call {
	_c.Call.Return(run)
	return _c
}

// StopRunning provides a mock function with given fields:
func (_m *MockLogLoader) StopRunning() formatters.LogDynamicContent {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for StopRunning")
	}

	var r0 formatters.LogDynamicContent
	if rf, ok := ret.Get(0).(func() formatters.LogDynamicContent); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(formatters.LogDynamicContent)
		}
	}

	return r0
}

// MockLogLoader_StopRunning_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'StopRunning'
type MockLogLoader_StopRunning_Call struct {
	*mock.Call
}

// StopRunning is a helper method to define mock.On call
func (_e *MockLogLoader_Expecter) StopRunning() *MockLogLoader_StopRunning_Call {
	return &MockLogLoader_StopRunning_Call{Call: _e.mock.On("StopRunning")}
}

func (_c *MockLogLoader_StopRunning_Call) Run(run func()) *MockLogLoader_StopRunning_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockLogLoader_StopRunning_Call) Return(_a0 formatters.LogDynamicContent) *MockLogLoader_StopRunning_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockLogLoader_StopRunning_Call) RunAndReturn(run func() formatters.LogDynamicContent) *MockLogLoader_StopRunning_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockLogLoader creates a new instance of MockLogLoader. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockLogLoader(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockLogLoader {
	mock := &MockLogLoader{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
