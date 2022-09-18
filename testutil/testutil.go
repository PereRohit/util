package testutil

import (
	"path/filepath"
	"reflect"
	"runtime"
	"strconv"
	"strings"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

// https://www.youtube.com/watch?v=Wncq9qtzgh0

func Diff[T any](got, want T) string {
	opts := cmp.Options{
		cmp.Exporter(func(p reflect.Type) bool { return true }),
		cmpopts.EquateEmpty(),
	}
	diff := cmp.Diff(got, want, opts...)
	if diff != "" {
		return "\n-got +want\n" + diff
	}
	return ""
}

// Callers prints the stack trace of everything up till the line where Callers()
func Callers() string {
	var pc [50]uintptr
	n := runtime.Callers(2, pc[:]) // skip runtime.Callers + Callers
	callsites := make([]string, 0, n)
	frames := runtime.CallersFrames(pc[:n])
	for frame, more := frames.Next(); more; frame, more = frames.Next() {
		callsites = append(callsites, frame.File+":"+strconv.Itoa(frame.Line))
	}
	callsites = callsites[:len(callsites)-1] // skip testing.tRunner
	if len(callsites) == 1 {
		return ""
	}
	var b strings.Builder
	for i := len(callsites) - 1; i >= 0; i-- {
		if b.Len() > 0 {
			b.WriteString(" -> ")
		}
		b.WriteString(filepath.Base(callsites[i]))
	}
	return "\n" + b.String() + ":"
}
