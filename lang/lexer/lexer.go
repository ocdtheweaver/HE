package lexer

import (
	"fmt"
	"strings"
	"unicode"

	"hunterlang/lang/token"
)

type Lexer struct {
	src    string
	i      int
	line   int
	col    int
	looked *token.Token
}

func New(src string) *Lexer {
	return &Lexer{
		src:  src,
		i:    0,
		line: 1,
		col:  1,
	}
}

func (l *Lexer) peekByte() byte {
	if l.i >= len(l.src) {
		return 0
	}
	return l.src[l.i]
}

func (l *Lexer) nextByte() byte {
	if l.i >= len(l.src) {
		return 0
	}
	b := l.src[l.i]
	l.i++
	if b == '\n' {
		l.line++
		l.col = 1
	} else {
		l.col++
	}
	return b
}

func (l *Lexer) skipWhitespace() {
	for {
		b := l.peekByte()
		if b == 0 {
			return
		}
		if b == ' ' || b == '\t' || b == '\r' || b == '\n' {
			l.nextByte()
			continue
		}
		return
	}
}

// Comments: "~" ~"~"* "~"
// i.e. a tilde, then any number of tildes, then a final tilde.
func (l *Lexer) trySkipComment() bool {
	if l.peekByte() != '~' {
		return false
	}

	// Consume opening '~'
	l.nextByte()

	// Consume until the next '~' (allows arbitrary content in between).
	// This matches `~ ... ~` usage in game.he.
	for {
		b := l.peekByte()
		if b == 0 {
			// Unterminated comment; let lexer fail later.
			return false
		}
		if b == '~' {
			l.nextByte() // consume closing '~'
			return true
		}
		l.nextByte()
	}
}

func (l *Lexer) skipWhitespaceAndComments() {
	for {
		l.skipWhitespace()
		if l.trySkipComment() {
			continue
		}
		return
	}
}

func isLetter(b byte) bool {
	return (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z') || b == '_'
}

func isIdentChar(b byte) bool {
	return isLetter(b) || (b >= '0' && b <= '9')
}

func (l *Lexer) readNumber() token.Token {
	startLine, startCol := l.line, l.col
	var sb strings.Builder

	for {
		b := l.peekByte()
		if b >= '0' && b <= '9' {
			sb.WriteByte(l.nextByte())
			continue
		}
		break
	}

	// Optional fractional part: '.' Digit+
	if l.peekByte() == '.' {
		// Need at least one digit after '.'
		if l.i+1 < len(l.src) {
			next := l.src[l.i+1]
			if next >= '0' && next <= '9' {
				sb.WriteByte(l.nextByte()) // '.'
				for {
					b := l.peekByte()
					if b >= '0' && b <= '9' {
						sb.WriteByte(l.nextByte())
						continue
					}
					break
				}
			}
		}
	}

	return token.Token{
		Type:   token.NUMBER,
		Lexeme: sb.String(),
		Line:   startLine,
		Column: startCol,
	}
}

func (l *Lexer) readString() (token.Token, error) {
	startLine, startCol := l.line, l.col
	_ = l.nextByte() // consume opening '"'

	var sb strings.Builder
	for {
		b := l.peekByte()
		if b == 0 {
			return token.Token{}, fmt.Errorf("unterminated string at %d:%d", startLine, startCol)
		}
		if b == '"' {
			l.nextByte() // closing
			break
		}
		// Grammar: '"' [^"]* '"'
		// We'll allow any char except quote.
		sb.WriteByte(l.nextByte())
	}
	return token.Token{Type: token.STRING, Lexeme: sb.String(), Line: startLine, Column: startCol}, nil
}

func (l *Lexer) readIdentifierOrKeyword() token.Token {
	startLine, startCol := l.line, l.col
	var sb strings.Builder

	b := l.peekByte()
	if !isLetter(b) {
		// shouldn't happen
		l.nextByte()
		return token.Token{Type: token.ILLEGAL, Lexeme: string(b), Line: startLine, Column: startCol}
	}

	sb.WriteByte(l.nextByte())
	for {
		b := l.peekByte()
		if b == 0 || !isIdentChar(b) {
			break
		}
		sb.WriteByte(l.nextByte())
	}

	lex := sb.String()
	low := strings.ToLower(lex)

	// boolean keywords
	switch low {
	case "true", "yes":
		return token.Token{Type: token.BOOLEAN, Lexeme: "true", Line: startLine, Column: startCol}
	case "false", "no":
		return token.Token{Type: token.BOOLEAN, Lexeme: "false", Line: startLine, Column: startCol}
	}

	// multi-word keywords exist in the grammar; lexer maps single-word keywords only.
	switch low {
	case "summon":
		return token.Token{Type: token.K_SUMMON, Lexeme: lex, Line: startLine, Column: startCol}
	case "named":
		return token.Token{Type: token.K_NAMED, Lexeme: lex, Line: startLine, Column: startCol}
	case "as":
		return token.Token{Type: token.K_AS, Lexeme: lex, Line: startLine, Column: startCol}
	case "with":
		return token.Token{Type: token.K_WITH, Lexeme: lex, Line: startLine, Column: startCol}
	case "create":
		return token.Token{Type: token.K_CREATE, Lexeme: lex, Line: startLine, Column: startCol}
	case "make":
		return token.Token{Type: token.K_MAKE, Lexeme: lex, Line: startLine, Column: startCol}
	case "like":
		return token.Token{Type: token.K_LIKE, Lexeme: lex, Line: startLine, Column: startCol}
	case "has":
		return token.Token{Type: token.K_HAS, Lexeme: lex, Line: startLine, Column: startCol}
	case "owns":
		return token.Token{Type: token.K_OWNS, Lexeme: lex, Line: startLine, Column: startCol}
	case "carries":
		return token.Token{Type: token.K_CARRIES, Lexeme: lex, Line: startLine, Column: startCol}
	case "can":
		return token.Token{Type: token.K_CAN, Lexeme: lex, Line: startLine, Column: startCol}
	case "remembers":
		return token.Token{Type: token.K_REMEMBERS, Lexeme: lex, Line: startLine, Column: startCol}
	case "is":
		return token.Token{Type: token.K_IS, Lexeme: lex, Line: startLine, Column: startCol}
	case "starts":
		return token.Token{Type: token.IDENT, Lexeme: lex, Line: startLine, Column: startCol}
	case "on":
		return token.Token{Type: token.K_ON, Lexeme: lex, Line: startLine, Column: startCol}
	case "when":
		return token.Token{Type: token.K_WHEN, Lexeme: lex, Line: startLine, Column: startCol}
	case "whenever":
		return token.Token{Type: token.K_WHENEVER, Lexeme: lex, Line: startLine, Column: startCol}
	// NOTE: "collision" is treated as an IDENT so it can be used as a method name
	// (e.g. phys.collision(...)). Reaction triggers still work because the parser
	// accepts IDENT triggers too.
	case "tell":
		return token.Token{Type: token.K_CAN_TELL, Lexeme: lex, Line: startLine, Column: startCol}
	case "to":
		return token.Token{Type: token.K_TO, Lexeme: lex, Line: startLine, Column: startCol}
	case "and":
		return token.Token{Type: token.K_AND, Lexeme: lex, Line: startLine, Column: startCol}
	case "or":
		return token.Token{Type: token.K_OR, Lexeme: lex, Line: startLine, Column: startCol}
	case "then":
		return token.Token{Type: token.K_THEN, Lexeme: lex, Line: startLine, Column: startCol}
	case "else":
		return token.Token{Type: token.K_ELSE, Lexeme: lex, Line: startLine, Column: startCol}
	case "repeat":
		return token.Token{Type: token.K_REPEAT, Lexeme: lex, Line: startLine, Column: startCol}
	case "while":
		return token.Token{Type: token.K_WHILE, Lexeme: lex, Line: startLine, Column: startCol}
	case "return":
		return token.Token{Type: token.K_RETURN, Lexeme: lex, Line: startLine, Column: startCol}
	case "wait":
		return token.Token{Type: token.K_WAIT, Lexeme: lex, Line: startLine, Column: startCol}
	case "seconds":
		return token.Token{Type: token.K_SECONDS, Lexeme: lex, Line: startLine, Column: startCol}
	case "frames":
		return token.Token{Type: token.K_FRAMES, Lexeme: lex, Line: startLine, Column: startCol}
	case "print":
		return token.Token{Type: token.K_PRINT, Lexeme: lex, Line: startLine, Column: startCol}
	case "say":
		return token.Token{Type: token.K_SAY, Lexeme: lex, Line: startLine, Column: startCol}
	case "set":
		return token.Token{Type: token.K_SET, Lexeme: lex, Line: startLine, Column: startCol}
	case "image":
		return token.Token{Type: token.K_IMAGE, Lexeme: lex, Line: startLine, Column: startCol}
	case "sound":
		return token.Token{Type: token.K_SOUND, Lexeme: lex, Line: startLine, Column: startCol}
	case "music":
		return token.Token{Type: token.K_MUSIC, Lexeme: lex, Line: startLine, Column: startCol}
	case "video":
		return token.Token{Type: token.K_VIDEO, Lexeme: lex, Line: startLine, Column: startCol}
	case "font":
		return token.Token{Type: token.K_FONT, Lexeme: lex, Line: startLine, Column: startCol}
	case "shader":
		return token.Token{Type: token.K_SHADER, Lexeme: lex, Line: startLine, Column: startCol}
	}

	return token.Token{Type: token.IDENT, Lexeme: lex, Line: startLine, Column: startCol}
}

func (l *Lexer) NextToken() (token.Token, error) {
	if l.looked != nil {
		t := *l.looked
		l.looked = nil
		return t, nil
	}

	l.skipWhitespaceAndComments()

	startLine, startCol := l.line, l.col
	b := l.peekByte()
	if b == 0 {
		return token.Token{Type: token.EOF, Lexeme: "", Line: startLine, Column: startCol}, nil
	}

	if b == '"' {
		return l.readString()
	}

	if b >= '0' && b <= '9' {
		return l.readNumber(), nil
	}

	if unicode.IsLetter(rune(b)) || b == '_' {
		return l.readIdentifierOrKeyword(), nil
	}

	switch b {
	case '=':
		l.nextByte()
		if l.peekByte() == '=' {
			l.nextByte()
			return token.Token{Type: token.EQEQ, Lexeme: "==", Line: startLine, Column: startCol}, nil
		}
		return token.Token{Type: token.ILLEGAL, Lexeme: "=", Line: startLine, Column: startCol}, nil
	case '!':
		l.nextByte()
		if l.peekByte() == '=' {
			l.nextByte()
			return token.Token{Type: token.NEQ, Lexeme: "!=", Line: startLine, Column: startCol}, nil
		}
		l.nextByte()
		return token.Token{Type: token.BANG, Lexeme: "!", Line: startLine, Column: startCol}, nil
	case '>':
		l.nextByte()
		if l.peekByte() == '=' {
			l.nextByte()
			return token.Token{Type: token.GTE, Lexeme: ">=", Line: startLine, Column: startCol}, nil
		}
		return token.Token{Type: token.GT, Lexeme: ">", Line: startLine, Column: startCol}, nil
	case '<':
		l.nextByte()
		if l.peekByte() == '=' {
			l.nextByte()
			return token.Token{Type: token.LTE, Lexeme: "<=", Line: startLine, Column: startCol}, nil
		}
		return token.Token{Type: token.LT, Lexeme: "<", Line: startLine, Column: startCol}, nil
	case '*':
		l.nextByte()
		if l.peekByte() == '*' {
			l.nextByte()
			return token.Token{Type: token.POW, Lexeme: "**", Line: startLine, Column: startCol}, nil
		}
		return token.Token{Type: token.ASTERISK, Lexeme: "*", Line: startLine, Column: startCol}, nil
	}

	l.nextByte()
	switch b {
	case '+':
		return token.Token{Type: token.PLUS, Lexeme: "+", Line: startLine, Column: startCol}, nil
	case '-':
		return token.Token{Type: token.MINUS, Lexeme: "-", Line: startLine, Column: startCol}, nil
	case '/':
		return token.Token{Type: token.SLASH, Lexeme: "/", Line: startLine, Column: startCol}, nil
	case '^':
		return token.Token{Type: token.CARET, Lexeme: "^", Line: startLine, Column: startCol}, nil
	case '!':
		return token.Token{Type: token.BANG, Lexeme: "!", Line: startLine, Column: startCol}, nil
	case '(':
		return token.Token{Type: token.LPAREN, Lexeme: "(", Line: startLine, Column: startCol}, nil
	case ')':
		return token.Token{Type: token.RPAREN, Lexeme: ")", Line: startLine, Column: startCol}, nil
	case '[':
		return token.Token{Type: token.LBRACK, Lexeme: "[", Line: startLine, Column: startCol}, nil
	case ']':
		return token.Token{Type: token.RBRACK, Lexeme: "]", Line: startLine, Column: startCol}, nil
	case '{':
		return token.Token{Type: token.LBRACE, Lexeme: "{", Line: startLine, Column: startCol}, nil
	case '}':
		return token.Token{Type: token.RBRACE, Lexeme: "}", Line: startLine, Column: startCol}, nil
	case ',':
		return token.Token{Type: token.COMMA, Lexeme: ",", Line: startLine, Column: startCol}, nil
	case ':':
		return token.Token{Type: token.COLON, Lexeme: ":", Line: startLine, Column: startCol}, nil
	case '.':
		return token.Token{Type: token.DOT, Lexeme: ".", Line: startLine, Column: startCol}, nil
	}

	return token.Token{Type: token.ILLEGAL, Lexeme: string(b), Line: startLine, Column: startCol}, nil
}

func (l *Lexer) Peek() (token.Token, error) {
	if l.looked != nil {
		return *l.looked, nil
	}
	t, err := l.NextToken()
	if err != nil {
		return token.Token{}, err
	}
	l.looked = &t
	return t, nil
}

func (l *Lexer) Expect(tt token.TokenType) (token.Token, error) {
	t, err := l.NextToken()
	if err != nil {
		return token.Token{}, err
	}
	if t.Type != tt {
		return token.Token{}, fmt.Errorf("expected %s at %d:%d, got %s (%q)", tt, t.Line, t.Column, t.Type, t.Lexeme)
	}
	return t, nil
}
