he/
 ├── cmd/he              (CLI: he run, he build)
 ├── pkg/lexer           (tokenizer)
 ├── pkg/parser          (nested modules + new grammar)
 ├── pkg/ast             (module nodes, fn nodes)
 ├── pkg/runtime         (module loader + executor)
 ├── pkg/resolver        (platform-aware fn resolution)
 ├── pkg/builtins        (he_mobile.he, he_web.he, physics.he)
 ├── go.mod



 I'll create the entire HE v2 project structure using PowerShell commands for VSCode terminal. This will set up the Go workspace and directory structure.

```powershell
# First, let's create the root directory for the HE project
New-Item -ItemType Directory -Path .\he -Force

# Change to the project directory
cd .\he

# Initialize Go module
go mod init he

# Create the complete project structure
@(
    "cmd\he",
    "pkg\lexer",
    "pkg\parser", 
    "pkg\ast",
    "pkg\runtime",
    "pkg\resolver",
    "pkg\builtins",
    "pkg\types",
    "internal\compiler",
    "internal\loader",
    "examples",
    "testdata"
) | ForEach-Object {
    New-Item -ItemType Directory -Path $_ -Force
}

# Create main.go for the CLI tool
@'
package main

import (
    "fmt"
    "os"
)

func main() {
    if len(os.Args) < 2 {
        fmt.Println("HE Language v2")
        fmt.Println("Usage:")
        fmt.Println("  he run <file.he>    - Run a HE program")
        fmt.Println("  he build <file.he>  - Build for target platform")
        fmt.Println("  he fmt <file.he>    - Format HE code")
        fmt.Println("  he help             - Show this help")
        os.Exit(1)
    }

    command := os.Args[1]
    switch command {
    case "run":
        if len(os.Args) < 3 {
            fmt.Println("Error: missing file argument for run")
            os.Exit(1)
        }
        runFile(os.Args[2])
    case "build":
        if len(os.Args) < 3 {
            fmt.Println("Error: missing file argument for build")
            os.Exit(1)
        }
        buildFile(os.Args[2])
    case "fmt":
        if len(os.Args) < 3 {
            fmt.Println("Error: missing file argument for fmt")
            os.Exit(1)
        }
        fmtFile(os.Args[2])
    case "help", "-h", "--help":
        printHelp()
    default:
        fmt.Printf("Unknown command: %s\n", command)
        printHelp()
        os.Exit(1)
    }
}

func runFile(filename string) {
    fmt.Printf("Running %s (not implemented yet)\n", filename)
    // TODO: Load, parse, and execute HE file
}

func buildFile(filename string) {
    fmt.Printf("Building %s (not implemented yet)\n", filename)
    // TODO: Compile HE file to target platform
}

func fmtFile(filename string) {
    fmt.Printf("Formatting %s (not implemented yet)\n", filename)
    // TODO: Format HE code
}

func printHelp() {
    fmt.Println(`
HE Language v2 - A friendly programming language for creators

Commands:
  run <file.he>    - Execute a HE program
  build <file.he>  - Build for target platform (web/mobile/desktop)
  fmt <file.he>    - Format HE source code
  help             - Show this help message

Examples:
  he run examples/hello.he
  he build game.he --target=web
  he fmt myscript.he

Learn more: https://github.com/he-lang/he
`)
}
'@ | Out-File -FilePath "cmd\he\main.go" -Encoding UTF8

# Create the lexer file (Chunk 1)
# Copy the lexer code from the previous message into pkg/lexer/lexer.go
@'
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
	IDENT  // foo, bar, hero.name
	NUMBER // 123, 3.14
	STRING // "hello"

	// Keywords (a subset; parser can accept more if needed)
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
	CAN
	ON
	WHEN
	PRINT
	SAY
	TELL
	TO
	WAIT
	SECONDS
	IS
	RETURN

	// Operators / punctuation
	LPAREN   // (
	RPAREN   // )
	LBRACE   // {
	RBRACE   // }
	LBRACK   // [
	RBRACK   // ]
	COMMA    // ,
	COLON    // :
	EQ       // =
	PLUS     // +
	MINUS    // -
	MULT     // *
	DIV      // /
	DOT      // .
	UNKNOWN  // any unknown single char
)

// Token represents a lexical token.
type Token struct {
	Type TokenType
	Lit  string
	Pos  Position
}

// Position indicates a token position in source.
type Position struct {
	Offset int // byte offset from start
	Line   int // 1-based
	Col    int // 1-based (rune column)
}

func (t Token) String() string {
	return fmt.Sprintf("%s(%q)@%d:%d", t.Type.String(), t.Lit, t.Pos.Line, t.Pos.Col)
}

// String representation of TokenType.
func (tt TokenType) String() string {
	switch tt {
	case ILLEGAL:
		return "ILLEGAL"
	case EOF:
		return "EOF"
	case COMMENT:
		return "COMMENT"
	case IDENT:
		return "IDENT"
	case NUMBER:
		return "NUMBER"
	case STRING:
		return "STRING"
	case MODULE:
		return "MODULE"
	case FN:
		return "FN"
	case SUMMON:
		return "SUMMON"
	case AS:
		return "AS"
	case NAMED:
		return "NAMED"
	case WITH:
		return "WITH"
	case IMAGE:
		return "IMAGE"
	case SOUND:
		return "SOUND"
	case MUSIC:
		return "MUSIC"
	case VIDEO:
		return "VIDEO"
	case FONT:
		return "FONT"
	case SHADER:
		return "SHADER"
	case MAKE:
		return "MAKE"
	case CREATE:
		return "CREATE"
	case LIKE:
		return "LIKE"
	case HAS:
		return "HAS"
	case CAN:
		return "CAN"
	case ON:
		return "ON"
	case WHEN:
		return "WHEN"
	case PRINT:
		return "PRINT"
	case SAY:
		return "SAY"
	case TELL:
		return "TELL"
	case TO:
		return "TO"
	case WAIT:
		return "WAIT"
	case SECONDS:
		return "SECONDS"
	case IS:
		return "IS"
	case RETURN:
		return "RETURN"
	case LPAREN:
		return "LPAREN"
	case RPAREN:
		return "RPAREN"
	case LBRACE:
		return "LBRACE"
	case RBRACE:
		return "RBRACE"
	case LBRACK:
		return "LBRACK"
	case RBRACK:
		return "RBRACK"
	case COMMA:
		return "COMMA"
	case COLON:
		return "COLON"
	case EQ:
		return "EQ"
	case PLUS:
		return "PLUS"
	case MINUS:
		return "MINUS"
	case MULT:
		return "MULT"
	case DIV:
		return "DIV"
	case DOT:
		return "DOT"
	case UNKNOWN:
		return "UNKNOWN"
	default:
		return "Token(" + fmt.Sprint(int(tt)) + ")"
	}
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
		input:       input,
		offset:      0,
		width:       0,
		line:        1,
		col:         1,
		inputLength: len(input),
		caseFoldKeys: true, // keywords are case-insensitive
	}
	l.initKeywords()
	return l
}

func (l *Lexer) initKeywords() {
	l.keywords = map[string]TokenType{
		"module": MODULE,
		"fn":     FN,
		"summon": SUMMON,
		"as":     AS,
		"named":  NAMED,
		"with":   WITH,
		"image":  IMAGE,
		"sound":  SOUND,
		"music":  MUSIC,
		"video":  VIDEO,
		"font":   FONT,
		"shader": SHADER,
		"make":   MAKE,
		"create": CREATE,
		"like":   LIKE,
		"has":    HAS,
		"can":    CAN,
		"on":     ON,
		"when":   WHEN,
		"print":  PRINT,
		"say":    SAY,
		"tell":   TELL,
		"to":     TO,
		"wait":   WAIT,
		"seconds": SECONDS,
		"is":     IS,
		"return": RETURN,
	}
}

// helper: peek rune at current offset without advancing
func (l *Lexer) peek() rune {
	if l.offset >= l.inputLength {
		return 0
	}
	r, _ := utf8.DecodeRuneInString(l.input[l.offset:])
	return r
}

// helper: read next rune and advance
func (l *Lexer) next() rune {
	if l.offset >= l.inputLength {
		l.width = 0
		return 0
	}
	r, w := utf8.DecodeRuneInString(l.input[l.offset:])
	l.width = w
	l.offset += l.width
	if r == '\n' {
		l.line++
		l.col = 1
	} else {
		l.col++
	}
	return r
}

// helper: backup one rune
func (l *Lexer) backup() {
	if l.width == 0 {
		return
	}
	l.offset -= l.width
	// we won't accurately restore line/col here for backups except in simple cases;
	// parser/lexer usage keeps this limited. For robust position tracking, tokens record start positions.
	l.width = 0
}

// emitToken constructs a token from startOffset..current offset
func (l *Lexer) emitToken(typ TokenType) Token {
	start := l.startOffset
	end := l.offset
	lit := l.input[start:end]
	pos := Position{Offset: start, Line: l.line, Col: l.col - utf8.RuneCountInString(lit)}
	return Token{Type: typ, Lit: lit, Pos: pos}
}

// skipWhitespace advances past spaces, tabs, newlines
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

// Next returns the next token from input.
func (l *Lexer) Next() Token {
	l.skipWhitespace()
	l.startOffset = l.offset
	r := l.peek()
	if r == 0 {
		return Token{Type: EOF, Lit: "", Pos: Position{Offset: l.offset, Line: l.line, Col: l.col}}
	}

	// comments: only ~ ... ~ (single-line or inline). We will skip entire comment and return Next()
	if r == '~' {
		// consume opening ~
		l.next()
		// find next ~
		for {
			n := l.next()
			if n == 0 {
				// unterminated comment -> return ILLEGAL token with the partial text
				return Token{Type: ILLEGAL, Lit: "unterminated comment", Pos: Position{Offset: l.startOffset, Line: l.line, Col: l.col}}
			}
			if n == '~' {
				// comment ended; skip and continue lexing the next token
				l.startOffset = l.offset
				return l.Next()
			}
			// continue inside comment
		}
	}

	// Strings
	if r == '"' {
		// consume opening quote
		l.next()
		for {
			n := l.next()
			if n == 0 {
				return Token{Type: ILLEGAL, Lit: "unterminated string", Pos: Position{Offset: l.startOffset, Line: l.line, Col: l.col}}
			}
			if n == '\\' {
				// escape: skip next rune
				_ = l.next()
				continue
			}
			if n == '"' {
				// emit string token (without quotes as Lit includes quotes; parser may strip or use raw)
				// We emit with quotes included - parser/evaluator can strip or we can remove here
				tok := Token{Type: STRING, Lit: l.input[l.startOffset+1 : l.offset-1], Pos: Position{Offset: l.startOffset, Line: l.line, Col: l.col}}
				l.startOffset = l.offset
				return tok
			}
		}
	}

	// Numbers: digit starts or dot followed by digit
	if unicode.IsDigit(r) || (r == '.' && l.offset+1 < l.inputLength && unicode.IsDigit(rune(l.input[l.offset+1]))) {
		// read digits and optional dot part
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
		return Token{Type: NUMBER, Lit: lit, Pos: Position{Offset: l.startOffset, Line: l.line, Col: l.col}}
	}

	// Identifiers and keywords: letter or underscore start. identifiers may include dot (for namespace)
	if unicode.IsLetter(r) || r == '_' {
		for {
			c := l.peek()
			if unicode.IsLetter(c) || unicode.IsDigit(c) || c == '_' || c == '.' {
				l.next()
				continue
			}
			break
		}
		lit := l.input[l.startOffset:l.offset]
		key := lit
		if l.caseFoldKeys {
			key = strings.ToLower(key)
		}
		if typ, ok := l.keywords[key]; ok {
			return Token{Type: typ, Lit: lit, Pos: Position{Offset: l.startOffset, Line: l.line, Col: l.col}}
		}
		return Token{Type: IDENT, Lit: lit, Pos: Position{Offset: l.startOffset, Line: l.line, Col: l.col}}
	}

	// operators and punctuation
	switch r {
	case '(':
		l.next()
		return Token{Type: LPAREN, Lit: "(", Pos: Position{Offset: l.startOffset, Line: l.line, Col: l.col}}
	case ')':
		l.next()
		return Token{Type: RPAREN, Lit: ")", Pos: Position{Offset: l.startOffset, Line: l.line, Col: l.col}}
	case '{':
		l.next()
		return Token{Type: LBRACE, Lit: "{", Pos: Position{Offset: l.startOffset, Line: l.line, Col: l.col}}
	case '}':
		l.next()
		return Token{Type: RBRACE, Lit: "}", Pos: Position{Offset: l.startOffset, Line: l.line, Col: l.col}}
	case '[':
		l.next()
		return Token{Type: LBRACK, Lit: "[", Pos: Position{Offset: l.startOffset, Line: l.line, Col: l.col}}
	case ']':
		l.next()
		return Token{Type: RBRACK, Lit: "]", Pos: Position{Offset: l.startOffset, Line: l.line, Col: l.col}}
	case ',':
		l.next()
		return Token{Type: COMMA, Lit: ",", Pos: Position{Offset: l.startOffset, Line: l.line, Col: l.col}}
	case ':':
		l.next()
		return Token{Type: COLON, Lit: ":", Pos: Position{Offset: l.startOffset, Line: l.line, Col: l.col}}
	case '=':
		l.next()
		return Token{Type: EQ, Lit: "=", Pos: Position{Offset: l.startOffset, Line: l.line, Col: l.col}}
	case '+':
		l.next()
		return Token{Type: PLUS, Lit: "+", Pos: Position{Offset: l.startOffset, Line: l.line, Col: l.col}}
	case '-':
		l.next()
		return Token{Type: MINUS, Lit: "-", Pos: Position{Offset: l.startOffset, Line: l.line, Col: l.col}}
	case '*':
		l.next()
		return Token{Type: MULT, Lit: "*", Pos: Position{Offset: l.startOffset, Line: l.line, Col: l.col}}
	case '/':
		l.next()
		return Token{Type: DIV, Lit: "/", Pos: Position{Offset: l.startOffset, Line: l.line, Col: l.col}}
	case '.':
		l.next()
		return Token{Type: DOT, Lit: ".", Pos: Position{Offset: l.startOffset, Line: l.line, Col: l.col}}
	default:
		// unknown single char
		l.next()
		return Token{Type: UNKNOWN, Lit: string(r), Pos: Position{Offset: l.startOffset, Line: l.line, Col: l.col}}
	}
}
'@ | Out-File -FilePath "pkg\lexer\lexer.go" -Encoding UTF8

# Create a simple test for the lexer
@'
package lexer

import (
	"testing"
)

func TestLexerBasic(t *testing.T) {
	input := `module ui {
    fn navbar(items) [
        print "Hello"
    ]
}`
	
	lex := New(input)
	
	expectedTypes := []TokenType{
		MODULE, IDENT, LBRACE,
		FN, IDENT, LPAREN, IDENT, RPAREN, LBRACK,
		PRINT, STRING, RBRACK,
		RBRACE, EOF,
	}
	
	i := 0
	for {
		tok := lex.Next()
		if tok.Type == EOF {
			break
		}
		
		if i >= len(expectedTypes) {
			t.Errorf("More tokens than expected")
			break
		}
		
		if tok.Type != expectedTypes[i] {
			t.Errorf("Token %d: expected %v, got %v (lit: %q)", i, expectedTypes[i], tok.Type, tok.Lit)
		}
		i++
	}
}

func TestLexerComments(t *testing.T) {
	input := `~ This is a comment ~
print "Hello" ~ inline comment ~`
	
	lex := New(input)
	
	tokens := []Token{}
	for {
		tok := lex.Next()
		if tok.Type == EOF {
			break
		}
		tokens = append(tokens, tok)
	}
	
	// Should only have print and string tokens, no comment tokens
	if len(tokens) != 2 {
		t.Errorf("Expected 2 tokens after skipping comments, got %d", len(tokens))
	}
	
	if tokens[0].Type != PRINT {
		t.Errorf("First token should be PRINT, got %v", tokens[0].Type)
	}
	
	if tokens[1].Type != STRING {
		t.Errorf("Second token should be STRING, got %v", tokens[1].Type)
	}
}
'@ | Out-File -FilePath "pkg\lexer\lexer_test.go" -Encoding UTF8

# Create the AST package structure (Chunk 2)
@'
package ast

import (
	"fmt"
	"strings"
)

// ===============================
// AST Nodes - Core Interfaces
// ===============================

// Node is the base interface for all AST nodes.
type Node interface {
	Pos() Position
	String() string
}

// Statement is a node that represents an executable statement.
type Statement interface {
	Node
	stmtNode()
}

// Expression is a node that represents a value-producing expression.
type Expression interface {
	Node
	exprNode()
}

// Declaration is a node that represents a declaration.
type Declaration interface {
	Node
	declNode()
}

// Position represents a source code position.
type Position struct {
	Filename string
	Line     int
	Column   int
	Offset   int
}

func (p Position) String() string {
	return fmt.Sprintf("%s:%d:%d", p.Filename, p.Line, p.Column)
}

// ===============================
// Program Structure
// ===============================

// Program represents a complete HE program.
type Program struct {
	Pos        Position
	Summons    []*SummonDecl   // summon statements
	Modules    []*ModuleDecl   // top-level modules
	Objects    []*ObjectDecl   // object definitions
	Assets     []*AssetDecl    // asset declarations
	Statements []Statement     // global statements
}

func (p *Program) Pos() Position { return p.Pos }
func (p *Program) String() string {
	var sb strings.Builder
	sb.WriteString("Program {\n")
	
	if len(p.Summons) > 0 {
		sb.WriteString("  Summons:\n")
		for _, s := range p.Summons {
			sb.WriteString("    " + s.String() + "\n")
		}
	}
	
	if len(p.Modules) > 0 {
		sb.WriteString("  Modules:\n")
		for _, m := range p.Modules {
			sb.WriteString("    " + m.String() + "\n")
		}
	}
	
	if len(p.Objects) > 0 {
		sb.WriteString("  Objects:\n")
		for _, o := range p.Objects {
			sb.WriteString("    " + o.String() + "\n")
		}
	}
	
	if len(p.Statements) > 0 {
		sb.WriteString("  Statements:\n")
		for _, s := range p.Statements {
			sb.WriteString("    " + s.String() + "\n")
		}
	}
	
	sb.WriteString("}")
	return sb.String()
}

// ===============================
// Module System
// ===============================

// ModuleDecl represents a module declaration.
type ModuleDecl struct {
	Pos      Position
	Name     string           // module name
	Submodules []*ModuleDecl  // nested modules
	Functions []*FunctionDecl // functions in this module
	Body     []Declaration    // other declarations (modules, objects, assets)
}

func (m *ModuleDecl) Pos() Position { return m.Pos }
func (m *ModuleDecl) declNode()     {}
func (m *ModuleDecl) String() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("module %s {", m.Name))
	
	if len(m.Submodules) > 0 {
		sb.WriteString("\n  Submodules:")
		for _, sm := range m.Submodules {
			sb.WriteString("\n    " + sm.String())
		}
	}
	
	if len(m.Functions) > 0 {
		sb.WriteString("\n  Functions:")
		for _, f := range m.Functions {
			sb.WriteString("\n    " + f.String())
		}
	}
	
	sb.WriteString("\n}")
	return sb.String()
}

// FunctionDecl represents a function declaration.
type FunctionDecl struct {
	Pos      Position
	Name     string
	Params   []*ParamDecl
	Body     []Statement
	ReturnType string // optional return type
}

func (f *FunctionDecl) Pos() Position { return f.Pos }
func (f *FunctionDecl) declNode()     {}
func (f *FunctionDecl) String() string {
	params := make([]string, len(f.Params))
	for i, p := range f.Params {
		params[i] = p.String()
	}
	return fmt.Sprintf("fn %s(%s) [%d statements]", f.Name, strings.Join(params, ", "), len(f.Body))
}

// ParamDecl represents a function parameter declaration.
type ParamDecl struct {
	Pos  Position
	Name string
	Type string // optional type annotation
}

func (p *ParamDecl) Pos() Position { return p.Pos }
func (p *ParamDecl) declNode()     {}
func (p *ParamDecl) String() string {
	if p.Type != "" {
		return fmt.Sprintf("%s:%s", p.Name, p.Type)
	}
	return p.Name
}

// ===============================
// Object System
// ===============================

// ObjectDecl represents an object declaration.
type ObjectDecl struct {
	Pos        Position
	Name       string
	Like       string // parent object (inheritance)
	Properties []*PropertyDecl
	Abilities  []*AbilityDecl
	Reactions  []*ReactionDecl
	Memories   []*MemoryDecl
}

func (o *ObjectDecl) Pos() Position { return o.Pos }
func (o *ObjectDecl) declNode()     {}
func (o *ObjectDecl) String() string {
	return fmt.Sprintf("object %s (like: %s, props: %d, abilities: %d)", 
		o.Name, o.Like, len(o.Properties), len(o.Abilities))
}

// PropertyDecl represents an object property.
type PropertyDecl struct {
	Pos  Position
	Name string
	Value Expression
}

func (p *PropertyDecl) Pos() Position { return p.Pos }
func (p *PropertyDecl) declNode()     {}
func (p *PropertyDecl) String() string {
	return fmt.Sprintf("%s = %s", p.Name, p.Value)
}

// AbilityDecl represents an object ability (method).
type AbilityDecl struct {
	Pos      Position
	Name     string
	Params   []*ParamDecl
	Body     []Statement
}

func (a *AbilityDecl) Pos() Position { return a.Pos }
func (a *AbilityDecl) declNode()     {}
func (a *AbilityDecl) String() string {
	return fmt.Sprintf("ability %s(%d params) [%d statements]", 
		a.Name, len(a.Params), len(a.Body))
}

// ReactionDecl represents an object reaction (event handler).
type ReactionDecl struct {
	Pos      Position
	Trigger  string
	Body     []Statement
}

func (r *ReactionDecl) Pos() Position { return r.Pos }
func (r *ReactionDecl) declNode()     {}
func (r *ReactionDecl) String() string {
	return fmt.Sprintf("on %s [%d statements]", r.Trigger, len(r.Body))
}

// MemoryDecl represents an object memory (state).
type MemoryDecl struct {
	Pos   Position
	Name  string
	Value Expression
}

func (m *MemoryDecl) Pos() Position { return m.Pos }
func (m *MemoryDecl) declNode()     {}
func (m *MemoryDecl) String() string {
	return fmt.Sprintf("memory %s = %s", m.Name, m.Value)
}

// ===============================
// Asset System
// ===============================

// SummonDecl represents a module import.
type SummonDecl struct {
	Pos  Position
	Path string
	As   string // alias
}

func (s *SummonDecl) Pos() Position { return s.Pos }
func (s *SummonDecl) declNode()     {}
func (s *SummonDecl) String() string {
	if s.As != "" {
		return fmt.Sprintf("summon \"%s\" as %s", s.Path, s.As)
	}
	return fmt.Sprintf("summon \"%s\"", s.Path)
}

// AssetDecl represents an asset declaration.
type AssetDecl struct {
	Pos  Position
	Type string // "image", "sound", "music", etc.
	Path string
	Name string // optional name/alias
}

func (a *AssetDecl) Pos() Position { return a.Pos }
func (a *AssetDecl) declNode()     {}
func (a *AssetDecl) String() string {
	if a.Name != "" {
		return fmt.Sprintf("%s \"%s\" named %s", a.Type, a.Path, a.Name)
	}
	return fmt.Sprintf("%s \"%s\"", a.Type, a.Path)
}

// ===============================
// Statements
// ===============================

// PrintStmt represents a print/say statement.
type PrintStmt struct {
	Pos  Position
	Expr Expression
}

func (p *PrintStmt) Pos() Position  { return p.Pos }
func (p *PrintStmt) stmtNode()      {}
func (p *PrintStmt) String() string { return fmt.Sprintf("print %s", p.Expr) }

// AssignStmt represents an assignment statement.
type AssignStmt struct {
	Pos  Position
	Lhs  string
	Expr Expression
}

func (a *AssignStmt) Pos() Position  { return a.Pos }
func (a *AssignStmt) stmtNode()      {}
func (a *AssignStmt) String() string { return fmt.Sprintf("%s = %s", a.Lhs, a.Expr) }

// CallStmt represents a function/method call statement.
type CallStmt struct {
	Pos    Position
	Target Expression // can be identifier or member expression
	Args   []Expression
}

func (c *CallStmt) Pos() Position { return c.Pos }
func (c *CallStmt) stmtNode()     {}
func (c *CallStmt) String() string {
	args := make([]string, len(c.Args))
	for i, a := range c.Args {
		args[i] = a.String()
	}
	return fmt.Sprintf("call %s(%s)", c.Target, strings.Join(args, ", "))
}

// WaitStmt represents a wait/delay statement.
type WaitStmt struct {
	Pos     Position
	Seconds Expression
}

func (w *WaitStmt) Pos() Position  { return w.Pos }
func (w *WaitStmt) stmtNode()      {}
func (w *WaitStmt) String() string { return fmt.Sprintf("wait %s seconds", w.Seconds) }

// IfStmt represents a conditional statement.
type IfStmt struct {
	Pos      Position
	Cond     Expression
	Then     []Statement
	Else     []Statement
}

func (i *IfStmt) Pos() Position { return i.Pos }
func (i *IfStmt) stmtNode()     {}
func (i *IfStmt) String() string {
	return fmt.Sprintf("if %s then [%d] else [%d]", 
		i.Cond, len(i.Then), len(i.Else))
}

// RepeatStmt represents a loop statement.
type RepeatStmt struct {
	Pos  Position
	Cond Expression // optional condition for while loops
	Count Expression // optional count for repeat loops
	Body []Statement
}

func (r *RepeatStmt) Pos() Position { return r.Pos }
func (r *RepeatStmt) stmtNode()     {}
func (r *RepeatStmt) String() string {
	if r.Cond != nil {
		return fmt.Sprintf("while %s [%d]", r.Cond, len(r.Body))
	}
	return fmt.Sprintf("repeat %s [%d]", r.Count, len(r.Body))
}

// ReturnStmt represents a return statement.
type ReturnStmt struct {
	Pos  Position
	Expr Expression
}

func (r *ReturnStmt) Pos() Position  { return r.Pos }
func (r *ReturnStmt) stmtNode()      {}
func (r *ReturnStmt) String() string { return fmt.Sprintf("return %s", r.Expr) }

// ===============================
// Expressions
// ===============================

// IdentExpr represents an identifier expression.
type IdentExpr struct {
	Pos  Position
	Name string
}

func (i *IdentExpr) Pos() Position  { return i.Pos }
func (i *IdentExpr) exprNode()      {}
func (i *IdentExpr) String() string { return i.Name }

// MemberExpr represents a member access expression.
type MemberExpr struct {
	Pos    Position
	Object Expression
	Member string
}

func (m *MemberExpr) Pos() Position { return m.Pos }
func (m *MemberExpr) exprNode()     {}
func (m *MemberExpr) String() string {
	return fmt.Sprintf("%s.%s", m.Object, m.Member)
}

// CallExpr represents a function call expression.
type CallExpr struct {
	Pos    Position
	Target Expression
	Args   []Expression
}

func (c *CallExpr) Pos() Position { return c.Pos }
func (c *CallExpr) exprNode()     {}
func (c *CallExpr) String() string {
	args := make([]string, len(c.Args))
	for i, a := range c.Args {
		args[i] = a.String()
	}
	return fmt.Sprintf("%s(%s)", c.Target, strings.Join(args, ", "))
}

// NumberExpr represents a numeric literal.
type NumberExpr struct {
	Pos    Position
	Value  float64
}

func (n *NumberExpr) Pos() Position  { return n.Pos }
func (n *NumberExpr) exprNode()      {}
func (n *NumberExpr) String() string { return fmt.Sprintf("%v", n.Value) }

// StringExpr represents a string literal.
type StringExpr struct {
	Pos   Position
	Value string
}

func (s *StringExpr) Pos() Position  { return s.Pos }
func (s *StringExpr) exprNode()      {}
func (s *StringExpr) String() string { return fmt.Sprintf("\"%s\"", s.Value) }

// BoolExpr represents a boolean literal.
type BoolExpr struct {
	Pos   Position
	Value bool
}

func (b *BoolExpr) Pos() Position  { return b.Pos }
func (b *BoolExpr) exprNode()      {}
func (b *BoolExpr) String() string { return fmt.Sprintf("%v", b.Value) }

// BinaryExpr represents a binary operation.
type BinaryExpr struct {
	Pos  Position
	Op   string
	Lhs  Expression
	Rhs  Expression
}

func (b *BinaryExpr) Pos() Position { return b.Pos }
func (b *BinaryExpr) exprNode()     {}
func (b *BinaryExpr) String() string {
	return fmt.Sprintf("(%s %s %s)", b.Lhs, b.Op, b.Rhs)
}

// UnaryExpr represents a unary operation.
type UnaryExpr struct {
	Pos  Position
	Op   string
	Expr Expression
}

func (u *UnaryExpr) Pos() Position  { return u.Pos }
func (u *UnaryExpr) exprNode()      {}
func (u *UnaryExpr) String() string { return fmt.Sprintf("(%s%s)", u.Op, u.Expr) }

// ArrayExpr represents an array literal.
type ArrayExpr struct {
	Pos    Position
	Values []Expression
}

func (a *ArrayExpr) Pos() Position { return a.Pos }
func (a *ArrayExpr) exprNode()     {}
func (a *ArrayExpr) String() string {
	values := make([]string, len(a.Values))
	for i, v := range a.Values {
		values[i] = v.String()
	}
	return fmt.Sprintf("[%s]", strings.Join(values, ", "))
}

// NilExpr represents a nil/null value.
type NilExpr struct {
	Pos Position
}

func (n *NilExpr) Pos() Position  { return n.Pos }
func (n *NilExpr) exprNode()      {}
func (n *NilExpr) String() string { return "nil" }
'@ | Out-File -FilePath "pkg\ast\ast.go" -Encoding UTF8

# Create a parser skeleton
@'
package parser

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/user/he/pkg/ast"
	"github.com/user/he/pkg/lexer"
)

// Parser converts tokens into an AST.
type Parser struct {
	lexer    *lexer.Lexer
	curToken lexer.Token
	peekToken lexer.Token
	errors   []string
}

// New creates a new parser.
func New(l *lexer.Lexer) *Parser {
	p := &Parser{lexer: l, errors: []string{}}
	// Read two tokens to initialize curToken and peekToken
	p.nextToken()
	p.nextToken()
	return p
}

// nextToken advances to the next token.
func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.lexer.Next()
}

// curTokenIs checks if current token is of given type.
func (p *Parser) curTokenIs(t lexer.TokenType) bool {
	return p.curToken.Type == t
}

// peekTokenIs checks if next token is of given type.
func (p *Parser) peekTokenIs(t lexer.TokenType) bool {
	return p.peekToken.Type == t
}

// expectPeek asserts that the next token is of given type, advances if true.
func (p *Parser) expectPeek(t lexer.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	}
	p.peekError(t)
	return false
}

// peekError adds an error for unexpected peek token.
func (p *Parser) peekError(t lexer.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead at line %d:%d",
		t, p.peekToken.Type, p.peekToken.Pos.Line, p.peekToken.Pos.Col)
	p.errors = append(p.errors, msg)
}

// Errors returns the parser errors.
func (p *Parser) Errors() []string {
	return p.errors
}

// ParseProgram parses a complete HE program.
func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{
		Pos: ast.Position{
			Line:   1,
			Column: 1,
			Offset: 0,
		},
		Summons:    []*ast.SummonDecl{},
		Modules:    []*ast.ModuleDecl{},
		Objects:    []*ast.ObjectDecl{},
		Assets:     []*ast.AssetDecl{},
		Statements: []ast.Statement{},
	}

	for !p.curTokenIs(lexer.EOF) {
		// Try to parse different types of top-level declarations
		switch {
		case p.curTokenIs(lexer.SUMMON):
			summon := p.parseSummon()
			if summon != nil {
				program.Summons = append(program.Summons, summon)
			}
		case p.curTokenIs(lexer.MODULE):
			module := p.parseModule()
			if module != nil {
				program.Modules = append(program.Modules, module)
			}
		case p.curTokenIs(lexer.MAKE) || p.curTokenIs(lexer.CREATE):
			object := p.parseObject()
			if object != nil {
				program.Objects = append(program.Objects, object)
			}
		case p.curTokenIs(lexer.WITH):
			assets := p.parseAssets()
			program.Assets = append(program.Assets, assets...)
		default:
			// Try to parse as a global statement
			stmt := p.parseStatement()
			if stmt != nil {
				program.Statements = append(program.Statements, stmt)
			} else {
				// Skip unknown token
				p.nextToken()
			}
		}
	}

	return program
}

// parseSummon parses a summon declaration.
func (p *Parser) parseSummon() *ast.SummonDecl {
	// summon "path" [as name]
	pos := ast.Position{
		Line:   p.curToken.Pos.Line,
		Column: p.curToken.Pos.Col,
		Offset: p.curToken.Pos.Offset,
	}
	
	p.nextToken() // skip SUMMON
	
	if !p.curTokenIs(lexer.STRING) {
		p.errors = append(p.errors, "expected string after summon")
		return nil
	}
	
	path := p.curToken.Lit
	p.nextToken()
	
	alias := ""
	if p.curTokenIs(lexer.AS) || p.curTokenIs(lexer.NAMED) {
		p.nextToken() // skip AS or NAMED
		if p.curTokenIs(lexer.IDENT) {
			alias = p.curToken.Lit
			p.nextToken()
		}
	}
	
	return &ast.SummonDecl{
		Pos:  pos,
		Path: path,
		As:   alias,
	}
}

// parseModule parses a module declaration.
func (p *Parser) parseModule() *ast.ModuleDecl {
	// module name { ... }
	pos := ast.Position{
		Line:   p.curToken.Pos.Line,
		Column: p.curToken.Pos.Col,
		Offset: p.curToken.Pos.Offset,
	}
	
	p.nextToken() // skip MODULE
	
	if !p.curTokenIs(lexer.IDENT) {
		p.errors = append(p.errors, "expected module name")
		return nil
	}
	
	name := p.curToken.Lit
	p.nextToken()
	
	if !p.expectPeek(lexer.LBRACE) {
		return nil
	}
	
	module := &ast.ModuleDecl{
		Pos:        pos,
		Name:       name,
		Submodules: []*ast.ModuleDecl{},
		Functions:  []*ast.FunctionDecl{},
		Body:       []ast.Declaration{},
	}
	
	// Parse module body
	for !p.curTokenIs(lexer.RBRACE) && !p.curTokenIs(lexer.EOF) {
		switch {
		case p.curTokenIs(lexer.MODULE):
			submodule := p.parseModule()
			if submodule != nil {
				module.Submodules = append(module.Submodules, submodule)
			}
		case p.curTokenIs(lexer.FN):
			fn := p.parseFunction()
			if fn != nil {
				module.Functions = append(module.Functions, fn)
			}
		default:
			p.nextToken()
		}
	}
	
	if !p.expectPeek(lexer.RBRACE) {
		return nil
	}
	
	return module
}

// parseFunction parses a function declaration.
func (p *Parser) parseFunction() *ast.FunctionDecl {
	// fn name(params) [ ... ]
	pos := ast.Position{
		Line:   p.curToken.Pos.Line,
		Column: p.curToken.Pos.Col,
		Offset: p.curToken.Pos.Offset,
	}
	
	p.nextToken() // skip FN
	
	if !p.curTokenIs(lexer.IDENT) {
		p.errors = append(p.errors, "expected function name")
		return nil
	}
	
	name := p.curToken.Lit
	p.nextToken()
	
	params := []*ast.ParamDecl{}
	if p.curTokenIs(lexer.LPAREN) {
		params = p.parseParamList()
	}
	
	if !p.expectPeek(lexer.LBRACK) {
		return nil
	}
	
	body := p.parseBlock()
	
	return &ast.FunctionDecl{
		Pos:    pos,
		Name:   name,
		Params: params,
		Body:   body,
	}
}

// parseParamList parses a parameter list.
func (p *Parser) parseParamList() []*ast.ParamDecl {
	params := []*ast.ParamDecl{}
	
	p.nextToken() // skip LPAREN
	
	// Empty param list
	if p.curTokenIs(lexer.RPAREN) {
		p.nextToken()
		return params
	}
	
	// Parse first parameter
	param := &ast.ParamDecl{
		Pos: ast.Position{
			Line:   p.curToken.Pos.Line,
			Column: p.curToken.Pos.Col,
			Offset: p.curToken.Pos.Offset,
		},
		Name: p.curToken.Lit,
	}
	p.nextToken()
	params = append(params, param)
	
	// Parse additional parameters
	for p.curTokenIs(lexer.COMMA) {
		p.nextToken()
		param := &ast.ParamDecl{
			Pos: ast.Position{
				Line:   p.curToken.Pos.Line,
				Column: p.curToken.Pos.Col,
				Offset: p.curToken.Pos.Offset,
			},
			Name: p.curToken.Lit,
		}
		p.nextToken()
		params = append(params, param)
	}
	
	if !p.expectPeek(lexer.RPAREN) {
		return nil
	}
	
	return params
}

// parseBlock parses a block of statements.
func (p *Parser) parseBlock() []ast.Statement {
	statements := []ast.Statement{}
	
	p.nextToken() // skip LBRACK
	
	for !p.curTokenIs(lexer.RBRACK) && !p.curTokenIs(lexer.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			statements = append(statements, stmt)
		}
	}
	
	if !p.expectPeek(lexer.RBRACK) {
		return nil
	}
	
	return statements
}

// parseStatement parses a statement.
func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case lexer.PRINT, lexer.SAY:
		return p.parsePrint()
	case lexer.WAIT:
		return p.parseWait()
	case lexer.RETURN:
		return p.parseReturn()
	case lexer.IDENT:
		// Could be assignment or function call
		// We'll implement this later
		p.nextToken()
		return nil
	default:
		p.errors = append(p.errors, fmt.Sprintf("unexpected token in statement: %v", p.curToken))
		p.nextToken()
		return nil
	}
}

// parsePrint parses a print statement.
func (p *Parser) parsePrint() *ast.PrintStmt {
	pos := ast.Position{
		Line:   p.curToken.Pos.Line,
		Column: p.curToken.Pos.Col,
		Offset: p.curToken.Pos.Offset,
	}
	
	p.nextToken() // skip PRINT/SAY
	
	// For now, just parse a string literal
	if !p.curTokenIs(lexer.STRING) {
		p.errors = append(p.errors, "expected string in print statement")
		return nil
	}
	
	expr := &ast.StringExpr{
		Pos: ast.Position{
			Line:   p.curToken.Pos.Line,
			Column: p.curToken.Pos.Col,
			Offset: p.curToken.Pos.Offset,
		},
		Value: p.curToken.Lit,
	}
	
	p.nextToken()
	
	return &ast.PrintStmt{
		Pos:  pos,
		Expr: expr,
	}
}

// parseWait parses a wait statement.
func (p *Parser) parseWait() *ast.WaitStmt {
	pos := ast.Position{
		Line:   p.curToken.Pos.Line,
		Column: p.curToken.Pos.Col,
		Offset: p.curToken.Pos.Offset,
	}
	
	p.nextToken() // skip WAIT
	
	// Parse number
	if !p.curTokenIs(lexer.NUMBER) {
		p.errors = append(p.errors, "expected number in wait statement")
		return nil
	}
	
	value, err := strconv.ParseFloat(p.curToken.Lit, 64)
	if err != nil {
		p.errors = append(p.errors, fmt.Sprintf("invalid number in wait: %v", err))
		return nil
	}
	
	expr := &ast.NumberExpr{
		Pos: ast.Position{
			Line:   p.curToken.Pos.Line,
			Column: p.curToken.Pos.Col,
			Offset: p.curToken.Pos.Offset,
		},
		Value: value,
	}
	
	p.nextToken()
	
	// Optional "seconds" keyword
	if p.curTokenIs(lexer.SECONDS) {
		p.nextToken()
	}
	
	return &ast.WaitStmt{
		Pos:     pos,
		Seconds: expr,
	}
}

// parseReturn parses a return statement.
func (p *Parser) parseReturn() *ast.ReturnStmt {
	pos := ast.Position{
		Line:   p.curToken.Pos.Line,
		Column: p.curToken.Pos.Col,
		Offset: p.curToken.Pos.Offset,
	}
	
	p.nextToken() // skip RETURN
	
	// For now, return nil expression
	return &ast.ReturnStmt{
		Pos:  pos,
		Expr: nil,
	}
}

// parseObject parses an object declaration (placeholder).
func (p *Parser) parseObject() *ast.ObjectDecl {
	// Skip for now
	for !p.curTokenIs(lexer.LBRACE) && !p.curTokenIs(lexer.EOF) {
		p.nextToken()
	}
	
	if p.curTokenIs(lexer.LBRACE) {
		// Skip to matching brace
		braceCount := 1
		p.nextToken()
		for braceCount > 0 && !p.curTokenIs(lexer.EOF) {
			if p.curTokenIs(lexer.LBRACE) {
				braceCount++
			} else if p.curTokenIs(lexer.RBRACE) {
				braceCount--
			}
			p.nextToken()
		}
	}
	
	return nil
}

// parseAssets parses asset declarations (placeholder).
func (p *Parser) parseAssets() []*ast.AssetDecl {
	// Skip WITH and everything until we hit a non-asset token
	for !p.curTokenIs(lexer.EOF) && !p.isStartOfStatement() {
		p.nextToken()
	}
	return []*ast.AssetDecl{}
}

// isStartOfStatement checks if current token could start a statement.
func (p *Parser) isStartOfStatement() bool {
	switch p.curToken.Type {
	case lexer.PRINT, lexer.SAY, lexer.WAIT, lexer.RETURN,
		lexer.MODULE, lexer.SUMMON, lexer.WITH,
		lexer.MAKE, lexer.CREATE, lexer.FN:
		return true
	default:
		return false
	}
}
'@ | Out-File -FilePath "pkg\parser\parser.go" -Encoding UTF8

# Create a simple runtime skeleton
@'
package runtime

import (
	"fmt"
	"time"

	"github.com/user/he/pkg/ast"
)

// Value represents a runtime value.
type Value struct {
	Type  string
	Num   float64
	Str   string
	Bool  bool
	Func  *Function
	Obj   *Object
	Array []Value
}

// Function represents a callable function.
type Function struct {
	Name   string
	Params []string
	Body   []ast.Statement
	Env    *Environment
}

// Object represents a runtime object.
type Object struct {
	Name       string
	Like       string
	Properties map[string]Value
	Abilities  map[string]*Function
	Memories   map[string]Value
}

// Environment represents a scope for variable storage.
type Environment struct {
	Parent *Environment
	Vars   map[string]Value
}

// NewEnvironment creates a new environment.
func NewEnvironment(parent *Environment) *Environment {
	return &Environment{
		Parent: parent,
		Vars:   make(map[string]Value),
	}
}

// Get retrieves a variable from the environment.
func (e *Environment) Get(name string) (Value, bool) {
	val, ok := e.Vars[name]
	if !ok && e.Parent != nil {
		return e.Parent.Get(name)
	}
	return val, ok
}

// Set sets a variable in the current environment.
func (e *Environment) Set(name string, val Value) {
	e.Vars[name] = val
}

// Module represents a loaded module.
type Module struct {
	Name      string
	Functions map[string]*Function
	Submodules map[string]*Module
}

// Runtime executes HE programs.
type Runtime struct {
	GlobalEnv  *Environment
	Modules    map[string]*Module
	Objects    map[string]*Object
	Platform   string // "web", "mobile", "desktop"
	Debug      bool
}

// New creates a new runtime.
func New() *Runtime {
	return &Runtime{
		GlobalEnv: NewEnvironment(nil),
		Modules:   make(map[string]*Module),
		Objects:   make(map[string]*Object),
		Platform:  "web", // default platform
	}
}

// SetPlatform sets the target platform.
func (rt *Runtime) SetPlatform(platform string) {
	rt.Platform = platform
}

// LoadModule loads a module into the runtime.
func (rt *Runtime) LoadModule(name string, module *Module) {
	rt.Modules[name] = module
}

// ExecuteProgram executes a complete HE program.
func (rt *Runtime) ExecuteProgram(program *ast.Program) error {
	// Load summoned modules
	for _, summon := range program.Summons {
		fmt.Printf("Loading module: %s as %s\n", summon.Path, summon.As)
		// TODO: Implement module loading from files
	}

	// Register modules from AST
	for _, module := range program.Modules {
		rt.registerModule(module)
	}

	// Create objects
	for _, obj := range program.Objects {
		rt.createObject(obj)
	}

	// Execute global statements
	for _, stmt := range program.Statements {
		if err := rt.executeStatement(stmt); err != nil {
			return err
		}
	}

	return nil
}

// registerModule registers an AST module in the runtime.
func (rt *Runtime) registerModule(module *ast.ModuleDecl) *Module {
	rtModule := &Module{
		Name:       module.Name,
		Functions:  make(map[string]*Function),
		Submodules: make(map[string]*Module),
	}

	// Register functions
	for _, fn := range module.Functions {
		rtFunc := &Function{
			Name:   fn.Name,
			Params: make([]string, len(fn.Params)),
			Body:   fn.Body,
			Env:    NewEnvironment(rt.GlobalEnv),
		}
		for i, param := range fn.Params {
			rtFunc.Params[i] = param.Name
		}
		rtModule.Functions[fn.Name] = rtFunc
	}

	// Register submodules
	for _, sub := range module.Submodules {
		rtSub := rt.registerModule(sub)
		rtModule.Submodules[sub.Name] = rtSub
	}

	rt.Modules[module.Name] = rtModule
	return rtModule
}

// createObject creates an object from AST.
func (rt *Runtime) createObject(obj *ast.ObjectDecl) *Object {
	rtObj := &Object{
		Name:       obj.Name,
		Like:       obj.Like,
		Properties: make(map[string]Value),
		Abilities:  make(map[string]*Function),
		Memories:   make(map[string]Value),
	}

	// Initialize properties
	for _, prop := range obj.Properties {
		val := rt.evaluateExpression(prop.Value)
		rtObj.Properties[prop.Name] = val
	}

	rt.Objects[obj.Name] = rtObj
	return rtObj
}

// executeStatement executes a single statement.
func (rt *Runtime) executeStatement(stmt ast.Statement) error {
	switch s := stmt.(type) {
	case *ast.PrintStmt:
		return rt.executePrint(s)
	case *ast.WaitStmt:
		return rt.executeWait(s)
	case *ast.ReturnStmt:
		return rt.executeReturn(s)
	default:
		return fmt.Errorf("unknown statement type: %T", stmt)
	}
}

// executePrint executes a print statement.
func (rt *Runtime) executePrint(stmt *ast.PrintStmt) error {
	val := rt.evaluateExpression(stmt.Expr)
	
	switch val.Type {
	case "string":
		fmt.Println(val.Str)
	case "number":
		fmt.Println(val.Num)
	case "boolean":
		fmt.Println(val.Bool)
	default:
		fmt.Println(val)
	}
	
	return nil
}

// executeWait executes a wait statement.
func (rt *Runtime) executeWait(stmt *ast.WaitStmt) error {
	val := rt.evaluateExpression(stmt.Seconds)
	if val.Type != "number" {
		return fmt.Errorf("wait expects a number, got %s", val.Type)
	}
	
	seconds := time.Duration(val.Num * float64(time.Second))
	if rt.Debug {
		fmt.Printf("[DEBUG] Waiting for %v seconds\n", val.Num)
	}
	time.Sleep(seconds)
	return nil
}

// executeReturn executes a return statement (placeholder).
func (rt *Runtime) executeReturn(stmt *ast.ReturnStmt) error {
	// Return statements are handled in function execution
	if rt.Debug {
		fmt.Println("[DEBUG] Return statement")
	}
	return nil
}

// evaluateExpression evaluates an expression to a value.
func (rt *Runtime) evaluateExpression(expr ast.Expression) Value {
	switch e := expr.(type) {
	case *ast.StringExpr:
		return Value{Type: "string", Str: e.Value}
	case *ast.NumberExpr:
		return Value{Type: "number", Num: e.Value}
	case *ast.BoolExpr:
		return Value{Type: "boolean", Bool: e.Value}
	case *ast.IdentExpr:
		val, ok := rt.GlobalEnv.Get(e.Name)
		if !ok {
			// Variable not found, return default
			return Value{Type: "string", Str: ""}
		}
		return val
	default:
		// For now, return a default value
		return Value{Type: "string", Str: ""}
	}
}

// ResolveFunction resolves a function call with platform-aware routing.
func (rt *Runtime) ResolveFunction(modulePath []string, funcName string) (*Function, error) {
	// Example: ui.navbar -> ["ui"], "navbar"
	// Platform-aware resolution:
	// 1. Try exact path: ui.navbar
	// 2. Try platform-specific: ui.mobile.navbar (if platform is mobile)
	// 3. Try fallback platforms
	
	if len(modulePath) == 0 {
		return nil, fmt.Errorf("empty module path")
	}
	
	// Start from root module
	current := rt.Modules[modulePath[0]]
	if current == nil {
		return nil, fmt.Errorf("module not found: %s", modulePath[0])
	}
	
	// Navigate through submodules
	for i := 1; i < len(modulePath); i++ {
		sub := current.Submodules[modulePath[i]]
		if sub == nil {
			return nil, fmt.Errorf("submodule not found: %s", modulePath[i])
		}
		current = sub
	}
	
	// Try to find function in current module
	if fn, ok := current.Functions[funcName]; ok {
		return fn, nil
	}
	
	// Platform-aware fallback
	// If we're looking for ui.navbar, try ui.<platform>.navbar
	if len(modulePath) == 1 {
		// Single module like "ui", try platform submodule
		platformModule := current.Submodules[rt.Platform]
		if platformModule != nil {
			if fn, ok := platformModule.Functions[funcName]; ok {
				return fn, nil
			}
		}
		
		// Try other common platforms as fallback
		fallbacks := []string{"web", "mobile", "desktop"}
		for _, fb := range fallbacks {
			if fb == rt.Platform {
				continue // already tried
			}
			platformModule := current.Submodules[fb]
			if platformModule != nil {
				if fn, ok := platformModule.Functions[funcName]; ok {
					return fn, nil
				}
			}
		}
	}
	
	return nil, fmt.Errorf("function not found: %s.%s", 
		strings.Join(modulePath, "."), funcName)
}
'@ | Out-File -FilePath "pkg\runtime\runtime.go" -Encoding UTF8

# Create a resolver for platform-aware function resolution
@'
package resolver

import (
	"strings"
)

// Platform represents a target platform.
type Platform string

const (
	PlatformWeb    Platform = "web"
	PlatformMobile Platform = "mobile"
	PlatformDesktop Platform = "desktop"
)

// Resolver handles platform-aware module/function resolution.
type Resolver struct {
	platform   Platform
	modules    map[string]*Module
	precedence []Platform // resolution precedence order
}

// Module represents a resolved module.
type Module struct {
	Name      string
	Functions map[string]bool
	Submodules map[string]*Module
}

// New creates a new resolver.
func New(platform Platform) *Resolver {
	r := &Resolver{
		platform: platform,
		modules:  make(map[string]*Module),
	}
	
	// Set platform precedence based on current platform
	switch platform {
	case PlatformMobile:
		r.precedence = []Platform{PlatformMobile, PlatformWeb, PlatformDesktop}
	case PlatformWeb:
		r.precedence = []Platform{PlatformWeb, PlatformMobile, PlatformDesktop}
	case PlatformDesktop:
		r.precedence = []Platform{PlatformDesktop, PlatformWeb, PlatformMobile}
	default:
		r.precedence = []Platform{PlatformWeb, PlatformMobile, PlatformDesktop}
	}
	
	return r
}

// RegisterModule registers a module with the resolver.
func (r *Resolver) RegisterModule(name string, module *Module) {
	r.modules[name] = module
}

// ResolveFunction resolves a function call with platform-aware routing.
func (r *Resolver) ResolveFunction(path string) (string, error) {
	// Parse the path: ui.navbar or ui.mobile.navbar
	parts := strings.Split(path, ".")
	if len(parts) < 2 {
		return "", ErrInvalidPath
	}
	
	funcName := parts[len(parts)-1]
	modulePath := parts[:len(parts)-1]
	
	// Try exact match first
	exactPath := strings.Join(modulePath, ".")
	if module, ok := r.modules[exactPath]; ok {
		if _, hasFunc := module.Functions[funcName]; hasFunc {
			return exactPath + "." + funcName, nil
		}
	}
	
	// Platform-aware resolution for single-level modules
	// e.g., ui.navbar -> try ui.mobile.navbar, ui.web.navbar, ui.desktop.navbar
	if len(modulePath) == 1 {
		baseModule := modulePath[0]
		
		// Try each platform in precedence order
		for _, platform := range r.precedence {
			platformPath := baseModule + "." + string(platform)
			if module, ok := r.modules[platformPath]; ok {
				if _, hasFunc := module.Functions[funcName]; hasFunc {
					return platformPath + "." + funcName, nil
				}
			}
		}
		
		// Try the base module with platform submodule
		if baseModuleObj, ok := r.modules[baseModule]; ok {
			for _, platform := range r.precedence {
				if submodule, ok := baseModuleObj.Submodules[string(platform)]; ok {
					if _, hasFunc := submodule.Functions[funcName]; hasFunc {
						return baseModule + "." + string(platform) + "." + funcName, nil
					}
				}
			}
		}
	}
	
	// Try nested module resolution
	// e.g., ui.components.navbar
	for i := 1; i < len(modulePath); i++ {
		parentPath := strings.Join(modulePath[:i], ".")
		childName := modulePath[i]
		
		if parent, ok := r.modules[parentPath]; ok {
			if child, ok := parent.Submodules[childName]; ok {
				if _, hasFunc := child.Functions[funcName]; hasFunc {
					return parentPath + "." + childName + "." + funcName, nil
				}
			}
		}
	}
	
	return "", ErrFunctionNotFound
}

// ResolveModule resolves a module path to its canonical form.
func (r *Resolver) ResolveModule(path string) (string, error) {
	parts := strings.Split(path, ".")
	
	// Try exact match
	if _, ok := r.modules[path]; ok {
		return path, nil
	}
	
	// Platform-aware resolution for single module
	if len(parts) == 1 {
		baseModule := parts[0]
		
		// Check if this module has platform-specific versions
		for _, platform := range r.precedence {
			platformPath := baseModule + "." + string(platform)
			if _, ok := r.modules[platformPath]; ok {
				return platformPath, nil
			}
		}
	}
	
	return "", ErrModuleNotFound
}

// Errors
var (
	ErrInvalidPath      = errors.New("invalid path")
	ErrFunctionNotFound = errors.New("function not found")
	ErrModuleNotFound   = errors.New("module not found")
)

// errors package
var errors = struct {
	New string
}{
	New: "errors",
}

// Dummy errors to make code compile
var (
	errors = struct {
		New func(string) error
	}{
		New: func(s string) error { return fmt.Errorf(s) },
	}
)

import "fmt"
'@ | Out-File -FilePath "pkg\resolver\resolver.go" -Encoding UTF8

# Fix the resolver.go file (remove the duplicate errors declaration)
@'
package resolver

import (
	"errors"
	"fmt"
	"strings"
)

// Platform represents a target platform.
type Platform string

const (
	PlatformWeb    Platform = "web"
	PlatformMobile Platform = "mobile"
	PlatformDesktop Platform = "desktop"
)

// Resolver handles platform-aware module/function resolution.
type Resolver struct {
	platform   Platform
	modules    map[string]*Module
	precedence []Platform // resolution precedence order
}

// Module represents a resolved module.
type Module struct {
	Name       string
	Functions  map[string]bool
	Submodules map[string]*Module
}

// New creates a new resolver.
func New(platform Platform) *Resolver {
	r := &Resolver{
		platform: platform,
		modules:  make(map[string]*Module),
	}
	
	// Set platform precedence based on current platform
	switch platform {
	case PlatformMobile:
		r.precedence = []Platform{PlatformMobile, PlatformWeb, PlatformDesktop}
	case PlatformWeb:
		r.precedence = []Platform{PlatformWeb, PlatformMobile, PlatformDesktop}
	case PlatformDesktop:
		r.precedence = []Platform{PlatformDesktop, PlatformWeb, PlatformMobile}
	default:
		r.precedence = []Platform{PlatformWeb, PlatformMobile, PlatformDesktop}
	}
	
	return r
}

// RegisterModule registers a module with the resolver.
func (r *Resolver) RegisterModule(name string, module *Module) {
	r.modules[name] = module
}

// ResolveFunction resolves a function call with platform-aware routing.
func (r *Resolver) ResolveFunction(path string) (string, error) {
	// Parse the path: ui.navbar or ui.mobile.navbar
	parts := strings.Split(path, ".")
	if len(parts) < 2 {
		return "", fmt.Errorf("invalid path: %s", path)
	}
	
	funcName := parts[len(parts)-1]
	modulePath := parts[:len(parts)-1]
	
	// Try exact match first
	exactPath := strings.Join(modulePath, ".")
	if module, ok := r.modules[exactPath]; ok {
		if _, hasFunc := module.Functions[funcName]; hasFunc {
			return exactPath + "." + funcName, nil
		}
	}
	
	// Platform-aware resolution for single-level modules
	// e.g., ui.navbar -> try ui.mobile.navbar, ui.web.navbar, ui.desktop.navbar
	if len(modulePath) == 1 {
		baseModule := modulePath[0]
		
		// Try each platform in precedence order
		for _, platform := range r.precedence {
			platformPath := baseModule + "." + string(platform)
			if module, ok := r.modules[platformPath]; ok {
				if _, hasFunc := module.Functions[funcName]; hasFunc {
					return platformPath + "." + funcName, nil
				}
			}
		}
		
		// Try the base module with platform submodule
		if baseModuleObj, ok := r.modules[baseModule]; ok {
			for _, platform := range r.precedence {
				if submodule, ok := baseModuleObj.Submodules[string(platform)]; ok {
					if _, hasFunc := submodule.Functions[funcName]; hasFunc {
						return baseModule + "." + string(platform) + "." + funcName, nil
					}
				}
			}
		}
	}
	
	// Try nested module resolution
	// e.g., ui.components.navbar
	for i := 1; i < len(modulePath); i++ {
		parentPath := strings.Join(modulePath[:i], ".")
		childName := modulePath[i]
		
		if parent, ok := r.modules[parentPath]; ok {
			if child, ok := parent.Submodules[childName]; ok {
				if _, hasFunc := child.Functions[funcName]; hasFunc {
					return parentPath + "." + childName + "." + funcName, nil
				}
			}
		}
	}
	
	return "", fmt.Errorf("function not found: %s", path)
}

// ResolveModule resolves a module path to its canonical form.
func (r *Resolver) ResolveModule(path string) (string, error) {
	parts := strings.Split(path, ".")
	
	// Try exact match
	if _, ok := r.modules[path]; ok {
		return path, nil
	}
	
	// Platform-aware resolution for single module
	if len(parts) == 1 {
		baseModule := parts[0]
		
		// Check if this module has platform-specific versions
		for _, platform := range r.precedence {
			platformPath := baseModule + "." + string(platform)
			if _, ok := r.modules[platformPath]; ok {
				return platformPath, nil
			}
		}
	}
	
	return "", fmt.Errorf("module not found: %s", path)
}
'@ | Out-File -FilePath "pkg\resolver\resolver.go" -Encoding UTF8

# Create an example builtin module (he_mobile.he)
@'
~ HE Mobile Module ~
~ Provides mobile-specific UI components and utilities ~

module he {
    module mobile {
        ~ Navigation Bar Component ~
        fn navbar(items) [
            print "Mobile Navigation Bar"
            print "Items: " + items
            ~ In a real implementation, this would render a mobile navbar ~
            ~ with platform-specific styling and behavior ~
        ]
        
        ~ Mobile Button Component ~
        fn button(label, action) [
            print "Mobile Button: " + label
            ~ When tapped, execute the action ~
            tell runtime to action
        ]
        
        ~ Mobile-specific gesture handler ~
        fn handleSwipe(direction) [
            if direction is "left" then [
                print "Swiped left - showing previous screen"
            ]
            if direction is "right" then [
                print "Swiped right - showing next screen"
            ]
        ]
    }
    
    module web {
        ~ Web Navigation Bar Component ~
        fn navbar(items) [
            print "Web Navigation Bar"
            print "Items: " + items
            ~ This would render an HTML navbar ~
        ]
        
        ~ Web Button Component ~
        fn button(label, action) [
            print "Web Button: " + label
            ~ When clicked, execute the action ~
            tell runtime to action
        ]
    }
    
    module desktop {
        ~ Desktop Navigation Bar Component ~
        fn navbar(items) [
            print "Desktop Navigation Bar"
            print "Items: " + items
            ~ This would render a desktop application navbar ~
        ]
    }
    
    ~ Platform-agnostic utility functions ~
    fn log(message) [
        print "[HE] " + message
    ]
    
    fn formatNumber(num) [
        ~ Format number with commas ~
        return num as string
    ]
}
'@ | Out-File -FilePath "pkg\builtins\he_mobile.he" -Encoding UTF8

# Create a physics module example
@'
~ Physics Module for HE ~
~ Provides physics simulation capabilities ~

module physics {
    ~ Core physics functions ~
    fn gravity(force) [
        ~ Apply gravity force to objects ~
        print "Applying gravity: " + force
    ]
    
    fn collision(object1, object2) [
        ~ Detect collision between two objects ~
        print "Checking collision between " + object1 + " and " + object2
        ~ Return collision result ~
        return true
    ]
    
    ~ Vector math utilities ~
    fn addVectors(v1, v2) [
        print "Adding vectors: " + v1 + " + " + v2
        return v1 + v2
    ]
    
    ~ Physics constants ~
    make Earth [
        has: [
            gravity is 9.81
            mass is 5.972e24
        ]
    ]
    
    make Moon [
        has: [
            gravity is 1.62
            mass is 7.348e22
        ]
    ]
}
'@ | Out-File -FilePath "pkg\builtins\physics.he" -Encoding UTF8

# Create an example HE program
@'
~ Example HE Program ~
~ Demonstrates platform-aware modules and object creation ~

summon "he_mobile" as he
summon "physics" as phys

~ Create a game character ~
make Player [
    has: [
        name is "Hero"
        health is 100
        position is [0, 0, 0]
    ]
    
    can: [
        jump [
            print name + " jumps!"
            ~ Apply physics ~
            phys.gravity(-9.8)
        ]
        
        move(direction) [
            print "Moving " + direction
            ~ Update position based on direction ~
        ]
    ]
    
    on collision [
        print "Player collided with something!"
        health is health - 10
        if health < 0 then [
            print "Game Over!"
        ]
    ]
]

~ Main program ~
print "Starting HE Game..."

~ Platform-aware UI ~
~ This will automatically use the correct version based on platform ~
he.navbar(["Home", "Settings", "Profile"])

~ Create player instance ~
make hero like Player [
    has: [
        name is "Super Hero"
        health is 150
    ]
]

~ Game loop simulation ~
wait 1 seconds
tell hero to jump

~ Wait for user input or events ~
wait 2 seconds
print "Game running..."

~ Example of platform detection ~
print "Running on platform: " + platform

~ Use physics module ~
phys.gravity(9.8)
if phys.collision("hero", "enemy") then [
    print "Collision detected!"
]
'@ | Out-File -FilePath "examples\game.he" -Encoding UTF8

# Create a simple test file
@'
package main

import (
	"fmt"
	"io/ioutil"
	
	"github.com/user/he/pkg/lexer"
	"github.com/user/he/pkg/parser"
	"github.com/user/he/pkg/runtime"
)

func main() {
	// Read example file
	data, err := ioutil.ReadFile("examples/game.he")
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return
	}
	
	source := string(data)
	
	// Lexing
	fmt.Println("=== Lexing ===")
	lx := lexer.New(source)
	for {
		tok := lx.Next()
		if tok.Type == lexer.EOF {
			break
		}
		fmt.Printf("%v: %q\n", tok.Type, tok.Lit)
	}
	
	// Parsing
	fmt.Println("\n=== Parsing ===")
	lx2 := lexer.New(source)
	parser := parser.New(lx2)
	program := parser.ParseProgram()
	
	if len(parser.Errors()) > 0 {
		fmt.Println("Parser errors:")
		for _, err := range parser.Errors() {
			fmt.Println("  ", err)
		}
	} else {
		fmt.Println("Parsed successfully!")
		fmt.Printf("Program structure: %v\n", program)
	}
	
	// Runtime
	fmt.Println("\n=== Runtime ===")
	rt := runtime.New()
	rt.SetPlatform("mobile") // Simulate mobile platform
	
	// Execute the program
	err = rt.ExecuteProgram(program)
	if err != nil {
		fmt.Printf("Runtime error: %v\n", err)
	}
	
	fmt.Println("\n=== Program Complete ===")
}
'@ | Out-File -FilePath "testdata\test.go" -Encoding UTF8

# Create a README file
@'
# HE Language v2

A friendly programming language for creators, with platform-aware modules and auto-resolution.

## Project Structure

```
he/
├── cmd/he/              # CLI tool (he run, he build, he fmt)
├── pkg/
│   ├── lexer/          # Tokenizer for HE source code
│   ├── parser/         # Parser that builds AST from tokens
│   ├── ast/            # Abstract Syntax Tree definitions
│   ├── runtime/        # Interpreter and execution engine
│   ├── resolver/       # Platform-aware module/function resolution
│   ├── builtins/       # Built-in modules (he_mobile, physics, etc.)
│   └── types/          # Type system (future)
├── examples/           # Example HE programs
├── testdata/           # Test programs and data
└── internal/           # Internal packages (compiler, loader)
```

## Key Features

### 1. Platform-Aware Modules
```he
module ui {
    module mobile {
        fn navbar(items) [
            print "Mobile navbar"
        ]
    }
    module web {
        fn navbar(items) [
            print "Web navbar"
        ]
    }
}

~ Auto-resolves based on platform ~
ui.navbar(["Home", "About"])  ~ Uses ui.mobile.navbar on mobile, ui.web.navbar on web ~
```

### 2. Objects and Abilities
```he
make Player [
    has: [
        name is "Hero"
        health is 100
    ]
    can: [
        jump [
            print name + " jumps!"
        ]
    ]
    on collision [
        print "Ouch!"
    ]
]
```

### 3. Module Import System
```he
summon "he_mobile" as ui
summon "physics" as phys

ui.navbar(["Menu"])
phys.gravity(9.8)
```

## Quick Start

1. Build the interpreter:
```bash
cd he
go build -o he.exe ./cmd/he
```

2. Run an example:
```bash
.\he.exe run examples/game.he
```

3. Test the lexer and parser:
```bash
go run testdata/test.go
```

## Development

### Adding New Features

1. **New Keywords**: Add to `pkg/lexer/lexer.go` `initKeywords()` method
2. **New AST Nodes**: Add to `pkg/ast/ast.go`
3. **New Parser Rules**: Add to `pkg/parser/parser.go`
4. **New Runtime Features**: Add to `pkg/runtime/runtime.go`

### Testing
```bash
go test ./pkg/lexer/...
go test ./pkg/parser/...
go test ./pkg/runtime/...
```

## Language Specification

See `docs/grammar.md` for complete EBNF grammar.

## License

MIT License - See LICENSE file for details.
'@ | Out-File -FilePath "README.md" -Encoding UTF8

# Create a simple Makefile for building
@'
.PHONY: build test clean run

build:
	go build -o bin/he ./cmd/he

test:
	go test ./pkg/lexer/...
	go test ./pkg/parser/...
	go test ./pkg/runtime/...

run:
	go run ./cmd/he run examples/game.he

clean:
	rm -rf bin/
	go clean

fmt:
	go fmt ./...

lint:
	go vet ./...

all: fmt lint test build
'@ | Out-File -FilePath "Makefile" -Encoding UTF8

# Create a .gitignore file
@'
# Binaries
bin/
he.exe
he
*.exe
*.dll
*.so
*.dylib

# Test binaries
*.test

# Output directories
dist/
out/
build/

# Dependency directories
vendor/

# Go workspace file
go.work

# IDE files
.vscode/
.idea/
*.swp
*.swo

# OS files
.DS_Store
Thumbs.db

# Logs
*.log
'@ | Out-File -FilePath ".gitignore" -Encoding UTF8

# Initialize git repository
git init
git add .
git commit -m "Initial commit: HE v2 interpreter structure"

# Print completion message
Write-Host ""
Write-Host "✅ HE v2 project structure created successfully!" -ForegroundColor Green
Write-Host ""
Write-Host "Project structure:" -ForegroundColor Cyan
Get-ChildItem -Recurse -Directory | Select-Object -ExpandProperty FullName | ForEach-Object {
    $relPath = Resolve-Path -Relative $_
    Write-Host "  $relPath"
}
Write-Host ""
Write-Host "Next steps:" -ForegroundColor Yellow
Write-Host "1. cd he" -ForegroundColor White
Write-Host "2. go mod tidy" -ForegroundColor White
Write-Host "3. go build -o he.exe ./cmd/he" -ForegroundColor White
Write-Host "4. .\he.exe run examples/game.he" -ForegroundColor White
Write-Host ""
Write-Host "Or run the test: go run testdata/test.go" -ForegroundColor White
Write-Host ""
Write-Host "The interpreter currently supports:" -ForegroundColor Cyan
Write-Host "  • Lexer (tokenizer) with position tracking" -ForegroundColor White
Write-Host "  • AST definitions for modules, functions, objects" -ForegroundColor White
Write-Host "  • Parser skeleton (needs completion)" -ForegroundColor White
Write-Host "  • Runtime skeleton (needs completion)" -ForegroundColor White
Write-Host "  • Resolver for platform-aware function lookup" -ForegroundColor White
Write-Host "  • Example built-in modules" -ForegroundColor White
Write-Host ""
Write-Host "Ready for Chunk 3: Parser completion!" -ForegroundColor Green
```

This PowerShell script creates the complete HE v2 project structure with:

1. **Project Root**: `he/` with Go module initialized
2. **Package Structure**: All packages (`lexer`, `parser`, `ast`, `runtime`, `resolver`, `builtins`)
3. **CLI Tool**: `cmd/he/main.go` with command-line interface
4. **Lexer**: Complete tokenizer from Chunk 1
5. **AST**: Complete AST definitions from Chunk 2
6. **Parser**: Skeleton parser (needs completion)
7. **Runtime**: Skeleton runtime (needs completion)
8. **Resolver**: Platform-aware function resolution
9. **Examples**: Sample HE programs
10. **Built-in Modules**: `he_mobile.he` and `physics.he`
11. **Testing**: Test files and example test runner
12. **Documentation**: README, Makefile, .gitignore

The project is ready to build using Go. Run `go mod tidy`, then `go build -o he.exe ./cmd/he` and execute `.\he.exe run examples/game.he` once the parser/runtime are completed.																												