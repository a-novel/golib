package grpc

import (
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ErrorHandler interface {
	Is(target error, code codes.Code) ErrorHandler
	IsW(target error, code codes.Code, wrap error) ErrorHandler
	IsWF(target error, code codes.Code, format string, args ...interface{}) ErrorHandler

	As(target interface{}, code codes.Code) ErrorHandler
	AsW(target interface{}, code codes.Code, wrap error) ErrorHandler
	AsWF(target interface{}, code codes.Code, format string, args ...interface{}) ErrorHandler

	Test(caseFn func(initialErr error) (err error, ok bool)) ErrorHandler

	Handle(err error) error
}

type errorHandlerImpl struct {
	cases       []func(err error) (error, bool)
	defaultCode codes.Code
}

func (e *errorHandlerImpl) Test(caseFn func(initialErr error) (err error, ok bool)) ErrorHandler {
	e.cases = append(e.cases, caseFn)
	return e
}

func (e *errorHandlerImpl) Default(code codes.Code) ErrorHandler {
	e.defaultCode = code
	return e
}

func (e *errorHandlerImpl) Is(target error, code codes.Code) ErrorHandler {
	e.cases = append(e.cases, func(err error) (error, bool) {
		if !errors.Is(err, target) {
			return nil, false
		}

		return status.Errorf(code, "%s", err), true
	})
	return e
}

func (e *errorHandlerImpl) IsW(target error, code codes.Code, wrap error) ErrorHandler {
	e.cases = append(e.cases, func(err error) (error, bool) {
		if !errors.Is(err, target) {
			return nil, false
		}

		return status.Errorf(code, "%s", errors.Join(wrap, err)), true
	})
	return e
}

func (e *errorHandlerImpl) IsWF(target error, code codes.Code, format string, args ...interface{}) ErrorHandler {
	e.cases = append(e.cases, func(err error) (error, bool) {
		if !errors.Is(err, target) {
			return nil, false
		}

		return status.Errorf(code, format, args...), true
	})
	return e
}

func (e *errorHandlerImpl) As(target interface{}, code codes.Code) ErrorHandler {
	e.cases = append(e.cases, func(err error) (error, bool) {
		if !errors.As(err, target) {
			return nil, false
		}

		return status.Errorf(code, "%s", err), true
	})
	return e
}

func (e *errorHandlerImpl) AsW(target interface{}, code codes.Code, wrap error) ErrorHandler {
	e.cases = append(e.cases, func(err error) (error, bool) {
		if !errors.As(err, target) {
			return nil, false
		}

		return status.Errorf(code, "%s", errors.Join(wrap, err)), true
	})
	return e
}

func (e *errorHandlerImpl) AsWF(target interface{}, code codes.Code, format string, args ...interface{}) ErrorHandler {
	e.cases = append(e.cases, func(err error) (error, bool) {
		if !errors.As(err, target) {
			return nil, false
		}

		return status.Errorf(code, format, args...), true
	})
	return e
}

func (e *errorHandlerImpl) Handle(err error) error {
	for _, caseFn := range e.cases {
		if res, ok := caseFn(err); ok {
			return res
		}
	}

	return status.Errorf(e.defaultCode, "%s", err)
}

func HandleError(defaultCode codes.Code) ErrorHandler {
	return &errorHandlerImpl{defaultCode: defaultCode}
}
