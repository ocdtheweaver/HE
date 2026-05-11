package types

import "hunterlang/lang/ast"

// Runtime is the execution context used by the evaluator.
// It is intentionally small for now; it will grow as we implement execution.
type Runtime struct{}

// Action represents an object action (ability method).
type Action struct {
	Name string
	Body []ast.Statement
}

// Reaction represents event-driven behavior attached to an object.
type Reaction struct {
	Trigger ast.Trigger
	Body    []ast.Statement
}
