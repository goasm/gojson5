package json5

import "io"

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
	stateValue
	stateArray
	stateObject
	stateString
	stateEscapeChar
)

type token struct {
	Type  tokenType
	Value interface{}
}

// Lexer reads and tokenizes a JSON string
type Lexer struct {
	str   string
	pos   int
	state lexerState
	buf   []byte
	ret   token
	err   error
}

// NewLexer creates a JSON5 Lexer
func NewLexer(str string) *Lexer {
	return &Lexer{str: str, buf: []byte{}}
}

func (l *Lexer) readDefault() {
	c := l.str[l.pos]
	switch c {
	case ' ', '\t', '\n', '\r':
		l.pos++
	case '/':
		l.state = stateComment
		l.pos++
	default:
		l.state = stateValue
	}
}

func (l *Lexer) readValue() {
	c := l.str[l.pos]
	switch c {
	case '[':
		l.state = stateArray
		l.pos++
	case '{':
		l.state = stateObject
		l.pos++
	case '"':
		l.state = stateString
		l.pos++
		// TODO: case number
		// TODO: case bool
		// TODO: case null
	default:
		l.err = badCharError(c, l.pos)
	}
}

func (l *Lexer) readString() {
	c := l.str[l.pos]
	switch c {
	case '\\':
		l.state = stateEscapeChar
		l.pos++
	case '"':
		value := string(l.buf)
		l.ret = token{typeString, value}
		l.pos++
	default:
		l.buf = append(l.buf, c)
		l.pos++
	}
}

// Token gets the next JSON token
func (l *Lexer) Token() (token, error) {
	l.state = stateDefault
	l.ret = token{typeNull, nil}
	for {
		if l.pos >= len(l.str) {
			l.err = io.EOF
		}
		if l.ret.Type != typeNull || l.err != nil {
			return l.ret, l.err
		}
		switch l.state {
		case stateDefault:
			l.readDefault()
		case stateComment:
			// TODO: read comment
		case stateValue:
			l.readValue()
		case stateArray:
			// TODO: read array
		case stateObject:
			// TODO: read object
		case stateString:
			l.readString()
		case stateEscapeChar:
			// TODO: read escape char
		}
	}
}
