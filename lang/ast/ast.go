package ast

import "hunterlang/lang/token"

// ── Top-level ────────────────────────────────────────────────────────────────

type Program struct {
	Lines []Line
}

// Line is anything that can appear at the top level of a program.
type Line interface{ lineNode() }

// SummonLine: summon "module" as alias
type SummonLine struct {
	ModuleName token.Token
	Alias      *Identifier
}

func (SummonLine) lineNode() {}

// ObjectLine: create/make Name [like Parent] { ... }
type ObjectLine struct {
	Kind       string // "create" | "make"
	Name       Identifier
	Like       *Identifier // optional parent
	Body       ObjectBody
	ProtectTag string // "" = unprotected; "protected", "protected1", etc.
}

func (*ObjectLine) lineNode() {}

// GlobalStatementLine wraps a statement at the top level.
type GlobalStatementLine struct {
	Statement Statement
}

func (GlobalStatementLine) lineNode() {}

// WithAssetsLine: with assets { ... }
type WithAssetsLine struct{}

func (WithAssetsLine) lineNode() {}

// ── Object body ──────────────────────────────────────────────────────────────

type ObjectBody struct {
	Sections      []BodySection
	NestedObjects []ObjectLine
}

type BodySection interface{ sectionNode() }

// PropertiesSection: has/owns/carries: [ name is value, ... ]
type PropertiesSection struct {
	Props []Property
}

func (PropertiesSection) sectionNode() {}

type Property struct {
	Name  string
	Kind  string // e.g. "has is", "owns starts as"
	Value Expression
}

// AbilitiesSection: can: [ actionName(params) [ body ] ]
type AbilitiesSection struct {
	Abilities []ActionDef
}

func (AbilitiesSection) sectionNode() {}

type ActionDef struct {
	Name       string
	Params     []Param
	HasParams  bool
	Returns    *TypeNode
	Body       []Statement
	ProtectTag string // "" = unprotected; "protected", "protected1", etc.
}

type Param struct {
	Name string
	Type TypeNode
}

type TypeNode struct {
	Name string // e.g. "number", "text", "boolean", "number[]"
}

// ReactionsSection: on/when/whenever trigger [ body ]
type ReactionsSection struct {
	Reactions []ReactionDef
}

func (ReactionsSection) sectionNode() {}

type ReactionDef struct {
	Trigger Trigger
	Body    []Statement
}

// Trigger: eventName [with "value" | with identifier]
type Trigger struct {
	Left        string  // event name
	Right       *string // optional qualifier (identifier or string value)
	RightIsStr  bool    // true if Right came from a string literal
}

// MemoriesSection: remembers: [ ... ]
type MemoriesSection struct{}

func (MemoriesSection) sectionNode() {}

// ── Identifiers ──────────────────────────────────────────────────────────────

type Identifier struct {
	Name string
	Tok  token.Token
}

// ── Statements ───────────────────────────────────────────────────────────────

type Statement interface{ stmtNode() }

// SayStmt: say/print/show expr
type SayStmt struct{ Expr Expression }

func (SayStmt) stmtNode() {}

// ChangeStmt: set/let/change/make name to/be expr
type ChangeStmt struct {
	Kind string // "set" | "let" | "change" | "make"
	Name string
	Expr Expression
}

func (ChangeStmt) stmtNode() {}

// GrowStmt: grow name by expr
type GrowStmt struct {
	Name string
	Expr Expression
}

func (GrowStmt) stmtNode() {}

// ShrinkStmt: shrink name by expr
type ShrinkStmt struct {
	Name string
	Expr Expression
}

func (ShrinkStmt) stmtNode() {}

// DecideStmt: if cond then [ body ] else [ body ]
type DecideStmt struct {
	Cond Expression
	Then []Statement
	Else []Statement
}

func (DecideStmt) stmtNode() {}

// RepeatStmt: repeat while cond [ body ] | repeat N times [ body ]
type RepeatStmt struct {
	Kind  string     // "while" | "times"
	Cond  Expression // condition for "while", count for "times"
	Body  []Statement
}

func (RepeatStmt) stmtNode() {}

// CallStmt: tell obj to action [with args]
type CallStmt struct {
	Object string
	Action string
	Args   []Expression
}

func (CallStmt) stmtNode() {}

// WaitStmt: wait N seconds | wait N frames
type WaitStmt struct {
	Expr Expression
	Unit string // "seconds" | "frames"
}

func (WaitStmt) stmtNode() {}

// ReturnStmt: return expr
type ReturnStmt struct{ Expr Expression }

func (ReturnStmt) stmtNode() {}

// ExprStmt: expression used as a statement (e.g. method call)
type ExprStmt struct{ Expr Expression }

func (ExprStmt) stmtNode() {}

// ── Expressions ──────────────────────────────────────────────────────────────

type Expression interface{ exprNode() }

type NumberLit struct{ Value float64 }

func (NumberLit) exprNode() {}

type StringLit struct{ Value string }

func (StringLit) exprNode() {}

type BooleanLit struct{ Value bool }

func (BooleanLit) exprNode() {}

type IdentifierExpr struct{ Name string }

func (IdentifierExpr) exprNode() {}

type ArrayLit struct{ Elems []Expression }

func (ArrayLit) exprNode() {}

// NamedArgLit: [ key is value, key: value, ... ]
// Used for tell ui to window with [title is "...", children: [...]]
type NamedArgLit struct {
	Pairs []NamedPair
}

func (NamedArgLit) exprNode() {}

type NamedPair struct {
	Key   string
	Value Expression
}

type ParenExpr struct{ X Expression }

func (ParenExpr) exprNode() {}

type UnaryExpr struct {
	Op string
	X  Expression
}

func (UnaryExpr) exprNode() {}

type BinaryExpr struct {
	Left  Expression
	Op    string
	Right Expression
}

func (BinaryExpr) exprNode() {}

type PowerExpr struct {
	Left  Expression
	Right Expression
}

func (PowerExpr) exprNode() {}

type CompareExpr struct {
	Left  Expression
	Op    string
	Right Expression
}

func (CompareExpr) exprNode() {}

type LogicAndExpr struct {
	Left  Expression
	Right Expression
}

func (LogicAndExpr) exprNode() {}

type LogicOrExpr struct {
	Left  Expression
	Right Expression
}

func (LogicOrExpr) exprNode() {}

// MethodCallExpr: receiver.method(args)
type MethodCallExpr struct {
	Receiver IdentifierExpr
	Method   string
	Args     []Expression
}

func (MethodCallExpr) exprNode() {}

// CallExpr: callee(args)
type CallExpr struct {
	Callee string
	Args   []Expression
}

func (CallExpr) exprNode() {}

// FieldAccessExpr: obj.field (no call parens)
type FieldAccessExpr struct {
	Receiver IdentifierExpr
	Field    string
}

func (FieldAccessExpr) exprNode() {}

// ForEachStmt: for each item in list [ body ]
type ForEachStmt struct {
	VarName string
	List    Expression
	Body    []Statement
}

func (ForEachStmt) stmtNode() {}

// AskStmt: ask "prompt" storing result in VarName
type AskStmt struct {
	Prompt  Expression
	VarName string
}

func (AskStmt) stmtNode() {}

// DotAssignStmt: set obj.field to expr
type DotAssignStmt struct {
	Object string
	Field  string
	Expr   Expression
}

func (DotAssignStmt) stmtNode() {}

// ── Pass 3 additions ──────────────────────────────────────────────────────────

// RangeLoopStmt: for each i from 1 to 10 [step 2] [ body ]
type RangeLoopStmt struct {
	VarName string
	From    Expression
	To      Expression
	Step    Expression // nil = 1
	Body    []Statement
}

func (RangeLoopStmt) stmtNode() {}

// TryStmt: try [ body ] or [ handler ] if it fails
type TryStmt struct {
	Body    []Statement
	Handler []Statement
}

func (TryStmt) stmtNode() {}

// RepeatUntilStmt: repeat until cond [ body ]
type RepeatUntilStmt struct {
	Cond Expression
	Body []Statement
}

func (RepeatUntilStmt) stmtNode() {}

// InterpStringExpr: "Hello, {name}! You have {score} points."
// Segments alternate: StringLit, IdentifierExpr, StringLit, ...
type InterpStringExpr struct {
	Segments []Expression // alternating StringLit and any Expression
}

func (InterpStringExpr) exprNode() {}

// ── Pass 4 additions ──────────────────────────────────────────────────────────

// BetweenExpr: X is between A and B  →  A <= X <= B
type BetweenExpr struct {
	Value Expression
	Low   Expression
	High  Expression
}

func (BetweenExpr) exprNode() {}

// AbilityLit: ability(params) [ body ]  — anonymous function value
type AbilityLit struct {
	Params []string
	Body   []Statement
}

func (AbilityLit) exprNode() {}

// MultiAssignStmt: set a, b to expr  |  set a, b to expr1, expr2
type MultiAssignStmt struct {
	Names []string
	Exprs []Expression
}

func (MultiAssignStmt) stmtNode() {}

// MultiReturnStmt: return a, b, c
type MultiReturnStmt struct {
	Exprs []Expression
}

func (MultiReturnStmt) stmtNode() {}

// RememberStmt: remember name [as "key"]  — persist to disk
type RememberStmt struct {
	Name string
	Key  string // optional alias key; defaults to Name
}

func (RememberStmt) stmtNode() {}

// ForgetStmt: forget name  — remove from persistent store
type ForgetStmt struct {
	Name string
}

func (ForgetStmt) stmtNode() {}

// RecallStmt: recall name [as "key"]  — load from persistent store
// "recall score"  or  "recall score as highscore"
type RecallStmt struct {
	Name string
	Key  string
}

func (RecallStmt) stmtNode() {}

// CallAbilityExpr: call a stored ability value with args
// e.g.  set result to greet("Hunter")  where greet is an AbilityLit value
type CallAbilityExpr struct {
	Ability Expression
	Args    []Expression
}

func (CallAbilityExpr) exprNode() {}

// ── Pass 5 additions ──────────────────────────────────────────────────────────

// TryWithVarStmt: try [...] or (errVar) [...]
// Named error capture — errVar holds the error message in the handler
type TryWithVarStmt struct {
	Body    []Statement
	ErrVar  string
	Handler []Statement
}

func (TryWithVarStmt) stmtNode() {}

// EachRangeStmt: each N from A to B [step S] [...]
// Short alias for for-each range loop
type EachRangeStmt = RangeLoopStmt // same AST node, different parse path

// TypeAnnotatedProp — property with explicit type hint
// has: [ health: number is 100 ]
type TypeAnnotatedProp struct {
	Name     string
	TypeHint string
	Value    Expression
}

// LoadModuleStmt: summon "path/to/mylib.he" [as alias]
// Always flat-merges the file's top-level names into the importing scope.
// The alias (if given) is additive — a qualified path on top, never a gate.
type LoadModuleStmt struct {
	FilePath string
	Alias    string
	HasAlias bool
}

func (LoadModuleStmt) stmtNode() {}

// ── Pass 6 additions ──────────────────────────────────────────────────────────

// MembershipExpr: X is one of [a, b, c]
type MembershipExpr struct {
	Value Expression
	List  Expression
}

func (MembershipExpr) exprNode() {}

// ForEachFieldStmt: for each field in obj.fields [body]
// Iterates key-value pairs of an object
type ForEachFieldStmt struct {
	KeyVar   string
	ValVar   string // optional — "for each key, val in obj.fields"
	Object   Expression
	Body     []Statement
}

func (ForEachFieldStmt) stmtNode() {}

// WithScopeStmt: with expr as alias [ body ]
// Binds expr to alias for the duration of body
type WithScopeStmt struct {
	Alias string
	Expr  Expression
	Body  []Statement
}

func (WithScopeStmt) stmtNode() {}

// MethodChainExpr: value.methodName(args) where value is any expression
// Enables "name".upper(), text.split(","), etc. via virtual dispatch
type MethodChainExpr struct {
	Recv   Expression
	Method string
	Args   []Expression
}

func (MethodChainExpr) exprNode() {}

// CaptureExpr: captures surrounding scope for an AbilityLit
// (internal — produced by closure analysis, not parsed directly)
type ClosureExpr struct {
	Params   []string
	Body     []Statement
	Captured map[string]Expression // name → value at capture time
}

func (ClosureExpr) exprNode() {}

// ── Pass 7 additions ──────────────────────────────────────────────────────────

// BecomeStmt: "name becomes value" — mutation shorthand
// Same semantics as ChangeStmt — sugar for readability
type BecomeStmt = ChangeStmt // alias

// TypedProperty — property with enforced type
type TypedProperty struct {
	Name     string
	TypeHint string // "number", "text", "boolean", "list", "nothing"
	Value    Expression
}

// CountedRepeatStmt: repeat N times [as i] [body]
// Exposes loop counter as variable i (0-based by default, 1-based with "as i from 1")
type CountedRepeatStmt struct {
	Count    Expression
	CountVar string // variable name for counter — empty = no counter exposed
	Body     []Statement
}

func (CountedRepeatStmt) stmtNode() {}
