package runtime

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/hunter/he/pkg/ast"
)

// Interpreter executes a HE program.
type Interpreter struct {
	global *Environment
	reader *bufio.Reader
}

// New creates an Interpreter with all builtins registered.
func New() *Interpreter {
	interp := &Interpreter{
		global: NewEnvironment(),
		reader: bufio.NewReader(os.Stdin),
	}
	interp.registerBuiltins()
	return interp
}

// GlobalEnv returns the interpreter's global environment (useful for testing).
func (interp *Interpreter) GlobalEnv() *Environment { return interp.global }

// Run executes a parsed program. Returns the first RuntimeError if any.
func (interp *Interpreter) Run(prog *ast.Program) (err error) {
	defer func() {
		if r := recover(); r != nil {
			switch sig := r.(type) {
			case *RuntimeError:
				err = sig
			case returnSignal:
				// top-level return: ignore value
			default:
				err = fmt.Errorf("unexpected panic: %v", r)
			}
		}
	}()
	interp.execStmts(prog.Statements, interp.global)
	return nil
}

// ─── builtins ─────────────────────────────────────────────────────────────────

func (interp *Interpreter) registerBuiltins() {
	builtins := map[string]*HeBuiltin{
		"len": {
			Name: "len",
			Fn: func(args []HeValue) (HeValue, error) {
				if len(args) != 1 {
					return nil, fmt.Errorf("len() takes 1 argument")
				}
				switch v := args[0].(type) {
				case *HeList:
					return HeNumber{V: float64(len(v.Elements))}, nil
				case HeString:
					return HeNumber{V: float64(len([]rune(v.V)))}, nil
				}
				return nil, fmt.Errorf("len() requires a list or string")
			},
		},
		"toNumber": {
			Name: "toNumber",
			Fn: func(args []HeValue) (HeValue, error) {
				if len(args) != 1 {
					return nil, fmt.Errorf("toNumber() takes 1 argument")
				}
				switch v := args[0].(type) {
				case HeNumber:
					return v, nil
				case HeString:
					f, err := strconv.ParseFloat(strings.TrimSpace(v.V), 64)
					if err != nil {
						return nil, fmt.Errorf("cannot convert %q to number", v.V)
					}
					return HeNumber{V: f}, nil
				case HeBool:
					if v.V {
						return HeNumber{V: 1}, nil
					}
					return HeNumber{V: 0}, nil
				}
				return nil, fmt.Errorf("cannot convert to number")
			},
		},
		"toString": {
			Name: "toString",
			Fn: func(args []HeValue) (HeValue, error) {
				if len(args) != 1 {
					return nil, fmt.Errorf("toString() takes 1 argument")
				}
				return HeString{V: args[0].Repr()}, nil
			},
		},
		"typeOf": {
			Name: "typeOf",
			Fn: func(args []HeValue) (HeValue, error) {
				if len(args) != 1 {
					return nil, fmt.Errorf("typeOf() takes 1 argument")
				}
				switch args[0].(type) {
				case HeNumber:
					return HeString{V: "number"}, nil
				case HeString:
					return HeString{V: "string"}, nil
				case HeBool:
					return HeString{V: "boolean"}, nil
				case HeNil:
					return HeString{V: "nil"}, nil
				case *HeList:
					return HeString{V: "list"}, nil
				case *HeObject:
					return HeString{V: "object"}, nil
				case *HeClass:
					return HeString{V: "class"}, nil
				case *HeBuiltin:
					return HeString{V: "builtin"}, nil
				}
				return HeString{V: "unknown"}, nil
			},
		},
		"append": {
			Name: "append",
			Fn: func(args []HeValue) (HeValue, error) {
				if len(args) < 2 {
					return nil, fmt.Errorf("append() takes at least 2 arguments")
				}
				list, ok := args[0].(*HeList)
				if !ok {
					return nil, fmt.Errorf("append() first argument must be a list")
				}
				list.Elements = append(list.Elements, args[1:]...)
				return list, nil
			},
		},
		"remove": {
			Name: "remove",
			Fn: func(args []HeValue) (HeValue, error) {
				if len(args) != 2 {
					return nil, fmt.Errorf("remove() takes 2 arguments")
				}
				list, ok := args[0].(*HeList)
				if !ok {
					return nil, fmt.Errorf("remove() first argument must be a list")
				}
				idx, ok2 := args[1].(HeNumber)
				if !ok2 {
					return nil, fmt.Errorf("remove() second argument must be a number (index)")
				}
				i := int(idx.V)
				if i < 0 || i >= len(list.Elements) {
					return nil, fmt.Errorf("remove() index %d out of range", i)
				}
				list.Elements = append(list.Elements[:i], list.Elements[i+1:]...)
				return list, nil
			},
		},
		"contains": {
			Name: "contains",
			Fn: func(args []HeValue) (HeValue, error) {
				if len(args) != 2 {
					return nil, fmt.Errorf("contains() takes 2 arguments")
				}
				switch col := args[0].(type) {
				case *HeList:
					for _, el := range col.Elements {
						if heEqual(el, args[1]) {
							return HeBool{V: true}, nil
						}
					}
					return HeBool{V: false}, nil
				case HeString:
					needle, ok := args[1].(HeString)
					if !ok {
						return HeBool{V: false}, nil
					}
					return HeBool{V: strings.Contains(col.V, needle.V)}, nil
				}
				return nil, fmt.Errorf("contains() requires a list or string")
			},
		},
		"sqrt": {
			Name: "sqrt",
			Fn: func(args []HeValue) (HeValue, error) {
				if len(args) != 1 {
					return nil, fmt.Errorf("sqrt() takes 1 argument")
				}
				n, ok := args[0].(HeNumber)
				if !ok {
					return nil, fmt.Errorf("sqrt() requires a number")
				}
				return HeNumber{V: math.Sqrt(n.V)}, nil
			},
		},
		"abs": {
			Name: "abs",
			Fn: func(args []HeValue) (HeValue, error) {
				if len(args) != 1 {
					return nil, fmt.Errorf("abs() takes 1 argument")
				}
				n, ok := args[0].(HeNumber)
				if !ok {
					return nil, fmt.Errorf("abs() requires a number")
				}
				return HeNumber{V: math.Abs(n.V)}, nil
			},
		},
		"floor": {
			Name: "floor",
			Fn: func(args []HeValue) (HeValue, error) {
				if len(args) != 1 {
					return nil, fmt.Errorf("floor() takes 1 argument")
				}
				n, ok := args[0].(HeNumber)
				if !ok {
					return nil, fmt.Errorf("floor() requires a number")
				}
				return HeNumber{V: math.Floor(n.V)}, nil
			},
		},
		"ceil": {
			Name: "ceil",
			Fn: func(args []HeValue) (HeValue, error) {
				if len(args) != 1 {
					return nil, fmt.Errorf("ceil() takes 1 argument")
				}
				n, ok := args[0].(HeNumber)
				if !ok {
					return nil, fmt.Errorf("ceil() requires a number")
				}
				return HeNumber{V: math.Ceil(n.V)}, nil
			},
		},
		"upper": {
			Name: "upper",
			Fn: func(args []HeValue) (HeValue, error) {
				if len(args) != 1 {
					return nil, fmt.Errorf("upper() takes 1 argument")
				}
				s, ok := args[0].(HeString)
				if !ok {
					return nil, fmt.Errorf("upper() requires a string")
				}
				return HeString{V: strings.ToUpper(s.V)}, nil
			},
		},
		"lower": {
			Name: "lower",
			Fn: func(args []HeValue) (HeValue, error) {
				if len(args) != 1 {
					return nil, fmt.Errorf("lower() takes 1 argument")
				}
				s, ok := args[0].(HeString)
				if !ok {
					return nil, fmt.Errorf("lower() requires a string")
				}
				return HeString{V: strings.ToLower(s.V)}, nil
			},
		},
		"trim": {
			Name: "trim",
			Fn: func(args []HeValue) (HeValue, error) {
				if len(args) != 1 {
					return nil, fmt.Errorf("trim() takes 1 argument")
				}
				s, ok := args[0].(HeString)
				if !ok {
					return nil, fmt.Errorf("trim() requires a string")
				}
				return HeString{V: strings.TrimSpace(s.V)}, nil
			},
		},
		"split": {
			Name: "split",
			Fn: func(args []HeValue) (HeValue, error) {
				if len(args) != 2 {
					return nil, fmt.Errorf("split() takes 2 arguments")
				}
				s, ok1 := args[0].(HeString)
				sep, ok2 := args[1].(HeString)
				if !ok1 || !ok2 {
					return nil, fmt.Errorf("split() requires two strings")
				}
				parts := strings.Split(s.V, sep.V)
				elems := make([]HeValue, len(parts))
				for i, p := range parts {
					elems[i] = HeString{V: p}
				}
				return &HeList{Elements: elems}, nil
			},
		},
		"join": {
			Name: "join",
			Fn: func(args []HeValue) (HeValue, error) {
				if len(args) != 2 {
					return nil, fmt.Errorf("join() takes 2 arguments")
				}
				list, ok1 := args[0].(*HeList)
				sep, ok2 := args[1].(HeString)
				if !ok1 || !ok2 {
					return nil, fmt.Errorf("join() requires a list and a string")
				}
				parts := make([]string, len(list.Elements))
				for i, e := range list.Elements {
					parts[i] = e.Repr()
				}
				return HeString{V: strings.Join(parts, sep.V)}, nil
			},
		},
	}

	for name, b := range builtins {
		interp.global.Define(name, b)
	}
}

// ─── statement execution ──────────────────────────────────────────────────────

func (interp *Interpreter) execStmts(stmts []ast.Stmt, env *Environment) {
	for _, s := range stmts {
		interp.execStmt(s, env)
	}
}

func (interp *Interpreter) execStmt(stmt ast.Stmt, env *Environment) {
	if stmt == nil {
		return
	}
	switch s := stmt.(type) {
	case *ast.ClassDecl:
		interp.execClassDecl(s, env)
	case *ast.ObjectDecl:
		interp.execObjectDecl(s, env)
	case *ast.ImportStmt:
		// Module system placeholder — silently ignored for now.
	case *ast.SetStmt:
		val := interp.evalExpr(s.Value, env)
		interp.assignTarget(s.Target, val, env)
	case *ast.CompoundAssignStmt:
		interp.execCompoundAssign(s, env)
	case *ast.PrintStmt:
		val := interp.evalExpr(s.Expr, env)
		fmt.Println(val.Repr())
	case *ast.InputStmt:
		interp.execInput(s, env)
	case *ast.IfStmt:
		interp.execIf(s, env)
	case *ast.WhileStmt:
		interp.execWhile(s, env)
	case *ast.ForEachStmt:
		interp.execForEach(s, env)
	case *ast.ReturnStmt:
		var val HeValue = HeNil{}
		if s.Value != nil {
			val = interp.evalExpr(s.Value, env)
		}
		panic(returnSignal{value: val})
	case *ast.BreakStmt:
		panic(breakSignal{})
	case *ast.ContinueStmt:
		panic(continueSignal{})
	case *ast.CallStmt:
		interp.execCallStmt(s, env)
	case *ast.ExprStmt:
		interp.evalExpr(s.Expr, env)
	}
}

func (interp *Interpreter) execClassDecl(cd *ast.ClassDecl, env *Environment) {
	var super *HeClass
	if cd.SuperClass != "" {
		v, ok := env.Get(cd.SuperClass)
		if !ok {
			panic(runtimeErr(cd.Pos, "class '%s' not found", cd.SuperClass))
		}
		cls, ok := v.(*HeClass)
		if !ok {
			panic(runtimeErr(cd.Pos, "'%s' is not a class", cd.SuperClass))
		}
		super = cls
	}

	methods := map[string]*ast.MethodDecl{}
	for _, m := range cd.Methods {
		methods[m.Name] = m
	}

	cls := &HeClass{
		Name:       cd.Name,
		SuperClass: super,
		Properties: cd.Properties,
		Methods:    methods,
	}
	env.Define(cd.Name, cls)
}

func (interp *Interpreter) execObjectDecl(od *ast.ObjectDecl, env *Environment) {
	obj := interp.instantiate(od.ClassName, od.InitFields, od.Pos, env)
	env.Define(od.Name, obj)
}

func (interp *Interpreter) instantiate(className string, initFields map[string]ast.Expr, pos ast.Pos, env *Environment) *HeObject {
	v, ok := env.Get(className)
	if !ok {
		panic(runtimeErr(pos, "class '%s' not found", className))
	}
	cls, ok := v.(*HeClass)
	if !ok {
		panic(runtimeErr(pos, "'%s' is not a class", className))
	}

	obj := &HeObject{Class: cls, Properties: map[string]HeValue{}}

	// Set property defaults (walk class hierarchy oldest → newest).
	chain := classChain(cls)
	for i := len(chain) - 1; i >= 0; i-- {
		c := chain[i]
		for _, prop := range c.Properties {
			if prop.Default != nil {
				obj.Properties[prop.Name] = interp.evalExpr(prop.Default, env)
			} else {
				if _, exists := obj.Properties[prop.Name]; !exists {
					obj.Properties[prop.Name] = HeNil{}
				}
			}
		}
	}

	// Apply init fields.
	for name, valExpr := range initFields {
		obj.Properties[name] = interp.evalExpr(valExpr, env)
	}

	// Call 'init' method if it exists (constructor convention).
	if m, _ := cls.LookupMethod("init"); m != nil {
		interp.callMethod(obj, m, nil, pos, env)
	}

	return obj
}

func classChain(cls *HeClass) []*HeClass {
	var chain []*HeClass
	for c := cls; c != nil; c = c.SuperClass {
		chain = append(chain, c)
	}
	return chain
}

func (interp *Interpreter) execCompoundAssign(s *ast.CompoundAssignStmt, env *Environment) {
	current := interp.evalExpr(s.Target, env)
	operand := interp.evalExpr(s.Value, env)
	var result HeValue
	switch s.Operator {
	case "+=":
		result = interp.opAdd(current, operand, s.Pos)
	case "-=":
		result = interp.numOp(current, operand, s.Pos, "-")
	case "*=":
		result = interp.numOp(current, operand, s.Pos, "*")
	case "/=":
		result = interp.numOp(current, operand, s.Pos, "/")
	}
	interp.assignTarget(s.Target, result, env)
}

func (interp *Interpreter) execInput(s *ast.InputStmt, env *Environment) {
	if s.Prompt != nil {
		p := interp.evalExpr(s.Prompt, env)
		fmt.Print(p.Repr())
	}
	line, _ := interp.reader.ReadString('\n')
	line = strings.TrimRight(line, "\r\n")
	env.Set(s.Target, HeString{V: line})
}

func (interp *Interpreter) execIf(s *ast.IfStmt, env *Environment) {
	cond := interp.evalExpr(s.Condition, env)
	if isTruthy(cond) {
		child := env.NewChild()
		interp.execStmts(s.Consequent, child)
	} else if len(s.Alternate) > 0 {
		child := env.NewChild()
		interp.execStmts(s.Alternate, child)
	}
}

func (interp *Interpreter) execWhile(s *ast.WhileStmt, env *Environment) {
	for {
		cond := interp.evalExpr(s.Condition, env)
		if !isTruthy(cond) {
			return
		}
		broken := false
		func() {
			defer func() {
				if r := recover(); r != nil {
					switch r.(type) {
					case breakSignal:
						broken = true
					case continueSignal:
						// end this iteration; outer loop continues
					default:
						panic(r)
					}
				}
			}()
			child := env.NewChild()
			interp.execStmts(s.Body, child)
		}()
		if broken {
			return
		}
	}
}

func (interp *Interpreter) execForEach(s *ast.ForEachStmt, env *Environment) {
	iterable := interp.evalExpr(s.Iterable, env)
	var elements []HeValue
	switch v := iterable.(type) {
	case *HeList:
		elements = v.Elements
	case HeString:
		runes := []rune(v.V)
		elements = make([]HeValue, len(runes))
		for i, r := range runes {
			elements[i] = HeString{V: string(r)}
		}
	default:
		panic(runtimeErr(s.Pos, "cannot iterate over %T", iterable))
	}

	for _, el := range elements {
		broken := false
		func() {
			defer func() {
				if r := recover(); r != nil {
					switch r.(type) {
					case breakSignal:
						broken = true
					case continueSignal:
					default:
						panic(r)
					}
				}
			}()
			child := env.NewChild()
			child.Define(s.VarName, el)
			interp.execStmts(s.Body, child)
		}()
		if broken {
			return
		}
	}
}

func (interp *Interpreter) execCallStmt(s *ast.CallStmt, env *Environment) {
	obj := interp.evalExpr(s.Object, env)
	heObj, ok := obj.(*HeObject)
	if !ok {
		panic(runtimeErr(s.Pos, "cannot call method on non-object (%s)", obj.Repr()))
	}
	m, _ := heObj.Class.LookupMethod(s.Method)
	if m == nil {
		panic(runtimeErr(s.Pos, "object of class '%s' has no method '%s'", heObj.Class.Name, s.Method))
	}
	args := make([]HeValue, len(s.Args))
	for i, a := range s.Args {
		args[i] = interp.evalExpr(a, env)
	}
	interp.callMethod(heObj, m, args, s.Pos, env)
}

// ─── method invocation ────────────────────────────────────────────────────────

func (interp *Interpreter) callMethod(receiver *HeObject, m *ast.MethodDecl, args []HeValue, pos ast.Pos, env *Environment) HeValue {
	if len(args) != len(m.Params) {
		panic(runtimeErr(pos, "method '%s' expects %d argument(s), got %d", m.Name, len(m.Params), len(args)))
	}

	frame := env.NewChild()
	frame.Define("this", receiver)

	// super: a proxy object with the parent class's methods.
	if receiver.Class.SuperClass != nil {
		superProxy := &HeObject{Class: receiver.Class.SuperClass, Properties: receiver.Properties}
		frame.Define("super", superProxy)
	}

	for i, param := range m.Params {
		frame.Define(param, args[i])
	}

	var result HeValue = HeNil{}
	func() {
		defer func() {
			if r := recover(); r != nil {
				if ret, ok := r.(returnSignal); ok {
					result = ret.value
				} else {
					panic(r)
				}
			}
		}()
		interp.execStmts(m.Body, frame)
	}()
	return result
}

// ─── assignment target ────────────────────────────────────────────────────────

func (interp *Interpreter) assignTarget(target ast.Expr, val HeValue, env *Environment) {
	switch t := target.(type) {
	case *ast.Identifier:
		env.Set(t.Name, val)
	case *ast.DotExpr:
		obj := interp.evalExpr(t.Object, env)
		heObj, ok := obj.(*HeObject)
		if !ok {
			panic(runtimeErr(t.Pos, "cannot set property on non-object"))
		}
		heObj.Set(t.Field, val)
	case *ast.IndexExpr:
		obj := interp.evalExpr(t.Object, env)
		idx := interp.evalExpr(t.Index, env)
		list, ok := obj.(*HeList)
		if !ok {
			panic(runtimeErr(t.Pos, "index assignment requires a list"))
		}
		n, ok := idx.(HeNumber)
		if !ok {
			panic(runtimeErr(t.Pos, "list index must be a number"))
		}
		i := int(n.V)
		if i < 0 || i >= len(list.Elements) {
			panic(runtimeErr(t.Pos, "list index %d out of range", i))
		}
		list.Elements[i] = val
	default:
		panic(runtimeErr(ast.Pos{}, "invalid assignment target"))
	}
}

// ─── expression evaluation ────────────────────────────────────────────────────

func (interp *Interpreter) evalExpr(expr ast.Expr, env *Environment) HeValue {
	switch e := expr.(type) {
	case *ast.NumberLit:
		return HeNumber{V: e.Value}
	case *ast.StringLit:
		return HeString{V: e.Value}
	case *ast.BoolLit:
		return HeBool{V: e.Value}
	case *ast.NilLit:
		return HeNil{}
	case *ast.Identifier:
		return interp.evalIdent(e, env)
	case *ast.DotExpr:
		return interp.evalDot(e, env)
	case *ast.IndexExpr:
		return interp.evalIndex(e, env)
	case *ast.BinaryExpr:
		return interp.evalBinary(e, env)
	case *ast.UnaryExpr:
		return interp.evalUnary(e, env)
	case *ast.CallExpr:
		return interp.evalCallExpr(e, env)
	case *ast.FuncCallExpr:
		return interp.evalFuncCall(e, env)
	case *ast.NewExpr:
		return interp.instantiate(e.ClassName, e.InitFields, e.Pos, env)
	case *ast.ListLiteral:
		elems := make([]HeValue, len(e.Elements))
		for i, el := range e.Elements {
			elems[i] = interp.evalExpr(el, env)
		}
		return &HeList{Elements: elems}
	}
	panic(runtimeErr(ast.Pos{}, "unknown expression type %T", expr))
}

func (interp *Interpreter) evalIdent(e *ast.Identifier, env *Environment) HeValue {
	v, ok := env.Get(e.Name)
	if !ok {
		panic(runtimeErr(e.Pos, "undefined variable '%s'", e.Name))
	}
	return v
}

func (interp *Interpreter) evalDot(e *ast.DotExpr, env *Environment) HeValue {
	obj := interp.evalExpr(e.Object, env)
	switch o := obj.(type) {
	case *HeObject:
		// property?
		if v, ok := o.Get(e.Field); ok {
			return v
		}
		// method reference? Return as callable.
		if m, _ := o.Class.LookupMethod(e.Field); m != nil {
			// Return a bound-method builtin.
			method := m
			receiver := o
			return &HeBuiltin{
				Name: e.Field,
				Fn: func(args []HeValue) (HeValue, error) {
					return interp.callMethod(receiver, method, args, e.Pos, env), nil
				},
			}
		}
		panic(runtimeErr(e.Pos, "object of class '%s' has no property or method '%s'", o.Class.Name, e.Field))
	case *HeList:
		return interp.listBuiltin(o, e.Field, e.Pos, env)
	case HeString:
		return interp.stringBuiltin(o, e.Field, e.Pos)
	}
	panic(runtimeErr(e.Pos, "cannot access field '%s' on %s", e.Field, obj.Repr()))
}

// listBuiltin provides .method() access on lists.
func (interp *Interpreter) listBuiltin(list *HeList, field string, pos ast.Pos, env *Environment) HeValue {
	switch field {
	case "length":
		return HeNumber{V: float64(len(list.Elements))}
	case "append":
		return &HeBuiltin{Name: "append", Fn: func(args []HeValue) (HeValue, error) {
			list.Elements = append(list.Elements, args...)
			return list, nil
		}}
	case "remove":
		return &HeBuiltin{Name: "remove", Fn: func(args []HeValue) (HeValue, error) {
			if len(args) != 1 {
				return nil, fmt.Errorf("remove() takes 1 argument")
			}
			n, ok := args[0].(HeNumber)
			if !ok {
				return nil, fmt.Errorf("remove() argument must be a number")
			}
			i := int(n.V)
			if i < 0 || i >= len(list.Elements) {
				return nil, fmt.Errorf("remove() index %d out of range", i)
			}
			list.Elements = append(list.Elements[:i], list.Elements[i+1:]...)
			return list, nil
		}}
	case "contains":
		return &HeBuiltin{Name: "contains", Fn: func(args []HeValue) (HeValue, error) {
			if len(args) != 1 {
				return nil, fmt.Errorf("contains() takes 1 argument")
			}
			for _, el := range list.Elements {
				if heEqual(el, args[0]) {
					return HeBool{V: true}, nil
				}
			}
			return HeBool{V: false}, nil
		}}
	case "first":
		if len(list.Elements) == 0 {
			return HeNil{}
		}
		return list.Elements[0]
	case "last":
		if len(list.Elements) == 0 {
			return HeNil{}
		}
		return list.Elements[len(list.Elements)-1]
	}
	panic(runtimeErr(pos, "list has no property or method '%s'", field))
}

// stringBuiltin provides .method access on strings.
func (interp *Interpreter) stringBuiltin(s HeString, field string, pos ast.Pos) HeValue {
	switch field {
	case "length":
		return HeNumber{V: float64(len([]rune(s.V)))}
	case "upper":
		return HeString{V: strings.ToUpper(s.V)}
	case "lower":
		return HeString{V: strings.ToLower(s.V)}
	case "trim":
		return HeString{V: strings.TrimSpace(s.V)}
	}
	panic(runtimeErr(pos, "string has no property or method '%s'", field))
}

func (interp *Interpreter) evalIndex(e *ast.IndexExpr, env *Environment) HeValue {
	obj := interp.evalExpr(e.Object, env)
	idx := interp.evalExpr(e.Index, env)
	switch o := obj.(type) {
	case *HeList:
		n, ok := idx.(HeNumber)
		if !ok {
			panic(runtimeErr(e.Pos, "list index must be a number"))
		}
		i := int(n.V)
		if i < 0 || i >= len(o.Elements) {
			panic(runtimeErr(e.Pos, "list index %d out of range (length %d)", i, len(o.Elements)))
		}
		return o.Elements[i]
	case HeString:
		n, ok := idx.(HeNumber)
		if !ok {
			panic(runtimeErr(e.Pos, "string index must be a number"))
		}
		runes := []rune(o.V)
		i := int(n.V)
		if i < 0 || i >= len(runes) {
			panic(runtimeErr(e.Pos, "string index %d out of range", i))
		}
		return HeString{V: string(runes[i])}
	}
	panic(runtimeErr(e.Pos, "cannot index into %T", obj))
}

func (interp *Interpreter) evalCallExpr(e *ast.CallExpr, env *Environment) HeValue {
	obj := interp.evalExpr(e.Object, env)
	args := make([]HeValue, len(e.Args))
	for i, a := range e.Args {
		args[i] = interp.evalExpr(a, env)
	}

	switch o := obj.(type) {
	case *HeObject:
		m, _ := o.Class.LookupMethod(e.Method)
		if m == nil {
			panic(runtimeErr(e.Pos, "object of class '%s' has no method '%s'", o.Class.Name, e.Method))
		}
		return interp.callMethod(o, m, args, e.Pos, env)
	case *HeBuiltin:
		// dot-called builtin (e.g. list.append(...))
		result, err := o.Fn(args)
		if err != nil {
			panic(runtimeErr(e.Pos, "%s", err.Error()))
		}
		return result
	}
	panic(runtimeErr(e.Pos, "cannot call method '%s' on %s", e.Method, obj.Repr()))
}

func (interp *Interpreter) evalFuncCall(e *ast.FuncCallExpr, env *Environment) HeValue {
	fn, ok := env.Get(e.Name)
	if !ok {
		panic(runtimeErr(e.Pos, "undefined function '%s'", e.Name))
	}
	args := make([]HeValue, len(e.Args))
	for i, a := range e.Args {
		args[i] = interp.evalExpr(a, env)
	}
	switch f := fn.(type) {
	case *HeBuiltin:
		result, err := f.Fn(args)
		if err != nil {
			panic(runtimeErr(e.Pos, "%s", err.Error()))
		}
		return result
	}
	panic(runtimeErr(e.Pos, "'%s' is not callable", e.Name))
}

// ─── binary / unary operators ─────────────────────────────────────────────────

func (interp *Interpreter) evalBinary(e *ast.BinaryExpr, env *Environment) HeValue {
	// Short-circuit logic operators.
	if e.Operator == "and" {
		left := interp.evalExpr(e.Left, env)
		if !isTruthy(left) {
			return left
		}
		return interp.evalExpr(e.Right, env)
	}
	if e.Operator == "or" {
		left := interp.evalExpr(e.Left, env)
		if isTruthy(left) {
			return left
		}
		return interp.evalExpr(e.Right, env)
	}

	left := interp.evalExpr(e.Left, env)
	right := interp.evalExpr(e.Right, env)

	switch e.Operator {
	case "+":
		return interp.opAdd(left, right, e.Pos)
	case "-":
		return interp.numOp(left, right, e.Pos, "-")
	case "*":
		return interp.numOp(left, right, e.Pos, "*")
	case "/":
		return interp.numOp(left, right, e.Pos, "/")
	case "%":
		return interp.numOp(left, right, e.Pos, "%")
	case "^":
		return interp.numOp(left, right, e.Pos, "^")
	case "==":
		return HeBool{V: heEqual(left, right)}
	case "!=", "<>":
		return HeBool{V: !heEqual(left, right)}
	case "<":
		return HeBool{V: heCompare(left, right, e.Pos) < 0}
	case ">":
		return HeBool{V: heCompare(left, right, e.Pos) > 0}
	case "<=":
		return HeBool{V: heCompare(left, right, e.Pos) <= 0}
	case ">=":
		return HeBool{V: heCompare(left, right, e.Pos) >= 0}
	}
	panic(runtimeErr(e.Pos, "unknown operator '%s'", e.Operator))
}

func (interp *Interpreter) opAdd(left, right HeValue, pos ast.Pos) HeValue {
	// number + number
	ln, lok := left.(HeNumber)
	rn, rok := right.(HeNumber)
	if lok && rok {
		return HeNumber{V: ln.V + rn.V}
	}
	// string concatenation (auto-coerce)
	return HeString{V: left.Repr() + right.Repr()}
}

func (interp *Interpreter) numOp(left, right HeValue, pos ast.Pos, op string) HeValue {
	ln, lok := left.(HeNumber)
	rn, rok := right.(HeNumber)
	if !lok || !rok {
		panic(runtimeErr(pos, "operator '%s' requires two numbers (got %s and %s)", op, left.Repr(), right.Repr()))
	}
	switch op {
	case "-":
		return HeNumber{V: ln.V - rn.V}
	case "*":
		return HeNumber{V: ln.V * rn.V}
	case "/":
		if rn.V == 0 {
			panic(runtimeErr(pos, "division by zero"))
		}
		return HeNumber{V: ln.V / rn.V}
	case "%":
		if rn.V == 0 {
			panic(runtimeErr(pos, "modulo by zero"))
		}
		return HeNumber{V: math.Mod(ln.V, rn.V)}
	case "^":
		return HeNumber{V: math.Pow(ln.V, rn.V)}
	}
	panic(runtimeErr(pos, "unknown numeric operator '%s'", op))
}

func (interp *Interpreter) evalUnary(e *ast.UnaryExpr, env *Environment) HeValue {
	operand := interp.evalExpr(e.Operand, env)
	switch e.Operator {
	case "not":
		return HeBool{V: !isTruthy(operand)}
	case "-":
		n, ok := operand.(HeNumber)
		if !ok {
			panic(runtimeErr(e.Pos, "unary '-' requires a number"))
		}
		return HeNumber{V: -n.V}
	}
	panic(runtimeErr(e.Pos, "unknown unary operator '%s'", e.Operator))
}

// ─── helpers ──────────────────────────────────────────────────────────────────

func isTruthy(v HeValue) bool {
	switch val := v.(type) {
	case HeBool:
		return val.V
	case HeNil:
		return false
	case HeNumber:
		return val.V != 0
	case HeString:
		return val.V != ""
	case *HeList:
		return len(val.Elements) > 0
	}
	return true
}

func heEqual(a, b HeValue) bool {
	switch av := a.(type) {
	case HeNumber:
		bv, ok := b.(HeNumber)
		return ok && av.V == bv.V
	case HeString:
		bv, ok := b.(HeString)
		return ok && av.V == bv.V
	case HeBool:
		bv, ok := b.(HeBool)
		return ok && av.V == bv.V
	case HeNil:
		_, ok := b.(HeNil)
		return ok
	}
	return a == b
}

func heCompare(a, b HeValue, pos ast.Pos) int {
	switch av := a.(type) {
	case HeNumber:
		bv, ok := b.(HeNumber)
		if !ok {
			panic(runtimeErr(pos, "cannot compare number with %T", b))
		}
		if av.V < bv.V {
			return -1
		}
		if av.V > bv.V {
			return 1
		}
		return 0
	case HeString:
		bv, ok := b.(HeString)
		if !ok {
			panic(runtimeErr(pos, "cannot compare string with %T", b))
		}
		return strings.Compare(av.V, bv.V)
	}
	panic(runtimeErr(pos, "cannot compare %T with %T", a, b))
}
