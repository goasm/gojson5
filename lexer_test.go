package json5_test

import (
	"testing"

	json5 "github.com/goasm/gojson5"
)

func expectToken(t *testing.T, token json5.Token, expected json5.TokenType) {
	t.Helper()
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

func TestReadEscapeChar(t *testing.T) {
	lexer := json5.NewLexer(` "foo\"bar" `)
	t0, err := lexer.Token()
	noError(t, err)
	expectToken(t, t0, json5.TypeString)
	equals(t, "foo\"bar", t0.Value)
	t1, err := lexer.Token()
	noError(t, err)
	expectToken(t, t1, json5.TypeEOF)
}

func TestReadValidEscapeChars(t *testing.T) {
	samples := []string{
		`"\""`, `"\\"`, `"\/"`, `"\b"`, `"\f"`, `"\n"`, `"\r"`, `"\t"`,
	}
	expectedValues := []string{
		"\"", "\\", "/", "\b", "\f", "\n", "\r", "\t",
	}
	for idx, sample := range samples {
		lexer := json5.NewLexer(sample)
		t0, err := lexer.Token()
		noError(t, err)
		expectToken(t, t0, json5.TypeString)
		equals(t, expectedValues[idx], t0.Value)
		t1, err := lexer.Token()
		noError(t, err)
		expectToken(t, t1, json5.TypeEOF)
	}
}

func TestReadInvalidEscapeChars(t *testing.T) {
	samples := []string{
		`"\a"`, `"\e"`, `"\v"`, `"\'"`, `"\?"`, `"\x"`,
	}
	for _, sample := range samples {
		lexer := json5.NewLexer(sample)
		t0, err := lexer.Token()
		hasError(t, err, "unexpected character")
		expectToken(t, t0, json5.TypeNone)
	}
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

func TestReadValidNumbers(t *testing.T) {
	samples := []string{
		"0", "1", "12", "1204",
		"0.2", "1.2", "12.4", "12.04",
		"0e2", "1e2", "12.4e4", "12.04e4",
	}
	expectedValues := []interface{}{
		int64(0), int64(1), int64(12), int64(1204),
		0.2, 1.2, 12.4, 12.04,
		0.0, 100.0, 124000.0, 120400.0,
	}
	for idx, sample := range samples {
		lexer := json5.NewLexer(sample)
		t0, err := lexer.Token()
		noError(t, err)
		expectToken(t, t0, json5.TypeNumber)
		equals(t, expectedValues[idx], t0.Value)
		t1, err := lexer.Token()
		noError(t, err)
		expectToken(t, t1, json5.TypeEOF)
	}
}

func TestReadInvalidNumber(t *testing.T) {
	lexer := json5.NewLexer(` 3.e8 `)
	t0, err := lexer.Token()
	hasError(t, err, "unexpected character")
	expectToken(t, t0, json5.TypeNone)
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

func TestReadInvalidLiteral(t *testing.T) {
	lexer := json5.NewLexer(` falsy `)
	t0, err := lexer.Token()
	hasError(t, err, "unexpected token")
	expectToken(t, t0, json5.TypeNone)
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
