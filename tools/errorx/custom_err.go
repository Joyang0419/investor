package errorx

import (
	"fmt"
	"runtime"
)

type CustomError struct {
	Name     string
	CausedBy error
	line     int
	file     string
}

func New(name string, causedBy error, detailMsg string) *CustomError {
	_, file, line, _ := runtime.Caller(1)
	return &CustomError{
		Name:     name,
		CausedBy: fmt.Errorf("%w ->\n%s", causedBy, detailMsg),
		file:     file,
		line:     line,
	}
}

func (e *CustomError) Unwrap() error {
	return e.CausedBy
}

func (e *CustomError) Error() string {
	return fmt.Sprintf("Name: %s, CausedBy: %v\nfile: %s , line: %d",
		e.Name,
		e.CausedBy,
		e.file,
		e.line,
	)
}
