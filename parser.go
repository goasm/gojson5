package json5

import "errors"

type parserState int

const (
	stateStart parserState = iota
	stateBeforeArrayValue
	stateAfterArrayValue
	stateObject
	stateEnd
)

type Parser struct {
	state parserState
	stack stateStack
	value interface{}
	lex   *Lexer
}

func NewParser(str string) *Parser {
	lexer := NewLexer(str)
	parser := &Parser{stack: stateStack{}, lex: lexer}
	lexer.ps = parser
	return parser
}

func (p *Parser) parseStart(tk Token) (err error) {
	switch tk.Type {
	case TypeArrayBegin:
		p.state = stateBeforeArrayValue
		p.stack.Push(make([]interface{}, 0))
	case TypeObjectBegin:
		p.state = stateObject
	case TypeString, TypeNumber, TypeBool, TypeNull:
		p.state = stateEnd
		p.stack.Push(tk.Value)
	default:
		err = errors.New("unexpected token")
	}
	return
}

func (p *Parser) parseBeforeArrayValue(tk Token) (err error) {
	switch tk.Type {
	case TypeArrayBegin:
		p.stack.Push(make([]interface{}, 0))
	case TypeObjectBegin:
		p.state = stateObject
	case TypeString, TypeNumber, TypeBool, TypeNull:
		p.state = stateAfterArrayValue
		arr, _ := p.stack.Top().([]interface{})
		p.stack.elements[p.stack.Size()-1] = append(arr, tk.Value)
	case TypeArrayEnd:
		// FIXME:(restore state)
		p.state = stateEnd
		p.value = p.stack.Pop()
	default:
		err = errors.New("unexpected token")
	}
	return
}

func (p *Parser) parseAfterArrayValue(tk Token) (err error) {
	switch tk.Type {
	case TypeValueSep:
		p.state = stateBeforeArrayValue
	case TypeArrayEnd:
		// FIXME:(restore state)
		p.state = stateEnd
		p.value = p.stack.Pop()
	default:
		err = errors.New("unexpected token")
	}
	return
}

func (p *Parser) parseObject(tk Token) (err error) {
	p.stack.Push(make(map[string]interface{}))
	return
}

func (p *Parser) parseEnd(tk Token) (err error) {
	if tk.Type != TypeEOF {
		err = errors.New("unexpected token")
	}
	return
}

func (p *Parser) Parse() (value interface{}, err error) {
	for {
		tk, e := p.lex.Token()
		if e != nil {
			err = e
			return
		}
		switch p.state {
		case stateStart:
			p.parseStart(tk)
		case stateBeforeArrayValue:
			p.parseBeforeArrayValue(tk)
		case stateAfterArrayValue:
			p.parseAfterArrayValue(tk)
		case stateObject:
			p.parseObject(tk)
		case stateEnd:
			p.parseEnd(tk)
			value = p.value
			return
		}
	}
}
