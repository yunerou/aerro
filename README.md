# aerro

- Make your error prettier 
- aerro name find out from mix betweenn `App error` and `Aero effect`
 
# Features

- No dependency with any packages
- Sure: it's implement `error` interface
- Generic with your own enum
- Allow custom error message with any `i18n` package
- Custom hook when error is created
- Tools for generating enum error
- Stacktrace
- Multi error
- Detailed error
- support context, that help you inject anything to control your error (enricher log, know exactly where error come from, filter, ...)

## Make more ultility function

### Usecase 1: Write your own validator error

```
package apperror

import (
	"context"

	detailtype "github.com/yunerou/aerro/aerror/detail_type"
	"git.ponos-tech.com/gophers/toys/providers/validation"
	"github.com/samber/lo"
)

func ValidationErrsAggregate(
	ctx context.Context,
	validationErrs []validation.ValidationErr,
) DetailAppError {
	detail := detailtype.ValidateError(
		lo.Map(
			validationErrs,
			func(vErr validation.ValidationErr, _ int) detailtype.ValidateErrorItem {
				return detailtype.ValidateErrorItem{
					Key:     vErr.Field,
					Message: vErr.Error(),
				}
			},
		))

	return NewWithDetail(ctx, ErrManyValidation, nil, detail, map[string]any{
		"Length": len(validationErrs),
	})
}
```



```

// get http status code which corresponds to error.
func (e *appError) HTTPStatusCode() int {
	x := int(e.errCode)
	switch {
	case x == int(ErrOrigin):
		return http.StatusBadRequest
	case x > int(start400) && x < int(end400):
		return http.StatusBadRequest
	case x > int(start401) && x < int(end401):
		return http.StatusUnauthorized
	case x > int(start403) && x < int(end403):
		return http.StatusForbidden
	case x > int(start500) && x < int(end500):
		return http.StatusInternalServerError
	case x > int(start500Trace) && x < int(end500Trace):
		return http.StatusInternalServerError
	}

	return http.StatusNotExtended
}

```