package assert

import (
	"fmt"
	"testing"
)

var numTestsRan = 0

// Asserts that args are equal. Will prepend optional "err" arg to error message.
func Equal[T comparable](t *testing.T, actual, expected T, err ...string) {
	// NOTE-1: Mark func as helper so line numbers in here aren't reported.
	t.Helper()

	if actual != expected {
		msg := fmt.Sprintf("Got \"%v\", expected \"%v\"", actual, expected)
		if err[0] != "" {
			msg = err[0] + ". " + msg
		}
		t.Errorf(msg)
	}

	numTestsRan += 1
}

// Asserts that the arg is nil.
func Nil(t *testing.T, val any, errMsg string) {
	t.Helper() // See NOTE-1

	if val != nil {
		t.Error(errMsg)
	}

	numTestsRan += 1
}

// Asserts that the arg is not nil.
func NotNil(t *testing.T, val any, errMsg string) {
	t.Helper() // See NOTE-1

	if val == nil {
		t.Error(errMsg)
	}

	numTestsRan += 1
}

func Finished() {
	fmt.Printf("Ran %d tests\n", numTestsRan)
}