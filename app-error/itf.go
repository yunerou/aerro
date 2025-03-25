package apperror

import "encoding/json"

type AppError[T ErrorCode] interface {
	error

	ErrorCode() T
	Origin() error
	Stacktrace() stackTraceT

	SetTag(key string, value string) AppError[T]
	GetTag(key string) (string, bool)
	Tags() map[string]string

	Is(target error) bool
	As(target interface{}) bool
	Unwrap() error

	MarshalJSON() ([]byte, error)
	ToJSON() json.RawMessage
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

	SetTag(key string, value string) MultiAppError[T]
	GetTag(key string) (string, bool)
	Tags() map[string]string

	Errors() []AppError[T]
}
