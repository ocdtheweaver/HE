// Package compiler implements HE's compilation pipeline.
// Pass 4 delivers:
//   - Type inference (infer types without annotations)
//   - Bytecode IR (intermediate representation before native codegen)
//
// The design goal: HE programs compile to typed bytecode,
// then bytecode compiles to native binary or WASM.
package compiler

import (
	"fmt"

	"hunterlang/lang/ast"
)

// ── Inferred types ────────────────────────────────────────────────────────────

type HEType int

const (
	TypeUnknown HEType = iota
	TypeNumber
	TypeText
	TypeBoolean
	TypeList
	TypeObject
	TypeAbility // anonymous function
	TypeNil
)

func (t HEType) String() string {
	switch t {
	case TypeNumber:
		return "number"
	case TypeText:
		return "text"
	case TypeBoolean:
		return "boolean"
	case TypeList:
		return "list"
	case TypeObject:
		return "object"
	case TypeAbility:
		return "ability"
	case TypeNil:
		return "nothing"
	default:
		return "unknown"
	}
}

// TypeInfo holds the inferred type of a named variable.
type TypeInfo struct {
	Name     string
	Type     HEType
	ElemType HEType // for lists
	ObjName  string // for objects — which object definition
}

// TypeError records a type mismatch found during inference.
type TypeError struct {
	Line    int
	Message string
}

func (e TypeError) Error() string {
	return fmt.Sprintf("line %d: %s", e.Line, e.Message)
}

// ── Type environment ──────────────────────────────────────────────────────────

type TypeEnv struct {
	vars   map[string]TypeInfo
	parent *TypeEnv
	errors []TypeError
}

func NewTypeEnv() *TypeEnv {
	env := &TypeEnv{vars: map[string]TypeInfo{}}
	// Built-in vars
	env.Set("pi", TypeInfo{Name: "pi", Type: TypeNumber})
	env.Set("nothing", TypeInfo{Name: "nothing", Type: TypeNil})
	env.Set("platform", TypeInfo{Name: "platform", Type: TypeText})
	env.Set("error", TypeInfo{Name: "error", Type: TypeText})
	return env
}

func (e *TypeEnv) child() *TypeEnv {
	return &TypeEnv{vars: map[string]TypeInfo{}, parent: e}
}

func (e *TypeEnv) Set(name string, info TypeInfo) {
	e.vars[name] = info
}

func (e *TypeEnv) Get(name string) (TypeInfo, bool) {
	if info, ok := e.vars[name]; ok {
		return info, true
	}
	if e.parent != nil {
		return e.parent.Get(name)
	}
	return TypeInfo{Type: TypeUnknown}, false
}

func (e *TypeEnv) AddError(line int, msg string) {
	e.errors = append(e.errors, TypeError{Line: line, Message: msg})
}

func (e *TypeEnv) Errors() []TypeError {
	return e.errors
}

// ── Type Inferencer ───────────────────────────────────────────────────────────

type Inferencer struct {
	env     *TypeEnv
	objects map[string]map[string]TypeInfo // object name → field types

	// protected records every #protected[N] tag found during inference,
	// for visibility in `he check` — no enforcement happens here, this is
	// purely a report of what's tagged and where.
	protected []ProtectedItem
}

// ProtectedItem describes a single #protected[N]-tagged declaration.
type ProtectedItem struct {
	Kind   string // "object" | "ability"
	Object string // owning object name (for abilities) or the object's own name
	Name   string // ability name, or "" for a protected object itself
	Tag    string // "protected", "protected1", "protected2", ...
}

func NewInferencer() *Inferencer {
	return &Inferencer{
		env:     NewTypeEnv(),
		objects: map[string]map[string]TypeInfo{},
	}
}

// InferProgram runs type inference over an entire program.
// Returns a map of variable→type and any type errors found.
func (inf *Inferencer) InferProgram(prog *ast.Program) (map[string]TypeInfo, []TypeError) {
	for _, line := range prog.Lines {
		inf.inferLine(line)
	}
	return inf.env.vars, inf.env.Errors()
}

func (inf *Inferencer) inferLine(line ast.Line) {
	switch l := line.(type) {
	case ast.SummonLine:
		// Module produces an object type
		if l.Alias != nil {
			inf.env.Set(l.Alias.Name, TypeInfo{
				Name:    l.Alias.Name,
				Type:    TypeObject,
				ObjName: l.ModuleName.Lexeme,
			})
		}
	case *ast.ObjectLine:
		fields := map[string]TypeInfo{}
		for _, sec := range l.Body.Sections {
			if ps, ok := sec.(ast.PropertiesSection); ok {
				for _, prop := range ps.Props {
					t := inf.inferExpr(prop.Value)
					fields[prop.Name] = TypeInfo{Name: prop.Name, Type: t}
				}
			}
			if as, ok := sec.(ast.AbilitiesSection); ok {
				for _, act := range as.Abilities {
					if act.ProtectTag != "" {
						inf.protected = append(inf.protected, ProtectedItem{
							Kind:   "ability",
							Object: l.Name.Name,
							Name:   act.Name,
							Tag:    act.ProtectTag,
						})
					}
				}
			}
		}
		if l.ProtectTag != "" {
			inf.protected = append(inf.protected, ProtectedItem{
				Kind:   "object",
				Object: l.Name.Name,
				Name:   "",
				Tag:    l.ProtectTag,
			})
		}
		inf.objects[l.Name.Name] = fields
		inf.env.Set(l.Name.Name, TypeInfo{
			Name:    l.Name.Name,
			Type:    TypeObject,
			ObjName: l.Name.Name,
		})
	case ast.GlobalStatementLine:
		inf.inferStmt(l.Statement)
	case *ast.GlobalStatementLine:
		inf.inferStmt(l.Statement)
	}
}

func (inf *Inferencer) inferStmt(s ast.Statement) {
	switch st := s.(type) {
	case ast.ChangeStmt:
		t := inf.inferExpr(st.Expr)
		inf.env.Set(st.Name, TypeInfo{Name: st.Name, Type: t})
	case ast.MultiAssignStmt:
		for i, name := range st.Names {
			if i < len(st.Exprs) {
				t := inf.inferExpr(st.Exprs[i])
				inf.env.Set(name, TypeInfo{Name: name, Type: t})
			}
		}
	case ast.GrowStmt:
		inf.env.Set(st.Name, TypeInfo{Name: st.Name, Type: TypeNumber})
	case ast.ShrinkStmt:
		inf.env.Set(st.Name, TypeInfo{Name: st.Name, Type: TypeNumber})
	case ast.DecideStmt:
		child := inf.env.child()
		childInf := &Inferencer{env: child, objects: inf.objects}
		for _, b := range st.Then {
			childInf.inferStmt(b)
		}
		for _, b := range st.Else {
			childInf.inferStmt(b)
		}
	case ast.ForEachStmt:
		// The loop variable takes the element type of the list
		listType := inf.inferExpr(st.List)
		elemType := TypeUnknown
		if listType == TypeList {
			elemType = TypeUnknown // would need element type tracking
		}
		child := inf.env.child()
		child.Set(st.VarName, TypeInfo{Name: st.VarName, Type: elemType})
		childInf := &Inferencer{env: child, objects: inf.objects}
		for _, b := range st.Body {
			childInf.inferStmt(b)
		}
	case ast.RangeLoopStmt:
		child := inf.env.child()
		child.Set(st.VarName, TypeInfo{Name: st.VarName, Type: TypeNumber})
		childInf := &Inferencer{env: child, objects: inf.objects}
		for _, b := range st.Body {
			childInf.inferStmt(b)
		}
	case ast.AskStmt:
		inf.env.Set(st.VarName, TypeInfo{Name: st.VarName, Type: TypeText})
	case ast.RecallStmt:
		// Type is unknown at compile time — would need store schema
		inf.env.Set(st.Name, TypeInfo{Name: st.Name, Type: TypeUnknown})
	}
}

func (inf *Inferencer) inferExpr(e ast.Expression) HEType {
	switch ex := e.(type) {
	case ast.NumberLit:
		return TypeNumber
	case ast.StringLit:
		return TypeText
	case ast.BooleanLit:
		return TypeBoolean
	case ast.InterpStringExpr:
		return TypeText
	case ast.ArrayLit:
		return TypeList
	case ast.NamedArgLit:
		return TypeObject
	case ast.AbilityLit:
		return TypeAbility
	case ast.IdentifierExpr:
		if info, ok := inf.env.Get(ex.Name); ok {
			return info.Type
		}
		return TypeUnknown
	case ast.BinaryExpr:
		left := inf.inferExpr(ex.Left)
		right := inf.inferExpr(ex.Right)
		if ex.Op == "+" {
			if left == TypeText || right == TypeText {
				return TypeText
			}
			return TypeNumber
		}
		return TypeNumber
	case ast.CompareExpr, ast.BetweenExpr, ast.LogicAndExpr, ast.LogicOrExpr:
		return TypeBoolean
	case ast.UnaryExpr:
		if ex.Op == "!" {
			return TypeBoolean
		}
		return TypeNumber
	case ast.PowerExpr:
		return TypeNumber
	case ast.ParenExpr:
		return inf.inferExpr(ex.X)
	case ast.FieldAccessExpr:
		if fields, ok := inf.objects[ex.Receiver.Name]; ok {
			if fi, ok := fields[ex.Field]; ok {
				return fi.Type
			}
		}
		return TypeUnknown
	case ast.MethodCallExpr:
		// Known stdlib returns
		switch ex.Receiver.Name {
		case "m": // math
			return TypeNumber
		case "t": // text
			switch ex.Method {
			case "upper", "lower", "replace", "trim", "join", "from":
				return TypeText
			case "contains", "starts", "ends":
				return TypeBoolean
			case "length":
				return TypeNumber
			case "split":
				return TypeList
			}
		case "lst": // list
			switch ex.Method {
			case "length":
				return TypeNumber
			case "contains":
				return TypeBoolean
			case "add", "remove":
				return TypeList
			default:
				return TypeUnknown
			}
		}
		return TypeUnknown
	default:
		return TypeUnknown
	}
}

// Protected returns every #protected[N]-tagged declaration found during
// the most recent InferProgram call.
func (inf *Inferencer) Protected() []ProtectedItem {
	return inf.protected
}

// Report formats inference results as a human-readable summary.
func (inf *Inferencer) Report() string {
	out := "── Type Inference Report ──────────────────\n\n"
	out += "  Variables:\n"
	for name, info := range inf.env.vars {
		if name == "pi" || name == "nothing" || name == "platform" || name == "error" {
			continue
		}
		objHint := ""
		if info.ObjName != "" {
			objHint = fmt.Sprintf(" (%s)", info.ObjName)
		}
		out += fmt.Sprintf("    %-16s : %s%s\n", name, info.Type, objHint)
	}
	if len(inf.objects) > 0 {
		out += "\n  Objects:\n"
		for obj, fields := range inf.objects {
			out += fmt.Sprintf("    %s\n", obj)
			for field, fi := range fields {
				out += fmt.Sprintf("      .%-14s : %s\n", field, fi.Type)
			}
		}
	}
	if len(inf.protected) > 0 {
		out += "\n  Protected (#tag):\n"
		for _, p := range inf.protected {
			switch p.Kind {
			case "object":
				out += fmt.Sprintf("    %-20s  #%s\n", p.Object, p.Tag)
			case "ability":
				out += fmt.Sprintf("    %s.%-17s #%s\n", p.Object, p.Name, p.Tag)
			}
		}
		out += "\n  (parsing only — no enforcement yet; see protected.he policy design)\n"
	}
	if errs := inf.env.Errors(); len(errs) > 0 {
		out += "\n  Type errors:\n"
		for _, e := range errs {
			out += fmt.Sprintf("    ✗ %s\n", e.Error())
		}
	} else {
		out += "\n  No type errors found.\n"
	}
	out += "──────────────────────────────────────────\n"
	return out
}
