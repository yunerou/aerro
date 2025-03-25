package apperror

import (
	"errors"
	"fmt"

	"encoding/json"
)

type multiAppError[T ErrorCode] struct {
	errs           []AppError[T]
	httpStatusCode *int
}

func (a Aerro[T]) Append(mErr MultiAppError[T], appErrs ...AppError[T]) MultiAppError[T] {
	if mErr == nil {
		mErr = &multiAppError[T]{
			errs: appErrs,
		}
	} else {
		mErr.Append(appErrs...)
	}
	return mErr
}

func (me *multiAppError[T]) Append(e ...AppError[T]) {
	me.errs = append(me.errs, e...)
}

func (me *multiAppError[T]) Errors() []AppError[T] {
	return me.errs
}

func (me *multiAppError[T]) Error() string {
	result := fmt.Sprintf("Total %d errors: \n", len(me.errs))
	for _, err := range me.errs {
		result += fmt.Sprintf("  %s\n", err.Error())
	}
	return result
}

func (me *multiAppError[T]) Is(target error) bool {
	for _, err := range me.errs {
		if errors.Is(err, target) {
			return true
		}
	}
	return false
}

// set http status code for multierror when return it to client or use default httpStatusCode
func (me *multiAppError[T]) SetHttpStatusCode(code int) {
	tmp := code
	me.httpStatusCode = &tmp
}

// make JSON format from error data.
func (me *multiAppError[T]) ToJSON() json.RawMessage {
	data := me.errs
	jsonB, _ := json.Marshal(data)
	return jsonB
}

func (me *multiAppError[T]) MarshalJSON() ([]byte, error) {
	return me.ToJSON(), nil
}
