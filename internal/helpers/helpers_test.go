package helpers

import (
	"bytes"
	"errors"
	"os"
	"testing"
)

func TestCheckPanic(t *testing.T) {
	testErr := errors.New("test: invalid string")
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()
	Check(testErr)
}

func TestCheckNoPanic(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("The code panicked unexpectedly")
		}
	}()
	Check(nil)
}

func TestDebugPrintsWhenDebugIsTrue(t *testing.T) {
	DEBUG = true
	expected := "---------\nTest Debug Message\n---------\n\n"

	// Redirect stdout for testing
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	Debug("Test Debug Message")

	// Reset stdout
	w.Close()
	os.Stdout = old

	// Read redirected output
	var buf bytes.Buffer
	buf.ReadFrom(r)
	actual := buf.String()

	if actual != expected {
		t.Errorf("Debug did not print the expected message. Got: %q, Expected: %q", actual, expected)
	}
}

func TestDebugDoesNotPrintWhenDebugIsFalse(t *testing.T) {
	DEBUG = false

	// Redirect stdout for testing
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	Debug("Test Debug Message")

	// Reset stdout
	w.Close()
	os.Stdout = old

	// Read redirected output
	var buf bytes.Buffer
	buf.ReadFrom(r)
	actual := buf.String()

	if actual != "" {
		t.Errorf("Debug printed when DEBUG was false. Got: %q, Expected: empty", actual)
	}
}
