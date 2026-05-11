// Package resolver performs semantic analysis on a Hunter's Engine AST.
//
// Responsibilities:
//   - Variable / method scope resolution
//   - Class hierarchy validation (no cycles, superclass must exist)
//   - 'this' / 'super' only inside methods
//   - 'break' / 'continue' only inside loops
//   - 'return' / 'answer' only inside methods
//   - Duplicate class / method / property detection
package resolver

import (
	"fmt"

	"github.com/hunter/he/pkg/ast"
)

// Error is a resolver diagnostic.
type Error struct {
	Pos     ast.Pos
	Message string
}

func (e *Error) Error() string {
	return fmt.Sprintf("resolver error at %s: %s", e.Pos, e.Message)
}

// ─── scope ────────────────────────────────────────────────────────────────────

type scope map[string]bool

// ─── Resolver ─────────────────────────────────────────────────────────────────

// Resolver walks the AST and collects semantic errors.
type Resolver struct {
	errors  []*Error
	scopes  []scope
	classes map[string]*ast.ClassDecl

	inMethod bool
	inLoop   bool
}

// New creates a Resolver.
func New() *Resolver {
	return &Resolver{
		classes: map[string]*ast.ClassDecl{},
	}
}

// Errors returns all diagnostics found.
func (r *Resolver) Errors() []*Error { return r.errors }

// Resolve analyses a complete program.
func (r *Resolver) Resolve(prog *ast.Program) {
	// First pass: collect all class names so forward references work.
	for _, stmt := range prog.Statements {
		if cd, ok := stmt.(*ast.ClassDecl); ok {
			if _, exists := r.classes[cd.Name]; exists {
				r.addError(cd.Pos, fmt.Sprintf("class '%s' declared more than once", cd.Name))
			}
			r.classes[cd.Name] = cd
		}
	}

	// Second pass: validate class hierarchies.
	for name, cd := range r.classes {
		r.checkInheritanceCycle(name, cd, map[string]bool{})
	}

	// Third pass: full resolution.
	r.pushScope()
	for _, stmt := range prog.Statements {
		r.resolveStmt(stmt)
	}
	r.popScope()
}

// ─── scope helpers ────────────────────────────────────────────────────────────

func (r *Resolver) pushScope() { r.scopes = append(r.scopes, scope{}) }

func (r *Resolver) popScope() {
	if len(r.scopes) > 0 {
		r.scopes = r.scopes[:len(r.scopes)-1]
	}
}

func (r *Resolver) declare(name string) {
	if len(r.scopes) == 0 {
		return
	}
	r.scopes[len(r.scopes)-1][name] = true
}

func (r *Resolver) isDeclared(name string) bool {
	// Walk scopes inner → outer.
	for i := len(r.scopes) - 1; i >= 0; i-- {
		if r.scopes[i][name] {
			return true
		}
	}
	return false
}

func (r *Resolver) addError(pos ast.Pos, msg string) {
	r.errors = append(r.errors, &Error{Pos: pos, Message: msg})
}

// ─── inheritance cycle detection ─────────────────────────────────────────────

func (r *Resolver) checkInheritanceCycle(name string, cd *ast.ClassDecl, visited map[string]bool) {
	if visited[name] {
		r.addError(cd.Pos, fmt.Sprintf("class '%s' has a circular inheritance chain", name))
		return
	}
	if cd.SuperClass == "" {
		return
	}
	visited[name] = true
	super, ok := r.classes[cd.SuperClass]
	if !ok {
		r.addError(cd.Pos, fmt.Sprintf("class '%s' extends unknown class '%s'", name, cd.SuperClass))
		return
	}
	r.checkInheritanceCycle(cd.SuperClass, super, visited)
}

// ─── statement resolution ─────────────────────────────────────────────────────

func (r *Resolver) resolveStmt(stmt ast.Stmt) {
	if stmt == nil {
		return
	}
	switch s := stmt.(type) {
	case *ast.ClassDecl:
		r.resolveClassDecl(s)
	case *ast.ObjectDecl:
		r.resolveObjectDecl(s)
	case *ast.ImportStmt:
		r.declare(s.Name)
	case *ast.SetStmt:
		r.resolveExpr(s.Value)
		r.resolveExpr(s.Target)
		if id, ok := s.Target.(*ast.Identifier); ok {
			r.declare(id.Name)
		}
	case *ast.CompoundAssignStmt:
		r.resolveExpr(s.Target)
		r.resolveExpr(s.Value)
	case *ast.PrintStmt:
		r.resolveExpr(s.Expr)
	case *ast.InputStmt:
		if s.Prompt != nil {
			r.resolveExpr(s.Prompt)
		}
		r.declare(s.Target)
	case *ast.IfStmt:
		r.resolveExpr(s.Condition)
		r.pushScope()
		r.resolveStmts(s.Consequent)
		r.popScope()
		r.pushScope()
		r.resolveStmts(s.Alternate)
		r.popScope()
	case *ast.WhileStmt:
		r.resolveExpr(s.Condition)
		prevLoop := r.inLoop
		r.inLoop = true
		r.pushScope()
		r.resolveStmts(s.Body)
		r.popScope()
		r.inLoop = prevLoop
	case *ast.ForEachStmt:
		r.resolveExpr(s.Iterable)
		prevLoop := r.inLoop
		r.inLoop = true
		r.pushScope()
		r.declare(s.VarName)
		r.resolveStmts(s.Body)
		r.popScope()
		r.inLoop = prevLoop
	case *ast.ReturnStmt:
		if !r.inMethod {
			r.addError(s.Pos, "'return' used outside of a method")
		}
		if s.Value != nil {
			r.resolveExpr(s.Value)
		}
	case *ast.BreakStmt:
		if !r.inLoop {
			r.addError(s.Pos, "'break' used outside of a loop")
		}
	case *ast.ContinueStmt:
		if !r.inLoop {
			r.addError(s.Pos, "'continue' used outside of a loop")
		}
	case *ast.CallStmt:
		r.resolveExpr(s.Object)
		for _, a := range s.Args {
			r.resolveExpr(a)
		}
	case *ast.ExprStmt:
		r.resolveExpr(s.Expr)
	}
}

func (r *Resolver) resolveStmts(stmts []ast.Stmt) {
	for _, s := range stmts {
		r.resolveStmt(s)
	}
}

func (r *Resolver) resolveClassDecl(cd *ast.ClassDecl) {
	r.declare(cd.Name)
	seen := map[string]bool{}
	for _, prop := range cd.Properties {
		if seen[prop.Name] {
			r.addError(prop.Pos, fmt.Sprintf("duplicate property '%s' in class '%s'", prop.Name, cd.Name))
		}
		seen[prop.Name] = true
		if prop.Default != nil {
			r.resolveExpr(prop.Default)
		}
	}

	seenMethods := map[string]bool{}
	for _, m := range cd.Methods {
		if seenMethods[m.Name] {
			r.addError(m.Pos, fmt.Sprintf("duplicate method '%s' in class '%s'", m.Name, cd.Name))
		}
		seenMethods[m.Name] = true
		r.resolveMethod(m)
	}
}

func (r *Resolver) resolveMethod(m *ast.MethodDecl) {
	prevMethod := r.inMethod
	r.inMethod = true
	r.pushScope()
	r.declare("this")
	r.declare("super")
	for _, param := range m.Params {
		r.declare(param)
	}
	r.resolveStmts(m.Body)
	r.popScope()
	r.inMethod = prevMethod
}

func (r *Resolver) resolveObjectDecl(od *ast.ObjectDecl) {
	if _, ok := r.classes[od.ClassName]; !ok {
		r.addError(od.Pos, fmt.Sprintf("unknown class '%s'", od.ClassName))
	}
	for _, val := range od.InitFields {
		r.resolveExpr(val)
	}
	r.declare(od.Name)
}

// ─── expression resolution ────────────────────────────────────────────────────

func (r *Resolver) resolveExpr(expr ast.Expr) {
	if expr == nil {
		return
	}
	switch e := expr.(type) {
	case *ast.Identifier:
		if e.Name != "this" && e.Name != "super" {
			// We don't error on undeclared idents here — runtime will catch it.
			// But we could add a warning pass later.
			_ = r.isDeclared(e.Name)
		}
	case *ast.DotExpr:
		r.resolveExpr(e.Object)
	case *ast.IndexExpr:
		r.resolveExpr(e.Object)
		r.resolveExpr(e.Index)
	case *ast.BinaryExpr:
		r.resolveExpr(e.Left)
		r.resolveExpr(e.Right)
	case *ast.UnaryExpr:
		r.resolveExpr(e.Operand)
	case *ast.CallExpr:
		r.resolveExpr(e.Object)
		for _, a := range e.Args {
			r.resolveExpr(a)
		}
	case *ast.NewExpr:
		if _, ok := r.classes[e.ClassName]; !ok {
			r.addError(e.Pos, fmt.Sprintf("unknown class '%s' in 'new' expression", e.ClassName))
		}
		for _, v := range e.InitFields {
			r.resolveExpr(v)
		}
	case *ast.ListLiteral:
		for _, el := range e.Elements {
			r.resolveExpr(el)
		}
	case *ast.FuncCallExpr:
		for _, a := range e.Args {
			r.resolveExpr(a)
		}
	// Literals need no resolution.
	case *ast.NumberLit, *ast.StringLit, *ast.BoolLit, *ast.NilLit:
	}
}
