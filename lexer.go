package json5

// TokenType represents an enum of token types
type TokenType int
type lexerState int

// Lexer token types
const (
	TypeNone TokenType = iota
	TypeString
	TypeNumber
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
	stateNumber
	stateSignedNumber
	stateDigitZero
	stateDecimalInteger
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
}

// NewLexer creates a JSON5 Lexer
func NewLexer(str string) *Lexer {
	return &Lexer{str: str}
}

func (l *Lexer) readDefault(c byte) (tk Token, err error) {
	switch c {
	case ' ', '\t', '\n', '\r':
		l.pos++
	case '/':
		l.state = stateComment
		l.pos++
	default:
		l.state = stateValue
	}
	return
}

func (l *Lexer) readComment(c byte) (tk Token, err error) {
	switch c {
	case '/':
		l.state = stateSingleLineComment
		l.pos++
	case '*':
		l.state = stateMultipleLineComment
		l.pos++
	default:
		err = badCharError(c, l.pos)
	}
	return
}

func (l *Lexer) readSingleLineComment(c byte) (tk Token, err error) {
	switch c {
	case '\n', '\r':
		l.state = stateDefault
		l.pos++
	default:
		l.pos++
	}
	return
}

func (l *Lexer) readMultipleLineComment(c byte) (tk Token, err error) {
	switch c {
	case '*':
		l.state = stateMultipleLineCommentEndAsterisk
		l.pos++
	default:
		l.pos++
	}
	return
}

func (l *Lexer) readMultipleLineCommentEndAsterisk(c byte) (tk Token, err error) {
	switch c {
	case '/':
		l.state = stateDefault
		l.pos++
	default:
		l.pos++
	}
	return
}

func (l *Lexer) readValue(c byte) (tk Token, err error) {
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
	case '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		l.state = stateNumber
	case 'f', 't':
		l.state = stateBool
	case 'n':
		l.state = stateNull
	default:
		err = badCharError(c, l.pos)
	}
	return
}

// ================================================================
// processing string {
// ================================================================

func (l *Lexer) readString(c byte) (tk Token, err error) {
	switch c {
	case '\\':
		l.state = stateEscapeChar
		l.pos++
	case '"':
		value := string(l.buf)
		tk = Token{TypeString, value}
		l.pos++
	default:
		l.buf = append(l.buf, c)
		l.pos++
	}
	return
}

// ================================================================
// }
// ================================================================

// ================================================================
// processing number {
// ================================================================

func (l *Lexer) readNumber(c byte) (tk Token, err error) {
	switch c {
	case '-':
		l.state = stateSignedNumber
		l.buf = append(l.buf, c)
		l.pos++
	case '0':
		l.state = stateDigitZero
		l.buf = append(l.buf, c)
		l.pos++
	case '1', '2', '3', '4', '5', '6', '7', '8', '9':
		l.state = stateDecimalInteger
		l.buf = append(l.buf, c)
		l.pos++
	default:
		err = badCharError(c, l.pos)
	}
	return
}

func (l *Lexer) readSignedNumber(c byte) (tk Token, err error) {
	switch c {
	case '0':
		l.state = stateDigitZero
		l.buf = append(l.buf, c)
		l.pos++
	case '1', '2', '3', '4', '5', '6', '7', '8', '9':
		l.state = stateDecimalInteger
		l.buf = append(l.buf, c)
		l.pos++
	default:
		err = badCharError(c, l.pos)
	}
	return
}

func (l *Lexer) readDigitZero(c byte) (tk Token, err error) {
	switch c {
	case '.':
		// TODO: float point
	default:
		// TODO: ???
	}
	return
}

func (l *Lexer) readDecimalInteger(c byte) (tk Token, err error) {
	switch c {
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		l.buf = append(l.buf, c)
		l.pos++
	case '.':
		// TODO: float point
	default:
		value := parseDecimalInteger(string(l.buf))
		tk = Token{TypeNumber, value}
	}
	return
}

// ================================================================
// }
// ================================================================

func (l *Lexer) readBool(c byte) (tk Token, err error) {
	p0 := l.pos
	switch c {
	case 'f':
		if expectLiteral(l, "false") {
			tk = Token{TypeBool, false}
			return
		}
	case 't':
		if expectLiteral(l, "true") {
			tk = Token{TypeBool, true}
			return
		}
	}
	err = badTokenError(l.str[p0:l.pos], p0)
	return
}

func (l *Lexer) readNull(c byte) (tk Token, err error) {
	p0 := l.pos
	if expectLiteral(l, "null") {
		tk = Token{TypeNull, nil}
		return
	}
	err = badTokenError(l.str[p0:l.pos], p0)
	return
}

// TODO: new state handler
// func (l *Lexer) readXxx() {}

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
}

// Token gets the next JSON token
func (l *Lexer) Token() (tk Token, err error) {
	l.Reset()
	for {
		var c byte
		if l.pos < len(l.str) {
			c = l.str[l.pos]
		} else {
			c = ' '
		}
		switch l.state {
		case stateDefault:
			tk, err = l.readDefault(c)
		case stateComment:
			tk, err = l.readComment(c)
		case stateSingleLineComment:
			tk, err = l.readSingleLineComment(c)
		case stateMultipleLineComment:
			tk, err = l.readMultipleLineComment(c)
		case stateMultipleLineCommentEndAsterisk:
			tk, err = l.readMultipleLineCommentEndAsterisk(c)
		case stateValue:
			tk, err = l.readValue(c)
		case stateArray:
			// TODO: read array
		case stateObject:
			// TODO: read object
		case stateString:
			tk, err = l.readString(c)
		case stateEscapeChar:
			// TODO: read escape char
		case stateNumber:
			tk, err = l.readNumber(c)
		case stateSignedNumber:
			tk, err = l.readSignedNumber(c)
		case stateDigitZero:
			tk, err = l.readDigitZero(c)
		case stateDecimalInteger:
			tk, err = l.readDecimalInteger(c)
		case stateBool:
			tk, err = l.readBool(c)
		case stateNull:
			tk, err = l.readNull(c)
		}
		// check EOF
		if l.pos > len(l.str) {
			// check state
			if e := l.checkEndState(); e != nil {
				err = e
			} else {
				// exit normally
				tk = Token{TypeEOF, nil}
			}
		}
		// check result and error
		if tk.Type != TypeNone || err != nil {
			return
		}
	}
}
