package lex

import "fmt"

var eof = rune(0)

type Token struct {
	typ  TokenType // The type of token
	pos  Pos       // The starting position, in bytes, of this token within input
	val  string    // The raw string of the token
	line int       // The line number at start of this Token
	col  int       // The column of this token within current line
}

func (t Token) String() string {
	switch {
	case t.typ == TokenEOF:
		return "EOF"
	case t.typ == TokenError:
		return t.val
	case t.typ == TokenImplicitObject:
		return "<implicit>"
	case len(t.val) > 10:
		return fmt.Sprintf("%.10q...", t.val)
	}
	return fmt.Sprintf("%q", t.val)
}

func (t Token) Type() TokenType {
	return t.typ
}

func (t Token) Pos() Pos {
	return t.pos
}

func (t Token) Val() string {
	return t.val
}

func (t Token) Line() int {
	return t.line
}

func (t Token) Col() int {
	return t.col
}

type TokenType int

const (
	TokenError TokenType = iota
	TokenEOF
	TokenWS

	// Literals
	TokenObjectName
	TokenImplicitObject
	TokenAttrKey
	TokenNumber
	TokenComplex
	TokenString
	TokenMultiString

	// Characters
	TokenLeftAttr  // [
	TokenRightAttr // ]
	TokenLeftCont  // {
	TokenRightCont // }
	TokenEquals    // =
)

var tokenChars = map[TokenType]string{
	TokenLeftAttr:  "[",
	TokenRightAttr: "]",
	TokenLeftCont:  "{",
	TokenRightCont: "}",
	TokenEquals:    "=",
}

func validateNest(want TokenType, got TokenType) string {
	if want == got {
		return ""
	}
	if ch, ok := tokenChars[got]; ok {
		return ch
	}
	return "eof"
}
