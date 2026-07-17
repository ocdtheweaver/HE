package token

type TokenType string

const (
	// Literals
	NUMBER  TokenType = "NUMBER"
	STRING  TokenType = "STRING"
	BOOLEAN       TokenType = "BOOLEAN"
	INTERP_STRING TokenType = "INTERP_STRING" // "Hello, {name}!"
	IDENT   TokenType = "IDENT"

	// Operators
	PLUS     TokenType = "+"
	MINUS    TokenType = "-"
	ASTERISK TokenType = "*"
	SLASH    TokenType = "/"
	CARET    TokenType = "^"
	POW      TokenType = "**"
	BANG     TokenType = "!"
	EQEQ     TokenType = "=="
	NEQ      TokenType = "!="
	GT       TokenType = ">"
	LT       TokenType = "<"
	GTE      TokenType = ">="
	LTE      TokenType = "<="

	// Delimiters
	LPAREN TokenType = "("
	RPAREN TokenType = ")"
	LBRACK TokenType = "["
	RBRACK TokenType = "]"
	LBRACE TokenType = "{"
	RBRACE TokenType = "}"
	COMMA  TokenType = ","
	COLON  TokenType = ":"
	DOT    TokenType = "."

	// Special
	EOF     TokenType = "EOF"
	ILLEGAL TokenType = "ILLEGAL"

	// ── Keywords ──────────────────────────────────────────────────────────

	// Module
	K_SUMMON TokenType = "summon"
	K_NAMED  TokenType = "named"
	K_AS     TokenType = "as"

	// Object creation
	K_CREATE TokenType = "create"
	K_MAKE   TokenType = "make"
	K_LIKE   TokenType = "like"

	// Property sections
	K_HAS      TokenType = "has"
	K_OWNS     TokenType = "owns"
	K_CARRIES  TokenType = "carries"
	K_REMEMBERS TokenType = "remembers"

	// Abilities
	K_CAN  TokenType = "can"
	K_KNOW TokenType = "know" // alias: "know how to"

	// Assignment / property
	K_IS      TokenType = "is"
	K_BECOMES TokenType = "becomes" // set X becomes Y (alias for set X to Y)

	// Events / reactions
	K_ON       TokenType = "on"
	K_WHEN     TokenType = "when"
	K_WHENEVER TokenType = "whenever"

	// Control flow
	K_IF     TokenType = "if"
	K_THEN   TokenType = "then"
	K_ELSE   TokenType = "else"
	K_REPEAT TokenType = "repeat"
	K_WHILE  TokenType = "while"
	K_TIMES  TokenType = "times"  // repeat N times
	K_RETURN TokenType = "return"

	// Tell / call
	K_CAN_TELL TokenType = "tell"
	K_TO       TokenType = "to"

	// Assignment keywords
	K_SET  TokenType = "set"
	K_LET  TokenType = "let"  // let X be Y
	K_BE   TokenType = "be"   // let X be Y

	// Change keywords (plain English mutation)
	K_CHANGE  TokenType = "change"  // change X to Y
	K_GROW    TokenType = "grow"    // grow X by N
	K_SHRINK  TokenType = "shrink"  // shrink X by N

	// Logic
	K_AND TokenType = "and"
	K_OR  TokenType = "or"
	K_NOT TokenType = "not" // "not X" as unary negation

	// Wait
	K_WAIT    TokenType = "wait"
	K_SECONDS TokenType = "seconds"
	K_FRAMES  TokenType = "frames"

	// Output
	K_PRINT TokenType = "print"
	K_SAY   TokenType = "say"
	K_SHOW  TokenType = "show" // alias for say/print

	// Check (alias for if — "check if X then")
	K_CHECK TokenType = "check"

	// With
	K_WITH TokenType = "with"

	// Asset types
	K_IMAGE  TokenType = "image"
	K_SOUND  TokenType = "sound"
	K_MUSIC  TokenType = "music"
	K_VIDEO  TokenType = "video"
	K_FONT   TokenType = "font"
	K_SHADER TokenType = "shader"

	// Collision (used as reaction trigger)
	K_COLLISION TokenType = "collision"

	// How (part of "knows how to")
	K_HOW TokenType = "how"

	// For-each loop
	K_FOR  TokenType = "for"
	K_EACH TokenType = "each"
	K_IN   TokenType = "in"

	// Ask / input
	K_ASK  TokenType = "ask"

	// Give (pass value to)
	K_GIVE TokenType = "give"

	// Done / stop (loop break)
	K_DONE TokenType = "done"
	K_STOP TokenType = "stop"
)

type Token struct {
	Type   TokenType
	Lexeme string
	Line   int
	Column int
}

const (
	// Numeric range loop: for each i from 1 to 10
	K_FROM  TokenType = "from"
	K_UPTO  TokenType = "upto"  // "up to" handled as two tokens: IDENT + K_TO

	// Try / error handling: try [...] or [...] if it fails
	K_TRY   TokenType = "try"
	K_FAILS TokenType = "fails" // "if it fails" → K_IF IDENT K_FAILS

	// String interpolation marker (internal — not a keyword users type)
	// Handled in lexer when it sees { inside a string

	// Step: for each i from 1 to 10 step 2
	K_STEP  TokenType = "step"

	// Until: repeat until cond
	K_UNTIL TokenType = "until"
)

const (
	// Range comparison: X is between A and B
	K_BETWEEN TokenType = "between"

	// Anonymous ability: ability(params) [body]
	K_ABILITY TokenType = "ability"

	// Persistence: remember X / forget X
	K_REMEMBER TokenType = "remember"
	K_FORGET   TokenType = "forget"

	// Multiple assignment: set a, b to expr
	// (comma already exists as COMMA)

	// Time module keyword (just a string summon, no keyword needed)

	// Compiler / type hints (optional annotations)
	K_AS_TYPE  TokenType = "as_type" // internal use only

	// Error capture: try [...] or (errVar) [...]
	K_CATCH TokenType = "catch" // or (e) — alias

	// Membership: X is one of [a, b, c]
	K_ONE  TokenType = "one"
	K_OF   TokenType = "of"

	// Scoped alias: with X as Y [...]
	K_WITH_AS TokenType = "with_as" // handled via existing K_WITH + K_AS

	// Object field iteration hint
	K_FIELDS TokenType = "fields"

	// Type enforcement
	K_TYPE   TokenType = "type"    // "type number" annotation keyword

	// Loop counter exposure: "repeat N times as i"
	// (handled via existing K_AS)

	// Becomes as mutation: "X becomes Y" standalone
	// K_BECOMES already exists — just needs parser case
)

const (
	// Protection tags: #protected, #protected1, #protected2, ...
	// Lexed as a single token: '#' followed immediately by an identifier
	// (letters/digits, no space). The tag name (e.g. "protected1") is
	// captured verbatim in the token's Lexeme for the parser to use.
	PROTECT_TAG TokenType = "PROTECT_TAG"
)
