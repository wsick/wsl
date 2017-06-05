package lex

import (
	"unicode"
)

// scan for attribute block begin '['
// push on nested attribute
func beginAttrBlock(l *Lexer) stateFn {
	if l.acceptWhitespace() {
		l.emit(TokenWS)
	}

	if !l.accept("[") {
		return l.errorf("expected '['")
	}
	l.emit(TokenLeftAttr)
	l.nest(TokenLeftAttr)

	return attrKey
}

// scan for attribute key
func attrKey(l *Lexer) stateFn {
	if l.acceptWhitespace() {
		l.emit(TokenWS)
	}

	// ensure the first character is alphabetic
	switch r := l.next(); {
	case unicode.IsLetter(r):
		// that's what we wanted, let lexer continue
	case r == ':':
		// this would imply a blank namespace alias, let lexer continue
	case r == eof:
		return l.errorf("unexpected end of file in attribute key")
	case r == ']':
		l.backup()
		return endAttrBlock(l)
	default:
		return l.errorf("unexpected attribute key")
	}

	// chomp through attribute key name
	for {
		switch r := l.next(); {
		case isAttributeTransition(r):
			l.backup()
			l.emit(TokenAttrKey)
			return attrValue(l)
		case r == eof:
			return l.errorf("unexpected end of file in attribute key")
		default:
			return l.errorf("unexpected character in attribute key")
		case isAttributeRune(r):
			// keep chomping through identifier
		}
	}
}

// scan for equals sign
func attrEquals(l *Lexer) stateFn {
	if l.acceptWhitespace() {
		l.emit(TokenWS)
	}

	if !l.accept("=") {
		return l.errorf("expected '='")
	}
	l.emit(TokenEquals)
	l.nest(TokenEquals)

	return attrValue
}

// scan for attribute value
// if we find '{' or '[', we are starting an implicit object
func attrValue(l *Lexer) stateFn {
	if l.acceptWhitespace() {
		l.emit(TokenWS)
	}

	switch l.peek() {
	case '=':
		return attrEquals(l)
	case '[':
		l.emit(TokenImplicitObject)
		return beginAttrBlock
	case '{':
		l.emit(TokenImplicitObject)
		return beginContentBlock
	default:
		return value(l)
	}
}

// scan for attribute block end ']'
// pop off nested attribute
func endAttrBlock(l *Lexer) stateFn {
	if l.acceptWhitespace() {
		l.emit(TokenWS)
	}

	if !l.accept("]") {
		return l.errorf("expected ']'")
	}
	l.emit(TokenRightAttr)

	last := l.unnest()
	if last == TokenEquals {
		last = l.unnest()
	}

	if got := validateNest(TokenLeftAttr, last); got != "" {
		return l.errorf("expected attribute termination, got %s", got)
	}

	return insideObject(true, false)
}
