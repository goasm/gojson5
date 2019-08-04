package json5

import (
	"bytes"
	"strconv"
)

// stringBuffer
type stringBuffer struct {
	buf bytes.Buffer
}

func (sb *stringBuffer) Append(c byte) {
	sb.buf.WriteByte(c)
}

func (sb *stringBuffer) Reset() {
	sb.buf.Reset()
}

func (sb *stringBuffer) String() string {
	return sb.buf.String()
}

// stateStack
type stateStack struct {
	elements []parserState
}

func (s *stateStack) Size() int {
	return len(s.elements)
}

func (s *stateStack) Top() parserState {
	return s.elements[len(s.elements)-1]
}

func (s *stateStack) Push(e parserState) {
	s.elements = append(s.elements, e)
}

func (s *stateStack) Pop() parserState {
	e := s.Top()
	s.elements = s.elements[:len(s.elements)-1]
	return e
}

// nameStack
type nameStack struct {
	elements []string
}

func (s *nameStack) Size() int {
	return len(s.elements)
}

func (s *nameStack) Top() string {
	return s.elements[len(s.elements)-1]
}

func (s *nameStack) Push(e string) {
	s.elements = append(s.elements, e)
}

func (s *nameStack) Pop() string {
	e := s.Top()
	s.elements = s.elements[:len(s.elements)-1]
	return e
}

// valueStack
type valueStack struct {
	elements []interface{}
}

func (s *valueStack) Size() int {
	return len(s.elements)
}

func (s *valueStack) Top() interface{} {
	return s.elements[len(s.elements)-1]
}

func (s *valueStack) Push(e interface{}) {
	s.elements = append(s.elements, e)
}

func (s *valueStack) Pop() interface{} {
	e := s.Top()
	s.elements = s.elements[:len(s.elements)-1]
	return e
}

// pair
type pair struct {
	name  string
	value interface{}
}

// helper functions

func expectLiteral(l *Lexer, expected string) bool {
	maxLen := len(l.str) - l.pos
	i := 0
	for ; i < len(expected) && i < maxLen; i++ {
		equal := l.str[l.pos] == expected[i]
		l.pos++
		if !equal {
			return false
		}
	}
	return i == len(expected)
}

func parseInteger(s string) (int64, error) {
	return strconv.ParseInt(s, 10, 64)
}

func parseFloat(s string) (float64, error) {
	return strconv.ParseFloat(s, 64)
}

func parseToken(tk Token) (interface{}, error) {
	switch tk.Type {
	case TypeString:
		return tk.Raw, nil
	case TypeInteger:
		return parseInteger(tk.Raw)
	case TypeFloat:
		return parseFloat(tk.Raw)
	case TypeFalse:
		return false, nil
	case TypeTrue:
		return true, nil
	case TypeNull:
		return nil, nil
	}
	panic("unreachable")
}
