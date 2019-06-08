package json5

type tokenType int
type lexerState int

// Lexer token types
const (
	TypeNull tokenType = iota
	TypeNumber
	TypeString
	TypeBool
	TypeEOF
)

const (
	stateDefault lexerState = iota
	stateComment
	stateValue
	stateArray
	stateObject
	stateString
	stateEscapeChar
	stateBool
	stateNull
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
	return &Lexer{str: str}
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
	case 'f', 't':
		l.state = stateBool
	case 'n':
		l.state = stateNull
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
		l.ret = token{TypeString, value}
		l.pos++
	default:
		l.buf = append(l.buf, c)
		l.pos++
	}
}

func (l *Lexer) readBool() {
	p0 := l.pos
	if expectLiteral(l, "false") {
		l.ret = token{TypeBool, false}
	} else if expectLiteral(l, "true") {
		l.ret = token{TypeBool, true}
	} else {
		l.err = badTokenError(l.str[p0:l.pos], p0)
	}
}

func (l *Lexer) readNull() {
}

// Reset resets the internals for next token
func (l *Lexer) Reset() {
	l.state = stateDefault
	l.buf = []byte{}
	l.ret = token{TypeNull, nil}
}

// Token gets the next JSON token
func (l *Lexer) Token() (token, error) {
	l.Reset()
	for {
		if l.pos >= len(l.str) {
			// check EOF
			l.ret = token{TypeEOF, nil}
		}
		if l.ret.Type != TypeNull || l.err != nil {
			// check result and error
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
		case stateBool:
			l.readBool()
		case stateNull:
			l.readNull()
		}
	}
}
