# HE Language — Changelog

All notable changes to HE are documented here.
Format: `[Pass N] — Description`

---

## [Pass 4] — Completeness + Compiler Foundation
*Anonymous abilities, multi-assign, is between, remember/recall, bytecode IR, type inference*

### Language features added
- **`X is between A and B`** — English range comparison, expands to `A <= X <= B`
- **Anonymous abilities** — first-class function values
  ```
  set greet to ability(name) [ return "Hello, {name}!" ]
  say greet("Hunter")
  set result to apply(double, 21)
  ```
- **`set a, b to expr1, expr2`** — multi-variable assignment in one line
  ```
  set x, y to 10, 20
  set x, y to y, x   ~ swap ~
  ```
- **`return a, b, c`** — return multiple values from an ability
  ```
  set quotient, remainder to Calculator.divmod(17, 5)
  set lo, hi to Calculator.minmax(numbers)
  ```
  When a multi-return result is unpacked into a multi-assign, values are distributed automatically.
- **`remember name`** — persist a variable to disk (`~/.he_memory.json`)
  - Optional key: `remember score as "highscore"`
- **`recall name`** — load a remembered variable back
  - Optional alias: `recall highscore as currentScore`
- **`forget name`** — remove a remembered variable
- **`summon` inside method bodies** — modules can now be summoned inside `can:` blocks
- **Full expression interpolation** — `{double(21)}`, `{apply(fn, x)}`, `{a + b}` all work inside strings

### New tooling
- **`he -file prog.he -types`** — type inference report
  - Infers: number, text, boolean, list, object, ability, unknown
  - Shows all variable types and object field types
  - Reports type errors (mismatches, undefined types)
- **`he -file prog.he -compile`** — bytecode disassembly (compiler preview)
  - Full bytecode IR with named opcodes
  - 40 opcodes covering: arithmetic, comparison, control flow, objects, arrays, persistence, I/O
  - Forward-jump patching for if/while/range loops
  - Foundation for native binary and WASM compilation

### New packages
- **`compiler/types.go`** — type inferencer and type environment
- **`compiler/bytecode.go`** — bytecode IR, opcode definitions, chunk compiler

### Fixed
- Anonymous abilities callable from interpolation: `{fn(arg)}`  
- Multi-return unpacking: array result automatically distributed across multi-assign names
- Ability calls in expression context now check Env before Objects

---

## [Pass 3] — Language Completeness + OS Foundation
*Range loops, error handling, string interpolation, net & WolfHead modules*

### Language features added
- **String interpolation** — `"Hello, {name}! Score: {score}."` — embed any variable or `obj.field` directly in strings; no more `"Hello, " + name`
- **`for each i from 1 to 10 [...]`** — numeric range loop
  - With step: `for each n from 2 to 10 step 2 [...]`
- **`repeat until cond [...]`** — loop that runs *until* a condition becomes true
  - Also: `repeat until` as sub-form of `repeat`
- **`try [...] or [...] if it fails`** — English error handling
  - On error, `error` variable is set to the error message
  - Handler is optional — `try [...]` alone silently ignores errors
  - Forms accepted: `or [...]`, `or if it fails [...]`, `or if something fails [...]`
- **`{obj.field}` in interpolation** — dot-access works inside `{...}` in strings

### New stdlib modules
- **`summon "net" as net`** — HTTP client
  - `net.get(url)` → response object with `.body`, `.status`, `.ok`
  - `net.post(url, body)` → response object
  - `net.status(response)` → status code number
- **`summon "wolfhead" as wh`** — WolfHead OS bindings
  - `wh.workspace()` → current workspace object (`.id`, `.name`, `.active`)
  - `wh.workspace(n)` → switch to workspace n
  - `wh.context()` → list all contexts, or `wh.context("Work")` to switch
  - `wh.notify(message)` or `wh.notify(title, message)` → OS notification
  - `wh.launch(appName)` → launch an application
  - `wh.gesture(name)` → trigger a gesture
  - `wh.platform()` → returns `"WolfHead/Linux"`
  - Also available as `summon "os" as os`

### Fixed
- `repeat until` correctly parsed as sub-form of `repeat` keyword
- String interpolation now handles `{obj.field}` dot-access expressions
- Empty string segments in interpolation handled correctly

---

## [Pass 2] — Evolution Pass
*Features, stdlib expansion, REPL, living docs*

### Added
- **`for each item in list [ ... ]`** — English for-each loop over any list
- **`is not`** — two-word not-equal operator (`mood is not "sad"`)
- **`not expr`** — English unary negation (alias for `!`)
- **`ask "prompt" as varName`** — read user input from stdin
  - Also: `ask "..." then set varName to answer`
  - Also: `ask "..." storing result in varName`
- **`set obj.field to expr`** — dot-assignment on object fields
- **`give value to obj.field`** — alternate dot-assignment syntax
- **`grow obj.field by N`** — field increment via dot notation
- **`show`** — third alias for `say` / `print`
- **`check if`** — alias for `if` (reads: `check if done is not true then [...]`)
- **`know how to:`** — alias for `can:` in object abilities section

### New stdlib modules
- **`summon "io" as io`**
  - `io.read(path)` — read file as text
  - `io.write(path, content)` — write text to file
  - `io.append(path, content)` — append text to file
  - `io.exists(path)` — boolean file existence check
  - `io.delete(path)` — remove a file

### New tools
- **REPL** (`go run ./cmd/repl`)
  - Persistent state across inputs
  - Multiline block detection (waits for `]` before executing)
  - `:help` — quick reference guide
  - `:state` — inspect all variables and objects
  - `:clear` — reset runtime state
  - `:run <file>` — load a `.he` file into current state
  - `:quit` / `:exit` — leave the REPL

### New docs
- **`SPEC.md`** — complete language specification with grammar
- **`CHANGELOG.md`** — this file

### Fixed (from Pass 1)
- `return` in ability bodies now correctly propagates return value
- `Interpreter.RunProg()` added for incremental execution (REPL foundation)
- `Interpreter.Run(nil)` safe for runtime initialisation without a program

---

## [Pass 1] — Foundation Pass
*Initial build: lexer, parser, AST, runtime, stdlib core*

### Language features implemented
- **`say` / `print`** — output to stdout
- **`set X to Y`** — variable assignment
- **`let X be Y`** — alias for `set`
- **`change X to Y`** — alias for `set`
- **`grow X by N`** — increment (`x = x + n`)
- **`shrink X by N`** — decrement (`x = x - n`)
- **`if cond then [...] else [...]`** — conditional
- **`repeat N times [...]`** — counted loop
- **`repeat while cond [...]`** — condition loop
- **`while cond [...]`** — alias for repeat while
- **`return expr`** — return value from ability
- **`wait N seconds`** / **`wait N frames`** — pause execution
- **`tell obj to action`** / **`tell obj to action with args`** — method call
- **`summon "module" as alias`** — import stdlib module
- **`create Name [...]`** / **`make Name [...]`** — object definition
- **`create Name like Parent [...]`** — object inheritance
- **`has: [name is val, ...]`** — property section
- **`owns:`** / **`carries:`** — aliases for `has:`
- **`can: [actionName [body]]`** — abilities section
- **`on trigger [body]`** — reaction (event handler)
- **`when trigger [body]`** / **`whenever trigger [body]`** — aliases for `on`
- **`on event with "qualifier" [body]`** — qualified reaction trigger
- **Dot notation**: `obj.method(args)` and `obj.field`
- **Named arg blocks**: `[title is "...", children: [...]]`
- **`yes`** / **`no`** — boolean aliases for `true` / `false`
- **`nothing`** — nil value
- **`not`** — unary negation keyword
- **`and`** / **`or`** — logical operators
- **`is`** — equality operator in expressions
- **`**`** — power operator
- **`~`** comment syntax (`~ comment ~`)

### Critical bugs fixed (vs original codebase)
- `registerBuiltinMethod` was discarding functions (`_ = impl`) — now wired correctly
- `if`, `return`, `not` were fragile string-checked hacks — now proper keyword tokens
- `parseType` consumed `[` meant for action bodies — fixed with one-token lookahead
- `on event with "string"` triggers rejected string literals — now accepted
- Builtin dispatch used hardcoded `recv.Name` strings — now uses `Builtins` map
- `BANG` token (`!`) had duplicate case in switch — cleaned up

### Stdlib modules
- **`summon "math" as m`** — abs, sqrt, floor, ceil, round, max, min, random, pow, sin, cos, pi
- **`summon "text" as t`** — upper, lower, length, contains, starts, ends, replace, trim, split, join, number, from
- **`summon "list" as lst`** — length, get, add, remove, contains, first, last
- **`summon "ui" as ui`** — window, button, text, navbar, renderDocs
- **`summon "physics" as phys`** — gravity, collision

### Architecture
- `lang/token` — all token type definitions
- `lang/lexer` — tokeniser with comment handling (`~ ... ~`)
- `lang/ast` — full AST node definitions
- `lang/parser` — recursive-descent parser
- `lang/types` — runtime value and object types
- `lang/eval/runtime.go` — statement/expression executor
- `lang/eval/stdlib.go` — all builtin module implementations
- `lang/eval/eval.go` — interpreter entry point
- `cmd/hunterlang/main.go` — CLI runner

---

## Planned (Pass 3 and beyond)

- **`each N from A to B [...]`** — numeric range loop (`for each i from 1 to 10`)
- **`X is between A and B`** — range comparison
- ~~**String interpolation**~~ ✓ Done in Pass 3
- **`remember`** / **`forget`** — persistent variable storage across runs
- ~~**`summon "net" as net`**~~ ✓ Done in Pass 3
- **`summon "time" as clock`** — date/time utilities
- **Error handling** — `try [...] or [...] if it fails`
- **Multiple return values** — `return name, age`
- **Anonymous abilities** — `set handler to ability(x) [...]`
- ~~**WolfHead OS bindings**~~ ✓ Done in Pass 3
- **Compiler** — compile `.he` to native binary or WASM
