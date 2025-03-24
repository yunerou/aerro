package apperror

import "context"

type ErrorCode interface {
	error
}

type Aerro[TCode ErrorCode, TDetail any] struct {
	// Inject function to create error message. If your context comes with a language key, use it to create an i18n message.
	BuildErrorMessage func(ctx context.Context,
		code TCode,
		origin error,
		templateData map[string]any) string

	// Configure reasonable conditions to reduce the workload on your computer.
	StacktraceEnabled func(code TCode) bool

	// Should be configured with conditions to reduce the workload on your computer.
	HookAfterCreated func(ctx context.Context, aerror AppError[TCode])
}
