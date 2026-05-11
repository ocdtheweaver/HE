package lexer

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

// TokenType identifies the kind of token.
type TokenType int

const (
	// Special
	ILLEGAL TokenType = iota
	EOF
	COMMENT

	// Identifiers + literals
	IDENT
	NUMBER
	STRING
	BOOLEAN

	// Keywords
	MODULE
	FN
	SUMMON
	AS
	NAMED
	WITH
	IMAGE
	SOUND
	MUSIC
	VIDEO
	FONT
	SHADER
	MAKE
	CREATE
	LIKE
	HAS
	OWNS
	CARRIES
	CAN
	KNOWS_HOW
	ON
	WHEN
	WHENEVER
	PRINT
	SAY
	TELL
	TO
	WAIT
	SECONDS
	FRAMES
	IS
	STARTS_AS
	RETURN
	IF
	THEN
	ELSE
	REPEAT
	WHILE

	// Operators / punctuation
	LPAREN  // (
	RPAREN  // )
	LBRACE  // {
	RBRACE  // }
	LBRACK  // [
	RBRACK  // ]
	COMMA   // ,
	COLON   // :
	EQ      // =
	EQEQ    // ==
	NOTEQ   // !=
	GT      // >
	LT      // <
	GTE     // >=
	LTE     // <=
	PLUS    // +
	MINUS   // -
	MULT    // *
	DIV     // /
	POWER   // ^ or **
	DOT     // .
	UNKNOWN // any unknown single char
)

// Token represents a lexical token.
type Token struct {
	Type TokenType
	Lit  string
	Pos  Position
}

// Position indicates a token position in source.
type Position struct {
	Offset int
	Line   int
	Col    int
}

func (t Token) String() string {
	return fmt.Sprintf("%s(%q)@%d:%d", t.Type.String(), t.Lit, t.Pos.Line, t.Pos.Col)
}

func (tt TokenType) String() string {
	names := map[TokenType]string{
		ILLEGAL:   "ILLEGAL",
		EOF:       "EOF",
		COMMENT:   "COMMENT",
		IDENT:     "IDENT",
		NUMBER:    "NUMBER",
		STRING:    "STRING",
		BOOLEAN:   "BOOLEAN",
		MODULE:    "MODULE",
		FN:        "FN",
		SUMMON:    "SUMMON",
		AS:        "AS",
		NAMED:     "NAMED",
		WITH:      "WITH",
		IMAGE:     "IMAGE",
		SOUND:     "SOUND",
		MUSIC:     "MUSIC",
		VIDEO:     "VIDEO",
		FONT:      "FONT",
		SHADER:    "SHADER",
		MAKE:      "MAKE",
		CREATE:    "CREATE",
		LIKE:      "LIKE",
		HAS:       "HAS",
		OWNS:      "OWNS",
		CARRIES:   "CARRIES",
		CAN:       "CAN",
		KNOWS_HOW: "KNOWS_HOW",
		ON:        "ON",
		WHEN:      "WHEN",
		WHENEVER:  "WHENEVER",
		PRINT:     "PRINT",
		SAY:       "SAY",
		TELL:      "TELL",
		TO:        "TO",
		WAIT:      "WAIT",
		SECONDS:   "SECONDS",
		FRAMES:    "FRAMES",
		IS:        "IS",
		STARTS_AS: "STARTS_AS",
		RETURN:    "RETURN",
		IF:        "IF",
		THEN:      "THEN",
		ELSE:      "ELSE",
		REPEAT:    "REPEAT",
		WHILE:     "WHILE",
		LPAREN:    "LPAREN",
		RPAREN:    "RPAREN",
		LBRACE:    "LBRACE",
		RBRACE:    "RBRACE",
		LBRACK:    "LBRACK",
		RBRACK:    "RBRACK",
		COMMA:     "COMMA",
		COLON:     "COLON",
		EQ:        "EQ",
		EQEQ:      "EQEQ",
		NOTEQ:     "NOTEQ",
		GT:        "GT",
		LT:        "LT",
		GTE:       "GTE",
		LTE:       "LTE",
		PLUS:      "PLUS",
		MINUS:     "MINUS",
		MULT:      "MULT",
		DIV:       "DIV",
		POWER:     "POWER",
		DOT:       "DOT",
		UNKNOWN:   "UNKNOWN",
	}
	if s, ok := names[tt]; ok {
		return s
	}
	return fmt.Sprintf("Token(%d)", int(tt))
}

// Lexer holds the state of the scanner.
type Lexer struct {
	input        string
	startOffset  int
	offset       int
	width        int
	line         int
	col          int
	inputLength  int
	keywords     map[string]TokenType
	caseFoldKeys bool
}

// New creates a new lexer for the given input string.
func New(input string) *Lexer {
	l := &Lexer{
		input:        input,
		offset:       0,
		width:        0,
		line:         1,
		col:          1,
		inputLength:  len(input),
		caseFoldKeys: true,
	}
	l.initKeywords()
	return l
}

func (l *Lexer) initKeywords() {
	// all HE keywords, case-insensitive
	l.keywords = map[string]TokenType{
		"module": MODULE, "fn": FN, "summon": SUMMON, "as": AS, "named": NAMED,
		"with": WITH, "image": IMAGE, "sound": SOUND, "music": MUSIC, "video": VIDEO,
		"font": FONT, "shader": SHADER, "make": MAKE, "create": CREATE, "like": LIKE,
		"has": HAS, "owns": OWNS, "carries": CARRIES, "can": CAN, "knows how": KNOWS_HOW,
		"on": ON, "when": WHEN, "whenever": WHENEVER, "print": PRINT, "say": SAY,
		"tell": TELL, "to": TO, "wait": WAIT, "seconds": SECONDS, "frames": FRAMES,
		"is": IS, "starts as": STARTS_AS, "return": RETURN, "if": IF, "then": THEN,
		"else": ELSE, "repeat": REPEAT, "while": WHILE,
		"true": BOOLEAN, "false": BOOLEAN, "yes": BOOLEAN, "no": BOOLEAN,
	}
}

// peek rune at current offset
func (l *Lexer) peek() rune {
	if l.offset >= l.inputLength {
		return 0
	}
	r, _ := utf8.DecodeRuneInString(l.input[l.offset:])
	return r
}

// read next rune
func (l *Lexer) next() rune {
	if l.offset >= l.inputLength {
		l.width = 0
		return 0
	}
	r, w := utf8.DecodeRuneInString(l.input[l.offset:])
	l.width = w
	l.offset += w
	if r == '\n' {
		l.line++
		l.col = 1
	} else {
		l.col++
	}
	return r
}

// backup one rune
func (l *Lexer) backup() {
	l.offset -= l.width
	l.width = 0
}

// current position
func (l *Lexer) currentPos() Position {
	return Position{Offset: l.startOffset, Line: l.line, Col: l.col}
}

// skipWhitespace skips spaces, tabs, newlines
func (l *Lexer) skipWhitespace() {
	for {
		r := l.peek()
		if r == 0 {
			return
		}
		if r == ' ' || r == '\t' || r == '\r' || r == '\n' {
			l.next()
			continue
		}
		return
	}
}

// Next returns the next token
func (l *Lexer) Next() Token {
	l.skipWhitespace()
	l.startOffset = l.offset
	r := l.peek()
	if r == 0 {
		return Token{Type: EOF, Lit: "", Pos: l.currentPos()}
	}

	// comments
	if r == '~' {
		l.next()
		for {
			n := l.next()
			if n == 0 {
				return Token{Type: ILLEGAL, Lit: "unterminated comment", Pos: l.currentPos()}
			}
			if n == '~' {
				l.startOffset = l.offset
				return l.Next()
			}
		}
	}

	// strings
	if r == '"' {
		l.next()
		for {
			n := l.next()
			if n == 0 {
				return Token{Type: ILLEGAL, Lit: "unterminated string", Pos: l.currentPos()}
			}
			if n == '\\' {
				l.next()
				continue
			}
			if n == '"' {
				tok := Token{Type: STRING, Lit: l.input[l.startOffset+1 : l.offset-1], Pos: l.currentPos()}
				l.startOffset = l.offset
				return tok
			}
		}
	}

	// numbers
	if unicode.IsDigit(r) || (r == '.' && l.offset+1 < l.inputLength && unicode.IsDigit(rune(l.input[l.offset+1]))) {
		hasDot := false
		for {
			c := l.peek()
			if c == '.' {
				if hasDot {
					break
				}
				hasDot = true
				l.next()
				continue
			}
			if unicode.IsDigit(c) {
				l.next()
				continue
			}
			break
		}
		lit := l.input[l.startOffset:l.offset]
		return Token{Type: NUMBER, Lit: lit, Pos: l.currentPos()}
	}

	// identifiers and keywords
	if unicode.IsLetter(r) || r == '_' {
		for {
			c := l.peek()
			if unicode.IsLetter(c) || unicode.IsDigit(c) || c == '_' || c == '.' || c == ' ' {
				l.next()
				continue
			}
			break
		}
		lit := strings.TrimSpace(l.input[l.startOffset:l.offset])
		key := lit
		if l.caseFoldKeys {
			key = strings.ToLower(key)
		}
		if typ, ok := l.keywords[key]; ok {
			if typ == BOOLEAN {
				return Token{Type: BOOLEAN, Lit: lit, Pos: l.currentPos()}
			}
			return Token{Type: typ, Lit: lit, Pos: l.currentPos()}
		}
		return Token{Type: IDENT, Lit: lit, Pos: l.currentPos()}
	}

	// operators & punctuation
	switch r {
	case '(':
		l.next()
		return Token{Type: LPAREN, Lit: "(", Pos: l.currentPos()}
	case ')':
		l.next()
		return Token{Type: RPAREN, Lit: ")", Pos: l.currentPos()}
	case '{':
		l.next()
		return Token{Type: LBRACE, Lit: "{", Pos: l.currentPos()}
	case '}':
		l.next()
		return Token{Type: RBRACE, Lit: "}", Pos: l.currentPos()}
	case '[':
		l.next()
		return Token{Type: LBRACK, Lit: "[", Pos: l.currentPos()}
	case ']':
		l.next()
		return Token{Type: RBRACK, Lit: "]", Pos: l.currentPos()}
	case ',':
		l.next()
		return Token{Type: COMMA, Lit: ",", Pos: l.currentPos()}
	case ':':
		l.next()
		return Token{Type: COLON, Lit: ":", Pos: l.currentPos()}
	case '+':
		l.next()
		return Token{Type: PLUS, Lit: "+", Pos: l.currentPos()}
	case '-':
		l.next()
		return Token{Type: MINUS, Lit: "-", Pos: l.currentPos()}
	case '*':
		l.next()
		if l.peek() == '*' {
			l.next()
			return Token{Type: POWER, Lit: "**", Pos: l.currentPos()}
		}
		return Token{Type: MULT, Lit: "*", Pos: l.currentPos()}
	case '/':
		l.next()
		return Token{Type: DIV, Lit: "/", Pos: l.currentPos()}
	case '=':
		l.next()
		if l.peek() == '=' {
			l.next()
			return Token{Type: EQEQ, Lit: "==", Pos: l.currentPos()}
		}
		return Token{Type: EQ, Lit: "=", Pos: l.currentPos()}
	case '!':
		l.next()
		if l.peek() == '=' {
			l.next()
			return Token{Type: NOTEQ, Lit: "!=", Pos: l.currentPos()}
		}
	case '>':
		l.next()
		if l.peek() == '=' {
			l.next()
			return Token{Type: GTE, Lit: ">=", Pos: l.currentPos()}
		}
		return Token{Type: GT, Lit: ">", Pos: l.currentPos()}
	case '<':
		l.next()
		if l.peek() == '=' {
			l.next()
			return Token{Type: LTE, Lit: "<=", Pos: l.currentPos()}
		}
		return Token{Type: LT, Lit: "<", Pos: l.currentPos()}
	case '^':
		l.next()
		return Token{Type: POWER, Lit: "^", Pos: l.currentPos()}
	case '.':
		l.next()
		return Token{Type: DOT, Lit: ".", Pos: l.currentPos()}
	}

	// unknown char
	l.next()
	return Token{Type: UNKNOWN, Lit: string(r), Pos: l.currentPos()}
}
