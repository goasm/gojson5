package json5

import "fmt"

// SyntaxError means JSON has an incorrect syntax
type SyntaxError struct {
	message string
	index   int
}

func (e *SyntaxError) Error() string {
	return fmt.Sprintf("json5: %s at position %d", e.message, e.index)
}

func badCharError(ch byte, index int) *SyntaxError {
	return &SyntaxError{fmt.Sprintf("unexpected character: %c", ch), index}
}

func badTokenError(token string, index int) *SyntaxError {
	return &SyntaxError{fmt.Sprintf("unexpected token: %s", token), index}
}

func badEOF(index int) *SyntaxError {
	return &SyntaxError{"unexpected end of JSON", index}
}
