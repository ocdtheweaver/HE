package eval

import (
	"fmt"
	"math"
	"strings"
	"time"

	"hunterlang/lang/ast"
	"hunterlang/lang/lexer"
	"hunterlang/lang/parser"
	"hunterlang/lang/types"
	"os"
	"path/filepath"
)

// returnSignal is a sentinel error used to implement `return` statements.
type returnSignal struct{ value types.Value }

func (r *returnSignal) Error() string { return "return" }

// ── Runtime ───────────────────────────────────────────────────────────────────

type Runtime struct {
	Env      map[string]types.Value
	Objects  map[string]*types.Object
	Platform string

	// currentRecv is non-nil while executing a user-defined action or reaction body.
	currentRecv *types.Object

	// persist handles remember/recall/forget
	persist *persistStore

	// callDepth tracks recursion depth for safety
	callDepth int

	// typeConstraints maps variable/field name to expected type
	typeConstraints map[string]string

	// baseDir is the directory of the currently-running .he file,
	// used to resolve relative summon paths (e.g. "lib/helpers.he")
	baseDir string

	// nameOrigin tracks which summoned file (or "" for locally-defined)
	// last contributed each flat-merged name. Used to detect collisions
	// between two different summoned files defining the same name.
	nameOrigin map[string]string

	// ambiguous marks names that exist from two or more *different* summoned
	// files. Bare access to an ambiguous name is a runtime error — the
	// person must use a qualified alias.Name form instead.
	ambiguous map[string]bool

	// entitlement answers "is this #protected[N] tag currently granted?"
	// Stage A, Step 2 of the protection roadmap: wired in here, but not
	// yet consulted anywhere. Enforcement (Step 3) is what will call this.
	// Defaults to AlwaysDeny — fail closed if nothing configures it.
	entitlement EntitlementChecker
}

func newRuntime() *Runtime {
	r := &Runtime{
		Env:      map[string]types.Value{},
		Objects:  map[string]*types.Object{},
		Platform: "unknown",
		persist:          newPersistStore(),
		typeConstraints:  map[string]string{},
		nameOrigin:       map[string]string{},
		ambiguous:        map[string]bool{},
		entitlement:      AlwaysDeny{},
	}
	r.Env["platform"] = types.FromString("unknown")
	r.Env["nothing"] = types.Nil()
	r.Env["pi"] = types.FromNumber(math.Pi)
	return r
}

// SetEntitlementChecker replaces the runtime's entitlement checker.
// Not called by anything in the standard CLI yet (Stage A, Step 2) —
// exposed so tests and future enforcement wiring can configure it.
func (r *Runtime) SetEntitlementChecker(c EntitlementChecker) {
	if c == nil {
		c = AlwaysDeny{}
	}
	r.entitlement = c
}

func (r *Runtime) setVar(name string, v types.Value) {
	// Type constraint enforcement
	if constraint, ok := r.typeConstraints[name]; ok {
		if !typeMatches(v, constraint) {
			// Soft enforcement: warn but don't crash (strict mode planned)
			fmt.Printf("  ⚠ type warning: %q expects %s, got %s\n", name, constraint, v.Type)
		}
	}
	r.Env[name] = v
	if v.Type == types.ObjectT && v.Object != nil {
		r.Objects[name] = v.Object
	}
}

func typeMatches(v types.Value, constraint string) bool {
	switch constraint {
	case "number":
		return v.Type == types.NumberT
	case "text":
		return v.Type == types.StringT
	case "boolean":
		return v.Type == types.BooleanT
	case "list":
		return v.Type == types.ArrayT
	case "nothing":
		return v.Type == types.NilT
	case "any":
		return true
	}
	return true // unknown constraint = permissive
}

func (r *Runtime) getVar(name string) (types.Value, bool) {
	// Prefer receiver fields when inside an action.
	if r.currentRecv != nil {
		if v, ok := r.currentRecv.Fields[name]; ok {
			return v, true
		}
	}
	v, ok := r.Env[name]
	return v, ok
}

func (r *Runtime) getObject(name string) (*types.Object, error) {
	if obj, ok := r.Objects[name]; ok && obj != nil {
		return obj, nil
	}
	if v, ok := r.Env[name]; ok && v.Type == types.ObjectT && v.Object != nil {
		return v.Object, nil
	}
	return nil, fmt.Errorf("I don't know anything called %q", name)
}

// ── Program execution ─────────────────────────────────────────────────────────

func (r *Runtime) execProgram(prog *ast.Program) error {
	for _, ln := range prog.Lines {
		if err := r.execLine(ln); err != nil {
			return err
		}
	}
	return nil
}

func (r *Runtime) execLine(ln ast.Line) error {
	switch l := ln.(type) {
	case ast.SummonLine:
		return r.execSummon(l)
	case *ast.ObjectLine:
		return r.execObjectLine(l)
	case ast.GlobalStatementLine:
		return r.execStatement(l.Statement)
	case *ast.GlobalStatementLine:
		return r.execStatement(l.Statement)
	case ast.WithAssetsLine:
		return nil
	default:
		return fmt.Errorf("unhandled top-level node %T", ln)
	}
}

func (r *Runtime) execSummon(l ast.SummonLine) error {
	obj := &types.Object{
		Name:     l.Alias.Name,
		Fields:   map[string]types.Value{},
		Actions:  map[string]*types.Action{},
		Builtins: map[string]types.BuiltinFn{},
	}
	registerModule(obj, l.ModuleName.Lexeme)
	r.setVar(l.Alias.Name, types.FromObject(obj))
	return nil
}

func (r *Runtime) execObjectLine(o *ast.ObjectLine) error {
	obj := &types.Object{
		Name:       o.Name.Name,
		Fields:     map[string]types.Value{},
		Actions:    map[string]*types.Action{},
		Builtins:   map[string]types.BuiltinFn{},
		ProtectTag: o.ProtectTag,
	}

	// Inherit from `like Parent`
	if o.Like != nil {
		if parent, ok := r.Objects[o.Like.Name]; ok && parent != nil {
			for k, v := range parent.Fields {
				obj.Fields[k] = v
			}
			for k, a := range parent.Actions {
				obj.Actions[k] = a
			}
			for k, b := range parent.Builtins {
				obj.Builtins[k] = b
			}
			obj.Reactions = append(obj.Reactions, parent.Reactions...)
		}
	}

	for _, sec := range o.Body.Sections {
		switch s := sec.(type) {
		case ast.PropertiesSection:
			for _, pr := range s.Props {
				v, err := r.evalExpr(pr.Value)
				if err != nil {
					return err
				}
				obj.Fields[pr.Name] = v
				// Register type constraint if kind contains a type hint
				// e.g. kind = "has is" — parse "number" from property name annotation
				// For now: infer constraint from actual value type
				switch v.Type {
				case types.NumberT:
					r.typeConstraints[o.Name.Name+"."+pr.Name] = "number"
				case types.StringT:
					r.typeConstraints[o.Name.Name+"."+pr.Name] = "text"
				case types.BooleanT:
					r.typeConstraints[o.Name.Name+"."+pr.Name] = "boolean"
				}
			}
		case ast.AbilitiesSection:
			for _, act := range s.Abilities {
				params := make([]string, len(act.Params))
				for i, p := range act.Params {
					params[i] = p.Name
				}
				obj.Actions[act.Name] = &types.Action{
					Name:       act.Name,
					Params:     params,
					Body:       act.Body,
					ProtectTag: act.ProtectTag,
				}
			}
		case ast.ReactionsSection:
			for _, rx := range s.Reactions {
				obj.Reactions = append(obj.Reactions, types.Reaction{
					Trigger: rx.Trigger,
					Body:    rx.Body,
				})
			}
		case ast.MemoriesSection:
			// placeholder
		}
	}

	r.setVar(o.Name.Name, types.FromObject(obj))
	return nil
}

// ── Statements ────────────────────────────────────────────────────────────────

func (r *Runtime) execStatement(s ast.Statement) error {
	switch st := s.(type) {

	case ast.SayStmt:
		v, err := r.evalExpr(st.Expr)
		if err != nil {
			return err
		}
		fmt.Println(v.String())
		return nil

	case ast.ChangeStmt:
		v, err := r.evalExpr(st.Expr)
		if err != nil {
			return err
		}
		if r.currentRecv != nil {
			if _, isField := r.currentRecv.Fields[st.Name]; isField {
				r.currentRecv.Fields[st.Name] = v
				return nil
			}
		}
		r.setVar(st.Name, v)
		return nil

	case ast.GrowStmt:
		delta, err := r.evalExpr(st.Expr)
		if err != nil {
			return err
		}
		cur, _ := r.getVar(st.Name)
		if cur.Type != types.NumberT {
			cur = types.FromNumber(0)
		}
		newVal := types.FromNumber(cur.Number + delta.Number)
		if r.currentRecv != nil {
			if _, isField := r.currentRecv.Fields[st.Name]; isField {
				r.currentRecv.Fields[st.Name] = newVal
				return nil
			}
		}
		r.setVar(st.Name, newVal)
		return nil

	case ast.ShrinkStmt:
		delta, err := r.evalExpr(st.Expr)
		if err != nil {
			return err
		}
		cur, _ := r.getVar(st.Name)
		if cur.Type != types.NumberT {
			cur = types.FromNumber(0)
		}
		newVal := types.FromNumber(cur.Number - delta.Number)
		if r.currentRecv != nil {
			if _, isField := r.currentRecv.Fields[st.Name]; isField {
				r.currentRecv.Fields[st.Name] = newVal
				return nil
			}
		}
		r.setVar(st.Name, newVal)
		return nil

	case ast.DecideStmt:
		cond, err := r.evalExpr(st.Cond)
		if err != nil {
			return err
		}
		branch := st.Else
		if isTruthy(cond) {
			branch = st.Then
		}
		for _, t := range branch {
			if err := r.execStatement(t); err != nil {
				return err
			}
		}
		return nil

	case ast.RepeatStmt:
		switch st.Kind {
		case "while":
			for i := 0; i < 100_000; i++ {
				cond, err := r.evalExpr(st.Cond)
				if err != nil {
					return err
				}
				if !isTruthy(cond) {
					break
				}
				for _, b := range st.Body {
					if err := r.execStatement(b); err != nil {
						if _, ok := err.(*returnSignal); ok {
							return err
						}
						return err
					}
				}
			}
		case "times":
			countV, err := r.evalExpr(st.Cond)
			if err != nil {
				return err
			}
			count := int(countV.Number)
			for i := 0; i < count; i++ {
				for _, b := range st.Body {
					if err := r.execStatement(b); err != nil {
						return err
					}
				}
			}
		}
		return nil

	case ast.CallStmt:
		obj, err := r.getObject(st.Object)
		if err != nil {
			return err
		}
		_, err = r.callMethod(obj, st.Action, st.Args)
		return err

	case ast.WaitStmt:
		v, err := r.evalExpr(st.Expr)
		if err != nil {
			return err
		}
		if v.Type != types.NumberT {
			return fmt.Errorf("wait expects a number")
		}
		switch st.Unit {
		case "seconds":
			time.Sleep(time.Duration(v.Number * float64(time.Second)))
		case "frames":
			time.Sleep(time.Duration(v.Number * (1.0 / 60.0) * float64(time.Second)))
		}
		return nil

	case ast.ReturnStmt:
		v, err := r.evalExpr(st.Expr)
		if err != nil {
			return err
		}
		return &returnSignal{value: v}



	case ast.RangeLoopStmt:
		fromVal, err := r.evalExpr(st.From)
		if err != nil {
			return err
		}
		toVal, err := r.evalExpr(st.To)
		if err != nil {
			return err
		}
		step := 1.0
		if st.Step != nil {
			stepVal, err := r.evalExpr(st.Step)
			if err != nil {
				return err
			}
			step = stepVal.Number
		}
		from, to := fromVal.Number, toVal.Number
		for i := from; i <= to; i += step {
			r.setVar(st.VarName, types.FromNumber(i))
			for _, b := range st.Body {
				if err := r.execStatement(b); err != nil {
					if _, ok := err.(*returnSignal); ok {
						return err
					}
					return err
				}
			}
		}
		return nil

	case ast.TryStmt:
		execErr := r.execBlock(st.Body)
		if execErr != nil {
			// Store error message in special variable
			r.setVar("error", types.FromString(execErr.Error()))
			for _, b := range st.Handler {
				if herr := r.execStatement(b); herr != nil {
					return herr
				}
			}
		}
		return nil

	case ast.RepeatUntilStmt:
		for i := 0; i < 100_000; i++ {
			cond, err := r.evalExpr(st.Cond)
			if err != nil {
				return err
			}
			if isTruthy(cond) {
				break
			}
			for _, b := range st.Body {
				if err := r.execStatement(b); err != nil {
					return err
				}
			}
		}
		return nil





	case ast.CountedRepeatStmt:
		countV, err := r.evalExpr(st.Count)
		if err != nil {
			return err
		}
		count := int(countV.Number)
		for i := 0; i < count; i++ {
			r.setVar(st.CountVar, types.FromNumber(float64(i+1))) // 1-based
			for _, b := range st.Body {
				if err := r.execStatement(b); err != nil {
					if _, ok := err.(*returnSignal); ok {
						return err
					}
					return err
				}
			}
		}
		return nil

	case ast.WithScopeStmt:
		// Evaluate the expression, bind to alias, execute body, then clean up
		v, err := r.evalExpr(st.Expr)
		if err != nil {
			return err
		}
		// Save any existing value under that name
		prev, hadPrev := r.Env[st.Alias]
		r.setVar(st.Alias, v)
		for _, b := range st.Body {
			if err := r.execStatement(b); err != nil {
				if !hadPrev {
					delete(r.Env, st.Alias)
				} else {
					r.Env[st.Alias] = prev
				}
				return err
			}
		}
		if !hadPrev {
			delete(r.Env, st.Alias)
		} else {
			r.Env[st.Alias] = prev
		}
		return nil

	case ast.ForEachFieldStmt:
		objVal, err := r.evalExpr(st.Object)
		if err != nil {
			return err
		}
		var fields map[string]types.Value
		if objVal.Type == types.ObjectT && objVal.Object != nil {
			fields = objVal.Object.Fields
		} else {
			return fmt.Errorf("'for each %s, %s in' expects an object", st.KeyVar, st.ValVar)
		}
		for k, v := range fields {
			r.setVar(st.KeyVar, types.FromString(k))
			if st.ValVar != "" {
				r.setVar(st.ValVar, v)
			}
			for _, b := range st.Body {
				if err := r.execStatement(b); err != nil {
					if _, ok := err.(*returnSignal); ok {
						return err
					}
					return err
				}
			}
		}
		return nil

	case ast.TryWithVarStmt:
		execErr := r.execBlock(st.Body)
		if execErr != nil {
			// Bind error to named variable
			r.setVar(st.ErrVar, types.FromString(execErr.Error()))
			for _, b := range st.Handler {
				if herr := r.execStatement(b); herr != nil {
					return herr
				}
			}
		}
		return nil

	case ast.LoadModuleStmt:
		return r.loadFileModule(st.FilePath, st.Alias, st.HasAlias)

	case ast.MultiAssignStmt:
		// Evaluate all right-hand expressions first (allows: set a, b to b, a)
		vals := make([]types.Value, len(st.Exprs))
		for i, ex := range st.Exprs {
			v, err := r.evalExpr(ex)
			if err != nil {
				return err
			}
			vals[i] = v
		}
		// If single expr returns an array (multi-return), unpack it
		if len(vals) == 1 && vals[0].Type == types.ArrayT && len(vals[0].Array) == len(st.Names) {
			vals = vals[0].Array
		}
		for i, name := range st.Names {
			if i < len(vals) {
				r.setVar(name, vals[i])
			} else {
				r.setVar(name, types.Nil())
			}
		}
		return nil

	case ast.MultiReturnStmt:
		if len(st.Exprs) == 1 {
			v, err := r.evalExpr(st.Exprs[0])
			if err != nil {
				return err
			}
			return &returnSignal{value: v}
		}
		vals := make([]types.Value, len(st.Exprs))
		for i, ex := range st.Exprs {
			v, err := r.evalExpr(ex)
			if err != nil {
				return err
			}
			vals[i] = v
		}
		return &returnSignal{value: types.FromArray(vals)}

	case ast.RememberStmt:
		v, _ := r.getVar(st.Name)
		return r.persist.save(st.Key, v)

	case ast.ForgetStmt:
		return r.persist.remove(st.Name)

	case ast.RecallStmt:
		v, err := r.persist.load(st.Key)
		if err != nil {
			r.setVar(st.Name, types.Nil())
			return nil
		}
		r.setVar(st.Name, v)
		return nil

	case ast.ForEachStmt:
		listVal, err := r.evalExpr(st.List)
		if err != nil {
			return err
		}
		var items []types.Value
		if listVal.Type == types.ArrayT {
			items = listVal.Array
		} else {
			items = []types.Value{listVal}
		}
		for _, item := range items {
			r.setVar(st.VarName, item)
			for _, b := range st.Body {
				if err := r.execStatement(b); err != nil {
					if _, ok := err.(*returnSignal); ok {
						return err
					}
					// stopSignal would go here
					return err
				}
			}
		}
		return nil

	case ast.AskStmt:
		promptVal, err := r.evalExpr(st.Prompt)
		if err != nil {
			return err
		}
		fmt.Print(promptVal.String() + " ")
		var input string
		fmt.Scanln(&input)
		r.setVar(st.VarName, types.FromString(input))
		return nil

	case ast.DotAssignStmt:
		obj, err := r.getObject(st.Object)
		if err != nil {
			return err
		}
		val, err := r.evalExpr(st.Expr)
		if err != nil {
			return err
		}
		obj.Fields[st.Field] = val
		return nil

	case ast.ExprStmt:
		// Detect inline summon: __summon__:module:alias
		if sl, ok := st.Expr.(ast.StringLit); ok && len(sl.Value) > 10 && sl.Value[:10] == "__summon__" {
			parts := strings.SplitN(sl.Value, ":", 3)
			if len(parts) == 3 {
				modName, alias := parts[1], parts[2]
				obj := &types.Object{
					Name:     alias,
					Fields:   map[string]types.Value{},
					Actions:  map[string]*types.Action{},
					Builtins: map[string]types.BuiltinFn{},
				}
				registerModule(obj, modName)
				r.setVar(alias, types.FromObject(obj))
				return nil
			}
		}
		_, err := r.evalExpr(st.Expr)
		return err

	default:
		return fmt.Errorf("unhandled statement %T", s)
	}
}

// ── Expressions ───────────────────────────────────────────────────────────────

func (r *Runtime) evalExpr(e ast.Expression) (types.Value, error) {
	switch ex := e.(type) {
	case ast.NumberLit:
		return types.FromNumber(ex.Value), nil
	case ast.StringLit:
		return types.FromString(ex.Value), nil
	case ast.BooleanLit:
		return types.FromBoolean(ex.Value), nil

	case ast.IdentifierExpr:
		if r.ambiguous[ex.Name] {
			return types.Nil(), fmt.Errorf(
				"%q is defined in more than one summoned file — use the file's alias to choose which one (e.g. summon \"...\" as theirname, then theirname.%s)",
				ex.Name, ex.Name,
			)
		}
		v, ok := r.getVar(ex.Name)
		if !ok {
			return types.Nil(), nil
		}
		return v, nil

	case ast.ArrayLit:
		out := make([]types.Value, 0, len(ex.Elems))
		for _, el := range ex.Elems {
			v, err := r.evalExpr(el)
			if err != nil {
				return types.Nil(), err
			}
			out = append(out, v)
		}
		return types.FromArray(out), nil

	case ast.NamedArgLit:
		// Evaluate named args into a pseudo-object so builtins can read by key name.
		obj := &types.Object{
			Name:   "__named__",
			Fields: map[string]types.Value{},
		}
		for _, pair := range ex.Pairs {
			v, err := r.evalExpr(pair.Value)
			if err != nil {
				return types.Nil(), err
			}
			obj.Fields[pair.Key] = v
		}
		return types.FromObject(obj), nil

	case ast.ParenExpr:
		return r.evalExpr(ex.X)

	case ast.UnaryExpr:
		v, err := r.evalExpr(ex.X)
		if err != nil {
			return types.Nil(), err
		}
		switch ex.Op {
		case "-":
			if v.Type != types.NumberT {
				return types.Nil(), fmt.Errorf("'-' expects a number, got %s", v.Type)
			}
			return types.FromNumber(-v.Number), nil
		case "!":
			return types.FromBoolean(!isTruthy(v)), nil
		}

	case ast.BinaryExpr:
		l, err := r.evalExpr(ex.Left)
		if err != nil {
			return types.Nil(), err
		}
		ri, err := r.evalExpr(ex.Right)
		if err != nil {
			return types.Nil(), err
		}
		return evalBinary(ex.Op, l, ri)

	case ast.PowerExpr:
		l, err := r.evalExpr(ex.Left)
		if err != nil {
			return types.Nil(), err
		}
		ri, err := r.evalExpr(ex.Right)
		if err != nil {
			return types.Nil(), err
		}
		if l.Type != types.NumberT || ri.Type != types.NumberT {
			return types.Nil(), fmt.Errorf("power (**) expects numbers")
		}
		return types.FromNumber(math.Pow(l.Number, ri.Number)), nil

	case ast.CompareExpr:
		l, err := r.evalExpr(ex.Left)
		if err != nil {
			return types.Nil(), err
		}
		ri, err := r.evalExpr(ex.Right)
		if err != nil {
			return types.Nil(), err
		}
		return evalCompare(ex.Op, l, ri)

	case ast.LogicAndExpr:
		lv, err := r.evalExpr(ex.Left)
		if err != nil {
			return types.Nil(), err
		}
		if !isTruthy(lv) {
			return types.FromBoolean(false), nil
		}
		rv, err := r.evalExpr(ex.Right)
		if err != nil {
			return types.Nil(), err
		}
		return types.FromBoolean(isTruthy(rv)), nil

	case ast.LogicOrExpr:
		lv, err := r.evalExpr(ex.Left)
		if err != nil {
			return types.Nil(), err
		}
		if isTruthy(lv) {
			return types.FromBoolean(true), nil
		}
		rv, err := r.evalExpr(ex.Right)
		if err != nil {
			return types.Nil(), err
		}
		return types.FromBoolean(isTruthy(rv)), nil

	case ast.MethodCallExpr:
		recvVal, err := r.evalExpr(ex.Receiver)
		if err != nil {
			return types.Nil(), err
		}
		// Virtual dispatch by value type (same as MethodChainExpr)
		switch recvVal.Type {
		case types.StringT:
			return r.callStringMethod(recvVal.Str, ex.Method, ex.Args)
		case types.ArrayT:
			return r.callListMethod(recvVal, ex.Method, ex.Args)
		case types.NumberT:
			return r.callNumberMethod(recvVal.Number, ex.Method, ex.Args)
		case types.ObjectT:
			if recvVal.Object != nil {
				return r.callMethod(recvVal.Object, ex.Method, ex.Args)
			}
		}
		return types.Nil(), fmt.Errorf("%q is not an object, can't call .%s on it", ex.Receiver.Name, ex.Method)

	case ast.FieldAccessExpr:
		recvVal, err := r.evalExpr(ex.Receiver)
		if err != nil {
			return types.Nil(), err
		}
		if recvVal.Type != types.ObjectT || recvVal.Object == nil {
			return types.Nil(), fmt.Errorf("%q is not an object", ex.Receiver.Name)
		}
		if v, ok := recvVal.Object.Fields[ex.Field]; ok {
			return v, nil
		}
		return types.Nil(), nil




	case ast.MembershipExpr:
		val, err := r.evalExpr(ex.Value)
		if err != nil {
			return types.Nil(), err
		}
		listVal, err := r.evalExpr(ex.List)
		if err != nil {
			return types.Nil(), err
		}
		if listVal.Type != types.ArrayT {
			return types.Nil(), fmt.Errorf("'is one of' expects a list")
		}
		needle := val.String()
		for _, item := range listVal.Array {
			if item.String() == needle {
				return types.FromBoolean(true), nil
			}
		}
		return types.FromBoolean(false), nil

	case ast.MethodChainExpr:
		recvVal, err := r.evalExpr(ex.Recv)
		if err != nil {
			return types.Nil(), err
		}
		// Evaluate args
		argVals := make([]ast.Expression, len(ex.Args))
		copy(argVals, ex.Args)
		// Virtual dispatch by value type
		switch recvVal.Type {
		case types.StringT:
			return r.callStringMethod(recvVal.Str, ex.Method, argVals)
		case types.ArrayT:
			return r.callListMethod(recvVal, ex.Method, argVals)
		case types.NumberT:
			return r.callNumberMethod(recvVal.Number, ex.Method, argVals)
		case types.ObjectT:
			if recvVal.Object != nil {
				return r.callMethod(recvVal.Object, ex.Method, argVals)
			}
		}
		return types.Nil(), fmt.Errorf("can't call .%s on %s", ex.Method, recvVal.Type)

	case ast.ClosureExpr:
		// Capture a snapshot of the current environment
		captured := make(map[string]types.Value, len(r.Env))
		for k, v := range r.Env {
			captured[k] = v
		}
		if r.currentRecv != nil {
			for k, v := range r.currentRecv.Fields {
				captured[k] = v
			}
		}
		obj := &types.Object{
			Name:     "__closure__",
			Fields:   captured,
			Actions:  map[string]*types.Action{},
			Builtins: map[string]types.BuiltinFn{},
		}
		obj.Actions["__call__"] = &types.Action{
			Name:   "__call__",
			Params: ex.Params,
			Body:   ex.Body,
		}
		return types.FromObject(obj), nil

	case ast.BetweenExpr:
		val, err := r.evalExpr(ex.Value)
		if err != nil {
			return types.Nil(), err
		}
		low, err := r.evalExpr(ex.Low)
		if err != nil {
			return types.Nil(), err
		}
		high, err := r.evalExpr(ex.High)
		if err != nil {
			return types.Nil(), err
		}
		if val.Type != types.NumberT || low.Type != types.NumberT || high.Type != types.NumberT {
			return types.Nil(), fmt.Errorf("'is between' requires numbers")
		}
		return types.FromBoolean(val.Number >= low.Number && val.Number <= high.Number), nil

	case ast.AbilityLit:
		// Store as a closure-like object in the env
		obj := &types.Object{
			Name:     "__ability__",
			Fields:   map[string]types.Value{},
			Actions:  map[string]*types.Action{},
			Builtins: map[string]types.BuiltinFn{},
		}
		obj.Actions["__call__"] = &types.Action{
			Name:   "__call__",
			Params: ex.Params,
			Body:   ex.Body,
		}
		return types.FromObject(obj), nil

	case ast.InterpStringExpr:
		var sb strings.Builder
		for _, seg := range ex.Segments {
			v, err := r.evalExpr(seg)
			if err != nil {
				return types.Nil(), err
			}
			sb.WriteString(v.String())
		}
		return types.FromString(sb.String()), nil

	case ast.CallExpr:
		if r.ambiguous[ex.Callee] {
			return types.Nil(), fmt.Errorf(
				"%q is defined in more than one summoned file — use the file's alias to choose which one (e.g. summon \"...\" as theirname, then theirname.%s(...))",
				ex.Callee, ex.Callee,
			)
		}
		// Try as a stored ability (anonymous function) first
		if v, ok := r.getVar(ex.Callee); ok && v.Type == types.ObjectT && v.Object != nil {
			if _, hasCall := v.Object.Actions["__call__"]; hasCall {
				return r.callMethod(v.Object, "__call__", ex.Args)
			}
		}
		// Also check current receiver fields (inside object methods)
		if r.currentRecv != nil {
			if v, ok := r.currentRecv.Fields[ex.Callee]; ok && v.Type == types.ObjectT && v.Object != nil {
				if _, hasCall := v.Object.Actions["__call__"]; hasCall {
					return r.callMethod(v.Object, "__call__", ex.Args)
				}
			}
		}
		obj, err := r.getObject(ex.Callee)
		if err != nil {
			return types.Nil(), fmt.Errorf("I don't know what %q is", ex.Callee)
		}
		return r.callMethod(obj, ex.Callee, ex.Args)
	}

	return types.Nil(), fmt.Errorf("unknown expression type %T", e)
}


// execBlock runs a slice of statements, returning first error.
func (r *Runtime) execBlock(stmts []ast.Statement) (err error) {
	for _, s := range stmts {
		if err = r.execStatement(s); err != nil {
			return
		}
	}
	return
}

// ── Method dispatch ───────────────────────────────────────────────────────────

func (r *Runtime) callMethod(recv *types.Object, method string, argExprs []ast.Expression) (types.Value, error) {
	// Evaluate args first
	evalArgs := make([]types.Value, 0, len(argExprs))
	for _, a := range argExprs {
		v, err := r.evalExpr(a)
		if err != nil {
			return types.Nil(), err
		}
		evalArgs = append(evalArgs, v)
	}

	// Stage A, Step 3: entitlement enforcement. A call is gated if either
	// the receiving object carries a whole-object tag, or the specific
	// ability being called carries its own tag. Either one denying is
	// enough to deny the call — protection never weakens by combining.
	if recv.ProtectTag != "" {
		if err := r.checkEntitlement(recv.ProtectTag); err != nil {
			return types.Nil(), err
		}
	}
	if recv.Actions != nil {
		if act, ok := recv.Actions[method]; ok && act != nil && act.ProtectTag != "" {
			if err := r.checkEntitlement(act.ProtectTag); err != nil {
				return types.Nil(), err
			}
		}
	}

	// Builtin (native Go function) takes priority
	if recv.Builtins != nil {
		if fn, ok := recv.Builtins[method]; ok {
			return fn(evalArgs)
		}
	}

	// User-defined action
	if recv.Actions != nil {
		if act, ok := recv.Actions[method]; ok && act != nil && len(act.Body) > 0 {
			return r.execAction(recv, act, evalArgs)
		}
	}

	// Field-as-callable fallback: summon "math.he" as maths wraps each
	// ability as a single-purpose object stored under maths.Fields[method].
	// maths.sqrt(16) means "call the wrapper object stored in field sqrt".
	if recv.Fields != nil {
		if fv, ok := recv.Fields[method]; ok && fv.Type == types.ObjectT && fv.Object != nil {
			inner := fv.Object
			if inner.ProtectTag != "" {
				if err := r.checkEntitlement(inner.ProtectTag); err != nil {
					return types.Nil(), err
				}
			}
			if act, ok := inner.Actions[method]; ok && act != nil && act.ProtectTag != "" {
				if err := r.checkEntitlement(act.ProtectTag); err != nil {
					return types.Nil(), err
				}
			}
			if bfn, ok := inner.Builtins[method]; ok {
				return bfn(evalArgs)
			}
			if act, ok := inner.Actions[method]; ok && act != nil {
				return r.execAction(inner, act, evalArgs)
			}
		}
	}

	return types.Nil(), fmt.Errorf(
		"%q doesn't know how to %q — is it defined in a 'can' block?",
		recv.Name, method,
	)
}

// checkEntitlement consults the runtime's EntitlementChecker for the given
// tag and returns a catchable error if access is denied. This is the only
// place enforcement decisions are made — Stage A, Step 3 of the protection
// roadmap. The error is a plain Go error, which means it flows through
// HE's existing try/or mechanism exactly like any other runtime error
// (division by zero, missing file, etc.) — no new control-flow concept
// was needed.
func (r *Runtime) checkEntitlement(tag string) error {
	granted, err := r.entitlement.Check(tag)
	if err != nil {
		return fmt.Errorf("couldn't verify access for %q: %v", tag, err)
	}
	if !granted {
		return fmt.Errorf("access denied — %q is not currently available to you", tag)
	}
	return nil
}

func (r *Runtime) execAction(recv *types.Object, act *types.Action, args []types.Value) (types.Value, error) {
	// Recursion depth guard
	r.callDepth++
	if r.callDepth > 500 {
		r.callDepth--
		return types.Nil(), fmt.Errorf("maximum recursion depth reached in %q — HE currently limits to 500 nested calls", act.Name)
	}
	defer func() { r.callDepth-- }()

	prevRecv := r.currentRecv
	r.currentRecv = recv
	defer func() { r.currentRecv = prevRecv }()

	// Bind params into a temporary local scope overlay
	// Save previous param values so recursive calls don't clobber them
	paramBackup := make(map[string]types.Value, len(act.Params))
	for i, paramName := range act.Params {
		if prev, hasPrev := recv.Fields[paramName]; hasPrev {
			paramBackup[paramName] = prev
		}
		if i < len(args) {
			recv.Fields[paramName] = args[i]
		} else {
			recv.Fields[paramName] = types.Nil()
		}
	}
	defer func() {
		// Restore param values after call returns
		for _, paramName := range act.Params {
			if prev, had := paramBackup[paramName]; had {
				recv.Fields[paramName] = prev
			} else {
				delete(recv.Fields, paramName)
			}
		}
	}()

	for _, st := range act.Body {
		if err := r.execStatement(st); err != nil {
			if sig, ok := err.(*returnSignal); ok {
				return sig.value, nil
			}
			return types.Nil(), err
		}
	}
	return types.Nil(), nil
}

// ── Reactions ─────────────────────────────────────────────────────────────────

func (r *Runtime) triggerReactions(recv *types.Object, event string, right *string) error {
	if recv == nil {
		return nil
	}
	for _, rx := range recv.Reactions {
		if !strings.EqualFold(rx.Trigger.Left, event) {
			continue
		}
		// If trigger has a qualifier, match it
		if rx.Trigger.Right != nil && right != nil {
			if !strings.EqualFold(*rx.Trigger.Right, *right) {
				continue
			}
		}
		prevRecv := r.currentRecv
		r.currentRecv = recv
		for _, st := range rx.Body {
			if err := r.execStatement(st); err != nil {
				r.currentRecv = prevRecv
				return err
			}
		}
		r.currentRecv = prevRecv
	}
	return nil
}

func (r *Runtime) triggerReactionsAll(event string, right *string) error {
	for _, obj := range r.Objects {
		if err := r.triggerReactions(obj, event, right); err != nil {
			return err
		}
	}
	return nil
}

// ── Operators ─────────────────────────────────────────────────────────────────

func isTruthy(v types.Value) bool {
	switch v.Type {
	case types.BooleanT:
		return v.Boolean
	case types.NilT:
		return false
	case types.NumberT:
		return v.Number != 0
	case types.StringT:
		return v.Str != ""
	default:
		return true
	}
}

func evalBinary(op string, l, r types.Value) (types.Value, error) {
	switch op {
	case "+":
		if l.Type == types.NumberT && r.Type == types.NumberT {
			return types.FromNumber(l.Number + r.Number), nil
		}
		// String concatenation
		return types.FromString(l.String() + r.String()), nil
	case "-":
		if l.Type != types.NumberT || r.Type != types.NumberT {
			return types.Nil(), fmt.Errorf("'-' expects numbers")
		}
		return types.FromNumber(l.Number - r.Number), nil
	case "*":
		if l.Type != types.NumberT || r.Type != types.NumberT {
			return types.Nil(), fmt.Errorf("'*' expects numbers")
		}
		return types.FromNumber(l.Number * r.Number), nil
	case "/":
		if l.Type != types.NumberT || r.Type != types.NumberT {
			return types.Nil(), fmt.Errorf("'/' expects numbers")
		}
		if r.Number == 0 {
			return types.Nil(), fmt.Errorf("can't divide by zero")
		}
		return types.FromNumber(l.Number / r.Number), nil
	}
	return types.Nil(), fmt.Errorf("unknown operator %q", op)
}

func evalCompare(op string, l, r types.Value) (types.Value, error) {
	switch op {
	case "==":
		return types.FromBoolean(l.String() == r.String()), nil
	case "!=":
		return types.FromBoolean(l.String() != r.String()), nil
	}
	if l.Type != types.NumberT || r.Type != types.NumberT {
		return types.Nil(), fmt.Errorf("'%s' comparison expects numbers", op)
	}
	switch op {
	case ">":
		return types.FromBoolean(l.Number > r.Number), nil
	case "<":
		return types.FromBoolean(l.Number < r.Number), nil
	case ">=":
		return types.FromBoolean(l.Number >= r.Number), nil
	case "<=":
		return types.FromBoolean(l.Number <= r.Number), nil
	}
	return types.Nil(), fmt.Errorf("unknown comparison %q", op)
}

// loadFileModule reads, parses, and executes a .he file. Its top-level
// names (variables and objects) are flat-merged directly into the
// importing scope — summon "math.he" lets you write sqrt(16) immediately,
// no alias required.
//
// If an alias is also given (summon "math.he" as m), the same contents
// are ADDITIONALLY exposed as a qualified object m.sqrt(16) — the alias
// never gates access, it's a second path to the same names.
//
// If two different summoned files define the same name, that name becomes
// ambiguous: bare access raises an error telling the person to alias one
// or both files, while qualified access (alias.name) keeps working.
func (r *Runtime) loadFileModule(filePath, alias string, hasAlias bool) error {
	resolvedPath := filePath
	data, err := os.ReadFile(resolvedPath)
	if err != nil && r.baseDir != "" && !filepath.IsAbs(filePath) {
		alt := filepath.Join(r.baseDir, filePath)
		if altData, altErr := os.ReadFile(alt); altErr == nil {
			resolvedPath = alt
			data = altData
			err = nil
		}
	}
	if err != nil {
		return fmt.Errorf("summon: can't load %q: %v", filePath, err)
	}

	lx := lexer.New(string(data))
	p := parser.New(lx)
	prog, err := p.ParseProgram()
	if err != nil {
		return fmt.Errorf("summon %q: parse error: %v", filePath, err)
	}

	sub := newRuntime()
	sub.baseDir = filepath.Dir(resolvedPath)
	if err := sub.execProgram(prog); err != nil {
		return fmt.Errorf("summon %q: %v", filePath, err)
	}

	// Collect the module's raw exports (top-level variables and objects).
	raw := map[string]types.Value{}
	for k, v := range sub.Env {
		if k == "platform" || k == "nothing" || k == "pi" {
			continue
		}
		raw[k] = v
	}
	for k, obj := range sub.Objects {
		raw[k] = types.FromObject(obj)
	}

	// flatNames is everything this summon makes available — the raw
	// top-level exports PLUS, for every exported object, a directly
	// callable wrapper per ability/builtin (so `sqrt(16)` works after
	// summoning a math.he that defines `create Math [can: [sqrt(n)...]]]`).
	// This single map drives both the flat global merge AND the alias
	// object's fields, so `theirname.sqrt(16)` and bare `sqrt(16)` always
	// agree on what's available.
	flatNames := map[string]types.Value{}
	for name, val := range raw {
		flatNames[name] = val
		if val.Type != types.ObjectT || val.Object == nil {
			continue
		}
		for abilityName, action := range val.Object.Actions {
			wrapper := &types.Object{
				Name:     abilityName,
				Fields:   val.Object.Fields, // shares state with the source object
				Actions:  map[string]*types.Action{abilityName: action},
				Builtins: map[string]types.BuiltinFn{},
			}
			if bfn, ok := val.Object.Builtins[abilityName]; ok {
				wrapper.Builtins[abilityName] = bfn
			}
			flatNames[abilityName] = types.FromObject(wrapper)
		}
		for builtinName, bfn := range val.Object.Builtins {
			if _, already := val.Object.Actions[builtinName]; already {
				continue
			}
			wrapper := &types.Object{
				Name:     builtinName,
				Fields:   map[string]types.Value{},
				Actions:  map[string]*types.Action{},
				Builtins: map[string]types.BuiltinFn{builtinName: bfn},
			}
			flatNames[builtinName] = types.FromObject(wrapper)
		}
	}

	// Flat-merge into the importing scope, tracking origin for collisions.
	// A name colliding with a *different* file is marked ambiguous and
	// left out of bare scope — qualified alias access still works for it.
	for name, val := range flatNames {
		if origin, exists := r.nameOrigin[name]; exists && origin != filePath {
			r.ambiguous[name] = true
			continue
		}
		r.nameOrigin[name] = filePath
		if !r.ambiguous[name] {
			r.setVar(name, val)
		}
	}

	// If an alias was given, ALSO expose every flat name as a qualified
	// object field — additive, never required, and works even when a name
	// is ambiguous at the bare/global level.
	if hasAlias {
		mod := &types.Object{
			Name:     alias,
			Fields:   flatNames,
			Actions:  make(map[string]*types.Action),
			Builtins: make(map[string]types.BuiltinFn),
		}
		r.setVar(alias, types.FromObject(mod))
	}

	return nil
}
// ── Virtual method dispatch (Pass 6) ─────────────────────────────────────────
// These enable value.method() syntax without needing a module variable.

func (r *Runtime) callStringMethod(s string, method string, argExprs []ast.Expression) (types.Value, error) {
	args, err := r.evalArgs(argExprs)
	if err != nil {
		return types.Nil(), err
	}
	switch method {
	case "upper":
		return types.FromString(strings.ToUpper(s)), nil
	case "lower":
		return types.FromString(strings.ToLower(s)), nil
	case "length", "len":
		return types.FromNumber(float64(len(s))), nil
	case "trim":
		return types.FromString(strings.TrimSpace(s)), nil
	case "split":
		sep := ","
		if len(args) > 0 {
			sep = args[0].Str
		}
		parts := strings.Split(s, sep)
		vals := make([]types.Value, len(parts))
		for i, p := range parts {
			vals[i] = types.FromString(p)
		}
		return types.FromArray(vals), nil
	case "contains":
		if len(args) == 0 {
			return types.Nil(), fmt.Errorf(".contains() expects a substring argument")
		}
		return types.FromBoolean(strings.Contains(s, args[0].Str)), nil
	case "starts", "startsWith":
		if len(args) == 0 {
			return types.Nil(), fmt.Errorf(".starts() expects a prefix argument")
		}
		return types.FromBoolean(strings.HasPrefix(s, args[0].Str)), nil
	case "ends", "endsWith":
		if len(args) == 0 {
			return types.Nil(), fmt.Errorf(".ends() expects a suffix argument")
		}
		return types.FromBoolean(strings.HasSuffix(s, args[0].Str)), nil
	case "replace":
		if len(args) < 2 {
			return types.Nil(), fmt.Errorf(".replace() expects (old, new)")
		}
		return types.FromString(strings.ReplaceAll(s, args[0].Str, args[1].Str)), nil
	case "number", "toNumber":
		var n float64
		if _, err := fmt.Sscanf(s, "%f", &n); err != nil {
			return types.Nil(), fmt.Errorf("can't convert %q to a number", s)
		}
		return types.FromNumber(n), nil
	case "repeat":
		if len(args) == 0 {
			return types.Nil(), fmt.Errorf(".repeat() expects a count")
		}
		return types.FromString(strings.Repeat(s, int(args[0].Number))), nil
	case "at", "charAt":
		if len(args) == 0 {
			return types.Nil(), fmt.Errorf(".at() expects an index")
		}
		idx := int(args[0].Number)
		runes := []rune(s)
		if idx < 0 || idx >= len(runes) {
			return types.Nil(), fmt.Errorf("string index %d out of range", idx)
		}
		return types.FromString(string(runes[idx])), nil
	case "isEmpty", "empty":
		return types.FromBoolean(len(s) == 0), nil
	}
	return types.Nil(), fmt.Errorf("text doesn't have a .%s() method", method)
}

func (r *Runtime) callListMethod(v types.Value, method string, argExprs []ast.Expression) (types.Value, error) {
	args, err := r.evalArgs(argExprs)
	if err != nil {
		return types.Nil(), err
	}
	arr := v.Array
	switch method {
	case "length", "len":
		return types.FromNumber(float64(len(arr))), nil
	case "first":
		if len(arr) == 0 {
			return types.Nil(), fmt.Errorf("list is empty")
		}
		return arr[0], nil
	case "last":
		if len(arr) == 0 {
			return types.Nil(), fmt.Errorf("list is empty")
		}
		return arr[len(arr)-1], nil
	case "add", "append", "push":
		if len(args) == 0 {
			return types.Nil(), fmt.Errorf(".add() expects an item")
		}
		return types.FromArray(append(append([]types.Value{}, arr...), args[0])), nil
	case "contains", "has", "includes":
		if len(args) == 0 {
			return types.Nil(), fmt.Errorf(".contains() expects an item")
		}
		needle := args[0].String()
		for _, item := range arr {
			if item.String() == needle {
				return types.FromBoolean(true), nil
			}
		}
		return types.FromBoolean(false), nil
	case "get", "at":
		if len(args) == 0 {
			return types.Nil(), fmt.Errorf(".get() expects an index")
		}
		idx := int(args[0].Number)
		if idx < 0 || idx >= len(arr) {
			return types.Nil(), fmt.Errorf("list index %d out of range", idx)
		}
		return arr[idx], nil
	case "remove":
		if len(args) == 0 {
			return types.Nil(), fmt.Errorf(".remove() expects an index")
		}
		idx := int(args[0].Number)
		if idx < 0 || idx >= len(arr) {
			return types.Nil(), fmt.Errorf("list index %d out of range", idx)
		}
		newArr := append(append([]types.Value{}, arr[:idx]...), arr[idx+1:]...)
		return types.FromArray(newArr), nil
	case "isEmpty", "empty":
		return types.FromBoolean(len(arr) == 0), nil
	case "join":
		sep := ""
		if len(args) > 0 {
			sep = args[0].Str
		}
		parts := make([]string, len(arr))
		for i, item := range arr {
			parts[i] = item.String()
		}
		return types.FromString(strings.Join(parts, sep)), nil
	case "reverse":
		newArr := make([]types.Value, len(arr))
		for i, item := range arr {
			newArr[len(arr)-1-i] = item
		}
		return types.FromArray(newArr), nil
	}
	return types.Nil(), fmt.Errorf("list doesn't have a .%s() method", method)
}

func (r *Runtime) callNumberMethod(n float64, method string, argExprs []ast.Expression) (types.Value, error) {
	args, err := r.evalArgs(argExprs)
	_ = args
	if err != nil {
		return types.Nil(), err
	}
	switch method {
	case "abs":
		if n < 0 {
			return types.FromNumber(-n), nil
		}
		return types.FromNumber(n), nil
	case "floor":
		return types.FromNumber(math.Floor(n)), nil
	case "ceil":
		return types.FromNumber(math.Ceil(n)), nil
	case "round":
		return types.FromNumber(math.Round(n)), nil
	case "sqrt":
		return types.FromNumber(math.Sqrt(n)), nil
	case "text", "toString":
		return types.FromString(types.FromNumber(n).String()), nil
	case "isPositive":
		return types.FromBoolean(n > 0), nil
	case "isNegative":
		return types.FromBoolean(n < 0), nil
	case "isZero":
		return types.FromBoolean(n == 0), nil
	}
	return types.Nil(), fmt.Errorf("number doesn't have a .%s() method", method)
}

// evalArgs evaluates a slice of argument expressions into Values
func (r *Runtime) evalArgs(argExprs []ast.Expression) ([]types.Value, error) {
	vals := make([]types.Value, len(argExprs))
	for i, a := range argExprs {
		v, err := r.evalExpr(a)
		if err != nil {
			return nil, err
		}
		vals[i] = v
	}
	return vals, nil
}
