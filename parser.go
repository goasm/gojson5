package json5

type parserState int

const (
	stateArrayBegin parserState = iota
	stateArrayEnd
	stateObjectBegin
	stateObjectEnd
)

type Parser struct {
	state parserState
	stack *stateStack
	lex   *Lexer
}

func NewParser(str string) *Parser {
	lexer := NewLexer(str)
	parser := &Parser{stack: &stateStack{}, lex: lexer}
	lexer.ps = parser
	return parser
}

func (p *Parser) Parse() (value interface{}, err error) {
	for {
		tk, e := p.lex.Token()
		if e != nil {
			err = e
			return
		}
		switch tk.Type {
		case TypeEOF:
			return
		case TypeArrayBegin:
			p.state = stateArrayBegin
		case TypeObjectBegin:
			p.state = stateObjectBegin
		default:
			value = tk.Value
		}
	}
}
