package json5

import "strconv"

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

func parseDecimalInteger(s string) int {
	value, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		panic(err)
	}
	return int(value)
}
