package apperror

import "encoding/json"

type AppError[T ErrorCode] interface {
	error

	ErrorCode() ErrorCode
	Origin() error
	Stacktrace() stackTraceT

	Is(target error) bool
	As(target interface{}) bool
	MarshalJSON() ([]byte, error)
	ToJSON() json.RawMessage
	Unwrap() error
	CastToDetail() (out DetailAppError[T], ok bool)
}

type DetailAppError[T ErrorCode] interface {
	AppError[T]
	Detail() any
}

type MultiAppError[T ErrorCode] interface {
	error

	Is(target error) bool
	MarshalJSON() ([]byte, error)
	ToJSON() json.RawMessage

	Append(e ...AppError[T])
	Errors() []AppError[T]
	SetHttpStatusCode(code int)
}
