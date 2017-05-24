package lex

func beginContentBlock(l *Lexer) stateFn {
	if l.acceptWhitespace() {
		l.emit(TokenWS)
	}

	if !l.accept("{") {
		return l.errorf("expected '{'")
	}
	l.emit(TokenLeftCont)
	l.nest(TokenLeftCont)

	return value
}

// scan for content block end '}'
// pop off nested attribute
func endContentBlock(l *Lexer) stateFn {
	if l.acceptWhitespace() {
		l.emit(TokenWS)
	}

	if !l.accept("}") {
		return l.errorf("expected '}'")
	}
	l.emit(TokenRightCont)

	last := l.unnest()
	if last == TokenEquals {
		last = l.unnest()
	}

	if got := validateNest(TokenLeftCont, last); got != "" {
		return l.errorf("expected content termination, got %s", got)
	}

	return insideObject(true, true)
}
