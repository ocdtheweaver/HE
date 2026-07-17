package compiler

import (
	"fmt"
	"strings"

	"hunterlang/lang/ast"
)

// ── Opcodes ───────────────────────────────────────────────────────────────────

type Opcode byte

const (
	OP_PUSH_NUM  Opcode = iota
	OP_PUSH_STR
	OP_PUSH_BOOL
	OP_PUSH_NIL
	OP_POP

	OP_LOAD
	OP_STORE

	OP_ADD
	OP_SUB
	OP_MUL
	OP_DIV
	OP_POW
	OP_NEG

	OP_EQ
	OP_NEQ
	OP_GT
	OP_LT
	OP_GTE
	OP_LTE
	OP_BETWEEN

	OP_AND
	OP_OR
	OP_NOT

	OP_CONCAT

	OP_JUMP
	OP_JUMP_IF_NOT
	OP_JUMP_IF

	OP_NEW_OBJECT
	OP_LOAD_FIELD
	OP_STORE_FIELD
	OP_CALL_METHOD
	OP_INIT_FIELD
	OP_DEF_METHOD
	OP_INHERIT
	OP_NEW_MODULE

	OP_CALL
	OP_RETURN
	OP_RETURN_VOID

	OP_SAY

	OP_REMEMBER
	OP_RECALL
	OP_FORGET

	OP_NEW_ARRAY
	OP_ARRAY_GET
	OP_ARRAY_LEN

	OP_CALL_VALUE_METHOD

	OP_TRY_START
	OP_TRY_END

	OP_CONTAINS // list, value → bool (true if list contains value)

	OP_FIELD_PAIRS // obj → array of [key, value] pairs (as 2-element arrays)

	OP_RUN_MODULE_CHUNK // [aliasIdx, subChunkIdx] → runs subchunk in isolated
	                     // VM, wraps its resulting vars as an object under alias

	OP_ASK // prompt on stack → prints prompt, reads a line from stdin, pushes it as text

	OP_SET_PROTECT_TAG // nameIdx → pops obj, sets obj.ProtectTag = Names[nameIdx], pushes obj back

	OP_NOP
	OP_HALT
)

var opNames = map[Opcode]string{
	OP_PUSH_NUM: "PUSH_NUM", OP_PUSH_STR: "PUSH_STR", OP_PUSH_BOOL: "PUSH_BOOL",
	OP_PUSH_NIL: "PUSH_NIL", OP_POP: "POP",
	OP_LOAD: "LOAD", OP_STORE: "STORE",
	OP_ADD: "ADD", OP_SUB: "SUB", OP_MUL: "MUL", OP_DIV: "DIV",
	OP_POW: "POW", OP_NEG: "NEG",
	OP_EQ: "EQ", OP_NEQ: "NEQ", OP_GT: "GT", OP_LT: "LT",
	OP_GTE: "GTE", OP_LTE: "LTE", OP_BETWEEN: "BETWEEN",
	OP_AND: "AND", OP_OR: "OR", OP_NOT: "NOT",
	OP_CONCAT: "CONCAT",
	OP_JUMP: "JUMP", OP_JUMP_IF_NOT: "JUMP_IF_NOT", OP_JUMP_IF: "JUMP_IF",
	OP_NEW_OBJECT: "NEW_OBJECT", OP_LOAD_FIELD: "LOAD_FIELD",
	OP_STORE_FIELD: "STORE_FIELD", OP_CALL_METHOD: "CALL_METHOD",
	OP_INIT_FIELD: "INIT_FIELD", OP_DEF_METHOD: "DEF_METHOD",
	OP_INHERIT: "INHERIT", OP_NEW_MODULE: "NEW_MODULE",
	OP_CALL: "CALL", OP_RETURN: "RETURN", OP_RETURN_VOID: "RETURN_VOID",
	OP_SAY: "SAY",
	OP_REMEMBER: "REMEMBER", OP_RECALL: "RECALL", OP_FORGET: "FORGET",
	OP_NEW_ARRAY: "NEW_ARRAY", OP_ARRAY_GET: "ARRAY_GET", OP_ARRAY_LEN: "ARRAY_LEN",
	OP_CALL_VALUE_METHOD: "CALL_VALUE_METHOD",
	OP_TRY_START: "TRY_START", OP_TRY_END: "TRY_END",
	OP_CONTAINS: "CONTAINS",
	OP_FIELD_PAIRS: "FIELD_PAIRS",
	OP_RUN_MODULE_CHUNK: "RUN_MODULE_CHUNK",
	OP_ASK: "ASK",
	OP_SET_PROTECT_TAG: "SET_PROTECT_TAG",
	OP_NOP: "NOP", OP_HALT: "HALT",
}

func (o Opcode) String() string {
	if s, ok := opNames[o]; ok {
		return s
	}
	return fmt.Sprintf("OP_%d", int(o))
}

// ── Instruction ───────────────────────────────────────────────────────────────

type Instruction struct {
	Op      Opcode
	Operand interface{}
}

func (ins Instruction) String() string {
	if ins.Operand != nil {
		return fmt.Sprintf("%-20s %v", ins.Op, ins.Operand)
	}
	return ins.Op.String()
}

// ── Chunk ─────────────────────────────────────────────────────────────────────

type Chunk struct {
	Name         string
	Instructions []Instruction
	Constants    []interface{}
	Names        []string
	SubChunks    []*MethodChunk
}

type MethodChunk struct {
	Name         string
	Params       []string
	Instructions []Instruction
	ProtectTag   string // "" = unprotected; "protected1", etc.
}

func NewChunk(name string) *Chunk { return &Chunk{Name: name} }

func (c *Chunk) Emit(op Opcode, operand interface{}) int {
	c.Instructions = append(c.Instructions, Instruction{Op: op, Operand: operand})
	return len(c.Instructions) - 1
}

func (c *Chunk) AddConst(v interface{}) int {
	c.Constants = append(c.Constants, v)
	return len(c.Constants) - 1
}

func (c *Chunk) AddName(name string) int {
	for i, n := range c.Names {
		if n == name {
			return i
		}
	}
	c.Names = append(c.Names, name)
	return len(c.Names) - 1
}

func (c *Chunk) AddSubChunk(mc *MethodChunk) int {
	c.SubChunks = append(c.SubChunks, mc)
	return len(c.SubChunks) - 1
}

func (c *Chunk) Patch(idx int, operand interface{}) {
	c.Instructions[idx].Operand = operand
}

func (c *Chunk) Disassemble() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("── Chunk: %s ─────────────────────────────\n", c.Name))
	sb.WriteString(fmt.Sprintf("  Constants : %v\n", c.Constants))
	sb.WriteString(fmt.Sprintf("  Names     : %v\n", c.Names))
	if len(c.SubChunks) > 0 {
		sb.WriteString(fmt.Sprintf("  Methods   : %d\n", len(c.SubChunks)))
		for i, mc := range c.SubChunks {
			sb.WriteString(fmt.Sprintf("    [%d] %s(%v)\n", i, mc.Name, mc.Params))
		}
	}
	sb.WriteString("\n")
	for i, ins := range c.Instructions {
		sb.WriteString(fmt.Sprintf("  %04d  %s\n", i, ins))
	}
	sb.WriteString("──────────────────────────────────────────\n")
	return sb.String()
}

// ── Compiler ──────────────────────────────────────────────────────────────────

type BytecodeCompiler struct {
	chunk   *Chunk
	errors  []string
	BaseDir string // directory used to resolve relative summon "file.he" paths
}

func NewBytecodeCompiler() *BytecodeCompiler {
	return &BytecodeCompiler{chunk: NewChunk("main")}
}

func (bc *BytecodeCompiler) Compile(prog *ast.Program) (*Chunk, []string) {
	for _, line := range prog.Lines {
		bc.compileLine(line)
	}
	bc.chunk.Emit(OP_HALT, nil)
	return bc.chunk, bc.errors
}

func (bc *BytecodeCompiler) compileLine(line ast.Line) {
	switch l := line.(type) {
	case *ast.ObjectLine:
		bc.compileObjectDef(l)
	case ast.GlobalStatementLine:
		bc.compileStmt(l.Statement)
	case *ast.GlobalStatementLine:
		bc.compileStmt(l.Statement)
	case ast.SummonLine:
		if l.Alias != nil {
			aliasIdx := bc.chunk.AddName(l.Alias.Name)
			modIdx := bc.chunk.AddName(l.ModuleName.Lexeme)
			bc.chunk.Emit(OP_NEW_MODULE, [2]int{aliasIdx, modIdx})
		}
	}
}

// ── Object compilation ────────────────────────────────────────────────────────

func (bc *BytecodeCompiler) compileObjectDef(o *ast.ObjectLine) {
	nameIdx := bc.chunk.AddName(o.Name.Name)
	bc.chunk.Emit(OP_NEW_OBJECT, nameIdx)
	if o.Like != nil {
		parentIdx := bc.chunk.AddName(o.Like.Name)
		bc.chunk.Emit(OP_INHERIT, parentIdx)
	}
	// Whole-object protection tag
	if o.ProtectTag != "" {
		tagIdx := bc.chunk.AddName(o.ProtectTag)
		bc.chunk.Emit(OP_SET_PROTECT_TAG, tagIdx)
	}
	for _, sec := range o.Body.Sections {
		switch s := sec.(type) {
		case ast.PropertiesSection:
			for _, prop := range s.Props {
				fieldIdx := bc.chunk.AddName(prop.Name)
				bc.compileExpr(prop.Value)
				bc.chunk.Emit(OP_INIT_FIELD, fieldIdx)
			}
		case ast.AbilitiesSection:
			for _, act := range s.Abilities {
				bc.compileMethodDef(act.Name, act.Params, act.Body, act.ProtectTag)
			}
		case ast.ReactionsSection:
			for _, rx := range s.Reactions {
				bc.compileMethodDef("__on_"+rx.Trigger.Left, nil, rx.Body, "")
			}
		}
	}
	bc.chunk.Emit(OP_STORE, nameIdx)
}

func (bc *BytecodeCompiler) compileMethodDef(name string, params []ast.Param, body []ast.Statement, protectTag string) {
	paramNames := make([]string, len(params))
	for i, p := range params {
		paramNames[i] = p.Name
	}
	sub := &BytecodeCompiler{chunk: &Chunk{
		Name:      name,
		Constants: bc.chunk.Constants,
		Names:     bc.chunk.Names,
		SubChunks: bc.chunk.SubChunks,
	}}
	for _, st := range body {
		sub.compileStmt(st)
	}
	sub.chunk.Emit(OP_RETURN_VOID, nil)
	bc.chunk.Constants = sub.chunk.Constants
	bc.chunk.Names = sub.chunk.Names
	bc.chunk.SubChunks = sub.chunk.SubChunks
	mc := &MethodChunk{
		Name:         name,
		Params:       paramNames,
		Instructions: sub.chunk.Instructions,
		ProtectTag:   protectTag,
	}
	chunkIdx := bc.chunk.AddSubChunk(mc)
	methodNameIdx := bc.chunk.AddName(name)
	bc.chunk.Emit(OP_DEF_METHOD, [2]int{methodNameIdx, chunkIdx})
}

func (bc *BytecodeCompiler) compileClosureLit(ex ast.ClosureExpr) {
	nameIdx := bc.chunk.AddName("__closure__")
	bc.chunk.Emit(OP_NEW_OBJECT, nameIdx)
	params := make([]ast.Param, len(ex.Params))
	for i, n := range ex.Params {
		params[i] = ast.Param{Name: n, Type: ast.TypeNode{Name: "any"}}
	}
	bc.compileMethodDef("__call__", params, ex.Body, "")
}

func (bc *BytecodeCompiler) compileMembership(ex ast.MembershipExpr) {
	// Works for both literal and dynamic lists: push value, push list,
	// let the VM do a proper equality scan via OP_CONTAINS.
	bc.compileExpr(ex.Value)
	bc.compileExpr(ex.List)
	bc.chunk.Emit(OP_CONTAINS, nil)
}

// compileFileModule compiles "summon path/to/lib.he [as alias]" by reading
// and compiling the target file into its own MethodChunk, then emitting
// OP_RUN_MODULE_CHUNK so the VM executes it in an isolated sub-VM, flat-merges
// its top-level names (and each exported object's abilities) into the
// importing scope, and — if an alias was given — additionally exposes
// everything as a qualified object too.
func (bc *BytecodeCompiler) compileFileModule(st ast.LoadModuleStmt) {
	path := st.FilePath
	data, err := fileModuleReader(path, bc.BaseDir)
	if err != nil {
		bc.errors = append(bc.errors, fmt.Sprintf("summon %q: %v", path, err))
		return
	}

	prog, err := parseModuleSource(data)
	if err != nil {
		bc.errors = append(bc.errors, fmt.Sprintf("summon %q: parse error: %v", path, err))
		return
	}

	chunkName := path
	if st.HasAlias {
		chunkName = st.Alias
	}
	sub := &BytecodeCompiler{
		chunk: &Chunk{
			Name:      chunkName,
			Constants: bc.chunk.Constants,
			Names:     bc.chunk.Names,
			SubChunks: bc.chunk.SubChunks,
		},
		BaseDir: dirOf(path, bc.BaseDir),
	}
	for _, line := range prog.Lines {
		sub.compileLine(line)
	}
	sub.chunk.Emit(OP_RETURN_VOID, nil)

	bc.chunk.Constants = sub.chunk.Constants
	bc.chunk.Names = sub.chunk.Names
	bc.chunk.SubChunks = sub.chunk.SubChunks
	bc.errors = append(bc.errors, sub.errors...)

	mc := &MethodChunk{Name: chunkName, Instructions: sub.chunk.Instructions}
	chunkIdx := bc.chunk.AddSubChunk(mc)
	pathIdx := bc.chunk.AddName(path)
	aliasIdx := -1
	if st.HasAlias {
		aliasIdx = bc.chunk.AddName(st.Alias)
	}
	bc.chunk.Emit(OP_RUN_MODULE_CHUNK, [3]int{pathIdx, aliasIdx, chunkIdx})
}

// ── Statement compilation ─────────────────────────────────────────────────────

func (bc *BytecodeCompiler) compileStmt(s ast.Statement) {
	switch st := s.(type) {

	case ast.SayStmt:
		bc.compileExpr(st.Expr)
		bc.chunk.Emit(OP_SAY, nil)

	case ast.ChangeStmt:
		bc.compileExpr(st.Expr)
		bc.chunk.Emit(OP_STORE, bc.chunk.AddName(st.Name))

	case ast.GrowStmt:
		idx := bc.chunk.AddName(st.Name)
		bc.chunk.Emit(OP_LOAD, idx)
		bc.compileExpr(st.Expr)
		bc.chunk.Emit(OP_ADD, nil)
		bc.chunk.Emit(OP_STORE, idx)

	case ast.ShrinkStmt:
		idx := bc.chunk.AddName(st.Name)
		bc.chunk.Emit(OP_LOAD, idx)
		bc.compileExpr(st.Expr)
		bc.chunk.Emit(OP_SUB, nil)
		bc.chunk.Emit(OP_STORE, idx)

	case ast.DotAssignStmt:
		objIdx := bc.chunk.AddName(st.Object)
		fieldIdx := bc.chunk.AddName(st.Field)
		bc.chunk.Emit(OP_LOAD, objIdx)
		bc.compileExpr(st.Expr)
		bc.chunk.Emit(OP_STORE_FIELD, fieldIdx)

	case ast.DecideStmt:
		bc.compileExpr(st.Cond)
		jmpFalse := bc.chunk.Emit(OP_JUMP_IF_NOT, 0)
		for _, b := range st.Then {
			bc.compileStmt(b)
		}
		if len(st.Else) > 0 {
			jmpEnd := bc.chunk.Emit(OP_JUMP, 0)
			bc.chunk.Patch(jmpFalse, len(bc.chunk.Instructions))
			for _, b := range st.Else {
				bc.compileStmt(b)
			}
			bc.chunk.Patch(jmpEnd, len(bc.chunk.Instructions))
		} else {
			bc.chunk.Patch(jmpFalse, len(bc.chunk.Instructions))
		}

	case ast.RepeatStmt:
		switch st.Kind {
		case "while":
			loopStart := len(bc.chunk.Instructions)
			bc.compileExpr(st.Cond)
			jmpOut := bc.chunk.Emit(OP_JUMP_IF_NOT, 0)
			for _, b := range st.Body {
				bc.compileStmt(b)
			}
			bc.chunk.Emit(OP_JUMP, loopStart)
			bc.chunk.Patch(jmpOut, len(bc.chunk.Instructions))
		case "times":
			cntIdx := bc.chunk.AddName("__cnt__")
			bc.compileExpr(st.Cond)
			bc.chunk.Emit(OP_STORE, cntIdx)
			iIdx := bc.chunk.AddName("__i__")
			bc.chunk.Emit(OP_PUSH_NUM, 0.0)
			bc.chunk.Emit(OP_STORE, iIdx)
			loopStart := len(bc.chunk.Instructions)
			bc.chunk.Emit(OP_LOAD, iIdx)
			bc.chunk.Emit(OP_LOAD, cntIdx)
			bc.chunk.Emit(OP_LT, nil)
			jmpOut := bc.chunk.Emit(OP_JUMP_IF_NOT, 0)
			for _, b := range st.Body {
				bc.compileStmt(b)
			}
			bc.chunk.Emit(OP_LOAD, iIdx)
			bc.chunk.Emit(OP_PUSH_NUM, 1.0)
			bc.chunk.Emit(OP_ADD, nil)
			bc.chunk.Emit(OP_STORE, iIdx)
			bc.chunk.Emit(OP_JUMP, loopStart)
			bc.chunk.Patch(jmpOut, len(bc.chunk.Instructions))
		}

	case ast.CountedRepeatStmt:
		varIdx := bc.chunk.AddName(st.CountVar)
		cntIdx := bc.chunk.AddName("__counted_n__")
		bc.compileExpr(st.Count)
		bc.chunk.Emit(OP_STORE, cntIdx)
		bc.chunk.Emit(OP_PUSH_NUM, 1.0)
		bc.chunk.Emit(OP_STORE, varIdx)
		loopStart := len(bc.chunk.Instructions)
		bc.chunk.Emit(OP_LOAD, varIdx)
		bc.chunk.Emit(OP_LOAD, cntIdx)
		bc.chunk.Emit(OP_LTE, nil)
		jmpOut := bc.chunk.Emit(OP_JUMP_IF_NOT, 0)
		for _, b := range st.Body {
			bc.compileStmt(b)
		}
		bc.chunk.Emit(OP_LOAD, varIdx)
		bc.chunk.Emit(OP_PUSH_NUM, 1.0)
		bc.chunk.Emit(OP_ADD, nil)
		bc.chunk.Emit(OP_STORE, varIdx)
		bc.chunk.Emit(OP_JUMP, loopStart)
		bc.chunk.Patch(jmpOut, len(bc.chunk.Instructions))

	case ast.RepeatUntilStmt:
		loopStart := len(bc.chunk.Instructions)
		bc.compileExpr(st.Cond)
		jmpOut := bc.chunk.Emit(OP_JUMP_IF, 0)
		for _, b := range st.Body {
			bc.compileStmt(b)
		}
		bc.chunk.Emit(OP_JUMP, loopStart)
		bc.chunk.Patch(jmpOut, len(bc.chunk.Instructions))

	case ast.RangeLoopStmt:
		varIdx := bc.chunk.AddName(st.VarName)
		bc.compileExpr(st.From)
		bc.chunk.Emit(OP_STORE, varIdx)
		loopStart := len(bc.chunk.Instructions)
		bc.chunk.Emit(OP_LOAD, varIdx)
		bc.compileExpr(st.To)
		bc.chunk.Emit(OP_LTE, nil)
		jmpOut := bc.chunk.Emit(OP_JUMP_IF_NOT, 0)
		for _, b := range st.Body {
			bc.compileStmt(b)
		}
		bc.chunk.Emit(OP_LOAD, varIdx)
		if st.Step != nil {
			bc.compileExpr(st.Step)
		} else {
			bc.chunk.Emit(OP_PUSH_NUM, 1.0)
		}
		bc.chunk.Emit(OP_ADD, nil)
		bc.chunk.Emit(OP_STORE, varIdx)
		bc.chunk.Emit(OP_JUMP, loopStart)
		bc.chunk.Patch(jmpOut, len(bc.chunk.Instructions))

	case ast.ForEachStmt:
		listIdx := bc.chunk.AddName("__feach_list__")
		idxIdx := bc.chunk.AddName("__feach_idx__")
		varIdx := bc.chunk.AddName(st.VarName)
		bc.compileExpr(st.List)
		bc.chunk.Emit(OP_STORE, listIdx)
		bc.chunk.Emit(OP_PUSH_NUM, 0.0)
		bc.chunk.Emit(OP_STORE, idxIdx)
		loopStart := len(bc.chunk.Instructions)
		bc.chunk.Emit(OP_LOAD, idxIdx)
		bc.chunk.Emit(OP_LOAD, listIdx)
		bc.chunk.Emit(OP_ARRAY_LEN, nil)
		bc.chunk.Emit(OP_LT, nil)
		jmpOut := bc.chunk.Emit(OP_JUMP_IF_NOT, 0)
		bc.chunk.Emit(OP_LOAD, listIdx)
		bc.chunk.Emit(OP_LOAD, idxIdx)
		bc.chunk.Emit(OP_ARRAY_GET, nil)
		bc.chunk.Emit(OP_STORE, varIdx)
		for _, b := range st.Body {
			bc.compileStmt(b)
		}
		bc.chunk.Emit(OP_LOAD, idxIdx)
		bc.chunk.Emit(OP_PUSH_NUM, 1.0)
		bc.chunk.Emit(OP_ADD, nil)
		bc.chunk.Emit(OP_STORE, idxIdx)
		bc.chunk.Emit(OP_JUMP, loopStart)
		bc.chunk.Patch(jmpOut, len(bc.chunk.Instructions))

	case ast.TryStmt:
		errVarIdx := bc.chunk.AddName("error")
		tryStart := bc.chunk.Emit(OP_TRY_START, [2]int{0, errVarIdx})
		for _, b := range st.Body {
			bc.compileStmt(b)
		}
		bc.chunk.Emit(OP_TRY_END, nil)
		jmpOver := bc.chunk.Emit(OP_JUMP, 0)
		handlerStart := len(bc.chunk.Instructions)
		bc.chunk.Patch(tryStart, [2]int{handlerStart, errVarIdx})
		for _, b := range st.Handler {
			bc.compileStmt(b)
		}
		bc.chunk.Patch(jmpOver, len(bc.chunk.Instructions))

	case ast.TryWithVarStmt:
		errVarIdx := bc.chunk.AddName(st.ErrVar)
		tryStart := bc.chunk.Emit(OP_TRY_START, [2]int{0, errVarIdx})
		for _, b := range st.Body {
			bc.compileStmt(b)
		}
		bc.chunk.Emit(OP_TRY_END, nil)
		jmpOver := bc.chunk.Emit(OP_JUMP, 0)
		handlerStart := len(bc.chunk.Instructions)
		bc.chunk.Patch(tryStart, [2]int{handlerStart, errVarIdx})
		for _, b := range st.Handler {
			bc.compileStmt(b)
		}
		bc.chunk.Patch(jmpOver, len(bc.chunk.Instructions))

	case ast.MultiAssignStmt:
		if len(st.Exprs) == 1 && len(st.Names) > 1 {
			tmpIdx := bc.chunk.AddName("__multi_tmp__")
			bc.compileExpr(st.Exprs[0])
			bc.chunk.Emit(OP_STORE, tmpIdx)
			for i, name := range st.Names {
				nameIdx := bc.chunk.AddName(name)
				bc.chunk.Emit(OP_LOAD, tmpIdx)
				bc.chunk.Emit(OP_PUSH_NUM, float64(i))
				bc.chunk.Emit(OP_ARRAY_GET, nil)
				bc.chunk.Emit(OP_STORE, nameIdx)
			}
		} else {
			for _, ex := range st.Exprs {
				bc.compileExpr(ex)
			}
			for i := len(st.Names) - 1; i >= 0; i-- {
				bc.chunk.Emit(OP_STORE, bc.chunk.AddName(st.Names[i]))
			}
		}

	case ast.MultiReturnStmt:
		for _, ex := range st.Exprs {
			bc.compileExpr(ex)
		}
		bc.chunk.Emit(OP_NEW_ARRAY, len(st.Exprs))
		bc.chunk.Emit(OP_RETURN, nil)

	case ast.WithScopeStmt:
		aliasIdx := bc.chunk.AddName(st.Alias)
		bc.compileExpr(st.Expr)
		bc.chunk.Emit(OP_STORE, aliasIdx)
		for _, b := range st.Body {
			bc.compileStmt(b)
		}

	case ast.ReturnStmt:
		bc.compileExpr(st.Expr)
		bc.chunk.Emit(OP_RETURN, nil)

	case ast.WaitStmt:
		bc.chunk.Emit(OP_NOP, nil)

	case ast.RememberStmt:
		bc.chunk.Emit(OP_LOAD, bc.chunk.AddName(st.Name))
		bc.chunk.Emit(OP_REMEMBER, bc.chunk.AddName(st.Key))

	case ast.RecallStmt:
		bc.chunk.Emit(OP_RECALL, bc.chunk.AddName(st.Key))
		bc.chunk.Emit(OP_STORE, bc.chunk.AddName(st.Name))

	case ast.ForgetStmt:
		bc.chunk.Emit(OP_FORGET, bc.chunk.AddName(st.Name))

	case ast.ExprStmt:
		if sl, ok := st.Expr.(ast.StringLit); ok && len(sl.Value) > 10 && sl.Value[:10] == "__summon__" {
			parts := strings.SplitN(sl.Value, ":", 3)
			if len(parts) == 3 {
				aliasIdx := bc.chunk.AddName(parts[2])
				modIdx := bc.chunk.AddName(parts[1])
				bc.chunk.Emit(OP_NEW_MODULE, [2]int{aliasIdx, modIdx})
				return
			}
		}
		bc.compileExpr(st.Expr)
		bc.chunk.Emit(OP_POP, nil)

	case ast.CallStmt:
		objIdx := bc.chunk.AddName(st.Object)
		methodIdx := bc.chunk.AddName(st.Action)
		bc.chunk.Emit(OP_LOAD, objIdx)
		for _, arg := range st.Args {
			bc.compileExpr(arg)
		}
		bc.chunk.Emit(OP_CALL_VALUE_METHOD, [2]int{methodIdx, len(st.Args)})
		bc.chunk.Emit(OP_POP, nil) // discard return value of tell/call

	case ast.ForEachFieldStmt:
		pairsIdx := bc.chunk.AddName("__field_pairs__")
		idxIdx := bc.chunk.AddName("__field_idx__")
		pairIdx := bc.chunk.AddName("__field_pair__")
		keyIdx := bc.chunk.AddName(st.KeyVar)

		bc.compileExpr(st.Object)
		bc.chunk.Emit(OP_FIELD_PAIRS, nil)
		bc.chunk.Emit(OP_STORE, pairsIdx)
		bc.chunk.Emit(OP_PUSH_NUM, 0.0)
		bc.chunk.Emit(OP_STORE, idxIdx)

		loopStart := len(bc.chunk.Instructions)
		bc.chunk.Emit(OP_LOAD, idxIdx)
		bc.chunk.Emit(OP_LOAD, pairsIdx)
		bc.chunk.Emit(OP_ARRAY_LEN, nil)
		bc.chunk.Emit(OP_LT, nil)
		jmpOut := bc.chunk.Emit(OP_JUMP_IF_NOT, 0)

		// pair = pairs[idx]; key = pair[0]
		bc.chunk.Emit(OP_LOAD, pairsIdx)
		bc.chunk.Emit(OP_LOAD, idxIdx)
		bc.chunk.Emit(OP_ARRAY_GET, nil)
		bc.chunk.Emit(OP_STORE, pairIdx)
		bc.chunk.Emit(OP_LOAD, pairIdx)
		bc.chunk.Emit(OP_PUSH_NUM, 0.0)
		bc.chunk.Emit(OP_ARRAY_GET, nil)
		bc.chunk.Emit(OP_STORE, keyIdx)

		if st.ValVar != "" {
			valIdx := bc.chunk.AddName(st.ValVar)
			bc.chunk.Emit(OP_LOAD, pairIdx)
			bc.chunk.Emit(OP_PUSH_NUM, 1.0)
			bc.chunk.Emit(OP_ARRAY_GET, nil)
			bc.chunk.Emit(OP_STORE, valIdx)
		}

		for _, b := range st.Body {
			bc.compileStmt(b)
		}

		bc.chunk.Emit(OP_LOAD, idxIdx)
		bc.chunk.Emit(OP_PUSH_NUM, 1.0)
		bc.chunk.Emit(OP_ADD, nil)
		bc.chunk.Emit(OP_STORE, idxIdx)
		bc.chunk.Emit(OP_JUMP, loopStart)
		bc.chunk.Patch(jmpOut, len(bc.chunk.Instructions))

	case ast.AskStmt:
		bc.compileExpr(st.Prompt)
		bc.chunk.Emit(OP_ASK, nil)
		bc.chunk.Emit(OP_STORE, bc.chunk.AddName(st.VarName))

	case ast.LoadModuleStmt:
		bc.compileFileModule(st)

	default:
		bc.errors = append(bc.errors, fmt.Sprintf("uncompilable statement: %T", s))
	}
}

// ── Expression compilation ────────────────────────────────────────────────────

func (bc *BytecodeCompiler) compileExpr(e ast.Expression) {
	switch ex := e.(type) {

	case ast.NumberLit:
		bc.chunk.Emit(OP_PUSH_NUM, ex.Value)

	case ast.StringLit:
		bc.chunk.Emit(OP_PUSH_STR, ex.Value)

	case ast.BooleanLit:
		bc.chunk.Emit(OP_PUSH_BOOL, ex.Value)

	case ast.InterpStringExpr:
		for i, seg := range ex.Segments {
			bc.compileExpr(seg)
			if i > 0 {
				bc.chunk.Emit(OP_CONCAT, nil)
			}
		}

	case ast.IdentifierExpr:
		bc.chunk.Emit(OP_LOAD, bc.chunk.AddName(ex.Name))

	case ast.FieldAccessExpr:
		bc.chunk.Emit(OP_LOAD, bc.chunk.AddName(ex.Receiver.Name))
		bc.chunk.Emit(OP_LOAD_FIELD, bc.chunk.AddName(ex.Field))

	case ast.ArrayLit:
		for _, el := range ex.Elems {
			bc.compileExpr(el)
		}
		bc.chunk.Emit(OP_NEW_ARRAY, len(ex.Elems))

	case ast.NamedArgLit:
		nameIdx := bc.chunk.AddName("__namedargs__")
		bc.chunk.Emit(OP_NEW_OBJECT, nameIdx)
		for _, pair := range ex.Pairs {
			fieldIdx := bc.chunk.AddName(pair.Key)
			bc.compileExpr(pair.Value)
			bc.chunk.Emit(OP_INIT_FIELD, fieldIdx)
		}

	case ast.BinaryExpr:
		bc.compileExpr(ex.Left)
		bc.compileExpr(ex.Right)
		switch ex.Op {
		case "+":
			bc.chunk.Emit(OP_ADD, nil)
		case "-":
			bc.chunk.Emit(OP_SUB, nil)
		case "*":
			bc.chunk.Emit(OP_MUL, nil)
		case "/":
			bc.chunk.Emit(OP_DIV, nil)
		}

	case ast.PowerExpr:
		bc.compileExpr(ex.Left)
		bc.compileExpr(ex.Right)
		bc.chunk.Emit(OP_POW, nil)

	case ast.UnaryExpr:
		bc.compileExpr(ex.X)
		if ex.Op == "-" {
			bc.chunk.Emit(OP_NEG, nil)
		} else {
			bc.chunk.Emit(OP_NOT, nil)
		}

	case ast.CompareExpr:
		bc.compileExpr(ex.Left)
		bc.compileExpr(ex.Right)
		switch ex.Op {
		case "==":
			bc.chunk.Emit(OP_EQ, nil)
		case "!=":
			bc.chunk.Emit(OP_NEQ, nil)
		case ">":
			bc.chunk.Emit(OP_GT, nil)
		case "<":
			bc.chunk.Emit(OP_LT, nil)
		case ">=":
			bc.chunk.Emit(OP_GTE, nil)
		case "<=":
			bc.chunk.Emit(OP_LTE, nil)
		}

	case ast.BetweenExpr:
		bc.compileExpr(ex.Value)
		bc.compileExpr(ex.Low)
		bc.compileExpr(ex.High)
		bc.chunk.Emit(OP_BETWEEN, nil)

	case ast.MembershipExpr:
		bc.compileMembership(ex)

	case ast.ClosureExpr:
		bc.compileClosureLit(ex)

	case ast.LogicAndExpr:
		bc.compileExpr(ex.Left)
		bc.compileExpr(ex.Right)
		bc.chunk.Emit(OP_AND, nil)

	case ast.LogicOrExpr:
		bc.compileExpr(ex.Left)
		bc.compileExpr(ex.Right)
		bc.chunk.Emit(OP_OR, nil)

	case ast.ParenExpr:
		bc.compileExpr(ex.X)

	case ast.MethodCallExpr:
		bc.chunk.Emit(OP_LOAD, bc.chunk.AddName(ex.Receiver.Name))
		for _, arg := range ex.Args {
			bc.compileExpr(arg)
		}
		bc.chunk.Emit(OP_CALL_VALUE_METHOD, [2]int{bc.chunk.AddName(ex.Method), len(ex.Args)})

	case ast.MethodChainExpr:
		bc.compileExpr(ex.Recv)
		for _, arg := range ex.Args {
			bc.compileExpr(arg)
		}
		bc.chunk.Emit(OP_CALL_VALUE_METHOD, [2]int{bc.chunk.AddName(ex.Method), len(ex.Args)})

	case ast.CallExpr:
		calleeIdx := bc.chunk.AddName(ex.Callee)
		for _, arg := range ex.Args {
			bc.compileExpr(arg)
		}
		bc.chunk.Emit(OP_CALL, [2]int{calleeIdx, len(ex.Args)})

	default:
		bc.chunk.Emit(OP_PUSH_NIL, nil)
		bc.errors = append(bc.errors, fmt.Sprintf("uncompilable expression: %T", e))
	}
}

// ── File module compile-time resolution ────────────────────────────────────────

func fileModuleReader(path, baseDir string) (string, error) {
	data, err := osReadFile(path)
	if err == nil {
		return data, nil
	}
	if baseDir != "" && !filepathIsAbs(path) {
		alt := filepathJoin(baseDir, path)
		if altData, altErr := osReadFile(alt); altErr == nil {
			return altData, nil
		}
	}
	return "", err
}

func dirOf(path, baseDir string) string {
	resolved := path
	if baseDir != "" && !filepathIsAbs(path) {
		if _, err := osStat(path); err != nil {
			resolved = filepathJoin(baseDir, path)
		}
	}
	return filepathDir(resolved)
}
