package ast

import (
	"fmt"
	"strings"
)

// ===============================
// AST Nodes - Core Interfaces
// ===============================

type Node interface {
	Pos() Position
	String() string
}

type Statement interface {
	Node
	stmtNode()
}

type Expression interface {
	Node
	exprNode()
}

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

type Program struct {
	Posn       Position
	Summons    []*SummonDecl // summon statements
	Modules    []*ModuleDecl // top-level modules
	Objects    []*ObjectDecl // object definitions
	Assets     []*AssetDecl  // asset declarations
	Statements []Statement   // global statements
}

func (p *Program) Pos() Position { return p.Posn }
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

type ModuleDecl struct {
	Posn       Position
	Name       string          // module name
	Submodules []*ModuleDecl   // nested modules
	Functions  []*FunctionDecl // functions in this module
	Body       []Declaration   // other declarations (modules, objects, assets)
}

func (m *ModuleDecl) Pos() Position { return m.Posn }
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

// ===============================
// Function Decl
// ===============================

type FunctionDecl struct {
	Posn       Position
	Name       string
	Params     []*ParamDecl
	Body       []Statement
	ReturnType string // optional return type
}

func (f *FunctionDecl) Pos() Position { return f.Posn }
func (f *FunctionDecl) declNode()     {}
func (f *FunctionDecl) String() string {
	params := make([]string, len(f.Params))
	for i, p := range f.Params {
		params[i] = p.String()
	}
	return fmt.Sprintf("fn %s(%s) [%d statements]", f.Name, strings.Join(params, ", "), len(f.Body))
}

type ParamDecl struct {
	Posn Position
	Name string
	Type string // optional type annotation
}

func (p *ParamDecl) Pos() Position { return p.Posn }
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

type ObjectDecl struct {
	Posn       Position
	Name       string
	Like       string // parent object (inheritance)
	Properties []*PropertyDecl
	Abilities  []*AbilityDecl
	Reactions  []*ReactionDecl
	Memories   []*MemoryDecl
}

func (o *ObjectDecl) Pos() Position { return o.Posn }
func (o *ObjectDecl) declNode()     {}
func (o *ObjectDecl) String() string {
	return fmt.Sprintf("object %s (like: %s, props: %d, abilities: %d)",
		o.Name, o.Like, len(o.Properties), len(o.Abilities))
}

type PropertyDecl struct {
	Posn  Position
	Name  string
	Value Expression
}

func (p *PropertyDecl) Pos() Position { return p.Posn }
func (p *PropertyDecl) declNode()     {}
func (p *PropertyDecl) String() string {
	return fmt.Sprintf("%s = %s", p.Name, p.Value)
}

type AbilityDecl struct {
	Posn   Position
	Name   string
	Params []*ParamDecl
	Body   []Statement
}

func (a *AbilityDecl) Pos() Position { return a.Posn }
func (a *AbilityDecl) declNode()     {}
func (a *AbilityDecl) String() string {
	return fmt.Sprintf("ability %s(%d params) [%d statements]",
		a.Name, len(a.Params), len(a.Body))
}

type ReactionDecl struct {
	Posn    Position
	Trigger string
	Body    []Statement
}

func (r *ReactionDecl) Pos() Position { return r.Posn }
func (r *ReactionDecl) declNode()     {}
func (r *ReactionDecl) String() string {
	return fmt.Sprintf("on %s [%d statements]", r.Trigger, len(r.Body))
}

type MemoryDecl struct {
	Posn  Position
	Name  string
	Value Expression
}

func (m *MemoryDecl) Pos() Position { return m.Posn }
func (m *MemoryDecl) declNode()     {}
func (m *MemoryDecl) String() string {
	return fmt.Sprintf("memory %s = %s", m.Name, m.Value)
}

// ===============================
// Asset System
// ===============================

type SummonDecl struct {
	Posn Position
	Path string
	As   string // alias
}

func (s *SummonDecl) Pos() Position { return s.Posn }
func (s *SummonDecl) declNode()     {}
func (s *SummonDecl) String() string {
	if s.As != "" {
		return fmt.Sprintf("summon \"%s\" as %s", s.Path, s.As)
	}
	return fmt.Sprintf("summon \"%s\"", s.Path)
}

type AssetDecl struct {
	Posn Position
	Type string // "image", "sound", "music", etc.
	Path string
	Name string // optional name/alias
}

func (a *AssetDecl) Pos() Position { return a.Posn }
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

type PrintStmt struct {
	Posn Position
	Expr Expression
}

func (p *PrintStmt) Pos() Position  { return p.Posn }
func (p *PrintStmt) stmtNode()      {}
func (p *PrintStmt) String() string { return fmt.Sprintf("print %s", p.Expr) }

type AssignStmt struct {
	Posn Position
	Lhs  string
	Expr Expression
}

func (a *AssignStmt) Pos() Position  { return a.Posn }
func (a *AssignStmt) stmtNode()      {}
func (a *AssignStmt) String() string { return fmt.Sprintf("%s = %s", a.Lhs, a.Expr) }

type CallStmt struct {
	Posn   Position
	Target Expression // can be identifier or member expression
	Args   []Expression
}

func (c *CallStmt) Pos() Position { return c.Posn }
func (c *CallStmt) stmtNode()     {}
func (c *CallStmt) String() string {
	args := make([]string, len(c.Args))
	for i, a := range c.Args {
		args[i] = a.String()
	}
	return fmt.Sprintf("call %s(%s)", c.Target, strings.Join(args, ", "))
}

// ExprStmt represents a standalone expression statement (e.g. a bare call).
type ExprStmt struct {
	Posn Position
	Expr Expression
}

func (e *ExprStmt) Pos() Position  { return e.Posn }
func (e *ExprStmt) stmtNode()      {}
func (e *ExprStmt) String() string { return e.Expr.String() }

type WaitStmt struct {
	Posn    Position
	Seconds Expression
}

func (w *WaitStmt) Pos() Position  { return w.Posn }
func (w *WaitStmt) stmtNode()      {}
func (w *WaitStmt) String() string { return fmt.Sprintf("wait %s seconds", w.Seconds) }

type IfStmt struct {
	Posn Position
	Cond Expression
	Then []Statement
	Else []Statement
}

func (i *IfStmt) Pos() Position { return i.Posn }
func (i *IfStmt) stmtNode()     {}
func (i *IfStmt) String() string {
	return fmt.Sprintf("if %s then [%d] else [%d]", i.Cond, len(i.Then), len(i.Else))
}

type RepeatStmt struct {
	Posn  Position
	Cond  Expression
	Count Expression
	Body  []Statement
}

func (r *RepeatStmt) Pos() Position { return r.Posn }
func (r *RepeatStmt) stmtNode()     {}
func (r *RepeatStmt) String() string {
	if r.Cond != nil {
		return fmt.Sprintf("while %s [%d]", r.Cond, len(r.Body))
	}
	return fmt.Sprintf("repeat %s [%d]", r.Count, len(r.Body))
}

type ReturnStmt struct {
	Posn Position
	Expr Expression
}

func (r *ReturnStmt) Pos() Position  { return r.Posn }
func (r *ReturnStmt) stmtNode()      {}
func (r *ReturnStmt) String() string { return fmt.Sprintf("return %s", r.Expr) }

// ===============================
// Expressions
// ===============================

type IdentExpr struct {
	Posn Position
	Name string
}

func (i *IdentExpr) Pos() Position  { return i.Posn }
func (i *IdentExpr) exprNode()      {}
func (i *IdentExpr) String() string { return i.Name }

type MemberExpr struct {
	Posn   Position
	Object Expression
	Member string
}

func (m *MemberExpr) Pos() Position { return m.Posn }
func (m *MemberExpr) exprNode()     {}
func (m *MemberExpr) String() string {
	return fmt.Sprintf("%s.%s", m.Object, m.Member)
}

type CallExpr struct {
	Posn   Position
	Target Expression
	Args   []Expression
}

func (c *CallExpr) Pos() Position { return c.Posn }
func (c *CallExpr) exprNode()     {}
func (c *CallExpr) String() string {
	args := make([]string, len(c.Args))
	for i, a := range c.Args {
		args[i] = a.String()
	}
	return fmt.Sprintf("%s(%s)", c.Target, strings.Join(args, ", "))
}

type NumberExpr struct {
	Posn  Position
	Value float64
}

func (n *NumberExpr) Pos() Position  { return n.Posn }
func (n *NumberExpr) exprNode()      {}
func (n *NumberExpr) String() string { return fmt.Sprintf("%v", n.Value) }

type StringExpr struct {
	Posn  Position
	Value string
}

func (s *StringExpr) Pos() Position  { return s.Posn }
func (s *StringExpr) exprNode()      {}
func (s *StringExpr) String() string { return fmt.Sprintf("\"%s\"", s.Value) }

type BoolExpr struct {
	Posn  Position
	Value bool
}

func (b *BoolExpr) Pos() Position  { return b.Posn }
func (b *BoolExpr) exprNode()      {}
func (b *BoolExpr) String() string { return fmt.Sprintf("%v", b.Value) }

type BinaryExpr struct {
	Posn Position
	Op   string
	Lhs  Expression
	Rhs  Expression
}

func (b *BinaryExpr) Pos() Position  { return b.Posn }
func (b *BinaryExpr) exprNode()      {}
func (b *BinaryExpr) String() string { return fmt.Sprintf("(%s %s %s)", b.Lhs, b.Op, b.Rhs) }

type UnaryExpr struct {
	Posn Position
	Op   string
	Expr Expression
}

func (u *UnaryExpr) Pos() Position  { return u.Posn }
func (u *UnaryExpr) exprNode()      {}
func (u *UnaryExpr) String() string { return fmt.Sprintf("(%s%s)", u.Op, u.Expr) }

type ArrayExpr struct {
	Posn   Position
	Values []Expression
}

func (a *ArrayExpr) Pos() Position { return a.Posn }
func (a *ArrayExpr) exprNode()     {}
func (a *ArrayExpr) String() string {
	values := make([]string, len(a.Values))
	for i, v := range a.Values {
		values[i] = v.String()
	}
	return fmt.Sprintf("[%s]", strings.Join(values, ", "))
}

type NilExpr struct {
	Posn Position
}

func (n *NilExpr) Pos() Position  { return n.Posn }
func (n *NilExpr) exprNode()      {}
func (n *NilExpr) String() string { return "nil" }
