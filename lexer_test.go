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
	equals(t, "foo", t0.Value)
	t1, err := lexer.Token()
	noError(t, err)
	expectToken(t, t1, json5.TypeEOF)
}

func TestReadBool(t *testing.T) {
	lexer := json5.NewLexer(` true false `)
	t0, err := lexer.Token()
	noError(t, err)
	expectToken(t, t0, json5.TypeBool)
	equals(t, true, t0.Value)
	t1, err := lexer.Token()
	noError(t, err)
	expectToken(t, t1, json5.TypeBool)
	equals(t, false, t1.Value)
	t2, err := lexer.Token()
	noError(t, err)
	expectToken(t, t2, json5.TypeEOF)
}

func TestReadNull(t *testing.T) {
	lexer := json5.NewLexer(` null `)
	t0, err := lexer.Token()
	noError(t, err)
	expectToken(t, t0, json5.TypeNull)
	equals(t, nil, t0.Value)
	t1, err := lexer.Token()
	noError(t, err)
	expectToken(t, t1, json5.TypeEOF)
}
