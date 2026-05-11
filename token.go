package lexer

import "fmt"

// TokenType identifies what kind of token this is.
type TokenType int

const (
	// ── Literals ────────────────────────────────────────────────────────────
	TOKEN_NUMBER TokenType = iota // 42, 3.14
	TOKEN_STRING                  // "hello"
	TOKEN_IDENT                   // identifiers

	// ── Keywords: declarations ───────────────────────────────────────────────
	TOKEN_CLASS    // class
	TOKEN_OBJECT   // object  (instance creation: "object of <class>")
	TOKEN_OF       // of
	TOKEN_EXTENDS  // extends
	TOKEN_METHOD   // method
	TOKEN_RETURNS  // returns
	TOKEN_PROPERTY // property
	TOKEN_IS       // is  (assignment shorthand: "x is 5")

	// ── Keywords: control flow ───────────────────────────────────────────────
	TOKEN_IF       // if
	TOKEN_ELSE     // else
	TOKEN_WHILE    // while
	TOKEN_FOR      // for
	TOKEN_EACH     // each  (for each x in list)
	TOKEN_IN       // in
	TOKEN_BREAK    // break
	TOKEN_CONTINUE // continue
	TOKEN_RETURN   // return

	// ── Keywords: boolean / nil ──────────────────────────────────────────────
	TOKEN_TRUE  // true
	TOKEN_FALSE // false
	TOKEN_NIL   // nil / nothing
	TOKEN_NOT   // not

	// ── Keywords: logic ─────────────────────────────────────────────────────
	TOKEN_AND // and
	TOKEN_OR  // or

	// ── Keywords: I/O ───────────────────────────────────────────────────────
	TOKEN_PRINT  // print
	TOKEN_INPUT  // input
	TOKEN_AS     // as  (input x as "prompt")
	TOKEN_ANSWER // answer  (return alias)

	// ── Keywords: misc ──────────────────────────────────────────────────────
	TOKEN_DO       // do  (block opener without condition)
	TOKEN_END      // end  (block closer)
	TOKEN_THEN     // then
	TOKEN_SET      // set  (set x to 5)
	TOKEN_TO       // to
	TOKEN_WITH     // with
	TOKEN_CALL     // call
	TOKEN_ON       // on   (call method on object)
	TOKEN_NEW      // new  (new object of ...)
	TOKEN_THIS     // this / self
	TOKEN_SUPER    // super
	TOKEN_IMPORT   // import
	TOKEN_FROM     // from
	TOKEN_NOTHING  // nothing  (nil alias)
	TOKEN_HAS      // has  (x has property ...)
	TOKEN_INHERITS // inherits

	// ── Arithmetic operators ─────────────────────────────────────────────────
	TOKEN_PLUS     // +
	TOKEN_MINUS    // -
	TOKEN_STAR     // *
	TOKEN_SLASH    // /
	TOKEN_PERCENT  // %
	TOKEN_CARET    // ^  (exponentiation)
	TOKEN_PLUS_EQ  // +=
	TOKEN_MINUS_EQ // -=
	TOKEN_STAR_EQ  // *=
	TOKEN_SLASH_EQ // /=

	// ── Comparison operators ──────────────────────────────────────────────────
	TOKEN_EQ     // ==
	TOKEN_NEQ    // !=  / <>
	TOKEN_LT     // <
	TOKEN_GT     // >
	TOKEN_LTE    // <=
	TOKEN_GTE    // >=
	TOKEN_ASSIGN // =

	// ── Punctuation ──────────────────────────────────────────────────────────
	TOKEN_LPAREN    // (
	TOKEN_RPAREN    // )
	TOKEN_LBRACKET  // [
	TOKEN_RBRACKET  // ]
	TOKEN_LBRACE    // {
	TOKEN_RBRACE    // }
	TOKEN_COMMA     // ,
	TOKEN_DOT       // .
	TOKEN_COLON     // :
	TOKEN_SEMICOLON // ;
	TOKEN_NEWLINE   // \n (significant in this language)

	// ── Special ──────────────────────────────────────────────────────────────
	TOKEN_EOF
	TOKEN_ILLEGAL
)

var tokenNames = map[TokenType]string{
	TOKEN_NUMBER:    "NUMBER",
	TOKEN_STRING:    "STRING",
	TOKEN_IDENT:     "IDENT",
	TOKEN_CLASS:     "class",
	TOKEN_OBJECT:    "object",
	TOKEN_OF:        "of",
	TOKEN_EXTENDS:   "extends",
	TOKEN_METHOD:    "method",
	TOKEN_RETURNS:   "returns",
	TOKEN_PROPERTY:  "property",
	TOKEN_IS:        "is",
	TOKEN_IF:        "if",
	TOKEN_ELSE:      "else",
	TOKEN_WHILE:     "while",
	TOKEN_FOR:       "for",
	TOKEN_EACH:      "each",
	TOKEN_IN:        "in",
	TOKEN_BREAK:     "break",
	TOKEN_CONTINUE:  "continue",
	TOKEN_RETURN:    "return",
	TOKEN_TRUE:      "true",
	TOKEN_FALSE:     "false",
	TOKEN_NIL:       "nil",
	TOKEN_NOT:       "not",
	TOKEN_AND:       "and",
	TOKEN_OR:        "or",
	TOKEN_PRINT:     "print",
	TOKEN_INPUT:     "input",
	TOKEN_AS:        "as",
	TOKEN_ANSWER:    "answer",
	TOKEN_DO:        "do",
	TOKEN_END:       "end",
	TOKEN_THEN:      "then",
	TOKEN_SET:       "set",
	TOKEN_TO:        "to",
	TOKEN_WITH:      "with",
	TOKEN_CALL:      "call",
	TOKEN_ON:        "on",
	TOKEN_NEW:       "new",
	TOKEN_THIS:      "this",
	TOKEN_SUPER:     "super",
	TOKEN_IMPORT:    "import",
	TOKEN_FROM:      "from",
	TOKEN_NOTHING:   "nothing",
	TOKEN_HAS:       "has",
	TOKEN_INHERITS:  "inherits",
	TOKEN_PLUS:      "+",
	TOKEN_MINUS:     "-",
	TOKEN_STAR:      "*",
	TOKEN_SLASH:     "/",
	TOKEN_PERCENT:   "%",
	TOKEN_CARET:     "^",
	TOKEN_PLUS_EQ:   "+=",
	TOKEN_MINUS_EQ:  "-=",
	TOKEN_STAR_EQ:   "*=",
	TOKEN_SLASH_EQ:  "/=",
	TOKEN_EQ:        "==",
	TOKEN_NEQ:       "!=",
	TOKEN_LT:        "<",
	TOKEN_GT:        ">",
	TOKEN_LTE:       "<=",
	TOKEN_GTE:       ">=",
	TOKEN_ASSIGN:    "=",
	TOKEN_LPAREN:    "(",
	TOKEN_RPAREN:    ")",
	TOKEN_LBRACKET:  "[",
	TOKEN_RBRACKET:  "]",
	TOKEN_LBRACE:    "{",
	TOKEN_RBRACE:    "}",
	TOKEN_COMMA:     ",",
	TOKEN_DOT:       ".",
	TOKEN_COLON:     ":",
	TOKEN_SEMICOLON: ";",
	TOKEN_NEWLINE:   "NEWLINE",
	TOKEN_EOF:       "EOF",
	TOKEN_ILLEGAL:   "ILLEGAL",
}

func (t TokenType) String() string {
	if s, ok := tokenNames[t]; ok {
		return s
	}
	return fmt.Sprintf("TOKEN(%d)", int(t))
}

// keywords maps reserved words to their token type.
var keywords = map[string]TokenType{
	"class":    TOKEN_CLASS,
	"object":   TOKEN_OBJECT,
	"of":       TOKEN_OF,
	"extends":  TOKEN_EXTENDS,
	"method":   TOKEN_METHOD,
	"returns":  TOKEN_RETURNS,
	"property": TOKEN_PROPERTY,
	"is":       TOKEN_IS,
	"if":       TOKEN_IF,
	"else":     TOKEN_ELSE,
	"while":    TOKEN_WHILE,
	"for":      TOKEN_FOR,
	"each":     TOKEN_EACH,
	"in":       TOKEN_IN,
	"break":    TOKEN_BREAK,
	"continue": TOKEN_CONTINUE,
	"return":   TOKEN_RETURN,
	"true":     TOKEN_TRUE,
	"false":    TOKEN_FALSE,
	"nil":      TOKEN_NIL,
	"nothing":  TOKEN_NOTHING,
	"not":      TOKEN_NOT,
	"and":      TOKEN_AND,
	"or":       TOKEN_OR,
	"print":    TOKEN_PRINT,
	"input":    TOKEN_INPUT,
	"as":       TOKEN_AS,
	"answer":   TOKEN_ANSWER,
	"do":       TOKEN_DO,
	"end":      TOKEN_END,
	"then":     TOKEN_THEN,
	"set":      TOKEN_SET,
	"to":       TOKEN_TO,
	"with":     TOKEN_WITH,
	"call":     TOKEN_CALL,
	"on":       TOKEN_ON,
	"new":      TOKEN_NEW,
	"this":     TOKEN_THIS,
	"self":     TOKEN_THIS,
	"super":    TOKEN_SUPER,
	"import":   TOKEN_IMPORT,
	"from":     TOKEN_FROM,
	"has":      TOKEN_HAS,
	"inherits": TOKEN_INHERITS,
}

// LookupIdent checks whether an identifier is a reserved keyword.
func LookupIdent(ident string) TokenType {
	if t, ok := keywords[ident]; ok {
		return t
	}
	return TOKEN_IDENT
}

// Token is a single lexical unit.
type Token struct {
	Type    TokenType
	Literal string
	Line    int
	Col     int
}

func (t Token) String() string {
	return fmt.Sprintf("Token{%s, %q, line:%d col:%d}", t.Type, t.Literal, t.Line, t.Col)
}
