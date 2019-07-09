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

// stack
type stack struct {
	elements []interface{}
}

func (s *stack) Size() int {
	return len(s.elements)
}

func (s *stack) Top() interface{} {
	return s.elements[len(s.elements)-1]
}

func (s *stack) Push(e interface{}) {
	s.elements = append(s.elements, e)
}

func (s *stack) Pop() interface{} {
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
