package json5

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
	Lexer
	state parserState
	stage stateStack
	paths nameStack
	stack valueStack
}

func (p *Parser) popValue() (err error) {
	p.state = p.stage.Pop()
	if p.stack.Size() == 1 {
		p.state = stateEnd
	} else {
		value := p.stack.Pop()
		switch p.state {
		case stateBeforeArrayItem:
			arr := p.stack.Top().([]interface{})
			p.stack.elements[p.stack.Size()-1] = append(arr, value)
			p.state = stateAfterArrayItem
		case stateBeforePropertyValue:
			name := p.paths.Pop()
			obj := p.stack.Top().(map[string]interface{})
			obj[name] = value
			p.state = stateAfterPropertyValue
		default:
			panic("unreachable")
		}
	}
	return
}

func (p *Parser) parseStart(tk Token) (err error) {
	switch tk.Type {
	case TypeArrayBegin:
		p.stage.Push(p.state)
		p.state = stateBeforeArrayItem
		p.stack.Push(make([]interface{}, 0))
	case TypeObjectBegin:
		p.stage.Push(p.state)
		p.state = stateBeforePropertyName
		p.stack.Push(make(map[string]interface{}))
	case TypeString, TypeInteger, TypeFloat, TypeFalse, TypeTrue, TypeNull:
		p.state = stateEnd
		value, e := parseToken(tk)
		if e != nil {
			err = e
			return
		}
		p.stack.Push(value)
	default:
		err = badTokenError(tk.Raw, p.pos)
	}
	return
}

func (p *Parser) parseBeforeArrayItem(tk Token) (err error) {
	switch tk.Type {
	case TypeArrayBegin:
		p.stage.Push(p.state)
		// p.state = stateBeforeArrayItem
		p.stack.Push(make([]interface{}, 0))
	case TypeObjectBegin:
		p.stage.Push(p.state)
		p.state = stateBeforePropertyName
		p.stack.Push(make(map[string]interface{}))
	case TypeString, TypeInteger, TypeFloat, TypeFalse, TypeTrue, TypeNull:
		p.state = stateAfterArrayItem
		value, e := parseToken(tk)
		if e != nil {
			err = e
			return
		}
		arr := p.stack.Top().([]interface{})
		p.stack.elements[p.stack.Size()-1] = append(arr, value)
	case TypeArrayEnd:
		err = p.popValue()
	default:
		err = badTokenError(tk.Raw, p.pos)
	}
	return
}

func (p *Parser) parseAfterArrayItem(tk Token) (err error) {
	switch tk.Type {
	case TypeValueSep:
		p.state = stateBeforeArrayItem
	case TypeArrayEnd:
		err = p.popValue()
	default:
		err = badTokenError(tk.Raw, p.pos)
	}
	return
}

func (p *Parser) parseBeforePropertyName(tk Token) (err error) {
	switch tk.Type {
	case TypeString:
		p.state = stateAfterPropertyName
		p.paths.Push(tk.Raw)
	case TypeObjectEnd:
		err = p.popValue()
	default:
		err = badTokenError(tk.Raw, p.pos)
	}
	return
}

func (p *Parser) parseAfterPropertyName(tk Token) (err error) {
	switch tk.Type {
	case TypePairSep:
		p.state = stateBeforePropertyValue
	default:
		err = badTokenError(tk.Raw, p.pos)
	}
	return
}

func (p *Parser) parseBeforePropertyValue(tk Token) (err error) {
	switch tk.Type {
	case TypeArrayBegin:
		p.stage.Push(p.state)
		p.state = stateBeforeArrayItem
		p.stack.Push(make([]interface{}, 0))
	case TypeObjectBegin:
		p.stage.Push(p.state)
		p.state = stateBeforePropertyName
		p.stack.Push(make(map[string]interface{}))
	case TypeString, TypeInteger, TypeFloat, TypeFalse, TypeTrue, TypeNull:
		p.state = stateAfterPropertyValue
		value, e := parseToken(tk)
		if e != nil {
			err = e
			return
		}
		name := p.paths.Pop()
		obj := p.stack.Top().(map[string]interface{})
		obj[name] = value
	default:
		err = badTokenError(tk.Raw, p.pos)
	}
	return
}

func (p *Parser) parseAfterPropertyValue(tk Token) (err error) {
	switch tk.Type {
	case TypeValueSep:
		p.state = stateBeforePropertyName
	case TypeObjectEnd:
		err = p.popValue()
	default:
		err = badTokenError(tk.Raw, p.pos)
	}
	return
}

func (p *Parser) parseEnd(tk Token) (err error) {
	if tk.Type != TypeEOF {
		err = badTokenError(tk.Raw, p.pos)
	}
	return
}

// Parse parses the JSON bytes
func (p *Parser) Parse(s []byte) (value interface{}, err error) {
	p.str = s
	for {
		tk, e := p.Token()
		if e != nil {
			err = e
			return
		}
		switch p.state {
		case stateStart:
			err = p.parseStart(tk)
		case stateBeforeArrayItem:
			err = p.parseBeforeArrayItem(tk)
		case stateAfterArrayItem:
			err = p.parseAfterArrayItem(tk)
		case stateBeforePropertyName:
			err = p.parseBeforePropertyName(tk)
		case stateAfterPropertyName:
			err = p.parseAfterPropertyName(tk)
		case stateBeforePropertyValue:
			err = p.parseBeforePropertyValue(tk)
		case stateAfterPropertyValue:
			err = p.parseAfterPropertyValue(tk)
		case stateEnd:
			err = p.parseEnd(tk)
			if err != nil {
				return
			}
			value = p.stack.Pop()
			return
		}
		if err != nil {
			return
		}
	}
}
