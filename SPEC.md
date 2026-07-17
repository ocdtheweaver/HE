# HE Language Specification
*Living document — updated with every evolution pass*
*Last updated: Evolution Pass 2*

---

## Philosophy

HE is a near-English programming language. The goal is that a person with no programming background can read an HE program and understand what it does. Every keyword is a real English word. Every structure reads like a sentence.

**Principles:**
1. **English first** — keywords are words, not symbols. `is not` not `!=`. `grow` not `+=`.
2. **Explicit is kind** — `set name to "Hunter"` not `name = "Hunter"`. The structure tells you what's happening.
3. **Blocks are brackets** — `[` and `]` open and close every block, like sentences in a paragraph.
4. **Comments are tildes** — `~ this is a comment ~`. Natural, unobtrusive.
5. **No surprises** — no implicit type coercion beyond string concatenation with `+`.

---

## Comments

```
~ This is a comment ~
~ Comments open and close with a tilde.
  They can span multiple lines. ~
```

---

## Values and Types

| Type      | Examples                          | Notes                          |
|-----------|-----------------------------------|--------------------------------|
| `number`  | `42`, `3.14`, `-9.8`             | All numbers are 64-bit floats  |
| `text`    | `"Hello"`, `"world"`             | Strings, always double-quoted  |
| `boolean` | `true`, `yes`, `false`, `no`     | `yes`/`no` are aliases         |
| `list`    | `["apple", "banana"]`, `[1,2,3]` | Mixed types allowed            |
| `nothing` | `nothing`                        | The nil/null value             |

---

## Variables

```
set score to 0
let name be "Hunter"
change score to 100

~ All three forms are equivalent. Use whichever reads most naturally. ~
```

### Mutation

```
grow score by 10        ~ score = score + 10 ~
shrink health by 5      ~ health = health - 5 ~
```

---


## String Interpolation

Embed values directly in strings using `{...}`:

```
set name to "Hunter"
set score to 42
say "Hello, {name}! Your score is {score}."
say "Pi is {pi}."
say "Workspace: {ws.name}"
```

Any variable, field access (`obj.field`), or built-in value works inside `{...}`.
For complex expressions, compute the value first and then interpolate.

## Output and Input

```
say "Hello, world"
print "Same as say"
show "Also the same"

ask "What is your name?" as myName
say "Hello, " + myName
```

---

## Operators

### Arithmetic
| Operator | Meaning        | Example              |
|----------|----------------|----------------------|
| `+`      | Add / join     | `5 + 3` or `"Hi" + name` |
| `-`      | Subtract       | `10 - 4`             |
| `*`      | Multiply       | `6 * 7`              |
| `/`      | Divide         | `10 / 2`             |
| `**`     | Power          | `2 ** 8`             |

### Comparison
| Operator   | Meaning              | Example              |
|------------|----------------------|----------------------|
| `==`       | Equal                | `x == 5`             |
| `is`       | Equal (English)      | `name is "Hunter"`   |
| `!=`       | Not equal            | `x != 0`             |
| `is not`   | Not equal (English)  | `mood is not "sad"`  |
| `>`        | Greater than         | `score > 100`        |
| `<`        | Less than            | `hp < 10`            |
| `>=`       | Greater or equal     | `level >= 5`         |
| `<=`       | Less or equal        | `time <= 60`         |

### Logic
| Operator | Meaning | Example                        |
|----------|---------|--------------------------------|
| `and`    | Both    | `x > 0 and x < 10`            |
| `or`     | Either  | `done or failed`               |
| `not`    | Negate  | `not ready`                    |

---

## Conditions

```
if score > 100 then [
  say "High score!"
] else [
  say "Keep going."
]

~ check if is an alias for if ~
check if name is "Hunter" then [
  say "Welcome back."
]
```

---

## Loops

### Count loop
```
repeat 5 times [
  say "Again!"
]
```

### Condition loop
```
repeat while health > 0 [
  shrink health by 10
]

~ 'while' alone also works ~
while alive [
  tell Player to update
]
```

### For-each loop
```
set fruits to ["apple", "banana", "cherry"]
for each fruit in fruits [
  say "I have a {fruit}"
]
```

### Range loop
```
~ count from 1 to 5 ~
for each i from 1 to 5 [
  say "Step {i}"
]

~ count by 2s ~
for each n from 0 to 20 step 2 [
  say "{n}"
]
```

---


## Error Handling

```
try [
  set data to io.read("config.txt")
  say "Loaded: {data}"
] or if it fails [
  say "Could not load config: {error}"
]

~ Shorter forms ~
try [
  tell server to connect
] or [
  say "Connection failed: {error}"
]

~ Try without handler — errors are silently ignored ~
try [
  io.delete("temp.txt")
]
```

When an error occurs, the special variable `error` holds the error message.

## Repeat Until

```
repeat until health == 0 [
  tell enemy to attack
  shrink health by enemy.damage
]
```

## Objects

Objects are the core of HE. They group data and behaviour together.

```
create Player [
  has: [
    name is "Hero"
    health starts as 100
    score starts as 0
  ]
  can: [
    greet [
      say "I am " + name
    ]
    heal(amount) [
      grow health by amount
    ]
    getHealth returns number [
      return health
    ]
  ]
  on damage [
    shrink health by 10
    if health < 0 then [
      set health to 0
    ]
  ]
]
```

### Property sections

| Keyword    | Meaning                                  |
|------------|------------------------------------------|
| `has`      | Properties this object starts with       |
| `owns`     | Same as `has`                            |
| `carries`  | Same as `has`                            |
| `remembers`| Persistent properties (planned)          |

### Initialising properties

```
name is "Hero"          ~ set immediately ~
health starts as 100    ~ alias for 'is' ~
```

### Abilities (`can`)

```
can: [
  greet [
    say "Hello!"
  ]
  add(a, b) [
    return a + b
  ]
]

~ 'know how to' is an alias for 'can' ~
know how to: [
  fly [
    say "Whoosh!"
  ]
]
```

### Reactions (`on` / `when` / `whenever`)

```
on click [
  say "Clicked!"
]

on click with "Increment" [
  grow count by 1
]

when collision [
  say "Ouch!"
]
```

---

## Calling objects

```
tell Player to greet
tell Player to heal with 20
```

### Dot notation

```
say Player.health          ~ read a field ~
set Player.name to "Max"   ~ write a field ~
give 50 to Player.score    ~ same as set Player.score to 50 ~
grow Player.score by 10    ~ field increment ~
```

---

## Inheritance

```
create Warrior like Player [
  has: [
    weapon is "sword"
  ]
  can: [
    strike [
      say name + " strikes with " + weapon
    ]
  ]
]
```

The `Warrior` inherits all fields, abilities, and reactions from `Player`.

---

## Modules (Standard Library)

```
summon "math" as m
summon "text" as t
summon "list" as lst
summon "io"   as io
summon "ui"   as ui
summon "physics" as phys
```

### `math`
| Method         | Description                   |
|----------------|-------------------------------|
| `m.abs(n)`     | Absolute value                |
| `m.sqrt(n)`    | Square root                   |
| `m.floor(n)`   | Round down                    |
| `m.ceil(n)`    | Round up                      |
| `m.round(n)`   | Round to nearest              |
| `m.max(a, b)`  | Larger of two numbers         |
| `m.min(a, b)`  | Smaller of two numbers        |
| `m.random()`   | Random 0–1                    |
| `m.random(a,b)`| Random between a and b        |
| `m.pow(a, b)`  | a to the power of b           |
| `m.sin(n)`     | Sine (radians)                |
| `m.cos(n)`     | Cosine (radians)              |
| `m.pi()`       | π (3.14159…)                  |

### `text`
| Method                   | Description                        |
|--------------------------|------------------------------------|
| `t.upper(s)`             | Uppercase                          |
| `t.lower(s)`             | Lowercase                          |
| `t.length(s)`            | Character count                    |
| `t.contains(s, sub)`     | True if s contains sub             |
| `t.starts(s, prefix)`    | True if s starts with prefix       |
| `t.ends(s, suffix)`      | True if s ends with suffix         |
| `t.replace(s, old, new)` | Replace all occurrences            |
| `t.trim(s)`              | Remove leading/trailing whitespace |
| `t.split(s, sep)`        | Split into a list                  |
| `t.join(list, sep)`      | Join list into text                |
| `t.number(s)`            | Convert text to number             |
| `t.from(val)`            | Convert any value to text          |

### `list`
| Method                | Description                    |
|-----------------------|--------------------------------|
| `lst.length(l)`       | Number of items                |
| `lst.get(l, i)`       | Item at index i (0-based)      |
| `lst.add(l, item)`    | New list with item appended    |
| `lst.remove(l, i)`    | New list with item at i removed|
| `lst.contains(l, x)`  | True if list contains x        |
| `lst.first(l)`        | First item                     |
| `lst.last(l)`         | Last item                      |

### `io`
| Method               | Description                        |
|----------------------|------------------------------------|
| `io.read(path)`      | Read file contents as text         |
| `io.write(path, s)`  | Write text to file (overwrite)     |
| `io.append(path, s)` | Append text to file                |
| `io.exists(path)`    | True if file exists                |
| `io.delete(path)`    | Delete a file                      |

### `ui`
| Method                        | Description                     |
|-------------------------------|---------------------------------|
| `ui.window([title is "..."])` | Render an HTML window           |
| `ui.button(label)`            | UI button element               |
| `ui.text(content)`            | UI text element                 |
| `ui.navbar(items)`            | Navigation bar                  |
| `ui.renderDocs(list)`         | Render a tabbed docs page       |

### `net`
| Method              | Description                              |
|---------------------|------------------------------------------|
| `net.get(url)`      | HTTP GET → response (`.body`, `.status`, `.ok`) |
| `net.post(url, body)` | HTTP POST with body → response         |
| `net.status(resp)`  | Status code as number                    |

### `wolfhead` / `os`
| Method                   | Description                               |
|--------------------------|-------------------------------------------|
| `wh.workspace()`         | Current workspace object                  |
| `wh.workspace(n)`        | Switch to workspace n                     |
| `wh.context()`           | List all contexts                         |
| `wh.context("Work")`     | Switch to named context                   |
| `wh.notify(msg)`         | Send OS notification                      |
| `wh.notify(title, msg)`  | Notification with title                   |
| `wh.launch(appName)`     | Launch an application                     |
| `wh.gesture(name)`       | Trigger a gesture                         |
| `wh.platform()`          | Returns `"WolfHead/Linux"`                |

### `physics`
| Method                    | Description                     |
|---------------------------|---------------------------------|
| `phys.gravity(n)`         | Set gravity value               |
| `phys.collision(a, b)`    | True if a and b collide         |

---

## The REPL

Run `go run ./cmd/repl` from the project root to start the interactive shell.

```
  → set name to "Hunter"
  → say "Hello, " + name
  Hello, Hunter
  → create Dog [ has: [name is "Rex"] can: [bark [say "Woof!"]] ]
  → tell Dog to bark
  Woof!
  → :state     ← show all variables and objects
  → :clear     ← reset everything
  → :run game.he  ← run a file in current state
  → :quit
```

---

## Running a File

```bash
go run ./cmd/hunterlang -file myprogram.he
```

---

## Reserved Words

```
and          ask         as          be          becomes
between         as          be          becomes
can          carries     change      check       collision
create       done        each        else        false
for          frames      give        grow        has
how          if          image       in          is
know         let         like        list        make
music        named       no          not         nothing
on           or          owns        print       remembers
repeat       return      say         seconds     set
shader       show        shrink      sound       stop
summon       tell        then        times       to          try         until
true        upto        video       wait        when        whenever
while       with        yes
```

---

## Grammar (simplified EBNF)

```ebnf
Program      = Line* EOF
Line         = SummonLine | ObjectLine | Statement

SummonLine   = "summon" STRING ("as" | "named") IDENT
ObjectLine   = ("create" | "make") IDENT ["like" IDENT] ("{" | "[") ObjectBody ("}" | "]")

ObjectBody   = Section*
Section      = PropsSection | AbilSection | ReactSection

PropsSection = ("has" | "owns" | "carries") ":" "[" Property* "]"
Property     = IDENT ("is" | "starts as" | "becomes") Expr ","?

AbilSection  = ("can" | "know how to") ":" "[" Action* "]"
Action       = IDENT ["(" Params ")"] ["returns" Type] "[" Statement* "]"

ReactSection = ("on" | "when" | "whenever") Trigger "[" Statement* "]"
Trigger      = IDENT ["with" (IDENT | STRING)]

Statement    = SayStmt | SetStmt | ChangeStmt | GrowStmt | ShrinkStmt
             | IfStmt  | RepeatStmt | ForEachStmt | WhileStmt
             | TellStmt | WaitStmt | ReturnStmt | AskStmt
             | DotAssignStmt | GiveStmt | ExprStmt

Expr         = LogicOr
LogicOr      = LogicAnd ("or" LogicAnd)*
LogicAnd     = Comparison ("and" Comparison)*
Comparison   = Arithmetic (Op Arithmetic)?
Op           = "==" | "!=" | ">" | "<" | ">=" | "<=" | "is" | "is not"
Arithmetic   = Term (("+"|"-") Term)*
Term         = Factor (("*"|"/") Factor)*
Factor       = Unary ("**" Unary)?
Unary        = ("-" | "!" | "not") Unary | Primary
Primary      = NUMBER | STRING | BOOLEAN | "nothing"
             | IDENT ["." IDENT ["(" Args ")"]]
             | IDENT "(" Args ")"
             | "(" Expr ")"
             | "[" (NamedArgs | ArrayItems) "]"
```
