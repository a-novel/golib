package grpc

import (
	"errors"
	"sync"

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

	mu sync.RWMutex
}

func (errorHandler *errorHandlerImpl) Test(caseFn func(initialErr error) (err error, ok bool)) ErrorHandler {
	errorHandler.mu.Lock()
	defer errorHandler.mu.Unlock()

	errorHandler.cases = append(errorHandler.cases, caseFn)
	return errorHandler
}

func (errorHandler *errorHandlerImpl) Is(target error, code codes.Code) ErrorHandler {
	errorHandler.mu.Lock()
	defer errorHandler.mu.Unlock()

	errorHandler.cases = append(errorHandler.cases, func(err error) (error, bool) {
		if !errors.Is(err, target) {
			return nil, false
		}

		return status.Errorf(code, "%s", err), true
	})
	return errorHandler
}

func (errorHandler *errorHandlerImpl) IsW(target error, code codes.Code, wrap error) ErrorHandler {
	errorHandler.mu.Lock()
	defer errorHandler.mu.Unlock()

	errorHandler.cases = append(errorHandler.cases, func(err error) (error, bool) {
		if !errors.Is(err, target) {
			return nil, false
		}

		return status.Errorf(code, "%s", errors.Join(wrap, err)), true
	})
	return errorHandler
}

func (errorHandler *errorHandlerImpl) IsWF(
	target error, code codes.Code, format string, args ...interface{},
) ErrorHandler {
	errorHandler.mu.Lock()
	defer errorHandler.mu.Unlock()

	errorHandler.cases = append(errorHandler.cases, func(err error) (error, bool) {
		if !errors.Is(err, target) {
			return nil, false
		}

		return status.Errorf(code, format, args...), true
	})
	return errorHandler
}

func (errorHandler *errorHandlerImpl) As(target interface{}, code codes.Code) ErrorHandler {
	errorHandler.mu.Lock()
	defer errorHandler.mu.Unlock()

	errorHandler.cases = append(errorHandler.cases, func(err error) (error, bool) {
		if !errors.As(err, target) {
			return nil, false
		}

		return status.Errorf(code, "%s", err), true
	})
	return errorHandler
}

func (errorHandler *errorHandlerImpl) AsW(target interface{}, code codes.Code, wrap error) ErrorHandler {
	errorHandler.mu.Lock()
	defer errorHandler.mu.Unlock()

	errorHandler.cases = append(errorHandler.cases, func(err error) (error, bool) {
		if !errors.As(err, target) {
			return nil, false
		}

		return status.Errorf(code, "%s", errors.Join(wrap, err)), true
	})
	return errorHandler
}

func (errorHandler *errorHandlerImpl) AsWF(
	target interface{}, code codes.Code, format string, args ...interface{},
) ErrorHandler {
	errorHandler.mu.Lock()
	defer errorHandler.mu.Unlock()

	errorHandler.cases = append(errorHandler.cases, func(err error) (error, bool) {
		if !errors.As(err, target) {
			return nil, false
		}

		return status.Errorf(code, format, args...), true
	})
	return errorHandler
}

func (errorHandler *errorHandlerImpl) Handle(err error) error {
	errorHandler.mu.RLock()
	defer errorHandler.mu.RUnlock()

	for _, caseFn := range errorHandler.cases {
		if res, ok := caseFn(err); ok {
			return res
		}
	}

	return status.Errorf(errorHandler.defaultCode, "%s", err)
}

// HandleError creates a new error handler. You can define different case depending on the error you want to handle.
// This handler is thread safe, and a single handler can be shared between multiple goroutines.
func HandleError(defaultCode codes.Code) ErrorHandler {
	return &errorHandlerImpl{defaultCode: defaultCode}
}
