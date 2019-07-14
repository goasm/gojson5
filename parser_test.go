package json5_test

import (
	"testing"

	json5 "github.com/goasm/gojson5"
)

func TestParseString(t *testing.T) {
	parser := json5.Parser{}
	val, err := parser.Parse([]byte(` "foo" `))
	noError(t, err)
	equals(t, "foo", val)
}

func TestParseNumber(t *testing.T) {
	parser := json5.Parser{}
	val, err := parser.Parse([]byte(` 100 `))
	noError(t, err)
	equals(t, int64(100), val)
}

func TestParseBool(t *testing.T) {
	parser := json5.Parser{}
	val, err := parser.Parse([]byte(` true `))
	noError(t, err)
	equals(t, true, val)
}

func TestParseNull(t *testing.T) {
	parser := json5.Parser{}
	val, err := parser.Parse([]byte(` null `))
	noError(t, err)
	equals(t, nil, val)
}

func TestParseArray(t *testing.T) {
	parser := json5.Parser{}
	raw, err := parser.Parse([]byte(` [1, 2, 3] `))
	noError(t, err)
	val, ok := raw.([]interface{})
	equals(t, true, ok)
	equals(t, 3, len(val))
	equals(t, int64(1), val[0])
	equals(t, int64(2), val[1])
	equals(t, int64(3), val[2])
}

func TestParseNestedArray(t *testing.T) {
	parser := json5.Parser{}
	raw, err := parser.Parse([]byte(` [[[1, 2, 3]]] `))
	noError(t, err)
	val, ok := raw.([]interface{})
	equals(t, true, ok)
	equals(t, 1, len(val))
	val, ok = val[0].([]interface{})
	equals(t, true, ok)
	equals(t, 1, len(val))
	val, ok = val[0].([]interface{})
	equals(t, true, ok)
	equals(t, 3, len(val))
	equals(t, int64(1), val[0])
	equals(t, int64(2), val[1])
	equals(t, int64(3), val[2])
}

func TestParseObject(t *testing.T) {
	parser := json5.Parser{}
	raw, err := parser.Parse([]byte(` { "foo": 1, "bar": 2, "baz": 3 } `))
	noError(t, err)
	val, ok := raw.(map[string]interface{})
	equals(t, true, ok)
	equals(t, int64(1), val["foo"])
	equals(t, int64(2), val["bar"])
	equals(t, int64(3), val["baz"])
}

func TestParseNestedObject(t *testing.T) {
	parser := json5.Parser{}
	raw, err := parser.Parse([]byte(` { "foo": { "bar": { "baz": 3 } } } `))
	noError(t, err)
	val, ok := raw.(map[string]interface{})
	equals(t, true, ok)
	equals(t, 1, len(val))
	val, ok = val["foo"].(map[string]interface{})
	equals(t, true, ok)
	equals(t, 1, len(val))
	val, ok = val["bar"].(map[string]interface{})
	equals(t, true, ok)
	equals(t, 1, len(val))
	equals(t, int64(3), val["baz"])
}
