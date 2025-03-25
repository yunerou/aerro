package apperror

import (
	"errors"
	"fmt"

	"encoding/json"
)

type multiAppError[T ErrorCode] struct {
	errs []AppError[T]
	tags map[string]string
}

func (a Aerro[T]) Append(mErr MultiAppError[T], appErrs ...AppError[T]) MultiAppError[T] {
	if mErr == nil {
		mErr = &multiAppError[T]{
			errs: appErrs,
		}
		return mErr
	}
	me, _ := mErr.(*multiAppError[T])
	me.errs = append(me.errs, appErrs...)
	return me
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

// make JSON format from error data.
func (me *multiAppError[T]) ToJSON() json.RawMessage {
	data := me.errs
	jsonB, _ := json.Marshal(data)
	return jsonB
}

func (me *multiAppError[T]) MarshalJSON() ([]byte, error) {
	return me.ToJSON(), nil
}

func (me *multiAppError[T]) SetTag(key string, value string) MultiAppError[T] {
	if me.tags == nil {
		me.tags = make(map[string]string)
	}
	me.tags[key] = value
	return me
}

func (me *multiAppError[T]) GetTag(key string) (string, bool) {
	if me.tags == nil {
		return "", false
	}
	value, ok := me.tags[key]
	return value, ok
}

func (me *multiAppError[T]) Tags() map[string]string {
	return me.tags
}
