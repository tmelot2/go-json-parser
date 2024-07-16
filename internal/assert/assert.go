package assert

import (
	"fmt"
	"testing"
)

// Asserts that args are equal. Will prepend optional "err" arg to error message.
func Equal[T comparable](t *testing.T, actual, expected T, err ...string) {
	t.Helper()
	if actual != expected {
		msg := fmt.Sprintf("Got \"%v\", expected \"%v\"", actual, expected)
		if err[0] != "" {
			msg = err[0] + ". " + msg
		}
		t.Errorf(msg)
	}
}

// Asserts that the arg is nil.
func Nil(t *testing.T, val any, errMsg string) {
	if val != nil {
		t.Error(errMsg)
	}
}

// Asserts that the arg is not nil.
func NotNil(t *testing.T, val any, errMsg string) {
	if val == nil {
		t.Error(errMsg)
	}
}
