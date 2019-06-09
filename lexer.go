package json5

// TokenType represents an enum of token types
type TokenType int
type lexerState int

// Lexer token types
const (
	TypeNone TokenType = iota
	TypeNumber
	TypeString
	TypeBool
	TypeNull
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

// Token represents an unit for syntax analysis
type Token struct {
	Type  TokenType
	Value interface{}
}

// Lexer reads and tokenizes a JSON string
type Lexer struct {
	str   string
	pos   int
	state lexerState
	buf   []byte
	ret   Token
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
		l.ret = Token{TypeString, value}
		l.pos++
	default:
		l.buf = append(l.buf, c)
		l.pos++
	}
}

func (l *Lexer) readBool() {
	p0 := l.pos
	if p1, ok := expectLiteral(l, "false"); ok {
		l.pos = p1
		l.ret = Token{TypeBool, false}
	} else if p1, ok := expectLiteral(l, "true"); ok {
		l.pos = p1
		l.ret = Token{TypeBool, true}
	} else {
		l.pos = p1
		l.err = badTokenError(l.str[p0:p1], p0)
	}
}

func (l *Lexer) readNull() {
	p0 := l.pos
	if p1, ok := expectLiteral(l, "null"); ok {
		l.pos = p1
		l.ret = Token{TypeNull, nil}
	} else {
		l.pos = p1
		l.err = badTokenError(l.str[p0:p1], p0)
	}
}

// Reset resets the internals for next token
func (l *Lexer) Reset() {
	l.state = stateDefault
	l.buf = []byte{}
	l.ret = Token{TypeNone, nil}
}

// Token gets the next JSON token
func (l *Lexer) Token() (Token, error) {
	l.Reset()
	for {
		if l.pos >= len(l.str) {
			// check EOF
			l.ret = Token{TypeEOF, nil}
		}
		if l.ret.Type != TypeNone || l.err != nil {
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
