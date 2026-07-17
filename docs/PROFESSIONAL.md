# HE Language — Professional Reference
*Architecture, internals, runtime semantics, and extension guide*
*Current: Pass 5*

---

## Table of Contents

1. [Architecture Overview](#1-architecture-overview)
2. [Lexer](#2-lexer)
3. [Parser](#3-parser)
4. [AST Reference](#4-ast-reference)
5. [Type System](#5-type-system)
6. [Runtime Semantics](#6-runtime-semantics)
7. [Standard Library Internals](#7-standard-library-internals)
8. [Compiler Pipeline](#8-compiler-pipeline)
9. [Module System](#9-module-system)
10. [Persistence Layer](#10-persistence-layer)
11. [Extension Guide](#11-extension-guide)
12. [Error Handling Model](#12-error-handling-model)
13. [Known Limitations](#13-known-limitations)
14. [Roadmap](#14-roadmap)

---

## 1. Architecture Overview

```
Source (.he)
    │
    ▼
┌─────────┐
│  Lexer  │  lang/lexer/lexer.go
│         │  Tokenises source; handles ~ comments ~, string interpolation
└────┬────┘
     │ []token.Token
     ▼
┌─────────┐
│ Parser  │  lang/parser/parser.go
│         │  Recursive-descent; produces AST
└────┬────┘
     │ *ast.Program
     ├────────────────────────────────────────────────────────┐
     ▼                                                        ▼
┌──────────────┐                                    ┌────────────────┐
│  Interpreter │  lang/eval/runtime.go              │   Compiler     │
│              │  Tree-walk evaluator               │  (preview)     │
│  + stdlib    │  lang/eval/stdlib.go               │                │
│  + persist   │  lang/eval/persist.go              │ compiler/      │
└──────────────┘                                    │ types.go       │
     │                                              │ bytecode.go    │
     ▼                                              └────────────────┘
   stdout / side effects                            bytecode IR / type report
```

### Package layout

```
hunterlang/
├── cmd/
│   ├── hunterlang/main.go   CLI runner: -file, -types, -compile
│   └── repl/main.go         Interactive REPL
├── lang/
│   ├── token/token.go       Token type definitions
│   ├── lexer/lexer.go       Tokeniser
│   ├── ast/ast.go           AST node types
│   ├── parser/parser.go     Parser
│   ├── types/types.go       Runtime value and object types
│   └── eval/
│       ├── eval.go          Interpreter entry point
│       ├── runtime.go       Statement/expression executor
│       ├── stdlib.go        Built-in module implementations
│       └── persist.go       remember/recall/forget store
├── compiler/
│   ├── types.go             Type inferencer
│   └── bytecode.go          Bytecode IR and compiler
└── docs/
    ├── HE_FOR_DUMMIES.md    Beginner guide
    ├── PROFESSIONAL.md      This document
    └── WOLFHEAD_INTEGRATION.md  OS binding guide
```

---


## 1.1 CLI Reference

The `he` binary uses a subcommand structure:

```
he run <file.he> [--event <name>]   Run a program
he check <file.he>                  Type inference report
he build <file.he>                  Bytecode disassembly
he vm <file.he>                     Run via bytecode VM
he bench <file.he>                  Interpreter vs VM benchmark
he repl                             Interactive shell
he new <name>                       Scaffold a new project
he version                          Version info
he help [command]                   Help
```

**Implicit run:** `he myprogram.he` is shorthand for `he run myprogram.he` — if the first argument ends in `.he` and isn't a known subcommand, it's treated as a file to run.

### Project scaffolding (`he new`)

```
he new myproject
```

Generates:
```
myproject/
├── main.he          — entry point, imports lib/helpers.he
├── lib/
│   └── helpers.he    — shared utilities, exports appName/version/Utils
└── README.md
```

### Module path resolution

`summon "lib/helpers.he" as helpers` resolves relative to the **importing file's directory**, not the working directory. This is tracked via `Runtime.baseDir`, set by `Interpreter.SetBaseDir()` before `Run()`. Nested file modules resolve relative to their own directory (`sub.baseDir = filepath.Dir(resolvedPath)`).

### Package layout addition

```
hunterlang/
├── repl/
│   └── repl.go       Shared REPL implementation (used by `he repl` and cmd/repl)
```


## 2. Lexer

**File:** `lang/lexer/lexer.go`

Single-pass, single-lookahead tokeniser. No regex — pure byte scanning.

### Comment syntax

```
~ This is a comment. It ends at the next tilde. ~
```

Comments can span multiple lines. Opening and closing `~` are consumed and discarded.

### String interpolation

When a string contains `{...}`, the lexer emits `INTERP_STRING` instead of `STRING`. The raw content (including braces) is preserved in the lexeme for the parser to split.

```
"Hello, {name}!"  →  INTERP_STRING: "Hello, {name}!"
```

The parser's `parseInterpString()` splits this into `InterpStringExpr` segments by creating a sub-lexer+parser for each `{...}` region. This means full expressions are valid inside `{}`:

```
"Result: {a + b}"
"Name: {person.name}"
"Double: {double(x)}"
```

### Keyword resolution

All keywords are resolved at lex time via the `keywords` map. Case-insensitive matching: `SET`, `Set`, `set` all produce `K_SET`.

Boolean literals `true`/`yes` → `BOOLEAN("true")`, `false`/`no` → `BOOLEAN("false")`.

---

## 3. Parser

**File:** `lang/parser/parser.go`

Recursive-descent, LL(1) with one token of lookahead (via `Lexer.Peek()`).

### Entry points

```go
p := parser.New(lx)
prog, err := p.ParseProgram()   // parses top-level lines
```

### Precedence hierarchy (lowest to highest)

```
LogicOr      →  "or"
LogicAnd     →  "and"
Comparison   →  == != > < >= <= "is" "is not" "is between"
Arithmetic   →  + -
Term         →  * /
Factor       →  ** (right-associative)
Unary        →  - ! "not"
Primary      →  literal, ident, (expr), [array/named], call, method-call, ability
```

### Named argument blocks

Inside `[...]`, if the first token pair is `IDENT (IS|COLON)`, the block is parsed as `NamedArgLit` rather than an array:

```
[title is "My App", children: [...]]  →  NamedArgLit
[1, 2, 3]                             →  ArrayLit
```

### Inline summon

When `summon` appears inside a method body (not at top level), it is parsed and encoded as a sentinel `StringLit`:

```
__summon__:module:alias
```

The runtime detects this in `ExprStmt` and registers the module correctly.

---

## 4. AST Reference

**File:** `lang/ast/ast.go`

All statement types implement `stmtNode()`. All expression types implement `exprNode()`.

### Statement nodes

| Node | Syntax |
|------|--------|
| `SayStmt` | `say expr` |
| `ChangeStmt` | `set/let/change name to expr` |
| `GrowStmt` | `grow name by expr` |
| `ShrinkStmt` | `shrink name by expr` |
| `DotAssignStmt` | `set obj.field to expr` |
| `GiveStmt` → `DotAssignStmt` | `give expr to obj.field` |
| `DecideStmt` | `if cond then [...] else [...]` |
| `RepeatStmt` | `repeat while/times cond [...]` |
| `RepeatUntilStmt` | `repeat until cond [...]` |
| `ForEachStmt` | `for each x in list [...]` |
| `RangeLoopStmt` | `for each/each i from A to B [...]` |
| `TryStmt` | `try [...] or [...]` |
| `TryWithVarStmt` | `try [...] or (e) [...]` |
| `CallStmt` | `tell obj to action [with args]` |
| `WaitStmt` | `wait N seconds/frames` |
| `ReturnStmt` | `return expr` |
| `MultiReturnStmt` | `return a, b, c` |
| `MultiAssignStmt` | `set a, b to expr1, expr2` |
| `AskStmt` | `ask "..." as var` |
| `RememberStmt` | `remember name [as "key"]` |
| `RecallStmt` | `recall name [as alias]` |
| `ForgetStmt` | `forget name` |
| `LoadModuleStmt` | `summon "file.he" as alias` |
| `ExprStmt` | expression used as statement |

### Expression nodes

| Node | Example |
|------|---------|
| `NumberLit` | `42`, `3.14` |
| `StringLit` | `"hello"` |
| `BooleanLit` | `true`, `yes` |
| `InterpStringExpr` | `"Hi {name}!"` |
| `IdentifierExpr` | `score` |
| `ArrayLit` | `[1, 2, 3]` |
| `NamedArgLit` | `[key is val, ...]` |
| `BinaryExpr` | `a + b`, `a * b` |
| `UnaryExpr` | `-x`, `not x` |
| `CompareExpr` | `x == 5`, `x > y` |
| `BetweenExpr` | `x is between 1 and 10` |
| `LogicAndExpr` | `a and b` |
| `LogicOrExpr` | `a or b` |
| `PowerExpr` | `2 ** 8` |
| `ParenExpr` | `(a + b)` |
| `FieldAccessExpr` | `obj.field` |
| `MethodCallExpr` | `obj.method(args)` |
| `CallExpr` | `name(args)` |
| `AbilityLit` | `ability(x) [...]` |

---


### Pass 6 statement nodes

| Node | Syntax |
|------|--------|
| `WithScopeStmt` | `with expr as alias [...]` |
| `ForEachFieldStmt` | `for each key, val in obj [...]` |
| `CountedRepeatStmt` | `repeat N times as i [...]` |

### Pass 6 expression nodes

| Node | Example |
|------|---------|
| `MembershipExpr` | `x is one of [a, b, c]` |
| `MethodChainExpr` | `word.upper()`, `items.reverse().first()` |
| `ClosureExpr` | `ability(x) [...]` with captured scope |

### Pass 7 additions

| Node | Syntax | Notes |
|------|--------|-------|
| `CountedRepeatStmt` | `repeat N times as i [...]` | Counter is 1-based |
| `ChangeStmt` (becomes) | `name becomes expr` | Sugar for set |


## 5. Type System

**File:** `lang/types/types.go`, `compiler/types.go`

### Runtime types

HE is dynamically typed at the interpreter level. All values are `types.Value`:

```go
type Value struct {
    Type    ValueType
    Number  float64
    Str     string
    Boolean bool
    Object  *Object
    Array   []Value
}
```

| `ValueType` | HE name   | Go representation |
|-------------|-----------|-------------------|
| `NumberT`   | `number`  | `float64`         |
| `StringT`   | `text`    | `string`          |
| `BooleanT`  | `boolean` | `bool`            |
| `ObjectT`   | `object`  | `*Object`         |
| `ArrayT`    | `list`    | `[]Value`         |
| `NilT`      | `nothing` | (zero value)      |

### Static inference (compiler)

The `compiler.Inferencer` walks the AST and infers types into a `TypeEnv`. This is used for `-types` reporting and will drive the bytecode compiler's type-annotated codegen.

**Inference rules:**
- `NumberLit`, `GrowStmt`, `ShrinkStmt`, `RangeLoopStmt` var → `TypeNumber`
- `StringLit`, `InterpStringExpr`, `AskStmt` result, `text.*` → `TypeText`
- `BooleanLit`, comparison expressions, logic expressions, `BetweenExpr` → `TypeBoolean`
- `ArrayLit`, `list.*` returns → `TypeList`
- Object definitions, `SummonLine`, `MethodCallExpr` on known types → `TypeObject`
- `AbilityLit` → `TypeAbility`
- Unknown: when inference can't determine the type (propagates through multi-return etc.)

---

## 6. Runtime Semantics

**File:** `lang/eval/runtime.go`

### Execution model

Tree-walk interpreter. Each `execStatement` call dispatches on the concrete statement type. Expressions are evaluated recursively via `evalExpr`.

### Environment

The `Runtime` struct holds:
- `Env map[string]types.Value` — flat variable store
- `Objects map[string]*types.Object` — object registry (subset of Env)
- `currentRecv *types.Object` — active receiver when inside an action/reaction body
- `persist *persistStore` — persistence backend

Variable lookup (`getVar`) checks `currentRecv.Fields` first (enabling `self` access without a keyword), then falls back to `Env`.

### Return propagation

`return` is implemented via a sentinel error type:

```go
type returnSignal struct{ value types.Value }
func (r *returnSignal) Error() string { return "return" }
```

`execAction` catches `returnSignal` and returns its value. All other callers propagate it. This avoids needing a separate control-flow mechanism.

### Object method dispatch

`callMethod(recv, method, argExprs)` evaluates args, then:
1. Checks `recv.Builtins[method]` — native Go functions
2. Checks `recv.Actions[method]` — user-defined abilities

For anonymous abilities stored as `ObjectT` values, `CallExpr` dispatch checks for `Actions["__call__"]`.

### Multi-return unpacking

When `MultiAssignStmt` has a single expr that evaluates to `ArrayT` with the same length as the name list, the array elements are distributed:

```
set a, b to Calculator.divmod(10, 3)
~ divmod returns [3, 1] as ArrayT ~
~ a = 3, b = 1 ~
```

---

## 7. Standard Library Internals

**File:** `lang/eval/stdlib.go`

All stdlib modules are registered via `registerModule(obj, moduleName)`. Each case wires `BuiltinFn` functions onto the module object using `registerBuiltin`.

### Adding a new module

```go
case "mymodule":
    registerBuiltin(obj, "myFunc", func(args []types.Value) (types.Value, error) {
        if len(args) != 1 {
            return types.Nil(), fmt.Errorf("mymodule.myFunc expects 1 argument")
        }
        result := doSomething(args[0].Str)
        return types.FromString(result), nil
    })
```

Then users can:
```
summon "mymodule" as m
say m.myFunc("hello")
```

### BuiltinFn signature

```go
type BuiltinFn func(args []types.Value) (types.Value, error)
```

**Conventions:**
- Always validate `len(args)` and `args[i].Type` before use
- Return `types.Nil(), fmt.Errorf(...)` for errors — these bubble to `try/or`
- Use `types.FromX()` constructors, never construct `Value` directly

---

## 8. Compiler Pipeline

**Files:** `compiler/types.go`, `compiler/bytecode.go`

### Type Inference pass

```go
inf := compiler.NewInferencer()
typeMap, errs := inf.InferProgram(prog)
fmt.Print(inf.Report())
```

### Bytecode IR

```go
bc := compiler.NewBytecodeCompiler()
chunk, errs := bc.Compile(prog)
fmt.Print(chunk.Disassemble())
```

### Opcode set (40 opcodes)

| Category | Opcodes |
|----------|---------|
| Stack | `PUSH_NUM`, `PUSH_STR`, `PUSH_BOOL`, `PUSH_NIL`, `POP` |
| Variables | `LOAD`, `STORE` |
| Arithmetic | `ADD`, `SUB`, `MUL`, `DIV`, `POW`, `NEG` |
| Comparison | `EQ`, `NEQ`, `GT`, `LT`, `GTE`, `LTE`, `BETWEEN` |
| Logic | `AND`, `OR`, `NOT` |
| Strings | `CONCAT` |
| Control | `JUMP`, `JUMP_IF_NOT`, `JUMP_IF` |
| Objects | `NEW_OBJECT`, `LOAD_FIELD`, `STORE_FIELD`, `CALL_METHOD` |
| Functions | `CALL`, `RETURN` |
| Arrays | `NEW_ARRAY`, `ARRAY_GET`, `ARRAY_LEN` |
| I/O | `SAY` |
| Persistence | `REMEMBER`, `RECALL`, `FORGET` |
| Control | `NOP`, `HALT` |

### Forward jump patching

```go
jmpOut := bc.chunk.Emit(OP_JUMP_IF_NOT, 0)  // placeholder
// ... compile body ...
bc.chunk.Patch(jmpOut, len(bc.chunk.Instructions))  // fill in target
```

### Planned: bytecode VM

The `Chunk` format is designed for a register-less stack VM. The VM will:
1. Load the `Chunk`
2. Maintain a value stack and instruction pointer
3. Dispatch each opcode via switch
4. Use the `Names` pool for variable lookup

---


## 8.1 Bytecode VM (Pass 8)

As of Pass 8, the bytecode VM (`compiler/vm.go`) is feature-complete for everyday HE programs — objects, methods, inheritance, closures, try/or, and all control flow compile and run correctly.

### VMObject

```go
type VMObject struct {
    Name    string
    Fields  map[string]VMValue
    Methods map[string]*MethodChunk
    Builtin map[string]func(args []VMValue) (VMValue, error)
}
```

Mirrors `types.Object` from the interpreter but operates on `VMValue` instead of `types.Value`.

### Method compilation

Each ability/reaction compiles to its own `MethodChunk` — an independent instruction stream sharing the parent chunk's `Constants`/`Names` pools (so name indices stay valid across method boundaries):

```go
type MethodChunk struct {
    Name         string
    Params       []string
    Instructions []Instruction
}
```

`OP_DEF_METHOD` attaches a compiled `MethodChunk` to an object on the stack. `OP_CALL_METHOD` / `OP_CALL_VALUE_METHOD` invoke it via `vm.invokeMethod`, which runs a **sub-VM** seeded with:
1. All of the parent VM's variables (so global objects remain callable from within their own methods — enables recursion: `Calculator.factorial(n-1)`)
2. The receiver's current field values
3. Bound parameters (backed up and restored per-call to keep recursion correct)

### Virtual dispatch

`OP_CALL_VALUE_METHOD` (used for all `.method()` calls, including `tell X to Y`) routes by the receiver's runtime tag — object, string, number, or list — exactly mirroring the interpreter's `callStringMethod`/`callListMethod`/`callNumberMethod`. This is what makes `"hello".upper()` and `Counter.increment` both work through the same opcode.

### Try/or error recovery

```go
type tryHandler struct {
    target    int // instruction index of the handler block
    errVarIdx int // name index for the error variable
}
```

`OP_TRY_START` pushes a handler frame; `OP_TRY_END` pops it on success. When `vm.exec()` returns an error, `Run()` and `runMethodBody()` both call `vm.tryRecover()`, which pops the nearest handler, binds the error message, clears the stack, and jumps to the handler block.

### .hbc binary format

`compiler/serial.go` implements a length-prefixed, little-endian binary encoding:

```
magic   "HE"
version 0x01
chunk:
  name       string
  constants  count + tagged values (n/s/b/0)
  names      count + strings
  subchunks  count + MethodChunks (name, params, instructions)
  instrs     count + Instructions (opcode byte + tagged operand)
```

Operand tags: `0` nil, `i` int32, `n` float64, `b` bool, `s` string, `2` [2]int (used for `[nameIdx, argCount]` pairs).

```go
data, _ := compiler.EncodeChunk(chunk)
os.WriteFile("prog.hbc", data, 0644)

chunk, _ := compiler.DecodeChunk(data)
compiler.RunTimed(chunk)
```



## 8.2 Pass 9 — VM Completeness

All four gaps documented at the end of Pass 8 are now closed:

| Gap | Opcode | Notes |
|-----|--------|-------|
| `is one of` with dynamic lists | `OP_CONTAINS` | Replaces literal-only unrolled equality chain |
| `for each key, val in obj` | `OP_FIELD_PAIRS` | Returns `[key, value]` pairs array, looped via standard indexed-loop pattern |
| File modules (`summon "lib.he"`) | `OP_RUN_MODULE_CHUNK` | Resolved and compiled at **compile time**, embedded as a `SubChunk` |
| `ask "..." as x` | `OP_ASK` | Reads a line from stdin via `fmt.Scanln` |

### File modules are compiled, not runtime-loaded

Unlike the interpreter (which reads and parses file modules at *run* time via `Runtime.loadFileModule`), the bytecode compiler resolves `summon "lib.he" as alias` at **compile* time: it reads the target file, parses it with the same `lang/lexer`/`lang/parser` pipeline, compiles it into a `MethodChunk`, and embeds that chunk directly in the parent chunk's `SubChunks`.

This means a compiled `.hbc` binary is **fully self-contained** — it does not need `lib.he` to exist on disk when later run via `he run prog.hbc`. The module's bytecode travels with the binary.

```go
// compiler/bytecode.go
func (bc *BytecodeCompiler) compileFileModule(st ast.LoadModuleStmt) {
    data, _ := fileModuleReader(st.FilePath, bc.BaseDir)
    prog, _ := parseModuleSource(data)
    sub := &BytecodeCompiler{...}
    for _, line := range prog.Lines {
        sub.compileLine(line)
    }
    // embed as SubChunk, emit OP_RUN_MODULE_CHUNK
}
```

`compiler/modreader.go` is a small bridge file that lets `compiler` import `lang/lexer`/`lang/parser` without creating an import cycle (verified: nothing in `lang/*` imports `compiler`).


## 9. Module System

### Built-in modules (summon by name)

```
summon "math"     as m
summon "text"     as t
summon "list"     as lst
summon "io"       as io
summon "clock"    as clock      (also: "time")
summon "net"      as net
summon "ui"       as ui
summon "physics"  as phys
summon "wolfhead" as wh         (also: "os")
```

### File modules (summon by path)

```
summon "path/to/mylib.he" as lib
```

When the module string ends in `.he`, `runtime.loadFileModule` is called:
1. Reads the file
2. Parses it
3. Executes it in a fresh sub-runtime
4. Wraps all exported variables and objects as fields of a module object
5. Registers the object under the given alias

**File module export:** everything at top level (`set`, `create`, `make`) is exported. No explicit `export` keyword needed.

**Example file module:**
```
~ mylib.he ~
set version to "1.0.0"

create Formatter [
  can: [
    pad(s, width) [
      ~ pad string to given width ~
      return s
    ]
  ]
]
```

### Inline summon (inside methods)

```
create Loader [
  can: [
    loadData [
      summon "io" as io
      return io.read("data.txt")
    ]
  ]
]
```

---

## 10. Persistence Layer

**File:** `lang/eval/persist.go`

### Storage format

JSON file at `~/.he_memory.json`:

```json
{
  "highscore": { "t": "n", "n": 750 },
  "playername": { "t": "s", "s": "Hunter" },
  "flags": { "t": "a", "a": [{"t": "b", "b": true}] }
}
```

### Type encoding

| HE type | JSON `t` field |
|---------|----------------|
| number  | `"n"` |
| text    | `"s"` |
| boolean | `"b"` |
| list    | `"a"` |
| nothing | `"nil"` |

Objects are not persistable (use file I/O for that).

### API

```
remember score           → saves Env["score"] to disk under key "score"
remember score as "best" → saves under key "best"
recall score             → loads key "score" into Env["score"]
recall best as score     → loads key "best" into Env["score"]
forget score             → removes key "score" from disk
```

---

## 11. Extension Guide

### Adding a keyword

1. Add `K_MYWORD TokenType = "myword"` to `lang/token/token.go`
2. Add `keywords["myword"] = token.K_MYWORD` to `lang/lexer/lexer.go` `init()`
3. Add `token.K_MYWORD` to `isStatementStart()` in parser if it starts a statement
4. Add `case token.K_MYWORD:` to `parseStatement()` or `parsePrimary()`
5. Add the AST node to `lang/ast/ast.go`
6. Add `case ast.MyNode:` to `execStatement()` or `evalExpr()` in runtime

### Adding an operator

1. Lex it in `NextToken()` (single or multi-char)
2. Insert it at the right precedence level in the parser expression hierarchy
3. Evaluate it in `evalExpr` or `evalBinary`/`evalCompare`

### Adding a builtin function (global, not on an object)

In `runtime.go`, add to `newRuntime()`:

```go
r.Env["myFunc"] = types.FromObject(&types.Object{
    Name: "myFunc",
    Actions: map[string]*types.Action{
        "__call__": {
            Name:   "__call__",
            Params: []string{"arg"},
        },
    },
    Builtins: map[string]types.BuiltinFn{
        "__call__": func(args []types.Value) (types.Value, error) {
            return types.FromString("hello " + args[0].Str), nil
        },
    },
})
```

---

## 12. Error Handling Model

### Parse errors

Format: `line N:M — expected X but got Y (token)`

The parser never panics. Errors are returned up the call stack and printed by `main`.

### Runtime errors

Errors propagate as Go `error` values. The `returnSignal` sentinel is NOT a real error — it is caught by `execAction` before escaping the action body.

`TryStmt` catches any error from its body:
```go
execErr := r.execBlock(st.Body)
if execErr != nil {
    r.setVar("error", types.FromString(execErr.Error()))
    // run handler
}
```

`TryWithVarStmt` additionally names the error:
```go
r.setVar(st.ErrVar, types.FromString(execErr.Error()))
```

### Friendly error messages

`main.go` passes errors through `friendly()` which strips Go type names before printing. Error messages are designed to reference HE concepts (`I don't know anything called "X"`) rather than Go internals.

---

## 13. Known Limitations

| Limitation | Status |
|------------|--------|
| Closures snapshot Env at creation time — mutations after capture not visible | Improving |
| `for each` over objects (not just lists) not yet supported | Planned |
| No tail-call optimisation — depth capped at 500 calls | Planned |
| Nested quotes in `{...}` interpolation require careful escaping | Known |
| File module exports are values only, not a proper namespace | Improving |
| `remember` doesn't support objects — only primitives and lists | Planned |
| `wait` in the interpreter is a real `time.Sleep` — blocks the process | By design for now |
| No concurrency primitives | Planned (WolfHead pass) |
| Bytecode VM implemented for arithmetic, control flow, arrays, persistence | Pass 8: full object support |

---

## 14. Roadmap

### Pass 6 — Bytecode VM
- Execute `.hbc` compiled files without full interpreter
- Stack-based VM consuming `Chunk`s
- Significant performance improvement for compute-heavy code

### Pass 7 — Native Compilation
- QBE or LLVM backend
- Compile `.he` → Linux ELF binary
- WASM target for browser/WolfHead apps

### Pass 8 — WolfHead Deep Integration
- Replace `wolfhead` module print stubs with real Wayland/compositor calls
- Workspace lifecycle management in HE
- Context-aware app scheduling
- Gesture bindings from compositor → HE reactions

### Pass 9 — Type System Hardening
- Optional type annotations enforced at runtime
- Generic list types: `list of number`
- Type-safe ability signatures
- Compile-time type errors

### Pass 10 — Package Manager
- `summon "github.com/user/lib"` fetches and caches HE packages
- `he.mod` dependency file
- Package versioning
