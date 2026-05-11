# How the HE Interpreter Was Built (Go)

This explains how the working `he run` interpreter was produced and how to extend it.

## High-level reality check
The repository contains an earlier “v2” implementation under `Go/he/pkg/*` (lexer/parser/ast/runtime). That code is currently **not fully wired together** (types and field names don’t match, plus some missing parser routines like `parseBlock`).

To deliver a functional interpreter quickly, the working solution uses the already-existing **prototype interpreter** logic from `Go/main.go` and embeds it into the CLI command that users actually run.

So today, the *actual* interpreter you execute via `he run` lives in the CLI file:
- `Go/he/cmd/he/main.go`

## What runs when you execute the CLI
Command:
```bash
cd Go/he
go run ./cmd/he run examples/game.he
```

### Entry point
`Go/he/cmd/he/main.go` is the `main` package.
It implements:
- argument parsing (`run`, `help`)
- `runFile(filename)` which:
  1. reads the `.he` source file
  2. parses it into an internal `Program`
  3. executes `Program.Globals`
  4. optionally triggers “start/collision/jump” demo hooks for objects if present

## Components inside the prototype interpreter

### 1) Parser (line-based)
Implemented directly in `Go/he/cmd/he/main.go` as a small, line-based parser:

- `type Parser struct { lines []string; pos int }`
- `ParseProgram() (*Program, error)`

It walks the file line-by-line and recognizes top-level constructs:
- `summon "..." as X` / `summon "..." named X`
- `with image "..." named ... and sound "..." ...` (multi-line support via `peekLine()`)
- `make <Name> [like <Parent>] { ... }` and `create ...`
- global statements like `print ...`, `say ...`, `wait ...`, and `x is ...` / `x = ...`

#### Parsing objects
`parseObject(firstLine)` handles:
- properties blocks (`has:` / `owns:` / `carries:`) by reading lines until `]`
- abilities blocks (`can:` / `knows how to`) by reading until `]`
- reactions (`on|when|whenever <trigger> [ ... ]`) by reading until the closing `]`

This design is intentionally “pragmatic”:
it mirrors the style of existing example scripts, rather than implementing the full EBNF grammar.

### 2) Runtime (execution)
Runtime data structures in the same file:

- `type Runtime struct { Program *Program; Objects ...; Assets ...; Vars map[string]Value }`

- `NewRuntime(p *Program)` initializes:
  - asset aliases into `rt.Assets`
  - object properties into `rt.Vars` using keys like:
    - `ObjectName.propertyName`

Execution:
- Each statement implements:
  - `Exec(rt *Runtime) error`

Global execution:
- after parsing, `runFile()` loops:
  - `for _, s := range prog.Globals { s.Exec(rt) }`

Additionally:
- it triggers reaction `"start"` for any object that defines it
- it triggers a demo for `Player`:
  - `Player.ability(jump)` if present
  - `Player.reaction(collision)` if present

### 3) Expression evaluation
Arithmetic and concatenation happen in functions:

- `tokenizeExpr(s string) ([]exprTok, error)`
- `shuntingYard(tokens []exprTok) ([]exprTok, error)`
- `evalRPN(rpn []exprTok, rt *Runtime) (Value, error)`
- `evaluateExpression(expr string, rt)`

Supported:
- numbers + / - * with parentheses
- quoted strings `"..."` supporting `+` concatenation
- identifiers resolved from `rt.Vars`
  - unknown identifiers default to `0` (for expression evaluation)

Unary minus:
- implemented in the shunting-yard pass by rewriting `-x` to `0 - x`.

### 4) Implemented statement types
Currently handled statement types include:
- `PrintStmt` / `say` (both map to the same struct)
- `WaitStmt`
- `TellStmt` (prints asset actions; can execute an ability on an object)
- `AssignStmt` (`x is expr` / `x = expr`)
- ability bodies and reaction bodies are executed by iterating their parsed statements

## Why the prototype approach was chosen
Because the v2 packages (`Go/he/pkg/lexer`, `Go/he/pkg/parser`, `Go/he/pkg/runtime`, etc.) currently fail to compile as a unified interpreter, embedding the proven prototype interpreter into the CLI guarantees:
- users can run scripts today (`he run`)
- the interpreter behavior is at least consistent with the existing prototype’s supported constructs

## How to extend the interpreter (the practical route)
Since `he run` is currently backed by the prototype interpreter:

### A) Add new statement syntax
Add:
1. Parsing in `parseStmt` and/or `parseGlobalStatement`
2. A new `type <Stmt> struct { ... }` that implements `Exec(rt *Runtime)`

Then update `TellStmt` or other existing statement execution as needed.

### B) Add new operators or expression features
Update:
- `tokenizeExpr` (if new tokens are required)
- precedence tables in `precedence` / shunting-yard logic
- `evalRPN` operator handling

### C) Add better object method invocation
Today:
- `tell Target to abilityName` looks up `rt.Objects[Target].Abilities[abilityName]`
- it ignores args at runtime (the parser captures args as raw strings but runtime doesn’t yet use them)

To implement args:
- extend `TellStmt.Exec` to pass values into a function/ability execution context
- store parameters in runtime environment / local scope

## Files you should modify (today)
- `Go/he/cmd/he/main.go`
  - prototype parser
  - prototype runtime
  - expression evaluation
  - CLI command routing

## Where the “v2” interpreter sits (not yet wired)
- `Go/he/pkg/lexer`
- `Go/he/pkg/parser`
- `Go/he/pkg/runtime`
- `Go/he/pkg/ast`
- `Go/he/pkg/resolver`

If you want to finish the “proper” interpreter using v2 AST + runtime, you will need to reconcile:
- AST field names (the `Pos` collision was fixed by renaming fields to `Posn`)
- missing parser routines (`parseBlock` etc.)
- lexer token names vs parser token references
- runtime execution support for statement kinds (currently limited)

## Quick usage recap
```bash
cd Go/he
go run ./cmd/he run examples/game.he
```

The interpreter prints:
- any `print` / `say`
- asset/tell logs for `tell <asset> to play ...`
- runtime-loaded summary:
  `HE program loaded. Objects: N Assets: M`
