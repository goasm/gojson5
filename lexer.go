package json5

type tokenType int
type lexerState int

const (
	typeNull tokenType = iota
	typeNumber
	typeString
)

const (
	stateDefault lexerState = iota
	stateComment
)

type token struct {
	Type tokenType
}

// Lexer reads and tokenizes a JSON string
type Lexer struct {
	str string
	pos int
}

// Token gets the next JSON token
func (l *Lexer) Token() (token, error) {
	state := stateDefault
	return
}
