package types

import (
	"fmt"
	"strings"

	"hunterlang/lang/ast"
)

type ValueType string

const (
	NilT     ValueType = "nil"
	NumberT  ValueType = "number"
	StringT  ValueType = "text"
	BooleanT ValueType = "boolean"
	ObjectT  ValueType = "object"
	ArrayT   ValueType = "list"
)

type Value struct {
	Type    ValueType
	Number  float64
	Str     string
	Boolean bool
	Object  *Object
	Array   []Value
}

func Nil() Value                    { return Value{Type: NilT} }
func FromNumber(n float64) Value    { return Value{Type: NumberT, Number: n} }
func FromString(s string) Value     { return Value{Type: StringT, Str: s} }
func FromBoolean(b bool) Value      { return Value{Type: BooleanT, Boolean: b} }
func FromObject(o *Object) Value    { return Value{Type: ObjectT, Object: o} }
func FromArray(a []Value) Value     { return Value{Type: ArrayT, Array: a} }

func (v Value) String() string {
	switch v.Type {
	case NilT:
		return "nothing"
	case NumberT:
		if v.Number == float64(int64(v.Number)) {
			return fmt.Sprintf("%d", int64(v.Number))
		}
		return fmt.Sprintf("%g", v.Number)
	case StringT:
		return v.Str
	case BooleanT:
		if v.Boolean {
			return "yes"
		}
		return "no"
	case ObjectT:
		if v.Object != nil {
			return fmt.Sprintf("[%s]", v.Object.Name)
		}
		return "[object]"
	case ArrayT:
		parts := make([]string, len(v.Array))
		for i, el := range v.Array {
			parts[i] = el.String()
		}
		return "[" + strings.Join(parts, ", ") + "]"
	}
	return ""
}

// ── Object ────────────────────────────────────────────────────────────────────

type Object struct {
	Name      string
	Fields    map[string]Value
	Actions   map[string]*Action
	Reactions []Reaction
	// Builtins maps method name to native Go function
	Builtins  map[string]BuiltinFn

	// ProtectTag is "" for unprotected objects, or a tag name like
	// "protected"/"protected1"/"protected2" for whole-object protection
	// (Stage A, Step 3 of the protection roadmap). Every ability called
	// on a protected object is checked against this tag, in addition to
	// any tag the ability itself carries.
	ProtectTag string
}

type BuiltinFn func(args []Value) (Value, error)

type Action struct {
	Name   string
	Params []string
	Body   []ast.Statement

	// ProtectTag is "" for an unprotected ability, or a tag name like
	// "protected1" for a per-ability protection (Stage A, Step 3).
	ProtectTag string
}

type Reaction struct {
	Trigger ast.Trigger
	Body    []ast.Statement
}
