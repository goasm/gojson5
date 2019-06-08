package json5_test

import (
	"testing"

	json5 "github.com/goasm/gojson5"
)

func expectToken(t *testing.T, token json5.Token, exp json5.TokenType) {
	if token.Type != exp {
		t.Fatal("Unexpected token:", token)
	}
}

func TestReadString(t *testing.T) {
	lexer := json5.NewLexer(` "foo" `)
	t0, err := lexer.Token()
	noError(t, err)
	expectToken(t, t0, json5.TypeString)
	t1, err := lexer.Token()
	noError(t, err)
	expectToken(t, t1, json5.TypeEOF)
}
