package json5

import "fmt"

// SyntaxError means JSON has an incorrect syntax
type SyntaxError struct {
	message string
	index   int
}

func newSyntaxError(message string, index int) *SyntaxError {
	return &SyntaxError{message, index}
}

func (e *SyntaxError) Error() string {
	return fmt.Sprintf("json5: %s at position %d", e.message, e.index)
}

func badCharError(ch byte, index int) *SyntaxError {
	return newSyntaxError(fmt.Sprintf("unexpected character: %c", ch), index)
}