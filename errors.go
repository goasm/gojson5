package json5

// SyntaxError means JSON has an incorrect syntax
type SyntaxError struct {
	message string
}

func newSyntaxError(message string) *SyntaxError {
	return &SyntaxError{message}
}

func (e *SyntaxError) Error() string {
	return "json5: " + e.message
}
