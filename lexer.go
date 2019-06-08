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
	stateValue
	stateArray
	stateObject
	stateString
)

type token struct {
	Type tokenType
}

// Lexer reads and tokenizes a JSON string
type Lexer struct {
	str   string
	pos   int
	state lexerState
	err   error
}

func (l *Lexer) readDefault() {
	// check boundary
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
		// case number
		// case bool
		// case null
	default:
		l.err = newSyntaxError("unexpected character")
	}
}

// Token gets the next JSON token
func (l *Lexer) Token() (token, error) {
	l.state = stateDefault
	for {
		switch l.state {
		case stateDefault:
			l.readDefault()
		case stateComment:
			// read comment
		case stateValue:
			l.readValue()
		}
	}
}
