package json5_test

import "testing"

func noError(t *testing.T, err error) {
	if err != nil {
		t.Fatal("Unexpected error:", err)
	}
}

func TestDecoding(t *testing.T) {
}
