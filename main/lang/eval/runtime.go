package eval

import (
	"fmt"
	"html"
	"math"
	"strings"
	"time"

	"hunterlang/lang/ast"
	"hunterlang/lang/types"
)

type RuntimeState struct {
	Env map[string]types.Value
	// Modules / objects created via summon/make.
	Objects map[string]*types.Object
}

func newRuntimeState() *RuntimeState {
	return &RuntimeState{
		Env:     map[string]types.Value{},
		Objects: map[string]*types.Object{},
	}
}

// Runtime executes the AST.
// receiver scoping is used while executing user-defined object actions.
type Runtime struct {
	*RuntimeState
	// Builtin behavior switches (platform-aware modules in a real engine).
	Platform string

	// currentRecv is non-nil only while executing a user-defined action body.
	currentRecv *types.Object
}

// newRuntime initializes the runtime and standard globals.
func newRuntime() *Runtime {
	rs := newRuntimeState()
	// Expose platform as a normal variable for `print "Running on " + platform`.
	rs.Env["platform"] = types.FromString("unknown")

	return &Runtime{
		RuntimeState: rs,
		Platform:     "unknown",
	}
}

func (r *Runtime) getVar(name string) (types.Value, bool) {
	v, ok := r.Env[name]
	return v, ok
}

func (r *Runtime) setVar(name string, v types.Value) {
	r.Env[name] = v
	if v.Type == types.ObjectT && v.Object != nil {
		r.Objects[name] = v.Object
	}
}

func (r *Runtime) getObject(name string) (*types.Object, error) {
	obj, ok := r.Objects[name]
	if !ok {
		if v, ok2 := r.Env[name]; ok2 && v.Type == types.ObjectT && v.Object != nil {
			return v.Object, nil
		}
		return nil, fmt.Errorf("unknown object %q", name)
	}
	return obj, nil
}

// ----------- Expressions -----------

func (r *Runtime) evalIdentifier(name string) (types.Value, bool) {
	// Prefer receiver fields when inside user-defined action execution.
	if r.currentRecv != nil && r.currentRecv.Fields != nil {
		if v, ok := r.currentRecv.Fields[name]; ok {
			return v, true
		}
	}
	// Fall back to global env.
	return r.getVar(name)
}

func (r *Runtime) evalExpr(e ast.Expression) (types.Value, error) {
	switch ex := e.(type) {
	case ast.NumberLit:
		return types.FromNumber(ex.Value), nil
	case ast.StringLit:
		return types.FromString(ex.Value), nil
	case ast.BooleanLit:
		return types.FromBoolean(ex.Value), nil
	case ast.IdentifierExpr:
		v, ok := r.evalIdentifier(ex.Name)
		if !ok {
			return types.Nil(), nil // undefined -> nil
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
				return types.Nil(), fmt.Errorf("unary - expects number, got %s", v.Type)
			}
			return types.FromNumber(-v.Number), nil
		case "!":
			return types.FromBoolean(!isTruthy(v)), nil
		default:
			return types.Nil(), fmt.Errorf("unknown unary op %q", ex.Op)
		}

	case ast.BinaryExpr:
		l, err := r.evalExpr(ex.Left)
		if err != nil {
			return types.Nil(), err
		}
		right, err := r.evalExpr(ex.Right)
		if err != nil {
			return types.Nil(), err
		}
		return r.evalBinary(ex.Op, l, right)

	case ast.PowerExpr:
		l, err := r.evalExpr(ex.Left)
		if err != nil {
			return types.Nil(), err
		}
		right, err := r.evalExpr(ex.Right)
		if err != nil {
			return types.Nil(), err
		}
		if l.Type != types.NumberT || right.Type != types.NumberT {
			return types.Nil(), fmt.Errorf("power expects numbers")
		}
		return types.FromNumber(math.Pow(l.Number, right.Number)), nil

	case ast.CompareExpr:
		l, err := r.evalExpr(ex.Left)
		if err != nil {
			return types.Nil(), err
		}
		right, err := r.evalExpr(ex.Right)
		if err != nil {
			return types.Nil(), err
		}
		return r.evalCompare(ex.Op, l, right)

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

	case ast.CallExpr:
		fn, err := r.getObject(ex.Callee)
		if err != nil {
			if v, ok := r.Env[ex.Callee]; ok && v.Type == types.ObjectT && v.Object != nil {
				fn = v.Object
			} else {
				return types.Nil(), fmt.Errorf("unknown function %q", ex.Callee)
			}
		}
		return r.callMethod(fn, ex.Callee, ex.Args)

	case ast.MethodCallExpr:
		recvVal, err := r.evalExpr(ex.Receiver)
		if err != nil {
			return types.Nil(), err
		}
		if recvVal.Type != types.ObjectT || recvVal.Object == nil {
			return types.Nil(), fmt.Errorf("method call receiver is not an object")
		}
		return r.callMethod(recvVal.Object, ex.Method, ex.Args)

	default:
		return types.Nil(), fmt.Errorf("unknown expression type %T", e)
	}
}

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

func (r *Runtime) evalBinary(op string, l, right types.Value) (types.Value, error) {
	switch op {
	case "+":
		if l.Type == types.NumberT && right.Type == types.NumberT {
			return types.FromNumber(l.Number + right.Number), nil
		}
		if l.Type == types.StringT {
			return types.FromString(l.Str + right.String()), nil
		}
		if right.Type == types.StringT {
			return types.FromString(l.String() + right.Str), nil
		}
		return types.Nil(), fmt.Errorf("+ expects number or string")
	case "-":
		if l.Type != types.NumberT || right.Type != types.NumberT {
			return types.Nil(), fmt.Errorf("- expects numbers")
		}
		return types.FromNumber(l.Number - right.Number), nil
	case "*":
		if l.Type != types.NumberT || right.Type != types.NumberT {
			return types.Nil(), fmt.Errorf("* expects numbers")
		}
		return types.FromNumber(l.Number * right.Number), nil
	case "/":
		if l.Type != types.NumberT || right.Type != types.NumberT {
			return types.Nil(), fmt.Errorf("/ expects numbers")
		}
		return types.FromNumber(l.Number / right.Number), nil
	default:
		return types.Nil(), fmt.Errorf("unknown binary op %q", op)
	}
}

func (r *Runtime) evalCompare(op string, l, right types.Value) (types.Value, error) {
	switch op {
	case "==":
		return types.FromBoolean(l.String() == right.String()), nil
	case "!=":
		return types.FromBoolean(l.String() != right.String()), nil
	}

	if l.Type != types.NumberT || right.Type != types.NumberT {
		return types.Nil(), fmt.Errorf("comparison %s expects numbers", op)
	}

	switch op {
	case ">":
		return types.FromBoolean(l.Number > right.Number), nil
	case "<":
		return types.FromBoolean(l.Number < right.Number), nil
	case ">=":
		return types.FromBoolean(l.Number >= right.Number), nil
	case "<=":
		return types.FromBoolean(l.Number <= right.Number), nil
	default:
		return types.Nil(), fmt.Errorf("unknown compare op %q", op)
	}
}

// ----------- Statements -----------

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
		obj := &types.Object{
			Name:    l.Alias.Name,
			Fields:  map[string]types.Value{},
			Actions: map[string]*types.Action{},
		}
		r.setVar(l.Alias.Name, types.FromObject(obj))

		switch l.ModuleName.Lexeme {
		case "ui":
			r.registerBuiltinMethod(obj, "navbar", builtinNavbar)
			r.registerBuiltinMethod(obj, "renderDocs", builtinRenderDocs)
			r.registerBuiltinMethod(obj, "window", builtinUIWindow)
		case "physics":
			r.registerBuiltinMethod(obj, "gravity", builtinGravity)
			r.registerBuiltinMethod(obj, "collision", builtinCollision)
		}
		return nil

	case ast.GlobalStatementLine:
		return r.execStatement(l.Statement)

	case *ast.GlobalStatementLine:
		return r.execStatement(l.Statement)

	case *ast.ObjectLine:
		return r.execObjectLine(l)

	case ast.WithAssetsLine:
		return nil

	default:
		return fmt.Errorf("unhandled line type %T", ln)
	}
}

func (r *Runtime) execObjectLine(o *ast.ObjectLine) error {
	obj := &types.Object{
		Name:    o.Name.Name,
		Fields:  map[string]types.Value{},
		Actions: map[string]*types.Action{},
	}

	// Inherit from `like <Identifier>` if present.
	if o.Like != nil {
		if parent, ok := r.Objects[o.Like.Name]; ok && parent != nil {
			for k, v := range parent.Fields {
				obj.Fields[k] = v
			}
			for k, act := range parent.Actions {
				obj.Actions[k] = act
			}
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
			}

		case ast.AbilitiesSection:
			for _, act := range s.Abilities {
				// params/returns are ignored for now; we execute body statements only.
				obj.Actions[act.Name] = &types.Action{Name: act.Name, Body: act.Body}
			}

		case ast.ReactionsSection:
			// Store reactions on the object so builtins/events can trigger them later.
			for _, rx := range s.Reactions {
				obj.Reactions = append(obj.Reactions, types.Reaction{
					Trigger: rx.Trigger,
					Body:    rx.Body,
				})
			}

		case ast.MemoriesSection:
			_ = s
		}
	}

	r.setVar(o.Name.Name, types.FromObject(obj))
	return nil
}

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
		// When inside a user-defined receiver action, assign into receiver fields.
		if r.currentRecv != nil && r.currentRecv.Fields != nil {
			r.currentRecv.Fields[st.Name] = v
		} else {
			r.setVar(st.Name, v)
		}
		return nil

	case ast.DecideStmt:
		cond, err := r.evalExpr(st.Cond)
		if err != nil {
			return err
		}
		if isTruthy(cond) {
			for _, t := range st.Then {
				if err := r.execStatement(t); err != nil {
					return err
				}
			}
		} else {
			for _, t := range st.Else {
				if err := r.execStatement(t); err != nil {
					return err
				}
			}
		}
		return nil

	case ast.RepeatStmt:
		max := 1000
		for i := 0; i < max; i++ {
			cond, err := r.evalExpr(st.Cond)
			if err != nil {
				return err
			}
			if st.Kind == "while" {
				if !isTruthy(cond) {
					break
				}
			} else {
				if !isTruthy(cond) && i > 0 {
					break
				}
			}
			for _, b := range st.Body {
				if err := r.execStatement(b); err != nil {
					return err
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
		secV, err := r.evalExpr(st.Expr)
		if err != nil {
			return err
		}
		if secV.Type != types.NumberT {
			return fmt.Errorf("wait expects number")
		}
		switch st.Unit {
		case "seconds":
			time.Sleep(time.Duration(secV.Number * float64(time.Second)))
		case "frames":
			time.Sleep(time.Duration(secV.Number * (1.0 / 60.0) * float64(time.Second)))
		}
		return nil

	case ast.ReturnStmt:
		_, err := r.evalExpr(st.Expr)
		return err

	case ast.ExprStmt:
		_, err := r.evalExpr(st.Expr)
		return err

	default:
		return fmt.Errorf("unhandled statement type %T", s)
	}
}

// ----------- Builtins / methods -----------

type builtinFunc func(i *Runtime, recv *types.Object, args []types.Value) (types.Value, error)

func (r *Runtime) registerBuiltinMethod(obj *types.Object, name string, impl builtinFunc) {
	if obj.Actions == nil {
		obj.Actions = map[string]*types.Action{}
	}
	obj.Actions[name] = &types.Action{Name: name}
	_ = impl
}

func (r *Runtime) callMethod(recv *types.Object, method string, args []ast.Expression) (types.Value, error) {
	evalArgs := make([]types.Value, 0, len(args))
	for _, a := range args {
		v, err := r.evalExpr(a)
		if err != nil {
			return types.Nil(), err
		}
		evalArgs = append(evalArgs, v)
	}
	_ = evalArgs

	// Builtins
	switch recv.Name {
	case "he":
		if method == "navbar" {
			return builtinNavbar(r, recv, evalArgs)
		}
	case "phys":
		if method == "gravity" {
			return builtinGravity(r, recv, evalArgs)
		}
		if method == "collision" {
			// If collision is true, trigger reactions on this receiver object.
			v, err := builtinCollision(r, recv, evalArgs)
			if err != nil {
				return types.Nil(), err
			}
			if v.Type == types.BooleanT && v.Boolean {
				if err := r.triggerReactions(recv, "collision", nil); err != nil {
					return types.Nil(), err
				}
			}
			return v, nil
		}
	}

	// User-defined actions (from `can: [...]`)
	if recv.Actions != nil {
		if act, ok := recv.Actions[method]; ok && act != nil {
			// Enable receiver scoping for the action body.
			prevRecv := r.currentRecv
			r.currentRecv = recv
			defer func() { r.currentRecv = prevRecv }()

			// For now we ignore params/returns; we just execute body statements.
			for _, st := range act.Body {
				if err := r.execStatement(st); err != nil {
					return types.Nil(), err
				}
			}
			return types.Nil(), nil
		}
	}

	return types.Nil(), fmt.Errorf("method %q not implemented on %q", method, recv.Name)
}

// triggerReactions runs matching reaction bodies for a given event on recv.
// For now we match on Trigger.Left == event.
func (r *Runtime) triggerReactions(recv *types.Object, event string, right *string) error {
	_ = right

	if recv == nil || len(recv.Reactions) == 0 {
		return nil
	}

	for _, rx := range recv.Reactions {
		if rx.Trigger.Left != event {
			continue
		}
		// Enable receiver scoping for reaction bodies.
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

// triggerReactionsAll runs matching reaction bodies for a given event across all runtime objects.
// For now we only match Trigger.Left == event (right matching will be added once the event payload is wired).
func (r *Runtime) triggerReactionsAll(event string, right *string) error {
	_ = right

	for _, obj := range r.Objects {
		if obj == nil {
			continue
		}
		if err := r.triggerReactions(obj, event, right); err != nil {
			return err
		}
	}
	return nil
}

func builtinNavbar(i *Runtime, recv *types.Object, args []types.Value) (types.Value, error) {
	_ = i
	_ = recv
	if len(args) != 1 || args[0].Type != types.ArrayT {
		return types.Nil(), fmt.Errorf("he.navbar expects one array arg")
	}
	fmt.Println("UI.navbar:", args[0].String())
	return types.Nil(), nil
}

func builtinUIWindow(i *Runtime, recv *types.Object, args []types.Value) (types.Value, error) {
	_ = i
	_ = recv

	// ui.window expects one array arg: [title, children]
	// children is an array of UI nodes encoded as arrays:
	//   ["text", "hello"] or ["button", "label", "onClickLabel"]
	if len(args) != 1 || args[0].Type != types.ArrayT {
		return types.Nil(), fmt.Errorf("ui.window expects one array arg [title, children]")
	}

	root := args[0].Array
	if len(root) < 2 {
		return types.Nil(), fmt.Errorf("ui.window expects [title, children]")
	}

	title := root[0].String()
	childrenV := root[1]
	if childrenV.Type != types.ArrayT {
		return types.Nil(), fmt.Errorf("ui.window children must be an array")
	}
	children := childrenV.Array

	// Render to HTML with minimal interactivity.
	// Clicks currently only log to console (we’ll wire HE event triggering later).
	html := fmt.Sprintf(`<!doctype html>
<html>
<head>
  <meta charset="utf-8"/>
  <meta name="viewport" content="width=device-width, initial-scale=1"/>
  <title>%s</title>
  <style>
    body{font-family:ui-monospace,SFMono-Regular,Menlo,Monaco,Consolas,monospace;margin:16px;background:#fff;color:#000;}
    .card{border:4px solid #000;background:#f8f8f8;box-shadow:6px 6px 0 #000;padding:14px;margin-top:12px;}
    button{border:4px solid #000;background:#fff;box-shadow:4px 4px 0 #000;padding:10px 14px;font-weight:900;cursor:pointer;}
    button:active{transform:translate(2px,2px);box-shadow:2px 2px 0 #000;}
    .muted{opacity:.78;}
  </style>
</head>
<body>
  <h1>%s</h1>
  <div class="card">
    %s
  </div>
<script>
  function heClick(label){
    console.log("HE click:", label);
  }
</script>
</body>
</html>`, title, title, renderChildrenHTML(children))
	fmt.Println(html)
	return types.Nil(), nil
}

func renderChildrenHTML(children []types.Value) string {
	out := ""
	for _, ch := range children {
		// UI node encoding: ["text", value] or ["button", label, clickLabel]
		if ch.Type != types.ArrayT || len(ch.Array) < 2 {
			continue
		}
		nodeType := ch.Array[0].String()
		switch nodeType {
		case "text":
			out += fmt.Sprintf(`<div>%s</div>`, escapeForHTML(ch.Array[1].String()))
		case "button":
			if len(ch.Array) < 3 {
				out += fmt.Sprintf(`<button onclick="heClick('')">%s</button>`, escapeForHTML(ch.Array[1].String()))
				continue
			}
			label := ch.Array[1].String()
			clickLabel := ch.Array[2].String()
			out += fmt.Sprintf(`<div style="margin-top:10px;"><button onclick="heClick('%s')">%s</button></div>`, escapeForHTML(clickLabel), escapeForHTML(label))
		}
	}
	return out
}

func escapeForHTML(s string) string {
	// Escape for embedding into HTML/JS string literals via onclick handlers.
	// We escape HTML-special characters and also ensure quotes/backticks are safe.
	// html.EscapeString escapes: &, <, >, ' ', " and / (where applicable).
	escaped := html.EscapeString(s)
	// Extra safety for JS string literals inside onclick="heClick('...')"
	escaped = strings.ReplaceAll(escaped, "\\", "\\\\")
	escaped = strings.ReplaceAll(escaped, "'", "\\'")
	return escaped
}

func builtinUIButton(i *Runtime, recv *types.Object, args []types.Value) (types.Value, error) {
	_ = i
	_ = recv

	// ui.button expects ["label", "clickLabel"] but we treat it as a node producer:
	// For now, since HE code will build children arrays manually, we just validate.
	if len(args) != 1 || args[0].Type != types.ArrayT {
		return types.Nil(), fmt.Errorf("ui.button expects one array arg [label, clickLabel]")
	}
	return types.Nil(), nil
}

func builtinRenderDocs(i *Runtime, recv *types.Object, args []types.Value) (types.Value, error) {
	_ = i
	_ = recv

	// args: [ [title, tab1, tab2, ...] ]
	if len(args) != 1 || args[0].Type != types.ArrayT {
		return types.Nil(), fmt.Errorf("ui.renderDocs expects one array arg [title, tab1, tab2, ...]")
	}
	arr := args[0].Array
	if len(arr) < 2 {
		return types.Nil(), fmt.Errorf("ui.renderDocs requires at least [title, tab1]")
	}

	// tab names (strings)
	title := arr[0].String()
	tabs := make([]string, 0, len(arr)-1)
	for _, v := range arr[1:] {
		tabs = append(tabs, v.String())
	}

	// Simple docs page generator.
	// Note: this is runtime-generated (HE intent only), so it’s cross-platform ready.
	html := fmt.Sprintf(`<!doctype html>
<html>
<head>
  <meta charset="utf-8"/>
  <meta name="viewport" content="width=device-width, initial-scale=1"/>
  <title>%s</title>
  <style>
    body{font-family: ui-monospace,SFMono-Regular,Menlo,Monaco,Consolas,monospace;margin:0;background:#fff;color:#000;}
    header{position:sticky;top:0;background:#fff;border-bottom:4px solid #000;padding:12px 16px;}
    .tabs{margin-top:10px;display:flex;gap:10px;flex-wrap:wrap;}
    .tab{border:3px solid #000;padding:8px 12px;box-shadow:4px 4px 0 #000;background:#fff;cursor:pointer;font-weight:900;user-select:none;}
    .tab[aria-selected='true']{background:#f8f8f8;}
    main{padding:16px;max-width:980px;margin:0 auto;}
    .card{border:4px solid #000;background:#f8f8f8;box-shadow:6px 6px 0 #000;padding:14px;margin-top:12px;}
    .section{display:none;}
    .section[aria-hidden='false']{display:block;}
    pre{background:#fff;border:3px solid #000;padding:12px;overflow:auto;box-shadow:4px 4px 0 #000;}
    button.big{border:4px solid #000;background:#fff;box-shadow:8px 8px 0 #000;padding:10px 14px;font-weight:900;cursor:pointer;margin-top:10px;}
    button.big:active{transform:translate(2px,2px);box-shadow:4px 4px 0 #000;}
    .muted{opacity:.78;}
  </style>
</head>
<body>
  <header>
    <div style="font-size:18px;font-weight:900;">%s <span class="muted">ALIVE</span></div>
    <div class="tabs">
%s
    </div>
  </header>
  <main>
%s
  </main>
<script>
  function setTab(tabName){
    document.querySelectorAll('.tab').forEach(t=>{
      const selected = t.getAttribute('data-tab') === tabName;
      t.setAttribute('aria-selected', selected ? 'true' : 'false');
    });
    document.querySelectorAll('.section').forEach(s=>{
      const name = s.getAttribute('data-tab');
      s.setAttribute('aria-hidden', name === tabName ? 'false' : 'true');
    });
  }
  document.querySelectorAll('.tab').forEach(tab=>{
    tab.addEventListener('click', ()=>setTab(tab.getAttribute('data-tab')));
  });
</script>
</body>
</html>`, title, title, renderTabsHTML(tabs), renderSectionsHTML(tabs))

	fmt.Println(html)
	return types.Nil(), nil
}

// renderTabsHTML renders the tab buttons.
func renderTabsHTML(tabs []string) string {
	out := ""
	for idx, t := range tabs {
		selected := "false"
		if idx == 0 {
			selected = "true"
		}
		out += fmt.Sprintf(`      <div class="tab" role="tab" tabindex="0" aria-selected='%s' data-tab='%s'>%s</div>
`, selected, t, t)
	}
	return out
}

// renderSectionsHTML renders placeholder content sections.
func renderSectionsHTML(tabs []string) string {
	out := ""
	for idx, t := range tabs {
		hidden := "true"
		if idx == 0 {
			hidden = "false"
		}
		out += fmt.Sprintf(`    <div class="card section" data-tab='%s' aria-hidden='%s'>
      <h3>%s</h3>
      <p class="muted">This section is generated by the runtime from HE intent (ui.renderDocs).</p>
    </div>
`, t, hidden, t)
	}
	return out
}

func builtinGravity(i *Runtime, recv *types.Object, args []types.Value) (types.Value, error) {
	_ = i
	_ = recv
	if len(args) != 1 || args[0].Type != types.NumberT {
		return types.Nil(), fmt.Errorf("phys.gravity expects one number")
	}
	fmt.Println("Physics.gravity:", args[0].String())
	return types.Nil(), nil
}

func builtinCollision(i *Runtime, recv *types.Object, args []types.Value) (types.Value, error) {
	_ = i
	_ = recv
	if len(args) != 2 {
		return types.Nil(), fmt.Errorf("phys.collision expects two args")
	}
	if args[0].String() == args[1].String() {
		return types.FromBoolean(true), nil
	}
	return types.FromBoolean(false), nil
}
