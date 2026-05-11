// Package runtime executes a Hunter's Engine AST.
//
// Value system
// ─────────────
//   HeNumber   → float64
//   HeString   → string
//   HeBool     → bool
//   HeNil      → nil sentinel
//   HeList     → []HeValue
//   HeObject   → instance of a user-defined class
//   HeClass    → the class descriptor itself (first-class)
//   HeBuiltin  → a Go function exposed to HE code
//
// Execution model
// ────────────────
// The interpreter uses a chain of *Environment frames (one per block/method
// call) and panics with signal values (returnSignal, breakSignal,
// continueSignal) to implement non-local control flow. The Execute* methods
// recover from these at the appropriate boundary.
package runtime

import (
	"fmt"
	"math"
	"strings"

	"github.com/hunter/he/pkg/ast"
)

// ─── value types ──────────────────────────────────────────────────────────────

// HeValue is any value that can live in a HE variable.
type HeValue interface {
	heValue()
	Repr() string // human-readable representation
}

// HeNumber wraps a float64.
type HeNumber struct{ V float64 }

func (n HeNumber) heValue() {}
func (n HeNumber) Repr() string {
	if n.V == math.Trunc(n.V) && !math.IsInf(n.V, 0) {
		return fmt.Sprintf("%g", n.V)
	}
	return fmt.Sprintf("%g", n.V)
}

// HeString wraps a Go string.
type HeString struct{ V string }

func (s HeString) heValue() {}
func (s HeString) Repr() string { return s.V }

// HeBool wraps a bool.
type HeBool struct{ V bool }

func (b HeBool) heValue() {}
func (b HeBool) Repr() string {
	if b.V {
		return "true"
	}
	return "false"
}

// HeNil is the nil/nothing value.
type HeNil struct{}

func (n HeNil) heValue()      {}
func (n HeNil) Repr() string  { return "nil" }

// HeList is an ordered, mutable list.
type HeList struct{ Elements []HeValue }

func (l *HeList) heValue() {}
func (l *HeList) Repr() string {
	parts := make([]string, len(l.Elements))
	for i, e := range l.Elements {
		parts[i] = e.Repr()
	}
	return "[" + strings.Join(parts, ", ") + "]"
}

// HeClass is a class descriptor.
type HeClass struct {
	Name       string
	SuperClass *HeClass
	Properties []*ast.PropertyDecl
	Methods    map[string]*ast.MethodDecl
}

func (c *HeClass) heValue() {}
func (c *HeClass) Repr() string { return "<class " + c.Name + ">" }

// LookupMethod searches the method on the class and its ancestors.
func (c *HeClass) LookupMethod(name string) (*ast.MethodDecl, *HeClass) {
	for cls := c; cls != nil; cls = cls.SuperClass {
		if m, ok := cls.Methods[name]; ok {
			return m, cls
		}
	}
	return nil, nil
}

// HeObject is an instance of a HeClass.
type HeObject struct {
	Class      *HeClass
	Properties map[string]HeValue
}

func (o *HeObject) heValue() {}
func (o *HeObject) Repr() string {
	return fmt.Sprintf("<object of %s>", o.Class.Name)
}

// Get fetches a property (walks the class hierarchy for defaults).
func (o *HeObject) Get(name string) (HeValue, bool) {
	if v, ok := o.Properties[name]; ok {
		return v, true
	}
	return nil, false
}

// Set stores a property value.
func (o *HeObject) Set(name string, v HeValue) { o.Properties[name] = v }

// HeBuiltin is a Go function callable from HE code.
type HeBuiltin struct {
	Name string
	Fn   func(args []HeValue) (HeValue, error)
}

func (b *HeBuiltin) heValue() {}
func (b *HeBuiltin) Repr() string { return "<builtin " + b.Name + ">" }

// ─── control-flow signals (panic values) ─────────────────────────────────────

type returnSignal struct{ value HeValue }
type breakSignal struct{}
type continueSignal struct{}

// ─── environment (scope chain) ────────────────────────────────────────────────

// Environment is a single scope frame in the variable lookup chain.
type Environment struct {
	vars   map[string]HeValue
	parent *Environment
}

// NewEnvironment creates a top-level environment.
func NewEnvironment() *Environment {
	return &Environment{vars: map[string]HeValue{}}
}

// NewChild creates a child scope.
func (e *Environment) NewChild() *Environment {
	return &Environment{vars: map[string]HeValue{}, parent: e}
}

// Get looks up a variable walking outer scopes.
func (e *Environment) Get(name string) (HeValue, bool) {
	if v, ok := e.vars[name]; ok {
		return v, true
	}
	if e.parent != nil {
		return e.parent.Get(name)
	}
	return nil, false
}

// Set assigns to an existing variable in the nearest enclosing scope that
// has it, or creates it in the current scope.
func (e *Environment) Set(name string, v HeValue) {
	if env := e.findOwner(name); env != nil {
		env.vars[name] = v
		return
	}
	e.vars[name] = v
}

// Define always creates the variable in the current scope.
func (e *Environment) Define(name string, v HeValue) { e.vars[name] = v }

func (e *Environment) findOwner(name string) *Environment {
	if _, ok := e.vars[name]; ok {
		return e
	}
	if e.parent != nil {
		return e.parent.findOwner(name)
	}
	return nil
}

// ─── runtime error ────────────────────────────────────────────────────────────

// RuntimeError is a fatal execution error.
type RuntimeError struct {
	Pos     ast.Pos
	Message string
}

func (e *RuntimeError) Error() string {
	return fmt.Sprintf("runtime error at %s: %s", e.Pos, e.Message)
}

func runtimeErr(pos ast.Pos, format string, args ...any) *RuntimeError {
	return &RuntimeError{Pos: pos, Message: fmt.Sprintf(format, args...)}
}
