package json5

import "errors"

type parserState int

const (
	stateStart parserState = iota
	stateBeforeArrayItem
	stateAfterArrayItem
	stateBeforePropertyName
	stateAfterPropertyName
	stateBeforePropertyValue
	stateAfterPropertyValue
	stateEnd
)

// Parser represents a JSON5 parser
type Parser struct {
	state parserState
	stack stack
	cache pair
	value interface{}
}

func (p *Parser) parseStart(tk Token) (err error) {
	switch tk.Type {
	case TypeArrayBegin:
		p.state = stateBeforeArrayItem
		p.stack.Push(make([]interface{}, 0))
	case TypeObjectBegin:
		p.state = stateBeforePropertyName
		p.stack.Push(make(map[string]interface{}))
	case TypeString, TypeNumber, TypeBool, TypeNull:
		p.state = stateEnd
		p.value = tk.Value
	default:
		err = errors.New("unexpected token")
	}
	return
}

func (p *Parser) parseBeforeArrayItem(tk Token) (err error) {
	switch tk.Type {
	case TypeArrayBegin:
		p.stack.Push(make([]interface{}, 0))
	case TypeObjectBegin:
		p.state = stateBeforePropertyName
		p.stack.Push(make(map[string]interface{}))
	case TypeString, TypeNumber, TypeBool, TypeNull:
		p.state = stateAfterArrayItem
		arr := p.stack.Top().([]interface{})
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

func (p *Parser) parseAfterArrayItem(tk Token) (err error) {
	switch tk.Type {
	case TypeValueSep:
		p.state = stateBeforeArrayItem
	case TypeArrayEnd:
		// FIXME:(restore state)
		p.state = stateEnd
		p.value = p.stack.Pop()
	default:
		err = errors.New("unexpected token")
	}
	return
}

func (p *Parser) parseBeforePropertyName(tk Token) (err error) {
	switch tk.Type {
	case TypeString:
		p.state = stateAfterPropertyName
		p.cache.name = tk.Value.(string)
	case TypeObjectEnd:
		// FIXME:(restore state)
		p.state = stateEnd
		p.value = p.stack.Pop()
	default:
		err = errors.New("unexpected token")
	}
	return
}

func (p *Parser) parseAfterPropertyName(tk Token) (err error) {
	switch tk.Type {
	case TypePairSep:
		p.state = stateBeforePropertyValue
	default:
		err = errors.New("unexpected token")
	}
	return
}

func (p *Parser) parseBeforePropertyValue(tk Token) (err error) {
	switch tk.Type {
	case TypeArrayBegin:
		p.state = stateBeforeArrayItem
		// TODO:(save state)
	case TypeObjectBegin:
		p.state = stateBeforePropertyName
		// TODO:(save state)
	case TypeString, TypeNumber, TypeBool, TypeNull:
		p.state = stateAfterPropertyValue
		p.cache.value = tk.Value
		obj := p.stack.Top().(map[string]interface{})
		obj[p.cache.name] = p.cache.value
	default:
		err = errors.New("unexpected token")
	}
	return
}

func (p *Parser) parseAfterPropertyValue(tk Token) (err error) {
	switch tk.Type {
	case TypeValueSep:
		p.state = stateBeforePropertyName
	case TypeObjectEnd:
		// FIXME:(restore state)
		p.state = stateEnd
		p.value = p.stack.Pop()
	default:
		err = errors.New("unexpected token")
	}
	return
}

func (p *Parser) parseEnd(tk Token) (err error) {
	if tk.Type != TypeEOF {
		err = errors.New("unexpected token")
	}
	return
}

// Parse parses the JSON bytes
func (p *Parser) Parse(s string) (value interface{}, err error) {
	lexer := NewLexer(s)
	for {
		tk, e := lexer.Token()
		if e != nil {
			err = e
			return
		}
		switch p.state {
		case stateStart:
			p.parseStart(tk)
		case stateBeforeArrayItem:
			p.parseBeforeArrayItem(tk)
		case stateAfterArrayItem:
			p.parseAfterArrayItem(tk)
		case stateBeforePropertyName:
			p.parseBeforePropertyName(tk)
		case stateAfterPropertyName:
			p.parseAfterPropertyName(tk)
		case stateBeforePropertyValue:
			p.parseBeforePropertyValue(tk)
		case stateAfterPropertyValue:
			p.parseAfterPropertyValue(tk)
		case stateEnd:
			p.parseEnd(tk)
			value = p.value
			return
		}
	}
}
