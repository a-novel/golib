// Code generated by mockery v2.46.3. DO NOT EDIT.

package adaptersmocks

import mock "github.com/stretchr/testify/mock"

// MockGRPC is an autogenerated mock type for the GRPC type
type MockGRPC struct {
	mock.Mock
}

type MockGRPC_Expecter struct {
	mock *mock.Mock
}

func (_m *MockGRPC) EXPECT() *MockGRPC_Expecter {
	return &MockGRPC_Expecter{mock: &_m.Mock}
}

// Report provides a mock function with given fields: service, err
func (_m *MockGRPC) Report(service string, err error) {
	_m.Called(service, err)
}

// MockGRPC_Report_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Report'
type MockGRPC_Report_Call struct {
	*mock.Call
}

// Report is a helper method to define mock.On call
//   - service string
//   - err error
func (_e *MockGRPC_Expecter) Report(service interface{}, err interface{}) *MockGRPC_Report_Call {
	return &MockGRPC_Report_Call{Call: _e.mock.On("Report", service, err)}
}

func (_c *MockGRPC_Report_Call) Run(run func(service string, err error)) *MockGRPC_Report_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(error))
	})
	return _c
}

func (_c *MockGRPC_Report_Call) Return() *MockGRPC_Report_Call {
	_c.Call.Return()
	return _c
}

func (_c *MockGRPC_Report_Call) RunAndReturn(run func(string, error)) *MockGRPC_Report_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockGRPC creates a new instance of MockGRPC. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockGRPC(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockGRPC {
	mock := &MockGRPC{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
