package lex

import "unicode"

// scan for object name
// after finding, we can pass off to lexPostObject to finish
func objectName(l *Lexer) stateFn {
	if l.acceptWhitespace() {
		l.emit(TokenWS)
	}

	// ensure the first character is alphabetic
	switch r := l.next(); {
	case unicode.IsLetter(r):
		// that's what we wanted, let lexer continue
	case r == eof:
		l.emit(TokenEOF)
		return nil
	default:
		return l.errorf("unexpected object name")
	}

	// chomp through object name
	for {
		switch r := l.next(); {
		case isWhitespace(r):
			l.backup()
			l.emit(TokenObjectName)
			return insideObject(false, false)
		case r == eof:
			l.emit(TokenObjectName)
			l.emit(TokenEOF)
			return nil
		case r == '[':
			l.backup()
			l.emit(TokenObjectName)
			return beginAttrBlock
		case r == '{':
			l.backup()
			l.emit(TokenObjectName)
			return beginContentBlock
		default:
			return l.errorf("unexpected character in object name")
		case isIdentifierRune(r):
			// keep chomping through identifier
		}
	}
}

// scanning within an object, can be after:
//   - object name
//   - after attribute block
//   - after content block
// we need to find our next state
//   - eof
//   - attribute block
//   - content block
//   - owning block termination
//   - new object
func insideObject(hitAttr, hitContent bool) stateFn {
	return func(l *Lexer) stateFn {
		if l.acceptWhitespace() {
			l.emit(TokenWS)
		}

		switch r := l.next(); {
		case r == eof:
			l.emit(TokenEOF)
			return nil
		case r == '[':
			if hitAttr {
				return l.errorf("duplicate attribute block found")
			}
			if hitContent {
				return l.errorf("attribute block must be listed before content block")
			}
			l.backup()
			return beginAttrBlock(l)
		case r == '{':
			if hitContent {
				return l.errorf("duplicate content block found")
			}
			l.backup()
			return beginContentBlock(l)
		case r == ']':
			l.backup()
			return endAttrBlock(l)
		case r == '}':
			l.backup()
			return endContentBlock(l)
		case unicode.IsLetter(r):
			l.backup()
			return objectName(l)
		default:
			return l.errorf("unexpected character after object name")
		}
	}
}
