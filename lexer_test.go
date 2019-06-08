package json5_test

import (
	"testing"

	json5 "github.com/goasm/gojson5"
)

func TestReadString(t *testing.T) {
	lexer := json5.NewLexer(` "foo" `)
	t0, err := lexer.Token()
	if err != nil {
		t.Fatal("Unexpected error:", err)
	}
	if t0.Type != json5.TypeString {
		t.Fatal("Unexpected token:", t0)
	}
	t1, err := lexer.Token()
	if err != nil {
		t.Fatal("Unexpected error:", err)
	}
	if t1.Type != json5.TypeEOF {
		t.Fatal("Unexpected token:", t1)
	}
}
