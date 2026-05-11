package types

import "fmt"

type ValueType string

const (
	NumberT  ValueType = "number"
	StringT  ValueType = "string"
	BooleanT ValueType = "boolean"
	ArrayT   ValueType = "array"
	ObjectT  ValueType = "object"
	FuncT    ValueType = "func"
	NilT     ValueType = "nil"
)

type Value struct {
	Type ValueType

	Number  float64
	Str     string
	Boolean bool
	Array   []Value
	Object  *Object
	Func    *Func
}

type Func struct {
	Name string
	// Builtins / closures.
	Impl func(i *Runtime, recv *Object, args []Value) (Value, error)
}

type Object struct {
	Name string
	// user-defined properties/fields
	Fields map[string]Value
	// actions: name -> action definition
	Actions map[string]*Action

	// event-driven reactions: parsed from `on/when/whenever ...`
	Reactions []Reaction
}

func (v Value) String() string {
	switch v.Type {
	case NumberT:
		return fmt.Sprintf("%g", v.Number)
	case StringT:
		return v.Str
	case BooleanT:
		if v.Boolean {
			return "true"
		}
		return "false"
	case ArrayT:
		return fmt.Sprintf("%v", v.Array)
	case ObjectT:
		if v.Object == nil {
			return "object(nil)"
		}
		return "object(" + v.Object.Name + ")"
	case FuncT:
		if v.Func == nil {
			return "func(nil)"
		}
		return "func(" + v.Func.Name + ")"
	case NilT:
		return "nil"
	default:
		return "<?>"
	}
}

func Nil() Value { return Value{Type: NilT} }

func FromNumber(n float64) Value { return Value{Type: NumberT, Number: n} }
func FromString(s string) Value  { return Value{Type: StringT, Str: s} }
func FromBoolean(b bool) Value   { return Value{Type: BooleanT, Boolean: b} }
func FromArray(a []Value) Value  { return Value{Type: ArrayT, Array: a} }
func FromObject(o *Object) Value {
	return Value{Type: ObjectT, Object: o}
}
