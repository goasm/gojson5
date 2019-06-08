package json5_test

import (
	"fmt"
	"testing"

	json5 "github.com/goasm/gojson5"
)

func TestLexer(t *testing.T) {
	lexer := json5.Lexer{`
	{
		"name": "target.json",
		"foo": "bar",
		"bar": 120,
		"baz": true,
		"qux": null,
		"list": [0, 2, 4, 8],
		"dict": {
			"a": 1,
			"b": 3,
			"c": 7
		}
	}
	`}
	for {
		token, err := lexer.Token()
		if err == nil {
			break
		}
		fmt.Println(token)
	}
}
