package json5_test

import (
	"testing"

	json5 "github.com/goasm/gojson5"
)

func expectToken(t *testing.T, token json5.Token, err error) {
	if err != nil {
		t.Fatal("Unexpected error:", err)
	}
	if token.Type != json5.TypeString {
		t.Fatal("Unexpected token:", token)
	}
}

func TestReadString(t *testing.T) {
	lexer := json5.NewLexer(` "foo" `)
	t0, err := lexer.Token()
	expectToken(t, t0, err)
	t1, err := lexer.Token()
	expectToken(t, t1, err)
}
