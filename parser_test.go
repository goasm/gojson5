package json5_test

import (
	"testing"

	json5 "github.com/goasm/gojson5"
)

func TestParseString(t *testing.T) {
	parser := json5.Parser{}
	val, err := parser.Parse(` "foo" `)
	noError(t, err)
	equals(t, "foo", val)
}

func TestParseArray(t *testing.T) {
	parser := json5.Parser{}
	raw, err := parser.Parse(` [1, 2, 3] `)
	noError(t, err)
	val, ok := raw.([]interface{})
	equals(t, true, ok)
	equals(t, 3, len(val))
	equals(t, int64(1), val[0])
	equals(t, int64(2), val[1])
	equals(t, int64(3), val[2])
}

func TestParseObject(t *testing.T) {
	parser := json5.Parser{}
	raw, err := parser.Parse(` { "foo": 1, "bar": 2, "baz": 3 } `)
	noError(t, err)
	val, ok := raw.(map[string]interface{})
	equals(t, true, ok)
	equals(t, int64(1), val["foo"])
	equals(t, int64(2), val["bar"])
	equals(t, int64(3), val["baz"])
}
