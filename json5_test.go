package json5_test

import "testing"

func noError(t *testing.T, err error) {
	if err != nil {
		t.Fatal("Unexpected error:", err)
	}
}

func equals(t *testing.T, expected, actual interface{}) {
	if actual != expected {
		t.Fatal("Not equal:", expected, "==", actual)
	}
}

func TestDecoding(t *testing.T) {
}
