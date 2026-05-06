package utils

import (
	"context"
	"fmt"
	"log"
	"runtime/debug"
)

// RecoverFunc is a function that handles panics
type RecoverFunc func(recovered interface{}, stack []byte)

// DefaultRecoveryHandler logs the panic and stack trace
func DefaultRecoveryHandler(recovered interface{}, stack []byte) {
	log.Printf("PANIC RECOVERED: %v\n%s", recovered, stack)
}

// Recover wraps a function with panic recovery
func Recover(fn func(), handler RecoverFunc) {
	defer func() {
		if r := recover(); r != nil {
			stack := debug.Stack()
			if handler != nil {
				handler(r, stack)
			} else {
				DefaultRecoveryHandler(r, stack)
			}
		}
	}()
	fn()
}

// RecoverWithError wraps a function with panic recovery and returns an error
func RecoverWithError(fn func() error) (err error) {
	defer func() {
		if r := recover(); r != nil {
			stack := debug.Stack()
			err = fmt.Errorf("panic recovered: %v\n%s", r, stack)
		}
	}()
	err = fn()
	return
}

// RecoverWithResult wraps a function with panic recovery and returns result with error
func RecoverWithResult[T any](fn func() (T, error)) (result T, err error) {
	defer func() {
		if r := recover(); r != nil {
			stack := debug.Stack()
			err = fmt.Errorf("panic recovered: %v\n%s", r, stack)
		}
	}()
	result, err = fn()
	return
}

// PanicHandler provides enhanced panic recovery with context
type PanicHandler struct {
	logger interface{} // Can be any logger interface
}

// NewPanicHandler creates a new panic handler
func NewPanicHandler(logger interface{}) *PanicHandler {
	return &PanicHandler{logger: logger}
}

// RecoverWithError recovers from panic and returns error
func (ph *PanicHandler) RecoverWithError(ctx context.Context, operation string, fn func()) (err error) {
	defer func() {
		if r := recover(); r != nil {
			stack := debug.Stack()
			err = fmt.Errorf("panic in %s: %v\n%s", operation, r, stack)
			if ph.logger != nil {
				// Log if logger is available
				log.Printf("PANIC in %s: %v\n%s", operation, r, stack)
			}
		}
	}()
	fn()
	return nil
}

// SafeGo launches a goroutine with panic recovery
func (ph *PanicHandler) SafeGo(ctx context.Context, operation string, fn func()) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				stack := debug.Stack()
				if ph.logger != nil {
					log.Printf("PANIC in goroutine %s: %v\n%s", operation, r, stack)
				} else {
					log.Printf("PANIC in goroutine %s: %v\n%s", operation, r, stack)
				}
			}
		}()
		fn()
	}()
}