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
