package ast

import "hunterlang/lang/token"

type Program struct {
	Lines []Line
}

type Line interface {
	isLine()
}

type (
	SummonLine struct {
		ModuleName token.Token // STRING lexeme
		Alias      *Identifier
	}
	WithAssetsLine struct {
		Assets []AssetSpec
		// In this grammar, `with` appears only inside object body, but we still model it as a line-like node.
	}
)

func (SummonLine) isLine()     {}
func (WithAssetsLine) isLine() {}

type EmptyLine struct{}

func (EmptyLine) isLine() {}

type CommentLine struct {
	Raw string
}

func (CommentLine) isLine() {}

type Identifier struct {
	Name string
	Tok  token.Token
}

type AssetSpec struct {
	Type  string
	Value string
	Named *Identifier
}

type ObjectLine struct {
	Kind   string // "create" | "make"
	Name   Identifier
	Like   *Identifier
	Body   ObjectBody
	Assets []AssetSpec // if we later support object-level with-assets
}

func (*ObjectLine) isLine() {}

type GlobalStatementLine struct {
	Statement Statement
}

func (GlobalStatementLine) isLine() {}

// ObjectBody models `{ ... }` contents.
type ObjectBody struct {
	Sections        []Section
	NestedObjects   []ObjectLine
	InnerStatements []Statement
	InnerWithAssets []AssetSpec
}

// Sections
type Section interface{ isSection() }

type PropertiesSection struct{ Props []Property }
type AbilitiesSection struct{ Abilities []ActionDef }
type ReactionsSection struct{ Reactions []ReactionDef }
type MemoriesSection struct{ Memories []MemoryDef }

func (PropertiesSection) isSection() {}
func (AbilitiesSection) isSection()  {}
func (ReactionsSection) isSection()  {}
func (MemoriesSection) isSection()   {}

type Property struct {
	Name  string
	Kind  string // "has" | "owns" | "carries" | "is" | "starts as" are parsed into Name/Kind/value by parser later
	Value Expression
}

type ActionDef struct {
	Name      string
	HasParams bool
	Params    []Param
	Returns   *TypeNode
	Body      []Statement
	BodyTok   token.Token
}

type Param struct {
	Name string
	Type TypeNode
}

type TypeNode struct {
	Name string // "number" | "string" | "boolean" | "number[]" etc.
}

type ReactionDef struct {
	Trigger Trigger
	Body    []Statement
}

type Trigger struct {
	Left  string  // identifier
	Right *string // optional identifier after `with`
}

type MemoryDef struct {
	Name  string
	Block []StatementOrProperty
}

// Memories grammar says: Identifier "remembers" ":" Block
type StatementOrProperty interface{ isStmtOrProp() }

type StmtAsAny struct{ Stmt Statement }

func (StmtAsAny) isStmtOrProp() {}

type PropAsAny struct{ Prop Property }

func (PropAsAny) isStmtOrProp() {}

// Statements
type Statement interface{ isStatement() }

type (
	SayStmt struct {
		Expr Expression
	}
	ChangeStmt struct {
		Kind string // "set" | "make"
		Name string
		Expr Expression
	}
	DecideStmt struct {
		Cond Expression
		Then []Statement
		Else []Statement
	}
	RepeatStmt struct {
		Kind string // "repeat" | "while"
		Cond Expression
		Body []Statement
	}
	CallStmt struct {
		Object string // identifier
		Action string // identifier
		Args   []Expression
	}
	WaitStmt struct {
		Expr Expression
		Unit string // "seconds" | "frames"
	}
	ReturnStmt struct {
		Expr Expression
	}
	ExprStmt struct {
		Expr Expression
	}
)

func (SayStmt) isStatement()    {}
func (ChangeStmt) isStatement() {}
func (DecideStmt) isStatement() {}
func (RepeatStmt) isStatement() {}
func (CallStmt) isStatement()   {}
func (WaitStmt) isStatement()   {}
func (ReturnStmt) isStatement() {}
func (ExprStmt) isStatement()   {}

type Block struct {
	Statements []Statement
}

// Expressions
type Expression interface{ isExpr() }

type (
	LogicOrExpr  struct{ Left, Right Expression }
	LogicAndExpr struct{ Left, Right Expression }

	CompareExpr struct {
		Left  Expression
		Op    string // == != > < >= <=
		Right Expression
	}

	BinaryExpr struct {
		Left  Expression
		Op    string // + - * /
		Right Expression
	}

	PowerExpr struct {
		Left  Expression
		Right Expression
	}

	UnaryExpr struct {
		Op string // - !
		X  Expression
	}

	NumberLit  struct{ Value float64 }
	StringLit  struct{ Value string }
	BooleanLit struct{ Value bool }

	IdentifierExpr struct{ Name string }

	CallExpr struct {
		Callee string // Identifier
		Args   []Expression
	}

	MethodCallExpr struct {
		Receiver Expression // IdentifierExpr or more complex later
		Method   string
		Args     []Expression
	}

	ArrayLit struct {
		Elems []Expression
	}

	ParenExpr struct{ X Expression }
)

func (LogicOrExpr) isExpr()    {}
func (LogicAndExpr) isExpr()   {}
func (CompareExpr) isExpr()    {}
func (BinaryExpr) isExpr()     {}
func (PowerExpr) isExpr()      {}
func (UnaryExpr) isExpr()      {}
func (NumberLit) isExpr()      {}
func (StringLit) isExpr()      {}
func (BooleanLit) isExpr()     {}
func (IdentifierExpr) isExpr() {}
func (CallExpr) isExpr()       {}
func (MethodCallExpr) isExpr() {}
func (ArrayLit) isExpr()       {}
func (ParenExpr) isExpr()      {}
