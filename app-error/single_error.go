package apperror

import (
	"context"
	"encoding/json"
	"errors"
)

type appError[T ErrorCode] struct {
	origin  error
	errCode T
	i18nMsg string
	tags    map[string]string
	*stack
}

func (a Aerro[T]) New(
	ctx context.Context,
	code T,
	origin error,
	templateData ...map[string]any,
) AppError[T] {
	aErr := &appError[T]{
		origin:  origin,
		errCode: code,
	}

	if len(templateData) > 0 {
		aErr.i18nMsg = a.BuildErrorMessage(ctx, code, origin, templateData[0])
	} else {
		aErr.i18nMsg = a.BuildErrorMessage(ctx, code, origin, nil)
	}

	if a.StacktraceEnabled(code) {
		aErr.stack = callers(3) // Skip 3 func runtime.Callers, callers, new, New
	}

	a.HookAfterCreated(ctx, aErr)

	return aErr
}

// method error interface.
func (e *appError[T]) Error() string {
	return e.i18nMsg
}

func (e *appError[T]) Unwrap() error {
	if e.origin != nil {
		return e.origin
	}
	return e.errCode
}

func (e *appError[T]) Is(target error) bool {
	return e.errCode.Error() == target.Error() ||
		errors.Is(e.origin, target)
}

func (e *appError[T]) As(target interface{}) bool {
	return errors.As(e.origin, target)
}

// get base error code.
func (e *appError[T]) ErrorCode() T {
	return e.errCode
}

// make JSON format from error data.
func (e *appError[T]) ToJSON() json.RawMessage {
	type responseError struct {
		ErrorCode string            `json:"error_code"`
		Message   string            `json:"message"`
		Tags      map[string]string `json:"tags,omitempty"`
	}

	data := responseError{
		ErrorCode: e.ErrorCode().Error(),
		Message:   e.Error(),
		Tags:      e.Tags(),
	}

	jsonB, _ := json.Marshal(data)
	return jsonB
}

func (e *appError[T]) MarshalJSON() ([]byte, error) {
	return e.ToJSON(), nil
}

func (e *appError[T]) CastToDetail() (out DetailAppError[T], ok bool) {
	return nil, false
}

func (e *appError[T]) Origin() error {
	return e.origin
}

func (e *appError[T]) Stacktrace() stackTraceT {
	return e.stacktrace()
}

func (e *appError[T]) SetTag(key string, value string) AppError[T] {
	if e.tags == nil {
		e.tags = make(map[string]string)
	}
	e.tags[key] = value
	return e
}

func (e *appError[T]) GetTag(key string) (string, bool) {
	if e.tags == nil {
		return "", false
	}
	value, ok := e.tags[key]
	return value, ok
}

func (e *appError[T]) Tags() map[string]string {
	return e.tags
}
