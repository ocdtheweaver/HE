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
