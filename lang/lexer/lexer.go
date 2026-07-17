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
	return &Lexer{src: src, i: 0, line: 1, col: 1}
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
		if b == ' ' || b == '\t' || b == '\r' || b == '\n' {
			l.nextByte()
		} else {
			return
		}
	}
}

// Comments: ~ anything ~
func (l *Lexer) trySkipComment() bool {
	if l.peekByte() != '~' {
		return false
	}
	l.nextByte() // consume opening ~
	for {
		b := l.peekByte()
		if b == 0 {
			return false
		}
		if b == '~' {
			l.nextByte()
			return true
		}
		l.nextByte()
	}
}

func (l *Lexer) skipWhitespaceAndComments() {
	for {
		l.skipWhitespace()
		if !l.trySkipComment() {
			return
		}
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
	for b := l.peekByte(); b >= '0' && b <= '9'; b = l.peekByte() {
		sb.WriteByte(l.nextByte())
	}
	if l.peekByte() == '.' && l.i+1 < len(l.src) && l.src[l.i+1] >= '0' && l.src[l.i+1] <= '9' {
		sb.WriteByte(l.nextByte())
		for b := l.peekByte(); b >= '0' && b <= '9'; b = l.peekByte() {
			sb.WriteByte(l.nextByte())
		}
	}
	return token.Token{Type: token.NUMBER, Lexeme: sb.String(), Line: startLine, Column: startCol}
}

// readString reads a string literal.
// If the string contains {varName} interpolation, the Lexeme will contain
// the raw source including braces, and Type will be INTERP_STRING.
// The parser then splits it into segments.
func (l *Lexer) readString() (token.Token, error) {
	startLine, startCol := l.line, l.col
	l.nextByte() // consume opening "
	var sb strings.Builder
	hasInterp := false
	for {
		b := l.peekByte()
		if b == 0 {
			return token.Token{}, fmt.Errorf("unterminated string at %d:%d", startLine, startCol)
		}
		if b == '"' {
			l.nextByte()
			break
		}
		if b == '{' {
			hasInterp = true
			sb.WriteByte(l.nextByte()) // {
			// Scan until matching } — track inner strings and parens
			// so expressions like {fn(", ")} work correctly.
			innerStr := false
			parenDepth := 0
			for {
				c := l.peekByte()
				if c == 0 {
					return token.Token{}, fmt.Errorf("unclosed interpolation at %d:%d", startLine, startCol)
				}
				if c == '"' {
					innerStr = !innerStr
					sb.WriteByte(l.nextByte())
					continue
				}
				if innerStr {
					sb.WriteByte(l.nextByte())
					continue
				}
				if c == '(' {
					parenDepth++
				}
				if c == ')' && parenDepth > 0 {
					parenDepth--
				}
				sb.WriteByte(l.nextByte())
				if c == '}' && parenDepth == 0 {
					break
				}
			}
			continue
		}
		if b == '\\' {
			l.nextByte()
			esc := l.nextByte()
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
				sb.WriteByte(esc)
			}
			continue
		}
		sb.WriteByte(l.nextByte())
	}
	tt := token.STRING
	if hasInterp {
		tt = token.INTERP_STRING
	}
	return token.Token{Type: tt, Lexeme: sb.String(), Line: startLine, Column: startCol}, nil
}

// keyword map — maps lowercase word to token type
var keywords = map[string]token.TokenType{
	// Module
	"summon": token.K_SUMMON,
	"named":  token.K_NAMED,
	"as":     token.K_AS,

	// Object
	"create": token.K_CREATE,
	"make":   token.K_MAKE,
	"like":   token.K_LIKE,

	// Sections
	"has":       token.K_HAS,
	"owns":      token.K_OWNS,
	"carries":   token.K_CARRIES,
	"remembers": token.K_REMEMBERS,

	// Abilities
	"can":  token.K_CAN,
	"know": token.K_KNOW,
	"how":  token.K_HOW,

	// Assignment
	"is":      token.K_IS,
	"becomes": token.K_BECOMES,
	"set":     token.K_SET,
	"let":     token.K_LET,
	"be":      token.K_BE,
	"change":  token.K_CHANGE,
	"grow":    token.K_GROW,
	"shrink":  token.K_SHRINK,

	// Events
	"on":       token.K_ON,
	"when":     token.K_WHEN,
	"whenever": token.K_WHENEVER,

	// Control flow
	"if":     token.K_IF,
	"then":   token.K_THEN,
	"else":   token.K_ELSE,
	"repeat": token.K_REPEAT,
	"while":  token.K_WHILE,
	"times":  token.K_TIMES,
	"return": token.K_RETURN,
	"check":  token.K_CHECK,

	// Tell
	"tell": token.K_CAN_TELL,
	"to":   token.K_TO,

	// Logic
	"and": token.K_AND,
	"or":  token.K_OR,
	"not": token.K_NOT,

	// Wait
	"wait":    token.K_WAIT,
	"seconds": token.K_SECONDS,
	"frames":  token.K_FRAMES,

	// Output
	"print": token.K_PRINT,
	"say":   token.K_SAY,
	"show":  token.K_SHOW,

	// With
	"with": token.K_WITH,

	// Assets
	"image":  token.K_IMAGE,
	"sound":  token.K_SOUND,
	"music":  token.K_MUSIC,
	"video":  token.K_VIDEO,
	"font":   token.K_FONT,
	"shader": token.K_SHADER,

	// Special triggers
	"collision": token.K_COLLISION,
}

// booleanWords maps to boolean literals
var booleanWords = map[string]string{
	"true":  "true",
	"yes":   "true",
	"false": "false",
	"no":    "false",
}

func (l *Lexer) readIdentifierOrKeyword() token.Token {
	startLine, startCol := l.line, l.col
	var sb strings.Builder
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

	// Boolean literals
	if bval, ok := booleanWords[low]; ok {
		return token.Token{Type: token.BOOLEAN, Lexeme: bval, Line: startLine, Column: startCol}
	}

	// Keywords
	if tt, ok := keywords[low]; ok {
		return token.Token{Type: tt, Lexeme: lex, Line: startLine, Column: startCol}
	}

	return token.Token{Type: token.IDENT, Lexeme: lex, Line: startLine, Column: startCol}
}

// readProtectTag reads "#protected", "#protected1", "#protected2", etc.
// The leading '#' is consumed and not included in the Lexeme — the
// Lexeme is just the tag name (e.g. "protected1"), so the parser/runtime
// can use it directly as a policy-lookup key.
func (l *Lexer) readProtectTag() token.Token {
	startLine, startCol := l.line, l.col
	l.nextByte() // consume '#'

	var sb strings.Builder
	for {
		b := l.peekByte()
		if b == 0 || !isIdentChar(b) {
			break
		}
		sb.WriteByte(l.nextByte())
	}
	return token.Token{Type: token.PROTECT_TAG, Lexeme: sb.String(), Line: startLine, Column: startCol}
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

	// Protection tags: #protected, #protected1, #protected2, ...
	if b == '#' {
		return l.readProtectTag(), nil
	}

	// Two-char operators
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

func init() {
	// Register new keywords added after initial build
	keywords["for"]  = token.K_FOR
	keywords["each"] = token.K_EACH
	keywords["in"]   = token.K_IN
	keywords["ask"]  = token.K_ASK
	keywords["give"] = token.K_GIVE
	keywords["done"] = token.K_DONE
	keywords["stop"] = token.K_STOP
}

func init() {
	// Pass 2 keywords
	keywords["for"]  = token.K_FOR
	keywords["each"] = token.K_EACH
	keywords["in"]   = token.K_IN
	keywords["ask"]  = token.K_ASK
	keywords["give"] = token.K_GIVE
	keywords["done"] = token.K_DONE
	keywords["stop"] = token.K_STOP
	// Pass 3 keywords
	keywords["from"]  = token.K_FROM
	keywords["upto"]  = token.K_UPTO
	keywords["try"]   = token.K_TRY
	keywords["fails"] = token.K_FAILS
	keywords["step"]  = token.K_STEP
	keywords["until"]   = token.K_UNTIL
	// Pass 4 keywords
	keywords["between"] = token.K_BETWEEN
	keywords["ability"] = token.K_ABILITY
	keywords["remember"] = token.K_REMEMBER
	keywords["forget"]  = token.K_FORGET
	// Pass 5 keywords
	keywords["catch"]   = token.K_CATCH
	// Pass 6 keywords
	keywords["one"]     = token.K_ONE
	keywords["of"]      = token.K_OF
	keywords["fields"]  = token.K_FIELDS
}
