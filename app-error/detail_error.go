package apperror

import (
	"context"
	"encoding/json"
)

type detailAppError[T ErrorCode] struct {
	*appError[T]
	detail any
}

// add detail info to error.
func (a Aerro[T]) NewWithDetail(
	ctx context.Context,
	code T,
	origin error,
	detail any,
	templateData ...map[string]any) DetailAppError[T] {
	aErr := &detailAppError[T]{
		appError: &appError[T]{
			origin:  origin,
			errCode: code,
		},
		detail: detail,
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

func (de *detailAppError[T]) Detail() any {
	return de.detail
}

func (de *detailAppError[T]) ToJSON() json.RawMessage {
	type responseError struct {
		ErrorCode string `json:"err_code"`
		Message   string `json:"err_msg"`
		Detail    any    `json:"err_detail"`
	}

	data := responseError{
		ErrorCode: de.ErrorCode().Error(),
		Message:   de.Error(),
		Detail:    de.Detail(),
	}

	jsonB, _ := json.Marshal(data)
	return jsonB
}

func (de *detailAppError[T]) MarshalJSON() ([]byte, error) {
	return de.ToJSON(), nil
}

func (de *detailAppError[T]) CastToDetail() (out DetailAppError[T], ok bool) {
	return de, true
}
