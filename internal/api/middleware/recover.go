package middleware

import (
	"context"
	"fmt"
	"net/http"
	"runtime"

	"connectrpc.com/connect"
)

func Recover(_ context.Context, _ connect.Spec, _ http.Header, p any) error {
	stack := make([]byte, 64<<10)
	stack = stack[:runtime.Stack(stack, false)]
	return &PanicError{Panic: p, Stack: stack}
}

type PanicError struct {
	Panic any
	Stack []byte
}

func (e *PanicError) Error() string {
	return fmt.Sprintf("panic caught: %v\n\n%s", e.Panic, e.Stack)
}
