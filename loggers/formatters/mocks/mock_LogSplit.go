// Code generated by mockery v2.46.3. DO NOT EDIT.

package formattersmocks

import (
	formatters "github.com/a-novel/golib/loggers/formatters"
	mock "github.com/stretchr/testify/mock"
)

// MockLogSplit is an autogenerated mock type for the LogSplit type
type MockLogSplit struct {
	mock.Mock
}

type MockLogSplit_Expecter struct {
	mock *mock.Mock
}

func (_m *MockLogSplit) EXPECT() *MockLogSplit_Expecter {
	return &MockLogSplit_Expecter{mock: &_m.Mock}
}

// RenderConsole provides a mock function with given fields:
func (_m *MockLogSplit) RenderConsole() string {
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

// MockLogSplit_RenderConsole_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'RenderConsole'
type MockLogSplit_RenderConsole_Call struct {
	*mock.Call
}

// RenderConsole is a helper method to define mock.On call
func (_e *MockLogSplit_Expecter) RenderConsole() *MockLogSplit_RenderConsole_Call {
	return &MockLogSplit_RenderConsole_Call{Call: _e.mock.On("RenderConsole")}
}

func (_c *MockLogSplit_RenderConsole_Call) Run(run func()) *MockLogSplit_RenderConsole_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockLogSplit_RenderConsole_Call) Return(_a0 string) *MockLogSplit_RenderConsole_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockLogSplit_RenderConsole_Call) RunAndReturn(run func() string) *MockLogSplit_RenderConsole_Call {
	_c.Call.Return(run)
	return _c
}

// RenderJSON provides a mock function with given fields:
func (_m *MockLogSplit) RenderJSON() interface{} {
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

// MockLogSplit_RenderJSON_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'RenderJSON'
type MockLogSplit_RenderJSON_Call struct {
	*mock.Call
}

// RenderJSON is a helper method to define mock.On call
func (_e *MockLogSplit_Expecter) RenderJSON() *MockLogSplit_RenderJSON_Call {
	return &MockLogSplit_RenderJSON_Call{Call: _e.mock.On("RenderJSON")}
}

func (_c *MockLogSplit_RenderJSON_Call) Run(run func()) *MockLogSplit_RenderJSON_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockLogSplit_RenderJSON_Call) Return(_a0 interface{}) *MockLogSplit_RenderJSON_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockLogSplit_RenderJSON_Call) RunAndReturn(run func() interface{}) *MockLogSplit_RenderJSON_Call {
	_c.Call.Return(run)
	return _c
}

// SetConsoleContent provides a mock function with given fields: content
func (_m *MockLogSplit) SetConsoleContent(content formatters.LogContent) formatters.LogSplit {
	ret := _m.Called(content)

	if len(ret) == 0 {
		panic("no return value specified for SetConsoleContent")
	}

	var r0 formatters.LogSplit
	if rf, ok := ret.Get(0).(func(formatters.LogContent) formatters.LogSplit); ok {
		r0 = rf(content)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(formatters.LogSplit)
		}
	}

	return r0
}

// MockLogSplit_SetConsoleContent_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SetConsoleContent'
type MockLogSplit_SetConsoleContent_Call struct {
	*mock.Call
}

// SetConsoleContent is a helper method to define mock.On call
//   - content formatters.LogContent
func (_e *MockLogSplit_Expecter) SetConsoleContent(content interface{}) *MockLogSplit_SetConsoleContent_Call {
	return &MockLogSplit_SetConsoleContent_Call{Call: _e.mock.On("SetConsoleContent", content)}
}

func (_c *MockLogSplit_SetConsoleContent_Call) Run(run func(content formatters.LogContent)) *MockLogSplit_SetConsoleContent_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(formatters.LogContent))
	})
	return _c
}

func (_c *MockLogSplit_SetConsoleContent_Call) Return(_a0 formatters.LogSplit) *MockLogSplit_SetConsoleContent_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockLogSplit_SetConsoleContent_Call) RunAndReturn(run func(formatters.LogContent) formatters.LogSplit) *MockLogSplit_SetConsoleContent_Call {
	_c.Call.Return(run)
	return _c
}

// SetConsoleMessage provides a mock function with given fields: content
func (_m *MockLogSplit) SetConsoleMessage(content string) formatters.LogSplit {
	ret := _m.Called(content)

	if len(ret) == 0 {
		panic("no return value specified for SetConsoleMessage")
	}

	var r0 formatters.LogSplit
	if rf, ok := ret.Get(0).(func(string) formatters.LogSplit); ok {
		r0 = rf(content)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(formatters.LogSplit)
		}
	}

	return r0
}

// MockLogSplit_SetConsoleMessage_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SetConsoleMessage'
type MockLogSplit_SetConsoleMessage_Call struct {
	*mock.Call
}

// SetConsoleMessage is a helper method to define mock.On call
//   - content string
func (_e *MockLogSplit_Expecter) SetConsoleMessage(content interface{}) *MockLogSplit_SetConsoleMessage_Call {
	return &MockLogSplit_SetConsoleMessage_Call{Call: _e.mock.On("SetConsoleMessage", content)}
}

func (_c *MockLogSplit_SetConsoleMessage_Call) Run(run func(content string)) *MockLogSplit_SetConsoleMessage_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *MockLogSplit_SetConsoleMessage_Call) Return(_a0 formatters.LogSplit) *MockLogSplit_SetConsoleMessage_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockLogSplit_SetConsoleMessage_Call) RunAndReturn(run func(string) formatters.LogSplit) *MockLogSplit_SetConsoleMessage_Call {
	_c.Call.Return(run)
	return _c
}

// SetConsoleRenderer provides a mock function with given fields: content
func (_m *MockLogSplit) SetConsoleRenderer(content func() string) formatters.LogSplit {
	ret := _m.Called(content)

	if len(ret) == 0 {
		panic("no return value specified for SetConsoleRenderer")
	}

	var r0 formatters.LogSplit
	if rf, ok := ret.Get(0).(func(func() string) formatters.LogSplit); ok {
		r0 = rf(content)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(formatters.LogSplit)
		}
	}

	return r0
}

// MockLogSplit_SetConsoleRenderer_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SetConsoleRenderer'
type MockLogSplit_SetConsoleRenderer_Call struct {
	*mock.Call
}

// SetConsoleRenderer is a helper method to define mock.On call
//   - content func() string
func (_e *MockLogSplit_Expecter) SetConsoleRenderer(content interface{}) *MockLogSplit_SetConsoleRenderer_Call {
	return &MockLogSplit_SetConsoleRenderer_Call{Call: _e.mock.On("SetConsoleRenderer", content)}
}

func (_c *MockLogSplit_SetConsoleRenderer_Call) Run(run func(content func() string)) *MockLogSplit_SetConsoleRenderer_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(func() string))
	})
	return _c
}

func (_c *MockLogSplit_SetConsoleRenderer_Call) Return(_a0 formatters.LogSplit) *MockLogSplit_SetConsoleRenderer_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockLogSplit_SetConsoleRenderer_Call) RunAndReturn(run func(func() string) formatters.LogSplit) *MockLogSplit_SetConsoleRenderer_Call {
	_c.Call.Return(run)
	return _c
}

// SetJSONContent provides a mock function with given fields: content
func (_m *MockLogSplit) SetJSONContent(content formatters.LogContent) formatters.LogSplit {
	ret := _m.Called(content)

	if len(ret) == 0 {
		panic("no return value specified for SetJSONContent")
	}

	var r0 formatters.LogSplit
	if rf, ok := ret.Get(0).(func(formatters.LogContent) formatters.LogSplit); ok {
		r0 = rf(content)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(formatters.LogSplit)
		}
	}

	return r0
}

// MockLogSplit_SetJSONContent_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SetJSONContent'
type MockLogSplit_SetJSONContent_Call struct {
	*mock.Call
}

// SetJSONContent is a helper method to define mock.On call
//   - content formatters.LogContent
func (_e *MockLogSplit_Expecter) SetJSONContent(content interface{}) *MockLogSplit_SetJSONContent_Call {
	return &MockLogSplit_SetJSONContent_Call{Call: _e.mock.On("SetJSONContent", content)}
}

func (_c *MockLogSplit_SetJSONContent_Call) Run(run func(content formatters.LogContent)) *MockLogSplit_SetJSONContent_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(formatters.LogContent))
	})
	return _c
}

func (_c *MockLogSplit_SetJSONContent_Call) Return(_a0 formatters.LogSplit) *MockLogSplit_SetJSONContent_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockLogSplit_SetJSONContent_Call) RunAndReturn(run func(formatters.LogContent) formatters.LogSplit) *MockLogSplit_SetJSONContent_Call {
	_c.Call.Return(run)
	return _c
}

// SetJSONMessage provides a mock function with given fields: content
func (_m *MockLogSplit) SetJSONMessage(content interface{}) formatters.LogSplit {
	ret := _m.Called(content)

	if len(ret) == 0 {
		panic("no return value specified for SetJSONMessage")
	}

	var r0 formatters.LogSplit
	if rf, ok := ret.Get(0).(func(interface{}) formatters.LogSplit); ok {
		r0 = rf(content)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(formatters.LogSplit)
		}
	}

	return r0
}

// MockLogSplit_SetJSONMessage_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SetJSONMessage'
type MockLogSplit_SetJSONMessage_Call struct {
	*mock.Call
}

// SetJSONMessage is a helper method to define mock.On call
//   - content interface{}
func (_e *MockLogSplit_Expecter) SetJSONMessage(content interface{}) *MockLogSplit_SetJSONMessage_Call {
	return &MockLogSplit_SetJSONMessage_Call{Call: _e.mock.On("SetJSONMessage", content)}
}

func (_c *MockLogSplit_SetJSONMessage_Call) Run(run func(content interface{})) *MockLogSplit_SetJSONMessage_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(interface{}))
	})
	return _c
}

func (_c *MockLogSplit_SetJSONMessage_Call) Return(_a0 formatters.LogSplit) *MockLogSplit_SetJSONMessage_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockLogSplit_SetJSONMessage_Call) RunAndReturn(run func(interface{}) formatters.LogSplit) *MockLogSplit_SetJSONMessage_Call {
	_c.Call.Return(run)
	return _c
}

// SetJSONRenderer provides a mock function with given fields: content
func (_m *MockLogSplit) SetJSONRenderer(content func() interface{}) formatters.LogSplit {
	ret := _m.Called(content)

	if len(ret) == 0 {
		panic("no return value specified for SetJSONRenderer")
	}

	var r0 formatters.LogSplit
	if rf, ok := ret.Get(0).(func(func() interface{}) formatters.LogSplit); ok {
		r0 = rf(content)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(formatters.LogSplit)
		}
	}

	return r0
}

// MockLogSplit_SetJSONRenderer_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SetJSONRenderer'
type MockLogSplit_SetJSONRenderer_Call struct {
	*mock.Call
}

// SetJSONRenderer is a helper method to define mock.On call
//   - content func() interface{}
func (_e *MockLogSplit_Expecter) SetJSONRenderer(content interface{}) *MockLogSplit_SetJSONRenderer_Call {
	return &MockLogSplit_SetJSONRenderer_Call{Call: _e.mock.On("SetJSONRenderer", content)}
}

func (_c *MockLogSplit_SetJSONRenderer_Call) Run(run func(content func() interface{})) *MockLogSplit_SetJSONRenderer_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(func() interface{}))
	})
	return _c
}

func (_c *MockLogSplit_SetJSONRenderer_Call) Return(_a0 formatters.LogSplit) *MockLogSplit_SetJSONRenderer_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockLogSplit_SetJSONRenderer_Call) RunAndReturn(run func(func() interface{}) formatters.LogSplit) *MockLogSplit_SetJSONRenderer_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockLogSplit creates a new instance of MockLogSplit. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockLogSplit(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockLogSplit {
	mock := &MockLogSplit{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
