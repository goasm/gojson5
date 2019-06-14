package json5_test

import (
	"testing"

	json5 "github.com/goasm/gojson5"
)

func expectToken(t *testing.T, token json5.Token, expected json5.TokenType) {
	if token.Type != expected {
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

func TestReadIntegerNumber(t *testing.T) {
	lexer := json5.NewLexer(` 5 `)
	t0, err := lexer.Token()
	noError(t, err)
	expectToken(t, t0, json5.TypeNumber)
	equals(t, int64(5), t0.Value)
	t1, err := lexer.Token()
	noError(t, err)
	expectToken(t, t1, json5.TypeEOF)
}

func TestReadNegativeIntegerNumber(t *testing.T) {
	lexer := json5.NewLexer(` -10 `)
	t0, err := lexer.Token()
	noError(t, err)
	expectToken(t, t0, json5.TypeNumber)
	equals(t, int64(-10), t0.Value)
	t1, err := lexer.Token()
	noError(t, err)
	expectToken(t, t1, json5.TypeEOF)
}

func TestReadFloatNumber(t *testing.T) {
	lexer := json5.NewLexer(` 12.566 `)
	t0, err := lexer.Token()
	noError(t, err)
	expectToken(t, t0, json5.TypeNumber)
	equals(t, 12.566, t0.Value)
	t1, err := lexer.Token()
	noError(t, err)
	expectToken(t, t1, json5.TypeEOF)
}

func TestReadExponentNumber(t *testing.T) {
	lexer := json5.NewLexer(` 3.14e8 `)
	t0, err := lexer.Token()
	noError(t, err)
	expectToken(t, t0, json5.TypeNumber)
	equals(t, 314000000.0, t0.Value)
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

func TestReadInWhitespaces(t *testing.T) {
	samples := [][]string{
		{`"foo"`, ` "foo"`, `"foo" `},
		{`100`, ` 100`, `100 `},
		{`true`, ` true`, `true `},
		{`null`, ` null`, `null `},
	}
	expectedTypes := []json5.TokenType{
		json5.TypeString,
		json5.TypeNumber,
		json5.TypeBool,
		json5.TypeNull,
	}
	for idx, line := range samples {
		expectedType := expectedTypes[idx]
		for _, sample := range line {
			lexer := json5.NewLexer(sample)
			t0, err := lexer.Token()
			noError(t, err)
			expectToken(t, t0, expectedType)
			t1, err := lexer.Token()
			noError(t, err)
			expectToken(t, t1, json5.TypeEOF)
		}
	}
}

func TestReadSingleLineComment(t *testing.T) {
	lexer := json5.NewLexer(`
	// This is a comment
	null
	`)
	t0, err := lexer.Token()
	noError(t, err)
	expectToken(t, t0, json5.TypeNull)
	equals(t, nil, t0.Value)
	t1, err := lexer.Token()
	noError(t, err)
	expectToken(t, t1, json5.TypeEOF)
}

func TestReadMultipleLineComment(t *testing.T) {
	lexer := json5.NewLexer(`
	/* =================
	 * This is a comment
	 * Ignore me
	 * =================
	 */
	null
	`)
	t0, err := lexer.Token()
	noError(t, err)
	expectToken(t, t0, json5.TypeNull)
	equals(t, nil, t0.Value)
	t1, err := lexer.Token()
	noError(t, err)
	expectToken(t, t1, json5.TypeEOF)
}

func TestReadUnclosedComment(t *testing.T) {
	lexer := json5.NewLexer(`
	null /* ==== Unclosed comment ====
	`)
	t0, err := lexer.Token()
	noError(t, err)
	expectToken(t, t0, json5.TypeNull)
	equals(t, nil, t0.Value)
	t1, err := lexer.Token()
	hasError(t, err, "unexpected end of JSON")
	expectToken(t, t1, json5.TypeNone)
}
