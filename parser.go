package json5

type parserState int

type Parser struct {
	lex   *Lexer
	stack *stateStack
}

func NewParser(str string) *Parser {
	lexer := NewLexer(str)
	parser := &Parser{lexer, &stateStack{}}
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
		case TypeObjectBegin:
		default:
			value = tk.Value
		}
	}
}
