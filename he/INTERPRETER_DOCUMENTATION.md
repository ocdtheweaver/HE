# HE Interpreter (Go) — Documentation

This repository includes a working HE interpreter that runs HE scripts via:

```bash
cd Go/he
go run ./cmd/he run <file.he>
```

## What this interpreter currently supports

### 1) Program structure
The interpreter is **line-based** and supports these top-level constructs:

- `summon "path" as Alias` / `summon "path" named Alias`
- `with image "..." named "..." and sound "..." ...`
- `make <ObjectName> [like <Parent>] { ... }`
- `create <ObjectName> [like <Parent>] { ... }`
- Global statements:
  - `print <expr>` / `say <expr>`
  - `wait <number>`

Inside objects it supports:

- `Object has: [ ... ]` properties
- `Object can: [ ... ]` abilities (methods)
- `on <trigger> [ ... ]` / `when <trigger> [ ... ]` / `whenever <trigger> [ ... ]` reactions

> Note: This is the prototype interpreter shipped as part of `Go/he/cmd/he/main.go`.  
> The more formal “v2” parser/runtime under `Go/he/pkg/*` is currently not fully aligned, but the CLI **does run** via the prototype interpreter.

---

### 2) Expressions (arithmetic + concatenation)
Expressions support:

- Numbers: `12`, `3.14`
- Variables / identifiers: `score`, `Player.jump` (identifier tokens are fairly permissive)
- Quoted strings: `"hello"`
- Operators: `+ - * /`
- Parentheses: `( ... )`
- Unary minus: `-x`

Operator behavior:
- `number + number` → number
- Any `+` involving non-numbers → coerces to strings and concatenates

Example valid expressions:
- `1 + 2 * 3`
- `"score: " + score`
- `-(1 + 2)`

---

### 3) Variable assignment
Assignments supported at global scope and inside abilities/reactions:

- `x is <expr>`
- `x = <expr>`

The interpreter stores variables in an internal map `rt.Vars`.

---

### 4) Printing
`print` / `say`:
- If the expression evaluates to a number, it prints it (dropping trailing `.0` when integer-like)
- If it evaluates to a string, it prints the string value
- Otherwise it prints the raw value (string representation fallback)

---

### 5) Waiting
`wait <number>` pauses the interpreter for that many seconds.

---

### 6) Assets and `tell`
Assets are declared with `with ...` blocks and stored by alias (or path if no alias).

`tell <Target> to <action> [with <args>]`:
- If `<Target>` matches a declared asset alias, it prints an asset action log.
- If `<Target>` matches an object name and `<action>` matches an ability on that object, it executes that ability’s statements.

---

### 7) Object reactions / triggers
Reactions are stored by trigger name.
The interpreter additionally simulates:

- A reaction called `"start"` for any object that defines it (runs at load time)
- A small demo hook for the `Player` object:
  - if `Player` defines ability `jump`, it runs `jump`
  - if `Player` defines reaction `collision`, it runs `collision`

These hooks are “demo-like” and may be changed in the future.

---

## How to use

### A) Run a script
```bash
cd Go/he
go run ./cmd/he run examples/game.he
```

### B) Help
```bash
go run ./cmd/he help
```

---

## Limitations / current behavior notes

- This interpreter is **not the complete “v2” AST runtime**; it is the prototype interpreter embedded into the CLI.
- The “v2” packages (`Go/he/pkg/parser`, `Go/he/pkg/runtime`) are currently not wired end-to-end due to mismatches.
- Parsing inside objects is pattern-driven and bracket-based; it supports the existing example style, but does not yet cover every grammar feature in `Go/grammar.ebnf`.

---

## Entry point files
- CLI + prototype interpreter implementation: `Go/he/cmd/he/main.go`
