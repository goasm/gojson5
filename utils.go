package json5

func expectLiteral(l *Lexer, exp string) (int, bool) {
	maxLen := len(l.str)
	p := l.pos
	for i := 0; i < len(exp); {
		if p >= maxLen {
			return p, false
		} else if l.str[p] != exp[i] {
			p++
			return p, false
		}
		i++
		p++
	}
	return p, true
}
