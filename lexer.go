// Package lexer tokenises Hunter's Engine source code.
//
// HE is a near-English OOP language. Its surface syntax looks like:
//
//	class Animal
//	    property name
//	    property sound
//
//	    method speak
//	        print this.name + " says " + this.sound
//	    end
//	end
//
//	object dog of Animal with name is "Rex", sound is "woof"
//	call speak on dog
//
// Significant newlines terminate most statements; blank lines and comment
// lines (-- ...) are discarded.  Indentation is advisory (not enforced by
// the grammar), but the lexer normalises it to make diagnostics cleaner.
package lexer

import (
	"fmt"
	"strings"
	"unicode"
)

// Error holds a lexer diagnostic.
type Error struct {
	Line    int
	Col     int
	Message string
}

func (e *Error) Error() string {
	return fmt.Sprintf("lexer error at line %d col %d: %s", e.Line, e.Col, e.Message)
}

// Lexer turns source text into a flat token slice.
type Lexer struct {
	src    []rune
	pos    int // current position in src
	line   int
	col    int
	errors []*Error
}

// New creates a Lexer from a source string.
func New(src string) *Lexer {
	return &Lexer{
		src:  []rune(src),
		line: 1,
		col:  1,
	}
}

// Errors returns all errors encountered during tokenisation.
func (l *Lexer) Errors() []*Error { return l.errors }

// Tokenise lexes the entire source and returns all tokens including EOF.
func (l *Lexer) Tokenise() []Token {
	var tokens []Token
	lastWasNewline := true // suppress leading blank lines

	for {
		tok := l.nextToken()
		switch tok.Type {
		case TOKEN_EOF:
			tokens = append(tokens, tok)
			return tokens
		case TOKEN_NEWLINE:
			if !lastWasNewline {
				tokens = append(tokens, tok)
				lastWasNewline = true
			}
			// collapse multiple consecutive newlines into one
		default:
			lastWasNewline = false
			tokens = append(tokens, tok)
		}
	}
}

// ─── internal helpers ────────────────────────────────────────────────────────

func (l *Lexer) peek() rune {
	if l.pos >= len(l.src) {
		return 0
	}
	return l.src[l.pos]
}

func (l *Lexer) peekAt(offset int) rune {
	idx := l.pos + offset
	if idx >= len(l.src) {
		return 0
	}
	return l.src[idx]
}

func (l *Lexer) advance() rune {
	if l.pos >= len(l.src) {
		return 0
	}
	ch := l.src[l.pos]
	l.pos++
	if ch == '\n' {
		l.line++
		l.col = 1
	} else {
		l.col++
	}
	return ch
}

func (l *Lexer) match(expected rune) bool {
	if l.peek() == expected {
		l.advance()
		return true
	}
	return false
}

func (l *Lexer) makeToken(t TokenType, lit string, line, col int) Token {
	return Token{Type: t, Literal: lit, Line: line, Col: col}
}

func (l *Lexer) addError(line, col int, msg string) {
	l.errors = append(l.errors, &Error{Line: line, Col: col, Message: msg})
}

// ─── core tokeniser ──────────────────────────────────────────────────────────

func (l *Lexer) nextToken() Token {
	// Skip horizontal whitespace and handle comments
	for l.pos < len(l.src) {
		ch := l.peek()
		switch {
		case ch == ' ' || ch == '\t' || ch == '\r':
			l.advance()
		case ch == '-' && l.peekAt(1) == '-':
			// line comment: skip until newline
			for l.pos < len(l.src) && l.peek() != '\n' {
				l.advance()
			}
		default:
			goto doneSkipping
		}
	}
doneSkipping:

	if l.pos >= len(l.src) {
		return l.makeToken(TOKEN_EOF, "", l.line, l.col)
	}

	line, col := l.line, l.col
	ch := l.advance()

	switch {
	// ── Newline ──────────────────────────────────────────────────────────────
	case ch == '\n':
		return l.makeToken(TOKEN_NEWLINE, "\n", line, col)

	// ── String literals ───────────────────────────────────────────────────────
	case ch == '"':
		return l.readString(line, col)

	// ── Numbers ───────────────────────────────────────────────────────────────
	case unicode.IsDigit(ch):
		return l.readNumber(ch, line, col)

	// ── Identifiers / keywords ────────────────────────────────────────────────
	case unicode.IsLetter(ch) || ch == '_':
		return l.readIdent(ch, line, col)

	// ── Operators ─────────────────────────────────────────────────────────────
	case ch == '+':
		if l.match('=') {
			return l.makeToken(TOKEN_PLUS_EQ, "+=", line, col)
		}
		return l.makeToken(TOKEN_PLUS, "+", line, col)
	case ch == '-':
		if l.match('=') {
			return l.makeToken(TOKEN_MINUS_EQ, "-=", line, col)
		}
		return l.makeToken(TOKEN_MINUS, "-", line, col)
	case ch == '*':
		if l.match('=') {
			return l.makeToken(TOKEN_STAR_EQ, "*=", line, col)
		}
		return l.makeToken(TOKEN_STAR, "*", line, col)
	case ch == '/':
		if l.match('=') {
			return l.makeToken(TOKEN_SLASH_EQ, "/=", line, col)
		}
		return l.makeToken(TOKEN_SLASH, "/", line, col)
	case ch == '%':
		return l.makeToken(TOKEN_PERCENT, "%", line, col)
	case ch == '^':
		return l.makeToken(TOKEN_CARET, "^", line, col)

	case ch == '=':
		if l.match('=') {
			return l.makeToken(TOKEN_EQ, "==", line, col)
		}
		return l.makeToken(TOKEN_ASSIGN, "=", line, col)
	case ch == '!':
		if l.match('=') {
			return l.makeToken(TOKEN_NEQ, "!=", line, col)
		}
		l.addError(line, col, "unexpected '!'; did you mean '!='?")
		return l.makeToken(TOKEN_ILLEGAL, "!", line, col)
	case ch == '<':
		if l.match('=') {
			return l.makeToken(TOKEN_LTE, "<=", line, col)
		}
		if l.match('>') {
			return l.makeToken(TOKEN_NEQ, "<>", line, col)
		}
		return l.makeToken(TOKEN_LT, "<", line, col)
	case ch == '>':
		if l.match('=') {
			return l.makeToken(TOKEN_GTE, ">=", line, col)
		}
		return l.makeToken(TOKEN_GT, ">", line, col)

	// ── Punctuation ───────────────────────────────────────────────────────────
	case ch == '(':
		return l.makeToken(TOKEN_LPAREN, "(", line, col)
	case ch == ')':
		return l.makeToken(TOKEN_RPAREN, ")", line, col)
	case ch == '[':
		return l.makeToken(TOKEN_LBRACKET, "[", line, col)
	case ch == ']':
		return l.makeToken(TOKEN_RBRACKET, "]", line, col)
	case ch == '{':
		return l.makeToken(TOKEN_LBRACE, "{", line, col)
	case ch == '}':
		return l.makeToken(TOKEN_RBRACE, "}", line, col)
	case ch == ',':
		return l.makeToken(TOKEN_COMMA, ",", line, col)
	case ch == '.':
		return l.makeToken(TOKEN_DOT, ".", line, col)
	case ch == ':':
		return l.makeToken(TOKEN_COLON, ":", line, col)
	case ch == ';':
		return l.makeToken(TOKEN_SEMICOLON, ";", line, col)
	}

	l.addError(line, col, fmt.Sprintf("unexpected character %q", ch))
	return l.makeToken(TOKEN_ILLEGAL, string(ch), line, col)
}

func (l *Lexer) readString(line, col int) Token {
	var sb strings.Builder
	for l.pos < len(l.src) {
		ch := l.advance()
		if ch == '"' {
			return l.makeToken(TOKEN_STRING, sb.String(), line, col)
		}
		if ch == '\\' {
			// escape sequences
			esc := l.advance()
			switch esc {
			case 'n':
				sb.WriteByte('\n')
			case 't':
				sb.WriteByte('\t')
			case '"':
				sb.WriteByte('"')
			case '\\':
				sb.WriteByte('\\')
			default:
				sb.WriteByte('\\')
				sb.WriteRune(esc)
			}
			continue
		}
		if ch == '\n' {
			l.addError(line, col, "unterminated string literal (newline in string)")
			return l.makeToken(TOKEN_STRING, sb.String(), line, col)
		}
		sb.WriteRune(ch)
	}
	l.addError(line, col, "unterminated string literal (reached EOF)")
	return l.makeToken(TOKEN_STRING, sb.String(), line, col)
}

func (l *Lexer) readNumber(first rune, line, col int) Token {
	var sb strings.Builder
	sb.WriteRune(first)
	hasDot := false
	for l.pos < len(l.src) {
		ch := l.peek()
		if unicode.IsDigit(ch) {
			sb.WriteRune(l.advance())
		} else if ch == '.' && !hasDot && unicode.IsDigit(l.peekAt(1)) {
			hasDot = true
			sb.WriteRune(l.advance()) // consume '.'
		} else {
			break
		}
	}
	return l.makeToken(TOKEN_NUMBER, sb.String(), line, col)
}

func (l *Lexer) readIdent(first rune, line, col int) Token {
	var sb strings.Builder
	sb.WriteRune(first)
	for l.pos < len(l.src) {
		ch := l.peek()
		if unicode.IsLetter(ch) || unicode.IsDigit(ch) || ch == '_' {
			sb.WriteRune(l.advance())
		} else {
			break
		}
	}
	lit := sb.String()
	tt := LookupIdent(lit)
	return l.makeToken(tt, lit, line, col)
}

