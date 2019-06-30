package json5

import "errors"

type parserState int

const (
	stateStart parserState = iota
	stateArray
	stateObject
)

type Parser struct {
	state parserState
	stack stateStack
	lex   *Lexer
}

func NewParser(str string) *Parser {
	lexer := NewLexer(str)
	parser := &Parser{stack: stateStack{}, lex: lexer}
	lexer.ps = parser
	return parser
}

func (p *Parser) parseValue(tk Token) (err error) {
	switch tk.Type {
	case TypeArrayBegin:
		p.state = stateArray
		p.stack.Push(make([]interface{}, 0))
	case TypeObjectBegin:
		p.state = stateObject
		p.stack.Push(make(map[string]interface{}))
	case TypeString, TypeNumber, TypeBool, TypeNull:
		p.stack.Push(tk.Value)
	default:
		err = errors.New("unexpected token")
	}
	return
}

func (p *Parser) parseArray(tk Token) {
}

func (p *Parser) parseObject(tk Token) {
}

func (p *Parser) Parse() (value interface{}, err error) {
	for {
		tk, e := p.lex.Token()
		if e != nil {
			err = e
			return
		}
		if tk.Type == TypeEOF {
			value = p.stack.Top()
			return
		}
		switch p.state {
		case stateStart:
			p.parseValue(tk)
		case stateArray:
			p.parseArray(tk)
		case stateObject:
			p.parseObject(tk)
		}
	}
}
