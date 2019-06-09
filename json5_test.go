package json5_test

import (
	"strings"
	"testing"
)

func noError(t *testing.T, err error) {
	if err != nil {
		t.Fatal("Unexpected error:", err)
	}
}

func hasError(t *testing.T, err error, errStr string) {
	if err == nil {
		t.Fatal("Expected an error")
	}
	errMsg := err.Error()
	if !strings.Contains(errMsg, errStr) {
		t.Fatal("Mismatched error message:", errMsg)
	}
}

func equals(t *testing.T, expected, actual interface{}) {
	if actual != expected {
		t.Fatal("Not equal:", expected, "==", actual)
	}
}

func TestDecoding(t *testing.T) {
}
