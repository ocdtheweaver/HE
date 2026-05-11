package token

type TokenType string

const (
	// Special
	ILLEGAL TokenType = "ILLEGAL"
	EOF     TokenType = "EOF"

	// Literals
	IDENT   TokenType = "IDENT"
	NUMBER  TokenType = "NUMBER"
	STRING  TokenType = "STRING"
	BOOLEAN TokenType = "BOOLEAN"

	// Keywords / special words
	K_SUMMON    TokenType = "summon"
	K_NAMED     TokenType = "named"
	K_AS        TokenType = "as"
	K_WITH      TokenType = "with"
	K_CREATE    TokenType = "create"
	K_MAKE      TokenType = "make"
	K_LIKE      TokenType = "like"
	K_HAS       TokenType = "has"
	K_OWNS      TokenType = "owns"
	K_CARRIES   TokenType = "carries"
	K_CAN       TokenType = "can"
	K_KNOWS     TokenType = "knows how to"
	K_REMEMBERS TokenType = "remembers"

	K_IS         TokenType = "is"
	K_STARTS_AS  TokenType = "starts as"
	K_REMEMBERS_ TokenType = "remembers:"

	K_ON        TokenType = "on"
	K_WHEN      TokenType = "when"
	K_WHENEVER  TokenType = "whenever"
	K_COLLISION TokenType = "collision"

	K_REMEMBERS_COLON TokenType = "remembers_colon"

	K_CAN_TELL TokenType = "tell"

	K_TO     TokenType = "to"
	K_OF     TokenType = "of"
	K_WITHIN TokenType = "with" // used for arg lists in calls; kept simple below

	K_AND TokenType = "and"
	K_OR  TokenType = "or"

	K_THEN TokenType = "then"
	K_ELSE TokenType = "else"

	K_REPEAT TokenType = "repeat"
	K_WHILE  TokenType = "while"

	K_RETURN  TokenType = "return"
	K_WAIT    TokenType = "wait"
	K_SECONDS TokenType = "seconds"
	K_FRAMES  TokenType = "frames"

	K_PRINT    TokenType = "print"
	K_SAY      TokenType = "say"
	K_SET      TokenType = "set"
	K_MAKE_SET TokenType = "make" // conflicts with object make; handled in parser by context

	// Asset types
	K_IMAGE  TokenType = "image"
	K_SOUND  TokenType = "sound"
	K_MUSIC  TokenType = "music"
	K_VIDEO  TokenType = "video"
	K_FONT   TokenType = "font"
	K_SHADER TokenType = "shader"

	// Operators / punctuation
	PLUS     TokenType = "+"
	MINUS    TokenType = "-"
	ASTERISK TokenType = "*"
	SLASH    TokenType = "/"
	CARET    TokenType = "^"
	POW      TokenType = "**"

	BANG TokenType = "!"
	EQEQ TokenType = "=="
	NEQ  TokenType = "!="
	GT   TokenType = ">"
	LT   TokenType = "<"
	GTE  TokenType = ">="
	LTE  TokenType = "<="

	ASSIGN_TO TokenType = "to" // for set? not used; left for parser
	LPAREN    TokenType = "("
	RPAREN    TokenType = ")"
	LBRACK    TokenType = "["
	RBRACK    TokenType = "]"
	LBRACE    TokenType = "{"
	RBRACE    TokenType = "}"

	COMMA TokenType = ","
	COLON TokenType = ":"
	DOT   TokenType = "."

	// Newline is mostly treated as whitespace by grammar (Line ::= ...), but we keep it for error messages.
	NEWLINE TokenType = "NEWLINE"
)

type Token struct {
	Type   TokenType
	Lexeme string
	Line   int
	Column int
}

func (t Token) String() string {
	return string(t.Type) + "('" + t.Lexeme + "')"
}
