package terr //TODO: need to be renamed

import (
	"errors"
	"fmt"
	"runtime"
	"strings"
)

type TraceError struct {
	error
	*stack
}

func New(msg string) error {
	return &TraceError{
		errors.New(msg),
		callers(),
	}
}

func Wrap(err error) error {
	if err == nil {
		return nil
	}

	if traceErr, ok := err.(*TraceError); ok {
		return traceErr
	}

	return &TraceError{
		err,
		callers(),
	}
}

func UnWrap(err error) error {
	if traceErr, ok := err.(*TraceError); ok {
		return traceErr.error
	}

	return err
}

func (e *TraceError) Error() string {
	return fmt.Sprint(e.error.Error(), "\tStackTrace: ", e.withStackMsg())
}

type stack []uintptr
type frame uintptr

func (f frame) pc() uintptr { return uintptr(f) - 1 }
func (f frame) fileLine() (file string, line int) {
	fn := runtime.FuncForPC(f.pc())
	if fn == nil {
		return "unknown", 0
	}

	return fn.FileLine(f.pc())
}

func callers() *stack {
	const depth = 32
	var pcs [depth]uintptr
	n := runtime.Callers(3, pcs[:])
	var st stack = pcs[0:n]
	return &st
}

func (e *TraceError) withStackMsg() string {
	var b []byte
	for i, frame := range e.Frames() {
		if i > 0 {
			b = append(b, " -> "...)
		}

		b = append(b, fmt.Sprintf("%s:%d", frame.File, frame.Line)...)
	}
	return string(b)
}

func withoutPath(filePath string) string {
	paths := strings.Split(filePath, "/")
	return paths[len(paths)-1]
}

func (te *TraceError) Frames() []*Frame {
	var frames []*Frame

	for _, pc := range *te.stack {
		fr := frame(pc)
		file, line := fr.fileLine()
		frames = append(frames, &Frame{
			Pc:   fr.pc(),
			File: withoutPath(file),
			Line: line,
		})
	}

	return frames
}

type Frame struct {
	Pc   uintptr
	File string
	Line int
}
