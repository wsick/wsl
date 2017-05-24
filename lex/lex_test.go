package lex

import "testing"

type lexTest struct {
	name   string
	input  string
	tokens []Token
}

var lexTests = []lexTest{
	{"empty", "", []Token{tEOF}},
	{"spaces", " \t\r\n", []Token{mkToken(TokenWS, " \t\r\n"), tEOF}},
	{"single object", "application", []Token{mkToken(TokenObjectName, "application"), tEOF}},
	{"qualified object name", "local:application", []Token{mkToken(TokenObjectName, "local:application"), tEOF}},
	{"single object with whitespace", "  application  \n", []Token{
		mkToken(TokenWS, "  "),
		mkToken(TokenObjectName, "application"),
		mkToken(TokenWS, "  \n"),
		tEOF,
	}},
	{"empty content block", "application { }", []Token{
		mkToken(TokenObjectName, "application"),
		tSpace, tLeftCont, tSpace, tRightCont,
		tEOF,
	}},
	{"empty attribute block", "application [ ]", []Token{
		mkToken(TokenObjectName, "application"),
		tSpace, tLeftAttr, tSpace, tRightAttr,
		tEOF,
	}},
	{"empty object blocks", "application[]{}", []Token{
		mkToken(TokenObjectName, "application"),
		tLeftAttr, tRightAttr,
		tLeftCont, tRightCont,
		tEOF,
	}},
	{"attribute number", "grid [ grid.row = 1 ]", []Token{
		mkToken(TokenObjectName, "grid"),
		tSpace, tLeftAttr,
		tSpace, mkToken(TokenAttrKey, "grid.row"),
		tSpace, tEquals,
		tSpace, mkToken(TokenNumber, "1"),
		tSpace, tRightAttr, tEOF,
	}},
	{"attribute string", `text-block[text="label"]`, []Token{
		mkToken(TokenObjectName, "text-block"),
		tLeftAttr,
		mkToken(TokenAttrKey, "text"), tEquals, mkToken(TokenString, `"label"`),
		tRightAttr, tEOF,
	}},
	{"content number", "number{0}", []Token{
		mkToken(TokenObjectName, "number"),
		tLeftCont, mkToken(TokenNumber, "0"), tRightCont, tEOF,
	}},
	{"content string", `text-block{"text"}`, []Token{
		mkToken(TokenObjectName, "text-block"),
		tLeftCont, mkToken(TokenString, `"text"`), tRightCont, tEOF,
	}},
	{"multiline content", "text-block{`first\nsecond\nthird`}", []Token{
		mkToken(TokenObjectName, "text-block"),
		tLeftCont, mkToken(TokenMultiString, "`first\nsecond\nthird`"), tRightCont, tEOF,
	}},
	{"implicit attribute", `grid[grid.row-definitions{row-definition{"auto"}}]`, []Token{
		mkToken(TokenObjectName, "grid"),
		tLeftAttr, mkToken(TokenAttrKey, "grid.row-definitions"), tImplicit,
		tLeftCont, mkToken(TokenObjectName, "row-definition"),
		tLeftCont, mkToken(TokenString, `"auto"`), tRightCont,
		tRightCont,
		tRightAttr, tEOF,
	}},
	// errors
	{"object name dot", "grid.row", []Token{mkToken(TokenError, "unexpected character in object name")}},
	{"duplicate attr block", "application[][]", []Token{
		mkToken(TokenObjectName, "application"),
		tLeftAttr, tRightAttr,
		mkToken(TokenError, "duplicate attribute block found"),
	}},
	{"duplicate content block", "application{}{}", []Token{
		mkToken(TokenObjectName, "application"),
		tLeftCont, tRightCont,
		mkToken(TokenError, "duplicate content block found"),
	}},
}

func mkToken(typ TokenType, text string) Token {
	return Token{
		typ: typ,
		val: text,
	}
}

var (
	tEOF       = mkToken(TokenEOF, "")
	tLeftAttr  = mkToken(TokenLeftAttr, "[")
	tRightAttr = mkToken(TokenRightAttr, "]")
	tLeftCont  = mkToken(TokenLeftCont, "{")
	tRightCont = mkToken(TokenRightCont, "}")
	tEquals    = mkToken(TokenEquals, "=")
	tSpace     = mkToken(TokenWS, " ")
	tImplicit  = mkToken(TokenImplicitObject, "")
)

// collect gathers the emitted items into a slice.
func collect(t *lexTest) (items []Token) {
	l := FromString(t.name, t.input)
	for {
		token := l.NextToken()
		items = append(items, token)
		if token.typ == TokenEOF || token.typ == TokenError {
			break
		}
	}
	return
}

func equal(i1, i2 []Token, checkPos bool) bool {
	if len(i1) != len(i2) {
		return false
	}
	for k := range i1 {
		if i1[k].typ != i2[k].typ {
			return false
		}
		if i1[k].val != i2[k].val {
			return false
		}
		if checkPos && i1[k].pos != i2[k].pos {
			return false
		}
	}
	return true
}

func TestLex(t *testing.T) {
	for _, test := range lexTests {
		items := collect(&test)
		if !equal(items, test.tokens, false) {
			t.Errorf("%s: got\n\t%+v\nexpected\n\t%v", test.name, items, test.tokens)
		}
	}
}
