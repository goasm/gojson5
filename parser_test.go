package json5_test

import (
	"testing"

	json5 "github.com/goasm/gojson5"
)

func TestParser(t *testing.T) {
	parser := json5.NewParser(`
	{
		// Use IntelliSense to learn about possible attributes.
		// Hover to view descriptions of existing attributes.
		// For more information, visit: https://go.microsoft.com/fwlink/?linkid=830387
		"version": "0.2.0",
		"configurations": [
			{
				"type": "node",
				"request": "launch",
				"name": "Launch Program",
				"program": "${workspaceFolder}/vendor/run.js"
			}
		]
	}
	`)
	parser.Parse()
}

func TestParseString(t *testing.T) {
	parser := json5.NewParser(` "foo" `)
	val, err := parser.Parse()
	noError(t, err)
	equals(t, "foo", val)
}

func TestParseArray(t *testing.T) {
	parser := json5.NewParser(` [1, 2, 3] `)
	raw, err := parser.Parse()
	noError(t, err)
	val, ok := raw.([]interface{})
	equals(t, true, ok)
	equals(t, 3, len(val))
	equals(t, int64(1), val[0])
	equals(t, int64(2), val[1])
	equals(t, int64(3), val[2])
}
