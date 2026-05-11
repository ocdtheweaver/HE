// main.go
package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

/*
HE interpreter prototype with arithmetic evaluation:
- supports + - * / and parentheses
- supports variables (identifiers) lookup
- supports quoted strings and concatenation via +
- accepts assignment expressed as `x is expr` or `x = expr`
*/

type TokenType int

const (
	TokEOF TokenType = iota
	TokIdentifier
	TokString
	TokNumber
	TokSymbol
	TokKeyword
)

type Token struct {
	Typ TokenType
	Val string
	Pos int
}

var keywords = map[string]bool{
	"summon": true, "as": true, "named": true,
	"with": true, "image": true, "sound": true, "music": true, "video": true, "font": true, "shader": true, "and": true,
	"make": true, "create": true, "like": true, "has": true, "can": true, "on": true, "when": true, "whenever": true,
	"print": true, "tell": true, "to": true, "wait": true, "seconds": true, "is": true,
}

type Asset struct {
	Kind string
	Path string
	Name string
}

type Ability struct {
	Name       string
	Params     []string
	Statements []Stmt
}

type ObjectDef struct {
	Name       string
	Like       string
	Properties map[string]Value
	Abilities  map[string]*Ability
	Reactions  map[string][]Stmt // trigger -> statements
}

type Program struct {
	Summons []struct{ Path, As string }
	Assets  []Asset
	Objects map[string]*ObjectDef
	Globals []Stmt
}

type Value struct {
	Number float64
	Str    string
	Bool   bool
	Type   string // "number", "string", "boolean"
}

type Stmt interface {
	Exec(rt *Runtime) error
	String() string
}

type Runtime struct {
	Program *Program
	// runtime state
	Objects map[string]*ObjectDef
	Assets  map[string]Asset
	Vars    map[string]Value
}

func NewRuntime(p *Program) *Runtime {
	rt := &Runtime{
		Program: p,
		Objects: map[string]*ObjectDef{},
		Assets:  map[string]Asset{},
		Vars:    map[string]Value{},
	}
	// register objects and assets
	for _, a := range p.Assets {
		n := a.Name
		if n == "" {
			n = a.Path
		}
		rt.Assets[n] = a
	}
	for k, v := range p.Objects {
		rt.Objects[k] = v
		// initialize properties into Vars with objectName.propertyName
		for pn, pv := range v.Properties {
			rt.Vars[k+"."+pn] = pv
		}
	}
	return rt
}

// ----------------- Expression evaluation -----------------

// Token for expression evaluation
type exprTok struct {
	typ string // "num","id","str","op","(" , ")"
	val string
}

func tokenizeExpr(s string) ([]exprTok, error) {
	out := []exprTok{}
	i := 0
	for i < len(s) {
		c := s[i]
		// skip whitespace
		if c == ' ' || c == '\t' || c == '\r' || c == '\n' {
			i++
			continue
		}
		// number (digit or dot)
		if (c >= '0' && c <= '9') || (c == '.' && i+1 < len(s) && s[i+1] >= '0' && s[i+1] <= '9') {
			start := i
			dot := false
			for i < len(s) && ((s[i] >= '0' && s[i] <= '9') || s[i] == '.') {
				if s[i] == '.' {
					if dot {
						break
					}
					dot = true
				}
				i++
			}
			out = append(out, exprTok{typ: "num", val: s[start:i]})
			continue
		}
		// identifier (letters, underscore, dot allowed)
		if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '_' {
			start := i
			i++
			for i < len(s) && ((s[i] >= 'a' && s[i] <= 'z') || (s[i] >= 'A' && s[i] <= 'Z') || (s[i] >= '0' && s[i] <= '9') || s[i] == '_' || s[i] == '.') {
				i++
			}
			out = append(out, exprTok{typ: "id", val: s[start:i]})
			continue
		}
		// quoted string
		if c == '"' {
			i++
			start := i
			for i < len(s) {
				if s[i] == '\\' && i+1 < len(s) {
					i += 2
					continue
				}
				if s[i] == '"' {
					break
				}
				i++
			}
			if i >= len(s) {
				return nil, errors.New("unterminated string")
			}
			val := s[start:i]
			i++ // skip closing quote
			out = append(out, exprTok{typ: "str", val: val})
			continue
		}
		// operators and parens
		// support + - * / ( )
		if strings.ContainsRune("+-*/()", rune(c)) {
			t := string(c)
			if t == "(" || t == ")" {
				out = append(out, exprTok{typ: t, val: t})
			} else {
				out = append(out, exprTok{typ: "op", val: t})
			}
			i++
			continue
		}
		// unknown char
		return nil, fmt.Errorf("unexpected char in expression: %c", c)
	}
	return out, nil
}

func precedence(op string) int {
	switch op {
	case "+", "-":
		return 1
	case "*", "/":
		return 2
	}
	return 0
}

func shuntingYard(tokens []exprTok) ([]exprTok, error) {
	out := []exprTok{}
	var stack []exprTok
	for i := 0; i < len(tokens); i++ {
		t := tokens[i]
		switch t.typ {
		case "num", "id", "str":
			out = append(out, t)
		case "op":
			// handle unary minus: if op is "-" and previous token is nil or operator or "(" then treat as unary by converting to (0 - x)
			if t.val == "-" {
				// unary if it's first token or previous is operator or "("
				if i == 0 || (tokens[i-1].typ == "op") || (tokens[i-1].typ == "(") {
					// push a zero then subtraction: we implement by appending num 0 before as if "0 - x"
					out = append(out, exprTok{typ: "num", val: "0"})
				}
			}
			for len(stack) > 0 {
				top := stack[len(stack)-1]
				if top.typ == "op" && precedence(top.val) >= precedence(t.val) {
					out = append(out, top)
					stack = stack[:len(stack)-1]
					continue
				}
				break
			}
			stack = append(stack, t)
		case "(":
			stack = append(stack, t)
		case ")":
			found := false
			for len(stack) > 0 {
				top := stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				if top.typ == "(" {
					found = true
					break
				}
				out = append(out, top)
			}
			if !found {
				return nil, errors.New("mismatched parentheses")
			}
		default:
			return nil, fmt.Errorf("unknown token type: %v", t.typ)
		}
	}
	for len(stack) > 0 {
		top := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if top.typ == "(" || top.typ == ")" {
			return nil, errors.New("mismatched parentheses")
		}
		out = append(out, top)
	}
	return out, nil
}

func evalRPN(rpn []exprTok, rt *Runtime) (Value, error) {
	var stack []Value
	push := func(v Value) { stack = append(stack, v) }
	pop := func() (Value, error) {
		if len(stack) == 0 {
			return Value{}, errors.New("stack underflow")
		}
		v := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		return v, nil
	}
	for _, t := range rpn {
		switch t.typ {
		case "num":
			f, _ := strconv.ParseFloat(t.val, 64)
			push(Value{Number: f, Type: "number"})
		case "str":
			push(Value{Str: t.val, Type: "string"})
		case "id":
			// lookup variable in runtime
			if v, ok := rt.Vars[t.val]; ok {
				push(v)
			} else {
				// if unknown identifier, push zero-number by default
				push(Value{Number: 0, Type: "number"})
			}
		case "op":
			// binary operators only
			b, err := pop()
			if err != nil {
				return Value{}, err
			}
			a, err := pop()
			if err != nil {
				return Value{}, err
			}
			switch t.val {
			case "+":
				// number+number -> number, string+string or mixed -> concat string
				if a.Type == "number" && b.Type == "number" {
					push(Value{Number: a.Number + b.Number, Type: "number"})
				} else {
					// coerce to strings
					as := a.Str
					if a.Type == "number" {
						as = fmt.Sprintf("%v", a.Number)
					}
					bs := b.Str
					if b.Type == "number" {
						bs = fmt.Sprintf("%v", b.Number)
					}
					push(Value{Str: as + bs, Type: "string"})
				}
			case "-":
				push(Value{Number: a.Number - b.Number, Type: "number"})
			case "*":
				push(Value{Number: a.Number * b.Number, Type: "number"})
			case "/":
				push(Value{Number: a.Number / b.Number, Type: "number"})
			default:
				return Value{}, fmt.Errorf("unknown operator %s", t.val)
			}
		default:
			return Value{}, fmt.Errorf("unexpected token in rpn: %v", t.typ)
		}
	}
	if len(stack) != 1 {
		return Value{}, fmt.Errorf("invalid expression (stack len %d)", len(stack))
	}
	return stack[0], nil
}

func evaluateExpression(expr string, rt *Runtime) (Value, error) {
	expr = strings.TrimSpace(expr)
	// quick: if it's a quoted string
	if strings.HasPrefix(expr, "\"") && strings.HasSuffix(expr, "\"") && len(expr) >= 2 {
		return Value{Str: expr[1 : len(expr)-1], Type: "string"}, nil
	}
	toks, err := tokenizeExpr(expr)
	if err != nil {
		return Value{}, err
	}
	rpn, err := shuntingYard(toks)
	if err != nil {
		return Value{}, err
	}
	val, err := evalRPN(rpn, rt)
	if err != nil {
		return Value{}, err
	}
	return val, nil
}

// ----------------- Statements -----------------

// PrintStmt now evaluates expressions
type PrintStmt struct {
	Expr string // expression text
}

func (p *PrintStmt) Exec(rt *Runtime) error {
	// if it's a raw quoted string, print directly
	trim := strings.TrimSpace(p.Expr)
	if strings.HasPrefix(trim, "\"") && strings.HasSuffix(trim, "\"") {
		fmt.Println(trim[1 : len(trim)-1])
		return nil
	}
	// evaluate expression
	v, err := evaluateExpression(p.Expr, rt)
	if err != nil {
		// fallback: try to print literal or identifier
		if vv, ok := rt.Vars[strings.TrimSpace(p.Expr)]; ok {
			if vv.Type == "number" {
				fmt.Println(vv.Number)
			} else {
				fmt.Println(vv.Str)
			}
			return nil
		}
		return fmt.Errorf("print eval error: %v", err)
	}
	if v.Type == "number" {
		// print without trailing decimal if integer
		if v.Number == float64(int64(v.Number)) {
			fmt.Printf("%d\n", int64(v.Number))
		} else {
			fmt.Println(v.Number)
		}
		return nil
	}
	if v.Type == "string" {
		fmt.Println(v.Str)
		return nil
	}
	// boolean not used often here
	fmt.Println(v)
	return nil
}
func (p *PrintStmt) String() string { return "print " + p.Expr }

// WaitStmt unchanged
type WaitStmt struct {
	Seconds float64
}

func (w *WaitStmt) Exec(rt *Runtime) error {
	d := time.Duration(w.Seconds * float64(time.Second))
	time.Sleep(d)
	return nil
}
func (w *WaitStmt) String() string { return fmt.Sprintf("wait %f seconds", w.Seconds) }

// TellStmt unchanged behavior
type TellStmt struct {
	Target string
	Action string
	Args   []string
}

func (t *TellStmt) Exec(rt *Runtime) error {
	// simplistic: support asset playAnimation and play sound names
	if a, ok := rt.Assets[t.Target]; ok {
		fmt.Printf("[asset:%s] %s %v\n", t.Target, t.Action, t.Args)
		if a.Kind == "sound" && (t.Action == "play" || t.Action == "playSound") {
			fmt.Printf("PLAY SOUND %s\n", a.Path)
		}
		return nil
	}
	// if target is an object and action is a method (ability), try to run ability body
	if obj, ok := rt.Objects[t.Target]; ok {
		if ab, ok := obj.Abilities[t.Action]; ok {
			for _, s := range ab.Statements {
				if err := s.Exec(rt); err != nil {
					return err
				}
			}
			return nil
		}
	}
	fmt.Printf("[tell] %s to %s %v\n", t.Target, t.Action, t.Args)
	return nil
}
func (t *TellStmt) String() string { return fmt.Sprintf("tell %s to %s", t.Target, t.Action) }

// AssignStmt now evaluates RHS expression fully
type AssignStmt struct {
	Lhs  string
	Expr string
}

func (a *AssignStmt) Exec(rt *Runtime) error {
	expr := strings.TrimSpace(a.Expr)
	// evaluate expression
	val, err := evaluateExpression(expr, rt)
	if err == nil {
		rt.Vars[a.Lhs] = val
		return nil
	}
	// fallback: if expression is identifier or numeric or quoted handled in evaluateExpression; error otherwise
	// to be safe, if evaluateExpression failed, try simple parse numeric or string
	if f, err2 := strconv.ParseFloat(expr, 64); err2 == nil {
		rt.Vars[a.Lhs] = Value{Number: f, Type: "number"}
		return nil
	}
	if strings.HasPrefix(expr, "\"") && strings.HasSuffix(expr, "\"") && len(expr) >= 2 {
		rt.Vars[a.Lhs] = Value{Str: expr[1 : len(expr)-1], Type: "string"}
		return nil
	}
	// identifier fallback
	if v, ok := rt.Vars[expr]; ok {
		rt.Vars[a.Lhs] = v
		return nil
	}
	// unknown, set empty string
	rt.Vars[a.Lhs] = Value{Str: expr, Type: "string"}
	return nil
}
func (a *AssignStmt) String() string { return fmt.Sprintf("%s is %s", a.Lhs, a.Expr) }

// ----------------- Parser (mostly unchanged) -----------------

type Parser struct {
	src   string
	lines []string
	pos   int
}

func NewParser(src string) *Parser {
	scanner := bufio.NewScanner(strings.NewReader(src))
	lines := []string{}
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return &Parser{src: src, lines: lines, pos: 0}
}

func trimCommentMarkers(line string) string {
	trim := strings.TrimSpace(line)
	if strings.HasPrefix(trim, "~") && strings.HasSuffix(trim, "~") {
		trim = strings.TrimSpace(trim[1 : len(trim)-1])
		return trim
	}
	return line
}

func (p *Parser) hasMore() bool { return p.pos < len(p.lines) }

func (p *Parser) peekLine() string {
	if p.hasMore() {
		return p.lines[p.pos]
	}
	return ""
}

func (p *Parser) nextLine() string {
	if p.hasMore() {
		l := p.lines[p.pos]
		p.pos++
		return l
	}
	return ""
}

func stripInlineComments(line string) string {
	out := ""
	i := 0
	for i < len(line) {
		if line[i] == '~' {
			j := i + 1
			for j < len(line) && line[j] != '~' {
				j++
			}
			if j < len(line) && line[j] == '~' {
				i = j + 1
				continue
			}
			out += string(line[i])
			i++
		} else {
			out += string(line[i])
			i++
		}
	}
	return out
}

func (p *Parser) ParseProgram() (*Program, error) {
	prg := &Program{Objects: map[string]*ObjectDef{}}
	for p.hasMore() {
		line := p.nextLine()
		trim := strings.TrimSpace(line)
		if trim == "" {
			continue
		}
		if strings.HasPrefix(strings.TrimSpace(trim), "~") && strings.HasSuffix(strings.TrimSpace(trim), "~") {
			continue
		}
		lower := strings.ToLower(trim)
		if strings.HasPrefix(lower, "summon") {
			path, asn, err := parseSummon(trim)
			if err != nil {
				return nil, err
			}
			prg.Summons = append(prg.Summons, struct{ Path, As string }{path, asn})
			continue
		}
		if strings.HasPrefix(lower, "with") {
			assets, err := p.parseWith(trim)
			if err != nil {
				return nil, err
			}
			prg.Assets = append(prg.Assets, assets...)
			continue
		}
		if strings.HasPrefix(lower, "make ") || strings.HasPrefix(lower, "create ") {
			obj, err := p.parseObject(trim)
			if err != nil {
				return nil, err
			}
			prg.Objects[obj.Name] = obj
			continue
		}
		stmt, err := parseGlobalStatement(trim)
		if err == nil && stmt != nil {
			prg.Globals = append(prg.Globals, stmt)
			continue
		}
	}
	return prg, nil
}

func parseSummon(line string) (string, string, error) {
	parts := strings.Fields(line)
	if len(parts) < 2 {
		return "", "", errors.New("invalid summon")
	}
	start := strings.Index(line, "\"")
	end := strings.LastIndex(line, "\"")
	if start == -1 || end == -1 || end == start {
		return "", "", errors.New("summon requires a string path")
	}
	path := strings.TrimSpace(line[start+1 : end])
	rest := strings.TrimSpace(line[end+1:])
	asname := ""
	if rest != "" {
		rparts := strings.Fields(rest)
		if len(rparts) >= 2 && (rparts[0] == "as" || rparts[0] == "named") {
			asname = rparts[1]
		}
	}
	return path, asname, nil
}

func (p *Parser) parseWith(firstLine string) ([]Asset, error) {
	assets := []Asset{}
	buf := firstLine
	for p.hasMore() {
		nl := p.peekLine()
		ts := strings.TrimSpace(nl)
		if ts == "" {
			break
		}
		if strings.HasPrefix(ts, "~") && strings.HasSuffix(ts, "~") {
			p.nextLine()
			continue
		}
		low := strings.ToLower(ts)
		if strings.HasPrefix(low, "and ") || strings.HasPrefix(low, ",") || strings.HasPrefix(low, "image ") || strings.HasPrefix(low, "sound ") ||
			strings.HasPrefix(low, "music ") || strings.HasPrefix(low, "video ") || strings.HasPrefix(low, "font ") || strings.HasPrefix(low, "shader ") {
			buf += " " + ts
			p.nextLine()
			continue
		}
		break
	}
	buf = strings.TrimSpace(buf)
	if strings.HasPrefix(strings.ToLower(buf), "with ") {
		buf = strings.TrimSpace(buf[5:])
	}
	parts := splitAssets(buf)
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		words := strings.Fields(part)
		if len(words) < 2 {
			continue
		}
		typ := strings.ToLower(words[0])
		start := strings.Index(part, "\"")
		end := strings.LastIndex(part, "\"")
		if start == -1 || end == -1 || end == start {
			continue
		}
		path := part[start+1 : end]
		after := strings.TrimSpace(part[end+1:])
		name := ""
		if after != "" {
			asParts := strings.Fields(after)
			for i := 0; i < len(asParts); i++ {
				if asParts[i] == "named" || asParts[i] == "as" {
					if i+1 < len(asParts) {
						name = asParts[i+1]
					}
				}
			}
		}
		a := Asset{Kind: typ, Path: path, Name: name}
		assets = append(assets, a)
	}
	return assets, nil
}

func splitAssets(s string) []string {
	out := []string{}
	cur := ""
	inQuote := false
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c == '"' {
			inQuote = !inQuote
			cur += string(c)
			continue
		}
		if !inQuote && i+5 <= len(s) && strings.ToLower(s[i:i+5]) == " and " {
			out = append(out, strings.TrimSpace(cur))
			cur = ""
			i += 4
			continue
		}
		if !inQuote && c == ',' {
			out = append(out, strings.TrimSpace(cur))
			cur = ""
			continue
		}
		cur += string(c)
	}
	if strings.TrimSpace(cur) != "" {
		out = append(out, strings.TrimSpace(cur))
	}
	return out
}

func (p *Parser) parseObject(firstLine string) (*ObjectDef, error) {
	l := strings.TrimSpace(firstLine)
	parts := strings.Fields(l)
	if len(parts) < 2 {
		return nil, errors.New("invalid make")
	}
	name := parts[1]
	like := ""
	if len(parts) >= 4 && parts[2] == "like" {
		like = parts[3]
	}
	if !strings.Contains(l, "{") {
		for p.hasMore() {
			nl := p.nextLine()
			if strings.Contains(nl, "{") {
				break
			}
		}
	}
	obj := &ObjectDef{
		Name:       name,
		Like:       like,
		Properties: map[string]Value{},
		Abilities:  map[string]*Ability{},
		Reactions:  map[string][]Stmt{},
	}
	for p.hasMore() {
		line := p.nextLine()
		trim := strings.TrimSpace(line)
		if trim == "" {
			continue
		}
		if strings.HasPrefix(trim, "~") && strings.HasSuffix(trim, "~") {
			continue
		}
		if strings.HasPrefix(trim, "}") {
			break
		}
		low := strings.ToLower(trim)
		if strings.Contains(low, " has:") || strings.Contains(low, " has :") || strings.Contains(low, " has: [") {
			for p.hasMore() {
				l2 := p.nextLine()
				t2 := strings.TrimSpace(l2)
				if t2 == "" {
					continue
				}
				if strings.HasPrefix(t2, "~") && strings.HasSuffix(t2, "~") {
					continue
				}
				if strings.HasPrefix(t2, "]") || strings.HasPrefix(t2, "]") {
					break
				}
				if strings.Contains(t2, " is ") {
					parts := strings.SplitN(t2, " is ", 2)
					prop := strings.TrimSpace(parts[0])
					expr := strings.TrimSpace(parts[1])
					val := Value{}
					if strings.HasPrefix(expr, "\"") && strings.HasSuffix(expr, "\"") {
						val.Type = "string"
						val.Str = expr[1 : len(expr)-1]
					} else {
						if f, err := strconv.ParseFloat(expr, 64); err == nil {
							val.Type = "number"
							val.Number = f
						} else {
							val.Type = "string"
							val.Str = expr
						}
					}
					obj.Properties[prop] = val
				}
			}
			continue
		}
		if strings.Contains(low, " can:") || strings.Contains(low, " can :") {
			for p.hasMore() {
				l2 := p.nextLine()
				t2 := strings.TrimSpace(l2)
				if t2 == "" {
					continue
				}
				if strings.HasPrefix(t2, "~") && strings.HasSuffix(t2, "~") {
					continue
				}
				if strings.HasPrefix(t2, "]") {
					break
				}
				if strings.HasSuffix(t2, "[") {
					head := strings.TrimSpace(t2[:len(t2)-1])
					nameEnd := strings.Index(head, "(")
					methodName := head
					if nameEnd != -1 {
						methodName = strings.TrimSpace(head[:nameEnd])
					}
					ab := &Ability{Name: methodName, Statements: []Stmt{}}
					for p.hasMore() {
						l3 := p.nextLine()
						t3 := strings.TrimSpace(l3)
						if t3 == "" {
							continue
						}
						if strings.HasPrefix(t3, "~") && strings.HasSuffix(t3, "~") {
							continue
						}
						if t3 == "]" {
							break
						}
						st, err := parseStmt(t3)
						if err == nil && st != nil {
							ab.Statements = append(ab.Statements, st)
						}
					}
					obj.Abilities[ab.Name] = ab
				}
			}
			continue
		}
		if strings.HasPrefix(low, "on ") || strings.HasPrefix(low, "when ") || strings.HasPrefix(low, "whenever ") {
			trigger := ""
			parts := strings.Fields(trim)
			if len(parts) >= 2 {
				trigger = parts[1]
			}
			stmts := []Stmt{}
			if strings.Contains(trim, "[") && strings.Contains(trim, "]") {
				inner := trim[strings.Index(trim, "[")+1 : strings.LastIndex(trim, "]")]
				lines := strings.Split(inner, "\n")
				for _, l := range lines {
					l = strings.TrimSpace(l)
					if l == "" {
						continue
					}
					if strings.HasPrefix(l, "~") && strings.HasSuffix(l, "~") {
						continue
					}
					if s, err := parseStmt(l); err == nil && s != nil {
						stmts = append(stmts, s)
					}
				}
			} else {
				for p.hasMore() {
					l2 := p.nextLine()
					t2 := strings.TrimSpace(l2)
					if t2 == "" {
						continue
					}
					if strings.HasPrefix(t2, "~") && strings.HasSuffix(t2, "~") {
						continue
					}
					if t2 == "]" {
						break
					}
					if s, err := parseStmt(t2); err == nil && s != nil {
						stmts = append(stmts, s)
					}
				}
			}
			obj.Reactions[trigger] = stmts
			continue
		}
	}
	return obj, nil
}

func parseGlobalStatement(line string) (Stmt, error) {
	lower := strings.ToLower(strings.TrimSpace(line))
	// print / say
	if strings.HasPrefix(lower, "print ") || strings.HasPrefix(lower, "say ") {
		idx := strings.Index(line, " ")
		return &PrintStmt{Expr: strings.TrimSpace(line[idx+1:])}, nil
	}
	// wait
	if strings.HasPrefix(lower, "wait ") {
		parts := strings.Fields(line)
		if len(parts) >= 2 {
			f, err := strconv.ParseFloat(parts[1], 64)
			if err == nil {
				return &WaitStmt{Seconds: f}, nil
			}
		}
	}
	// assignment: supports " is " and " = "
	if strings.Contains(lower, " is ") || stringsContainsOutsideQuotes(line, "=") {
		// handle "=" specially to avoid matching inside strings
		if stringsContainsOutsideQuotes(line, "=") && !strings.Contains(lower, " is ") {
			parts := strings.SplitN(line, "=", 2)
			lhs := strings.TrimSpace(parts[0])
			expr := strings.TrimSpace(parts[1])
			return &AssignStmt{Lhs: lhs, Expr: expr}, nil
		}
		parts := strings.SplitN(line, " is ", 2)
		if len(parts) == 2 {
			lhs := strings.TrimSpace(parts[0])
			expr := strings.TrimSpace(parts[1])
			return &AssignStmt{Lhs: lhs, Expr: expr}, nil
		}
	}
	return nil, errors.New("not a global stmt")
}

func stringsContainsOutsideQuotes(s string, sub string) bool {
	inQuote := false
	for i := 0; i+len(sub) <= len(s); i++ {
		c := s[i]
		if c == '"' {
			inQuote = !inQuote
		}
		if inQuote {
			continue
		}
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}

func parseStmt(line string) (Stmt, error) {
	low := strings.ToLower(strings.TrimSpace(line))
	if strings.HasPrefix(low, "print ") || strings.HasPrefix(low, "say ") {
		idx := strings.Index(line, " ")
		return &PrintStmt{Expr: strings.TrimSpace(line[idx+1:])}, nil
	}
	if strings.HasPrefix(low, "wait ") {
		parts := strings.Fields(line)
		if len(parts) >= 2 {
			f, err := strconv.ParseFloat(parts[1], 64)
			if err == nil {
				return &WaitStmt{Seconds: f}, nil
			}
		}
	}
	if strings.HasPrefix(low, "tell ") {
		rest := strings.TrimSpace(line[5:])
		idx := strings.Index(strings.ToLower(rest), " to ")
		if idx != -1 {
			target := strings.TrimSpace(rest[:idx])
			after := strings.TrimSpace(rest[idx+4:])
			action := after
			args := []string{}
			widx := strings.Index(strings.ToLower(after), " with ")
			if widx != -1 {
				action = strings.TrimSpace(after[:widx])
				argstr := strings.TrimSpace(after[widx+6:])
				rawargs := strings.Split(argstr, ",")
				for _, a := range rawargs {
					args = append(args, strings.TrimSpace(a))
				}
			}
			return &TellStmt{Target: target, Action: action, Args: args}, nil
		}
	}
	// assignment inside objects: accept " is " or "="
	if strings.Contains(low, " is ") || stringsContainsOutsideQuotes(line, "=") {
		if stringsContainsOutsideQuotes(line, "=") && !strings.Contains(low, " is ") {
			parts := strings.SplitN(line, "=", 2)
			return &AssignStmt{Lhs: strings.TrimSpace(parts[0]), Expr: strings.TrimSpace(parts[1])}, nil
		}
		parts := strings.SplitN(line, " is ", 2)
		return &AssignStmt{Lhs: strings.TrimSpace(parts[0]), Expr: strings.TrimSpace(parts[1])}, nil
	}
	return nil, errors.New("unknown statement")
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: heinterp <file.he>")
		return
	}
	srcPath := os.Args[1]
	b, err := os.ReadFile(srcPath)
	if err != nil {
		fmt.Println("Error reading:", err)
		return
	}
	parser := NewParser(string(b))
	prog, err := parser.ParseProgram()
	if err != nil {
		fmt.Println("Parse error:", err)
		return
	}
	rt := NewRuntime(prog)

	// Run global statements
	for _, s := range prog.Globals {
		if err := s.Exec(rt); err != nil {
			fmt.Println("Runtime error:", err)
			return
		}
	}

	// Simulate "on start" if any object has reaction named "start"
	for _, obj := range rt.Objects {
		if stmts, ok := obj.Reactions["start"]; ok {
			for _, s := range stmts {
				_ = s.Exec(rt)
			}
		}
	}

	fmt.Println("HE program loaded. Objects:", len(rt.Objects), "Assets:", len(rt.Assets))
	// For demo, trigger collision reaction if exists for Player or hero
	if obj, ok := rt.Objects["Player"]; ok {
		if ab, ok := obj.Abilities["jump"]; ok {
			fmt.Println("Running ability Player.jump() ...")
			for _, s := range ab.Statements {
				_ = s.Exec(rt)
			}
		}
		if r, ok := obj.Reactions["collision"]; ok {
			fmt.Println("Triggering collision reaction ...")
			for _, s := range r {
				_ = s.Exec(rt)
			}
		}
	}
}
