package apperror

import (
	"fmt"
	"os"
	"strings"

	"encoding/json"
)

type frame struct {
	File     string
	Line     int
	Function string
}

type stackTraceT []*frame

func (e *appError[T]) stacktrace() (st stackTraceT) {
	if e.stack == nil {
		return nil
	}

	pwd, _ := os.Getwd()

	frames := e.stack.Frames()
	st = make([]*frame, 0, maximumStackDeep)
	for {
		fr, ok := frames.Next()
		// skip go runtime functions
		if strings.HasPrefix(fr.Function, "runtime") {
			break
		}
		if !strings.HasPrefix(fr.File, pwd) {
			continue
		}
		st = append(st, &frame{
			File:     strings.Replace(fr.File, pwd, ".", 1),
			Line:     fr.Line,
			Function: fr.Function,
		})
		if !ok {
			break
		}
	}
	return st
}

func (st stackTraceT) String() string {
	stacktraceStr := ""
	for idx, fr := range st {
		stacktraceStr += fmt.Sprintf("[%d %s:%d#%s]", idx, fr.File, fr.Line, fr.Function)
	}
	return stacktraceStr
}

// make JSON format from error data.
func (st stackTraceT) toJSON() json.RawMessage {
	var result []string = make([]string, len(st))

	for idx, fr := range st {
		result[idx] = fmt.Sprintf("%d %s:%d#%s", idx, fr.File, fr.Line, fr.Function)
	}

	jsonB, _ := json.Marshal(result)
	return jsonB
}

func (st stackTraceT) MarshalJSON() ([]byte, error) {
	return st.toJSON(), nil
}
