package lex

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

// stateFn represents the state of the scanner as a function that returns the next state.
type stateFn func(*Lexer) stateFn

// Lexer holds the state of the scanner.
type Lexer struct {
	name     string      // the name of the input; used only for error reports
	input    string      // the string being scanned
	state    stateFn     // the next lexing function to enter
	pos      Pos         // current position in the input
	start    Pos         // start position of this item
	width    Pos         // width of last rune read from input
	lastPos  Pos         // position of most recent item returned by NextToken
	startCol int         // start column of this token
	col      int         // current rune's column
	lastCol  int         // last rune's column
	tokens   chan Token  // channel of scanned items
	line     int         // 1+number of newlines seen
	nestType []TokenType // nested type (content or attribute)
}

// next returns the next rune in the input.
func (l *Lexer) next() rune {
	if int(l.pos) >= len(l.input) {
		l.width = 0
		return eof
	}
	r, w := utf8.DecodeRuneInString(l.input[l.pos:])
	l.width = Pos(w)
	l.pos += l.width
	l.lastCol = l.col
	l.col++
	if r == '\n' {
		l.col = 0
		l.line++
	}
	return r
}

// peek returns but does not consume the next rune in the input.
func (l *Lexer) peek() rune {
	r := l.next()
	l.backup()
	return r
}

// backup steps back one rune. Can only be called once per call of next.
func (l *Lexer) backup() {
	l.pos -= l.width
	l.col = l.lastCol
	// Correct newline count.
	if l.width == 1 && l.input[l.pos] == '\n' {
		l.line--
	}
}

// emit passes an item back to the client.
func (l *Lexer) emit(t TokenType) {
	l.tokens <- Token{t, l.start, l.input[l.start:l.pos], l.line, l.startCol}
	l.start = l.pos
	l.startCol = l.col
}

// ignore skips over the pending input before this point.
func (l *Lexer) ignore() {
	l.start = l.pos
}

// accept consumes the next rune if it's from the valid set.
func (l *Lexer) accept(valid string) bool {
	if strings.ContainsRune(valid, l.next()) {
		return true
	}
	l.backup()
	return false
}

// acceptRun consumes a run of runes from the valid set.
func (l *Lexer) acceptRun(valid string) {
	for strings.ContainsRune(valid, l.next()) {
	}
	l.backup()
}

// acceptsWhitespace consumes a run of whitespace
func (l *Lexer) acceptWhitespace() bool {
	flag := false
	for strings.ContainsRune(" \t\r\n", l.next()) {
		flag = true
	}
	l.backup()
	return flag
}

func (l *Lexer) nest(typ TokenType) {
	l.nestType = append(l.nestType, typ)
}

// unnest a nested block and returns this block type
func (l *Lexer) unnest() TokenType {
	lastType := TokenEOF
	if len(l.nestType) > 0 {
		lastType = l.nestType[len(l.nestType)-1]
	}
	l.nestType = l.nestType[:len(l.nestType)-1]
	return lastType
}

// curNest returns the active nest block type
func (l *Lexer) curNest() TokenType {
	if len(l.nestType) > 0 {
		return l.nestType[len(l.nestType)-1]
	}
	return TokenEOF
}

// errorf returns an error token and terminates the scan by passing
// back a nil pointer that will be the next state, terminating l.NextToken.
func (l *Lexer) errorf(format string, args ...interface{}) stateFn {
	l.tokens <- Token{TokenError, l.start, fmt.Sprintf(format, args...), l.line, l.startCol}
	return nil
}

// NextToken returns the next item from the input.
// Called by the parser, not in the lexing goroutine.
func (l *Lexer) NextToken() Token {
	item := <-l.tokens
	l.lastPos = item.pos
	return item
}

// drain drains the output so the lexing goroutine will exit.
// Called by the parser, not in the lexing goroutine.
func (l *Lexer) drain() {
	for range l.tokens {
	}
}

// run runs the state machine for the lexer.
func (l *Lexer) run(startingLexFn stateFn) {
	for l.state = startingLexFn; l.state != nil; {
		l.state = l.state(l)
	}
	close(l.tokens)
}

// isWhitespace reports whether r is whitespace
func isWhitespace(r rune) bool {
	return r == ' ' || r == '\t' || r == '\r' || r == '\n'
}

// isAlphaNumeric reports whether r is alphabetic or numeric
func isAlphaNumeric(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r)
}

// isIdentifierRune reports whether r matches valid characters for an identifier
func isIdentifierRune(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r) || strings.ContainsRune(":_-", r)
}

// isAttributeRune reports whether r matches valid characters for an identifier
func isAttributeRune(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r) || strings.ContainsRune(":._-", r)
}

// isAttributeTransition reports whether r is a transition character between key and value
func isAttributeTransition(r rune) bool {
	return isWhitespace(r) || r == '=' || r == '{' || r == '['
}
