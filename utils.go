package json5

func expectLiteral(l *Lexer, exp string) bool {
	maxLen := len(l.str)
	for i := 0; i < len(exp); {
		if l.pos >= maxLen || l.str[l.pos] != exp[i] {
			return false
		}
		i++
		l.pos++
	}
	return true
}
