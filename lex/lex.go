package lex

func FromString(name, input string) *Lexer {
	l := &Lexer{
		name:     name,
		input:    input,
		tokens:   make(chan Token),
		line:     1,
		nestType: make([]TokenType, 0),
	}
	go l.run(objectName)
	return l
}
