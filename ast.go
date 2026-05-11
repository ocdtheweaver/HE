// Package ast defines the Abstract Syntax Tree nodes for Hunter's Engine.
//
// The grammar is near-English OOP.  Every construct maps to one of three
// categories: declarations, statements, and expressions.
package ast

import (
	"fmt"
	"strings"
)

// ─── position ────────────────────────────────────────────────────────────────

// Pos records a source location (for diagnostics).
type Pos struct {
	Line int
	Col  int
}

func (p Pos) String() string { return fmt.Sprintf("%d:%d", p.Line, p.Col) }

// ─── node interfaces ──────────────────────────────────────────────────────────

// Node is the root interface of all AST nodes.
type Node interface {
	nodePos() Pos
	String() string
}

// Stmt is a statement node.
type Stmt interface {
	Node
	stmtNode()
}

// Expr is an expression node.
type Expr interface {
	Node
	exprNode()
}

// ─── program ──────────────────────────────────────────────────────────────────

// Program is the top-level container.
type Program struct {
	Statements []Stmt
}

func (p *Program) nodePos() Pos { return Pos{} }
func (p *Program) String() string {
	var parts []string
	for _, s := range p.Statements {
		parts = append(parts, s.String())
	}
	return strings.Join(parts, "\n")
}

// ─── declarations ─────────────────────────────────────────────────────────────

// ClassDecl  →  class <Name> [extends <Super>] NEWLINE body end
type ClassDecl struct {
	Pos        Pos
	Name       string
	SuperClass string // empty if none
	Properties []*PropertyDecl
	Methods    []*MethodDecl
}

func (c *ClassDecl) nodePos() Pos   { return c.Pos }
func (c *ClassDecl) stmtNode()      {}
func (c *ClassDecl) String() string { return fmt.Sprintf("class %s", c.Name) }

// PropertyDecl  →  property <Name> [is <Expr>]
type PropertyDecl struct {
	Pos     Pos
	Name    string
	Default Expr // nil means no default
}

func (p *PropertyDecl) nodePos() Pos   { return p.Pos }
func (p *PropertyDecl) stmtNode()      {}
func (p *PropertyDecl) String() string { return "property " + p.Name }

// MethodDecl  →  method <Name> [with <params>] [returns] NEWLINE body end
type MethodDecl struct {
	Pos    Pos
	Name   string
	Params []string
	Body   []Stmt
}

func (m *MethodDecl) nodePos() Pos   { return m.Pos }
func (m *MethodDecl) stmtNode()      {}
func (m *MethodDecl) String() string { return fmt.Sprintf("method %s", m.Name) }

// ImportStmt  →  import <name> [from <path>]
type ImportStmt struct {
	Pos  Pos
	Name string
	Path string
}

func (i *ImportStmt) nodePos() Pos   { return i.Pos }
func (i *ImportStmt) stmtNode()      {}
func (i *ImportStmt) String() string { return "import " + i.Name }

// ─── statements ───────────────────────────────────────────────────────────────

// ExprStmt wraps a bare expression used as a statement.
type ExprStmt struct {
	Pos  Pos
	Expr Expr
}

func (e *ExprStmt) nodePos() Pos   { return e.Pos }
func (e *ExprStmt) stmtNode()      {}
func (e *ExprStmt) String() string { return e.Expr.String() }

// SetStmt  →  set <Target> to <Value>
// Also covers:  <ident> is <expr>   and   <ident> = <expr>
type SetStmt struct {
	Pos    Pos
	Target Expr // Identifier or DotExpr
	Value  Expr
}

func (s *SetStmt) nodePos() Pos   { return s.Pos }
func (s *SetStmt) stmtNode()      {}
func (s *SetStmt) String() string { return fmt.Sprintf("set %s to %s", s.Target, s.Value) }

// CompoundAssignStmt  →  <target> += / -= / *= /= <expr>
type CompoundAssignStmt struct {
	Pos      Pos
	Operator string
	Target   Expr
	Value    Expr
}

func (c *CompoundAssignStmt) nodePos() Pos   { return c.Pos }
func (c *CompoundAssignStmt) stmtNode()      {}
func (c *CompoundAssignStmt) String() string { return fmt.Sprintf("%s %s %s", c.Target, c.Operator, c.Value) }

// ObjectDecl  →  object <Name> of <Class> [with <prop> is <val>, ...]
type ObjectDecl struct {
	Pos        Pos
	Name       string
	ClassName  string
	InitFields map[string]Expr
}

func (o *ObjectDecl) nodePos() Pos   { return o.Pos }
func (o *ObjectDecl) stmtNode()      {}
func (o *ObjectDecl) String() string { return fmt.Sprintf("object %s of %s", o.Name, o.ClassName) }

// PrintStmt  →  print <Expr>
type PrintStmt struct {
	Pos  Pos
	Expr Expr
}

func (p *PrintStmt) nodePos() Pos   { return p.Pos }
func (p *PrintStmt) stmtNode()      {}
func (p *PrintStmt) String() string { return "print " + p.Expr.String() }

// InputStmt  →  input <Ident> [as <Prompt>]
type InputStmt struct {
	Pos    Pos
	Target string
	Prompt Expr // optional
}

func (i *InputStmt) nodePos() Pos   { return i.Pos }
func (i *InputStmt) stmtNode()      {}
func (i *InputStmt) String() string { return "input " + i.Target }

// IfStmt  →  if <Cond> then NEWLINE body [else body] end
type IfStmt struct {
	Pos        Pos
	Condition  Expr
	Consequent []Stmt
	Alternate  []Stmt // nil if no else
}

func (i *IfStmt) nodePos() Pos   { return i.Pos }
func (i *IfStmt) stmtNode()      {}
func (i *IfStmt) String() string { return "if " + i.Condition.String() }

// WhileStmt  →  while <Cond> do NEWLINE body end
type WhileStmt struct {
	Pos       Pos
	Condition Expr
	Body      []Stmt
}

func (w *WhileStmt) nodePos() Pos   { return w.Pos }
func (w *WhileStmt) stmtNode()      {}
func (w *WhileStmt) String() string { return "while " + w.Condition.String() }

// ForEachStmt  →  for each <Var> in <Iterable> do NEWLINE body end
type ForEachStmt struct {
	Pos      Pos
	VarName  string
	Iterable Expr
	Body     []Stmt
}

func (f *ForEachStmt) nodePos() Pos   { return f.Pos }
func (f *ForEachStmt) stmtNode()      {}
func (f *ForEachStmt) String() string { return "for each " + f.VarName }

// ReturnStmt  →  return <Expr> / answer <Expr>
type ReturnStmt struct {
	Pos   Pos
	Value Expr // nil → return nothing
}

func (r *ReturnStmt) nodePos() Pos   { return r.Pos }
func (r *ReturnStmt) stmtNode()      {}
func (r *ReturnStmt) String() string { return "return" }

// BreakStmt  →  break
type BreakStmt struct{ Pos Pos }

func (b *BreakStmt) nodePos() Pos   { return b.Pos }
func (b *BreakStmt) stmtNode()      {}
func (b *BreakStmt) String() string { return "break" }

// ContinueStmt  →  continue
type ContinueStmt struct{ Pos Pos }

func (c *ContinueStmt) nodePos() Pos   { return c.Pos }
func (c *ContinueStmt) stmtNode()      {}
func (c *ContinueStmt) String() string { return "continue" }

// CallStmt  →  call <Method> on <Object> [with <args>]
type CallStmt struct {
	Pos    Pos
	Method string
	Object Expr
	Args   []Expr
}

func (c *CallStmt) nodePos() Pos   { return c.Pos }
func (c *CallStmt) stmtNode()      {}
func (c *CallStmt) String() string { return fmt.Sprintf("call %s on %s", c.Method, c.Object) }

// ─── expressions ──────────────────────────────────────────────────────────────

// NumberLit holds a numeric literal.
type NumberLit struct {
	Pos   Pos
	Value float64
	Raw   string
}

func (n *NumberLit) nodePos() Pos   { return n.Pos }
func (n *NumberLit) exprNode()      {}
func (n *NumberLit) String() string { return n.Raw }

// StringLit holds a string literal.
type StringLit struct {
	Pos   Pos
	Value string
}

func (s *StringLit) nodePos() Pos   { return s.Pos }
func (s *StringLit) exprNode()      {}
func (s *StringLit) String() string { return fmt.Sprintf("%q", s.Value) }

// BoolLit holds true / false.
type BoolLit struct {
	Pos   Pos
	Value bool
}

func (b *BoolLit) nodePos() Pos   { return b.Pos }
func (b *BoolLit) exprNode()      {}
func (b *BoolLit) String() string {
	if b.Value {
		return "true"
	}
	return "false"
}

// NilLit represents nil / nothing.
type NilLit struct{ Pos Pos }

func (n *NilLit) nodePos() Pos   { return n.Pos }
func (n *NilLit) exprNode()      {}
func (n *NilLit) String() string { return "nil" }

// Identifier references a variable, parameter, or class name.
type Identifier struct {
	Pos  Pos
	Name string
}

func (i *Identifier) nodePos() Pos   { return i.Pos }
func (i *Identifier) exprNode()      {}
func (i *Identifier) String() string { return i.Name }

// DotExpr  →  <Object>.<Field>
type DotExpr struct {
	Pos    Pos
	Object Expr
	Field  string
}

func (d *DotExpr) nodePos() Pos   { return d.Pos }
func (d *DotExpr) exprNode()      {}
func (d *DotExpr) String() string { return d.Object.String() + "." + d.Field }

// IndexExpr  →  <Object>[<Index>]
type IndexExpr struct {
	Pos    Pos
	Object Expr
	Index  Expr
}

func (i *IndexExpr) nodePos() Pos   { return i.Pos }
func (i *IndexExpr) exprNode()      {}
func (i *IndexExpr) String() string { return fmt.Sprintf("%s[%s]", i.Object, i.Index) }

// BinaryExpr  →  <Left> <Op> <Right>
type BinaryExpr struct {
	Pos      Pos
	Operator string
	Left     Expr
	Right    Expr
}

func (b *BinaryExpr) nodePos() Pos   { return b.Pos }
func (b *BinaryExpr) exprNode()      {}
func (b *BinaryExpr) String() string { return fmt.Sprintf("(%s %s %s)", b.Left, b.Operator, b.Right) }

// UnaryExpr  →  <Op> <Operand>
type UnaryExpr struct {
	Pos      Pos
	Operator string
	Operand  Expr
}

func (u *UnaryExpr) nodePos() Pos   { return u.Pos }
func (u *UnaryExpr) exprNode()      {}
func (u *UnaryExpr) String() string { return fmt.Sprintf("(%s %s)", u.Operator, u.Operand) }

// CallExpr  →  <object>.<method>(<args>)  or  call <method> on <object>
type CallExpr struct {
	Pos    Pos
	Object Expr
	Method string
	Args   []Expr
}

func (c *CallExpr) nodePos() Pos   { return c.Pos }
func (c *CallExpr) exprNode()      {}
func (c *CallExpr) String() string { return fmt.Sprintf("%s.%s(...)", c.Object, c.Method) }

// NewExpr  →  new <ClassName> [with <prop> is <val>, ...]
type NewExpr struct {
	Pos        Pos
	ClassName  string
	InitFields map[string]Expr
}

func (n *NewExpr) nodePos() Pos   { return n.Pos }
func (n *NewExpr) exprNode()      {}
func (n *NewExpr) String() string { return "new " + n.ClassName }

// ListLiteral  →  [<expr>, <expr>, ...]
type ListLiteral struct {
	Pos      Pos
	Elements []Expr
}

func (l *ListLiteral) nodePos() Pos   { return l.Pos }
func (l *ListLiteral) exprNode()      {}
func (l *ListLiteral) String() string { return "[...]" }

// FuncCallExpr  →  bare function call: <name>(<args>)
type FuncCallExpr struct {
	Pos  Pos
	Name string
	Args []Expr
}

func (f *FuncCallExpr) nodePos() Pos   { return f.Pos }
func (f *FuncCallExpr) exprNode()      {}
func (f *FuncCallExpr) String() string { return f.Name + "(...)" }
