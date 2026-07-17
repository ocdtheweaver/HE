# HE Language — How To Use
## Version 5 — Complete Setup and Usage Guide

---

## What is HE?

HE is a programming language that reads like English. You don't need to memorize symbols, semicolons, or curly braces. If you can write a sentence, you can write HE.

```he
set name to "Hunter"
say "Hello, {name}!"
```

That's a complete program.

---

## Installation

### Requirements

- Go 1.22 or later — download from https://go.dev/dl/

### Build from source

```bash
# Unzip the HE folder, then:
cd HE

# Build the 'he' command-line tool
go build -o he ./cmd/hunterlang

# Move it somewhere on your PATH (Linux/Mac)
mv he /usr/local/bin/

# Windows: move he.exe to any folder on your PATH
```

### Verify it works

```bash
he version
```

Expected output:

```
HE language  v0.7.0
Build:       Pass 14
Runtime:     tree-walk interpreter + bytecode VM
```

---

## The `he` command

```
he <command> [file] [flags]

Commands:
  run     <file>          Run a .he program (or .hex compiled binary)
  check   <file>          Show type inference and protection tag report
  build   <file> --o out  Compile to a .hex bundle
  vm      <file>          Run via the bytecode VM
  bench   <file>          Benchmark interpreter vs VM
  repl                    Interactive shell
  new     <name>          Create a new HE project
  version                 Show version info
  help    [command]       Help for any command
```

---

## Your first program

Create `hello.he`:

```he
say "Hello, world!"
```

Run it:

```bash
he run hello.he
```

---

## Starting a new project

```bash
he new myproject
cd myproject
he run main.he
```

Creates:

```
myproject/
├── main.he          ← entry point
├── lib/
│   └── helpers.he   ← shared utilities
└── README.md
```

---

## The genesis.he convention

Name your entry point `genesis.he` and run with no filename:

```bash
he run
```

HE finds `genesis.he` automatically — same role as `index.html` for websites or `main()` elsewhere.

---

## Bringing in other files — `summon`

### Flat import (no alias needed)

```he
summon "math.he"

say sqrt(16)      ~ works immediately ~
say double(21)    ~ no prefix required ~
```

### Aliased import (more control)

```he
summon "UI.he" as ui

ui.nav("Home")
ui.button("Sign In")
```

Both at once — alias never gates flat access, just adds a second path:

```he
summon "math.he" as maths

say sqrt(16)         ~ flat access ~
say maths.sqrt(16)   ~ qualified access ~
```

### Built-in modules (alias always required)

```he
summon "math"   as m
summon "text"   as t
summon "list"   as lst
summon "io"     as io
summon "clock"  as clock
```

### Collision rule

Two summoned files defining the same name → error on bare access, must alias:

```he
summon "physics.he" as phys
summon "ui.he"      as ui

phys.render()   ~ unambiguous ~
ui.render()     ~ unambiguous ~
```

---

## Compiling to a .hex bundle

```bash
he build genesis.he --o myapp.hex
he run myapp.hex
```

A `.hex` file is fully self-contained — all summoned modules embedded. No source files needed to run it.

---

## The interactive shell

```bash
he repl
```

| Command | What it does |
|---------|-------------|
| `:help` | Show REPL help |
| `:state` | Show current variables |
| `:clear` | Reset state |
| `:quit` | Exit |

---

## Protection tags

Mark features as requiring entitlement:

```he
~ Whole object tagged ~
create Standard #protected1 [
  can: [
    watchEpisode(title) [ say "Streaming {title}" ]
  ]
]

~ Individual ability tagged ~
create Replay [
  can: [
    watchLive(event) #protected3 [ say "Live: {event}" ]
    browseSchedule [ say "Free to browse" ]
  ]
]
```

Define tag meanings in `protected.he`:

```he
create protected1 [
  has: [
    planName is "Standard"
    mode is "cached"
    endpoint is "https://api.myapp.com/entitlement/standard"
    graceSeconds is 604800
  ]
]

create protected3 [
  has: [
    planName is "Live Events"
    mode is "live"
    endpoint is "https://api.myapp.com/entitlement/live"
    graceSeconds is 0
  ]
]
```

Summon it:

```he
summon "protected.he"
```

Wrap protected calls in `try/or` — denial is a catchable error:

```he
try [
  tell Standard to watchEpisode with "Episode 1"
] or (err) [
  say "Upgrade required: {err}"
]
```

Check what's tagged:

```bash
he check genesis.he
```

**Current enforcement status:**
- Default: deny all (fail closed)
- Enforced across `he run`, `he vm`, and `.hex` binaries
- Network entitlement checks (real HTTP calls) are the next planned step

---

## Language quick reference

### Variables and assignment

```he
set score to 0
let name be "Hunter"
score becomes 42
grow score by 5
shrink score by 2
set x, y to 10, 20        ~ multi-assign ~
```

### Output and input

```he
say "Score: {score}"
print "Same thing"
ask "Your name?" as name
```

### Conditions

```he
if score > 100 then [
  say "High score!"
] else [
  say "Keep going."
]

if age is between 13 and 17 then [ say "Teenager" ]
if color is one of ["red", "green", "blue"] then [ say "Primary" ]
```

### Loops

```he
repeat 3 times [ say "Again" ]
repeat 5 times as i [ say "Round {i}" ]
repeat while lives > 0 [ shrink lives by 1 ]
for each item in list [ say item ]
each i from 1 to 10 [ say i ]
each n from 0 to 100 step 5 [ say n ]
```

### Objects

```he
create Dog [
  has: [
    name is "Rex"
    breed is "Labrador"
  ]
  can: [
    bark [ say "{name} says Woof!" ]
    greet(person) [ say "Hello {person}, I'm {name}" ]
  ]
  on click [ tell Dog to bark ]
]

tell Dog to bark
tell Dog to greet with "Hunter"
say Dog.name
set Dog.name to "Buddy"
```

### Inheritance

```he
create Cat like Dog [
  has: [ name is "Whiskers" ]
  can: [ bark [ say "Meow!" ] ]
]
```

### Anonymous abilities

```he
set double to ability(n) [ return n * 2 ]
say double(21)
```

### Multiple return values

```he
create Math [
  can: [
    divmod(a, b) [
      return a / b, a - ((a / b) * b)
    ]
  ]
]
set quotient, remainder to Math.divmod(17, 5)
```

### Error handling

```he
try [
  set data to io.read("config.txt")
] or (err) [
  say "Couldn't load: {err}"
]
```

### Persistence

```he
set highscore to 750
remember highscore

recall highscore
say "Best: {highscore}"
forget highscore
```

### Method chains

```he
say "hello".upper()
say "a,b,c".split(",").length()
say (-42).abs()
say items.reverse().first()
```

### String interpolation

```he
say "Score: {score}"
say "Double: {double(21)}"
say "Sum: {a + b}"
```

### Field iteration

```he
for each key, val in Config [
  say "{key} = {val}"
]
```

---

## Standard library

| Module | Key methods |
|--------|-------------|
| `math` | `sqrt`, `abs`, `round`, `floor`, `ceil`, `max`, `min`, `pow`, `pi`, `random` |
| `text` | `upper`, `lower`, `length`, `trim`, `split`, `replace`, `contains` |
| `list` | `length`, `first`, `last`, `add`, `remove`, `contains`, `reverse`, `join` |
| `io` | `read`, `write`, `append`, `exists`, `delete` |
| `clock` | `now`, `today`, `format`, `since`, `sleep`, `timestamp` |
| `net` | `get`, `post` |
| `wolfhead` | `workspace`, `context`, `notify`, `launch`, `gesture`, `platform` |

---

## What's coming next

| # | Feature | Status |
|---|---------|--------|
| 1 | `protected.he` policy format | ✅ Done |
| 2 | Stub entitlement check | ✅ Done |
| 3 | Full enforcement (interpreter + VM + .hex) | ✅ Done |
| 4 | Offline cache/grace-period logic | 🔜 Next |
| 5 | Real `.hbc` hardened distribution | 🔜 Planned |
| 6 | Network entitlement check (HTTP) | 🔜 Planned |

---

*HE is built with Go. Source at github.com/ocdtheweaver/HE*
