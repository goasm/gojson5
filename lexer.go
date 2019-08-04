package json5

// TokenType represents an enum of token types
type TokenType int
type lexerState int

// Lexer token types
const (
	TypeNone TokenType = iota
	TypeArrayBegin
	TypeArrayEnd
	TypeObjectBegin
	TypeObjectEnd
	TypeValueSep
	TypePairSep
	TypeString
	TypeInteger
	TypeFloat
	TypeFalse
	TypeTrue
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
	statePunctuator
	stateString
	stateEscapeChar
	stateNumber
	stateUnsignedNumber
	stateZero
	stateDecimalInteger
	statePoint
	stateDecimalFraction
	stateDecimalExponent
	stateUnsignedDecimalExponent
	stateLiteral
)

// Token represents an unit for syntax analysis
type Token struct {
	Type TokenType
	Raw  string
}

// Lexer reads and tokenizes a JSON string
type Lexer struct {
	str   []byte
	pos   int
	state lexerState
	buf   stringBuffer
	ps    *Parser
}

// Scan creates a Lexer that scans the give string
func Scan(s string) *Lexer {
	return &Lexer{str: []byte(s)}
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
		l.state = stateMultipleLineComment
		l.pos++
	}
	return
}

func (l *Lexer) readValue(c byte) (tk Token, err error) {
	switch c {
	// TODO: only this case '[', '{':
	case '[', ']', '{', '}', ',', ':':
		l.state = statePunctuator
	case '"':
		l.state = stateString
		l.pos++
	case '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		l.state = stateNumber
	case 'f', 't', 'n':
		l.state = stateLiteral
	default:
		err = badCharError(c, l.pos)
	}
	return
}

func (l *Lexer) readPunctuator(c byte) (tk Token, err error) {
	switch c {
	case '[':
		tk = Token{TypeArrayBegin, "["}
	case ']':
		tk = Token{TypeArrayEnd, "]"}
	case '{':
		tk = Token{TypeObjectBegin, "{"}
	case '}':
		tk = Token{TypeObjectEnd, "}"}
	case ',':
		tk = Token{TypeValueSep, ","}
	case ':':
		tk = Token{TypePairSep, ":"}
	default:
		err = badCharError(c, l.pos)
		return
	}
	l.pos++
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
		value := l.buf.String()
		tk = Token{TypeString, value}
		l.pos++
	default:
		l.buf.Append(c)
		l.pos++
	}
	return
}

func (l *Lexer) readEscapeChar(c byte) (tk Token, err error) {
	var value byte
	switch c {
	case '"', '\\', '/':
		value = c
	case 'b':
		value = '\b'
	case 'f':
		value = '\f'
	case 'n':
		value = '\n'
	case 'r':
		value = '\r'
	case 't':
		value = '\t'
	case 'u':
		value = 0 // TODO: support unicode
	default:
		err = badCharError(c, l.pos)
		return
	}
	l.state = stateString
	l.buf.Append(value)
	l.pos++
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
		l.state = stateUnsignedNumber
		l.buf.Append(c)
		l.pos++
	default:
		tk, err = l.readUnsignedNumber(c)
	}
	return
}

func (l *Lexer) readUnsignedNumber(c byte) (tk Token, err error) {
	switch c {
	case '0':
		l.state = stateZero
		l.buf.Append(c)
		l.pos++
	case '1', '2', '3', '4', '5', '6', '7', '8', '9':
		l.state = stateDecimalInteger
		l.buf.Append(c)
		l.pos++
	default:
		err = badCharError(c, l.pos)
	}
	return
}

func (l *Lexer) readZero(c byte) (tk Token, err error) {
	switch c {
	case '.':
		l.state = statePoint
		l.buf.Append(c)
		l.pos++
	case 'e', 'E':
		l.state = stateDecimalExponent
		l.buf.Append(c)
		l.pos++
	// case 'x', 'X': TODO: support hexadecimal number
	default:
		tk = Token{TypeInteger, l.buf.String()}
	}
	return
}

func (l *Lexer) readDecimalInteger(c byte) (tk Token, err error) {
	switch c {
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		l.buf.Append(c)
		l.pos++
	case '.':
		l.state = statePoint
		l.buf.Append(c)
		l.pos++
	case 'e', 'E':
		l.state = stateDecimalExponent
		l.buf.Append(c)
		l.pos++
	default:
		tk = Token{TypeInteger, l.buf.String()}
	}
	return
}

func (l *Lexer) readPoint(c byte) (tk Token, err error) {
	switch c {
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		l.state = stateDecimalFraction
		l.buf.Append(c)
		l.pos++
	default:
		err = badCharError(c, l.pos)
	}
	return
}

func (l *Lexer) readDecimalFraction(c byte) (tk Token, err error) {
	switch c {
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		l.buf.Append(c)
		l.pos++
	case 'e', 'E':
		l.state = stateDecimalExponent
		l.buf.Append(c)
		l.pos++
	default:
		tk = Token{TypeFloat, l.buf.String()}
	}
	return
}

func (l *Lexer) readDecimalExponent(c byte) (tk Token, err error) {
	switch c {
	case '+', '-':
		l.state = stateUnsignedDecimalExponent
		l.buf.Append(c)
		l.pos++
	default:
		tk, err = l.readUnsignedDecimalExponent(c)
	}
	return
}

func (l *Lexer) readUnsignedDecimalExponent(c byte) (tk Token, err error) {
	switch c {
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		l.buf.Append(c)
		l.pos++
	default:
		tk = Token{TypeFloat, l.buf.String()}
	}
	return
}

// ================================================================
// }
// ================================================================

func (l *Lexer) readLiteral(c byte) (tk Token, err error) {
	p0 := l.pos
	switch c {
	case 'f':
		if expectLiteral(l, "false") {
			tk = Token{TypeFalse, "false"}
			return
		}
	case 't':
		if expectLiteral(l, "true") {
			tk = Token{TypeTrue, "true"}
			return
		}
	case 'n':
		if expectLiteral(l, "null") {
			tk = Token{TypeNull, "null"}
			return
		}
	}
	err = badTokenError(string(l.str[p0:l.pos]), p0)
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
	l.buf.Reset()
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
		case statePunctuator:
			tk, err = l.readPunctuator(c)
		case stateString:
			tk, err = l.readString(c)
		case stateEscapeChar:
			tk, err = l.readEscapeChar(c)
		case stateNumber:
			tk, err = l.readNumber(c)
		case stateUnsignedNumber:
			tk, err = l.readUnsignedNumber(c)
		case stateZero:
			tk, err = l.readZero(c)
		case stateDecimalInteger:
			tk, err = l.readDecimalInteger(c)
		case statePoint:
			tk, err = l.readPoint(c)
		case stateDecimalFraction:
			tk, err = l.readDecimalFraction(c)
		case stateDecimalExponent:
			tk, err = l.readDecimalExponent(c)
		case stateUnsignedDecimalExponent:
			tk, err = l.readUnsignedDecimalExponent(c)
		case stateLiteral:
			tk, err = l.readLiteral(c)
		}
		// check EOF
		if l.pos > len(l.str) {
			// check state
			if e := l.checkEndState(); e != nil {
				err = e
			} else {
				// exit normally
				tk = Token{TypeEOF, ""}
			}
		}
		// check result and error
		if tk.Type != TypeNone || err != nil {
			return
		}
	}
}
