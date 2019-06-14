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

func parseInteger(s string) (int64, error) {
	return strconv.ParseInt(s, 10, 64)
}

func parseFloat(s string) (float64, error) {
	return strconv.ParseFloat(s, 64)
}
