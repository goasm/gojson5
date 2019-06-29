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
	elements []interface{}
}

func (ss *stateStack) Size() int {
	return len(ss.elements)
}

func (ss *stateStack) Top() interface{} {
	return ss.elements[len(ss.elements)-1]
}

func (ss *stateStack) Push(e interface{}) {
	ss.elements = append(ss.elements, e)
}

func (ss *stateStack) Pop() interface{} {
	e := ss.Top()
	ss.elements = ss.elements[:len(ss.elements)-1]
	return e
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
