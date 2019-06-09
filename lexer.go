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
	stateSingleLineComment
	stateMultipleLineComment
	stateMultipleLineCommentEndAsterisk
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

func (l *Lexer) readComment() {
	c := l.str[l.pos]
	switch c {
	case '/':
		l.state = stateSingleLineComment
		l.pos++
	case '*':
		l.state = stateMultipleLineComment
		l.pos++
	default:
		l.err = badCharError(c, l.pos)
	}
}

func (l *Lexer) readSingleLineComment() {
	c := l.str[l.pos]
	switch c {
	case '\n', '\r':
		l.state = stateDefault
		l.pos++
	default:
		l.pos++
	}
}

func (l *Lexer) readMultipleLineComment() {
	c := l.str[l.pos]
	switch c {
	case '*':
		l.state = stateMultipleLineCommentEndAsterisk
		l.pos++
	default:
		l.pos++
	}
}

func (l *Lexer) readMultipleLineCommentEndAsterisk() {
	c := l.str[l.pos]
	switch c {
	case '/':
		l.state = stateDefault
		l.pos++
	default:
		l.pos++
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
	c := l.str[l.pos]
	switch c {
	case 'f':
		if expectLiteral(l, "false") {
			l.ret = Token{TypeBool, false}
			return
		}
	case 't':
		if expectLiteral(l, "true") {
			l.ret = Token{TypeBool, true}
			return
		}
	}
	l.err = badTokenError(l.str[p0:l.pos], p0)
}

func (l *Lexer) readNull() {
	p0 := l.pos
	if expectLiteral(l, "null") {
		l.ret = Token{TypeNull, nil}
		return
	}
	l.err = badTokenError(l.str[p0:l.pos], p0)
}

func (l *Lexer) checkEndState() error {
	switch l.state {
	case stateMultipleLineComment, stateMultipleLineCommentEndAsterisk:
		return badEOF(l.pos)
	default:
		return nil
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
			if err := l.checkEndState(); err != nil {
				// check state
				l.err = err
			} else {
				// exit normally
				l.ret = Token{TypeEOF, nil}
			}
		}
		if l.ret.Type != TypeNone || l.err != nil {
			// check result and error
			return l.ret, l.err
		}
		switch l.state {
		case stateDefault:
			l.readDefault()
		case stateComment:
			l.readComment()
		case stateSingleLineComment:
			l.readSingleLineComment()
		case stateMultipleLineComment:
			l.readMultipleLineComment()
		case stateMultipleLineCommentEndAsterisk:
			l.readMultipleLineCommentEndAsterisk()
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
		if l.ret.Type != TypeNone || l.err != nil {
			// check result and error
			return l.ret, l.err
		}
	}
}
