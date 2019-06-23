package json5

type parserState int

type Parser struct {
	lex *Lexer
}

func NewParser(str string) *Parser {
	lexer := NewLexer(str)
	parser := &Parser{lexer}
	lexer.ps = parser
	return parser
}

func (p *Parser) Parse() interface{} {
	return nil
}
