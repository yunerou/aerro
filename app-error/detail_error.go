package apperror

import (
	"context"
	"encoding/json"
)

type detailAppError[T ErrorCode, D any] struct {
	*appError[T]
	detail D
}

// add detail info to error.
func (a Aerro[T, D]) NewWithDetail(
	ctx context.Context,
	code T,
	origin error,
	detail D,
	templateData ...map[string]any) DetailAppError[T] {
	aErr := &detailAppError[T, D]{
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

func (de *detailAppError[T, D]) Detail() any {
	return de.detail
}

func (de *detailAppError[T, D]) ToJSON() json.RawMessage {
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

func (de *detailAppError[T, D]) MarshalJSON() ([]byte, error) {
	return de.ToJSON(), nil
}

func (de *detailAppError[T, D]) CastToDetail() (out DetailAppError[T], ok bool) {
	return de, true
}
