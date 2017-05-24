package lex

import (
	"strings"
	"unicode"
)

// scan for attribute value
//   - object (alphabetic)
//   - number (digit)
//   - string '"'
//   - multiline string '`'
//   - extension '{{' (not implemented yet)
func value(l *Lexer) stateFn {
	if l.acceptWhitespace() {
		l.emit(TokenWS)
	}

	switch r := l.next(); {
	case unicode.IsLetter(r):
		l.backup()
		return objectName(l)
	case unicode.IsDigit(r):
		l.backup()
		return number(l)
	case r == '"':
		return quote(l)
	case r == '`':
		return multilineQuote(l)
	case r == '}':
		l.backup()
		return endContentBlock(l)
	case r == '{':
		if strings.HasPrefix(l.input[l.pos-l.width:], "{{") {
			return extension(l)
		}
		return l.errorf("unexpected character in value")
	default:
		return l.errorf("unexpected character in value")
	}

	return nil
}

// This ends a value
// After a number, string, or extension, we need to discover our next state
// Objects will naturally transition
func endValue(l *Lexer) stateFn {
	if l.acceptWhitespace() {
		l.emit(TokenWS)
	}

	switch l.curNest() {
	case TokenLeftCont:
		return endContentBlock(l)
	case TokenEquals:
		l.unnest()
		fallthrough
	case TokenLeftAttr:
		return attrKey(l)
	default:
		return objectName(l)
	}
}

// lexNumber scans a number: decimal, octal, hex, float, or imaginary. This
// isn't a perfect number scanner - for instance it accepts "." and "0x0.2"
// and "089" - but when it's wrong the input is invalid and the parser (via
// strconv) will notice.
func number(l *Lexer) stateFn {
	if !l.scanNumber() {
		return l.errorf("bad number syntax: %q", l.input[l.start:l.pos])
	}
	if sign := l.peek(); sign == '+' || sign == '-' {
		// Complex: 1+2i. No spaces, must end in 'i'.
		if !l.scanNumber() || l.input[l.pos-1] != 'i' {
			return l.errorf("bad number syntax: %q", l.input[l.start:l.pos])
		}
		l.emit(TokenComplex)
	} else {
		l.emit(TokenNumber)
	}
	return endValue
}

func (l *Lexer) scanNumber() bool {
	// Optional leading sign.
	l.accept("+-")
	// Is it hex?
	digits := "0123456789"
	if l.accept("0") && l.accept("xX") {
		digits = "0123456789abcdefABCDEF"
	}
	l.acceptRun(digits)
	if l.accept(".") {
		l.acceptRun(digits)
	}
	if l.accept("eE") {
		l.accept("+-")
		l.acceptRun("0123456789")
	}
	// Is it imaginary?
	l.accept("i")
	// Next thing must not be alphanumeric.
	if isAlphaNumeric(l.peek()) {
		l.next()
		return false
	}
	return true
}

func quote(l *Lexer) stateFn {
Loop:
	for {
		switch l.next() {
		case '\\':
			if r := l.next(); r != eof && r != '\n' {
				break
			}
			fallthrough
		case eof, '\n':
			return l.errorf("unterminated quoted string")
		case '"':
			break Loop
		}
	}
	l.emit(TokenString)
	return endValue
}

func multilineQuote(l *Lexer) stateFn {
Loop:
	for {
		switch l.next() {
		case '\\':
			if r := l.next(); r != eof {
				break
			}
			fallthrough
		case eof:
			return l.errorf("unterminated quoted string")
		case '`':
			break Loop
		}
	}
	l.emit(TokenMultiString)
	return endValue
}

func extension(l *Lexer) stateFn {
	return l.errorf("extensions not implemented")
}
