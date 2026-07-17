# HE for Dummies
## The Complete Beginner's Guide to Programming in Plain English

*If you've never written code before, this guide is for you.*
*If you have — it's still a quick read.*

---

## What is HE?

HE is a programming language that reads like English. Instead of:

```python
if x >= 18 and x <= 65:
    print("Working age")
```

You write:

```
if age is between 18 and 65 then [
  say "Working age"
]
```

That's it. Real sentences. Real words. No semicolons, no curly braces, no symbols you've never seen before.

---

## Getting Started

### Running your first program

1. Create a file called `hello.he`
2. Type this inside it:
   ```
   say "Hello, world!"
   ```
3. Run it:
   ```
   he run hello.he
   ```

You'll see:
```
Hello, world!
```

Congratulations. You just wrote a program.

### Starting a new project

```
he new myproject
```

This creates:
```
myproject/
├── main.he
├── lib/
│   └── helpers.he
└── README.md
```

Run it with `he run myproject/main.he`.

### The `he` command

| Command | What it does |
|---------|---------------|
| `he run file.he` | Run a program |
| `he repl` | Start an interactive shell |
| `he new myproject` | Create a new project |
| `he check file.he` | Check what types your variables are |
| `he version` | Show version info |
| `he help` | Show all commands |
| `he build file.he --o file.hbc` | Compile to a fast-loading file |
| `he vm file.he` | Run via the experimental fast VM |
| `he bench file.he` | Compare speed: normal vs fast VM |

---

## Part 1: The Basics

### Saying things

Use `say` to print something to the screen.

```
say "Hello!"
say "I can say anything."
say 42
say true
```

`print` and `show` do the same thing. Use whichever feels natural.

---

### Variables — storing information

A variable is a named box that holds a value.

```
set name to "Hunter"
set age to 25
set ready to true
```

You can also write:
```
let name be "Hunter"
```

Both mean the same thing. Use whichever reads better for your sentence.

**Changing a variable:**
```
change name to "Alex"
```

**Growing and shrinking numbers:**
```
set score to 0
grow score by 10      ~ score is now 10 ~
shrink score by 3     ~ score is now 7 ~
```

---

### Comments — notes to yourself

Use tildes `~` to write notes that HE ignores:

```
~ This is a comment — HE won't run this ~
set speed to 100   ~ this sets speed to 100 ~
```

---

### String interpolation — putting values inside text

Instead of joining strings together, put variable names inside `{` `}`:

```
set name to "Hunter"
set score to 42
say "Hello, {name}! Your score is {score}."
```

Output: `Hello, Hunter! Your score is 42.`

---

### Types of values

| What you want to store | Example           | Called     |
|------------------------|-------------------|------------|
| A whole or decimal number | `42`, `3.14`   | number     |
| Words or sentences     | `"Hello there"`   | text       |
| Yes or no              | `true`, `false`   | boolean    |
| A collection of things | `[1, 2, 3]`      | list       |
| Nothing at all         | `nothing`         | nothing    |

`yes` means the same as `true`. `no` means the same as `false`.

---

## Part 2: Making Decisions

### If / then / else

```
set temp to 22

if temp > 30 then [
  say "It's hot!"
] else [
  say "Nice weather."
]
```

The `[` and `]` brackets group the instructions inside each branch.

**English alternatives:**
```
check if temp > 30 then [
  say "Hot!"
]
```

---

### Comparing things

| What you write          | What it means             |
|-------------------------|---------------------------|
| `x == 5`                | x equals 5                |
| `x is 5`                | x equals 5 (English form) |
| `x != 5`                | x does not equal 5        |
| `x is not 5`            | x does not equal 5        |
| `x > 5`                 | x is greater than 5       |
| `x < 5`                 | x is less than 5          |
| `x >= 5`                | x is 5 or more            |
| `x <= 5`                | x is 5 or less            |
| `x is between 1 and 10` | x is at least 1, at most 10 |

---

### Combining conditions

```
if age > 18 and ready then [
  say "You can proceed."
]

if tired or hungry then [
  say "Take a break."
]

if not finished then [
  say "Keep going."
]
```

---

## Part 3: Loops — doing things repeatedly

### Repeat a fixed number of times

```
repeat 3 times [
  say "Again!"
]
```

Output:
```
Again!
Again!
Again!
```

---

### Repeat while something is true

```
set lives to 3
repeat while lives > 0 [
  say "You have {lives} lives"
  shrink lives by 1
]
```

---

### Repeat until something is true

```
set count to 0
repeat until count == 10 [
  grow count by 1
]
say "Done! Count is {count}"
```

---

### For each — go through a list

```
set fruits to ["apple", "banana", "cherry"]
for each fruit in fruits [
  say "I have a {fruit}"
]
```

---

### Count from one number to another

```
for each i from 1 to 5 [
  say "Step {i}"
]
```

Short form (no `for` needed):
```
each i from 1 to 5 [
  say "{i}"
]
```

Count by twos:
```
each n from 2 to 10 step 2 [
  say "{n}"
]
```

---

## Part 4: Asking for Input

```
ask "What is your name?" as myName
say "Hello, {myName}!"

ask "How old are you?" as age
if age is between 13 and 17 then [
  say "You're a teenager."
]
```

---

## Part 5: Objects — grouping things together

An object bundles related data and actions into one named thing.

```
create Dog [
  has: [
    name is "Rex"
    breed is "Labrador"
    age is 3
  ]
  can: [
    bark [
      say "{name} says: Woof!"
    ]
    describe [
      say "I'm {name}, a {breed}, {age} years old."
    ]
  ]
]
```

**Using an object:**
```
tell Dog to bark
tell Dog to describe
```

Output:
```
Rex says: Woof!
I'm Rex, a Labrador, 3 years old.
```

**Reading a value:**
```
say Dog.name
say Dog.age
```

**Changing a value:**
```
set Dog.name to "Buddy"
give 4 to Dog.age
```

---

### Actions with inputs

```
create Calculator [
  can: [
    add(a, b) [
      return a + b
    ]
    multiply(a, b) [
      return a * b
    ]
  ]
]

set result to Calculator.add(5, 3)
say "5 + 3 = {result}"
```

---

### Reactions — responding to events

```
create Button [
  has: [
    label is "Click me"
    clicks starts as 0
  ]
  on click [
    grow clicks by 1
    say "{label} has been clicked {clicks} times"
  ]
]
```

---

### Inheritance — one object based on another

```
create Animal [
  has: [
    name is "Animal"
    sound is "..."
  ]
  can: [
    speak [
      say "{name} says {sound}"
    ]
  ]
]

create Cat like Animal [
  has: [
    name is "Whiskers"
    sound is "Meow"
  ]
]

tell Cat to speak
```

Output: `Whiskers says Meow`

---

## Part 6: Math

HE has a built-in math toolbox:

```
summon "math" as m

say m.sqrt(16)        ~ 4 ~
say m.abs(-5)         ~ 5 ~
say m.round(3.7)      ~ 4 ~
say m.max(10, 20)     ~ 20 ~
say m.min(10, 20)     ~ 10 ~
say m.random(1, 100)  ~ a random number between 1 and 100 ~
say m.pow(2, 10)      ~ 1024 ~
```

Arithmetic works directly too:
```
set x to 10 + 3      ~ 13 ~
set x to 10 - 3      ~ 7 ~
set x to 10 * 3      ~ 30 ~
set x to 10 / 3      ~ 3.333... ~
set x to 2 ** 8      ~ 256 ~
```

---

## Part 7: Working with Text

```
summon "text" as t

say t.upper("hello")           ~ HELLO ~
say t.lower("WORLD")           ~ world ~
say t.length("hello")          ~ 5 ~
say t.contains("hello", "ell") ~ yes ~
say t.replace("I love cats", "cats", "dogs")  ~ I love dogs ~
say t.trim("  spaces  ")       ~ spaces ~

set parts to t.split("one,two,three", ",")
for each p in parts [
  say p
]
```

---

## Part 8: Working with Lists

```
summon "list" as lst

set items to ["apple", "banana", "cherry"]

say lst.length(items)      ~ 3 ~
say lst.first(items)       ~ apple ~
say lst.last(items)        ~ cherry ~

set items to lst.add(items, "date")
say lst.length(items)      ~ 4 ~

say lst.contains(items, "banana")  ~ yes ~
```

---

## Part 9: Files

```
summon "io" as io

~ Write to a file ~
io.write("notes.txt", "Hello from HE!")

~ Read it back ~
set content to io.read("notes.txt")
say content

~ Append to a file ~
io.append("log.txt", "Something happened\n")

~ Check if a file exists ~
if io.exists("notes.txt") then [
  say "File found!"
]

~ Delete a file ~
io.delete("notes.txt")
```

---

## Part 10: Handling Errors

Sometimes things go wrong. HE lets you handle that gracefully:

```
summon "io" as io

try [
  set data to io.read("config.txt")
  say "Loaded: {data}"
] or (problem) [
  say "Couldn't load config: {problem}"
  say "Using defaults instead."
]
```

The `or (problem)` part runs only if something goes wrong. The word `problem` (you can call it anything) holds the error message.

---

## Part 11: Remembering Things Between Runs

```
set highscore to 0
set playerName to "Hunter"

~ Save to disk ~
remember highscore
remember playerName

~ Next time the program runs, load them back ~
recall highscore
recall playerName
say "Welcome back, {playerName}! Your best score: {highscore}"

~ Remove a saved value ~
forget highscore
```

---

## Part 12: Functions You Can Store

```
~ Create a reusable ability ~
set greet to ability(name) [
  return "Hello, {name}! Welcome."
]

say greet("Hunter")
say greet("Alex")
```

---

## Part 13: Time and Dates

```
summon "clock" as clock

set today to clock.today()
say "Today is {today}"

set now to clock.now()
say "It's {now.hour}:{now.minute} on {now.weekday}"

set formatted to clock.format(now, "DD/MM/YYYY")
say "Date: {formatted}"
```

---


## Part 14: Method Chains — Doing Things Directly on Values

You can call methods directly on any value without needing a module:

```
~ On text ~
set word to "hello world"
say word.upper()              ~ HELLO WORLD ~
say word.length()             ~ 11 ~
say word.contains("world")    ~ yes ~
say word.replace("world","HE") ~ hello HE ~
say word.split(" ").length()   ~ 2 ~

~ On a string literal directly ~
say "  spaces  ".trim()
say "one,two,three".split(",").first()

~ On numbers ~
set n to -42
say n.abs()      ~ 42 ~
say n.floor()    ~ -42 ~
say (16.0).sqrt() ~ 4 ~

~ On lists ~
set items to [3, 1, 4, 1, 5]
say items.length()       ~ 5 ~
say items.first()        ~ 3 ~
say items.reverse()      ~ [5, 1, 4, 1, 3] ~
say items.contains(4)    ~ yes ~
```

---

## Part 15: Checking Membership

Check if a value is in a list without a loop:

```
set color to "blue"
if color is one of ["red", "green", "blue"] then [
  say "It's a primary color"
]

set day to "Saturday"
if day is one of ["Saturday", "Sunday"] then [
  say "Weekend!"
]
```

---

## Part 16: Scoped Aliases — `with`

Bind a value to a short name for a block:

```
summon "math" as m

with m.pi() as pi [
  set circumference to 2 * pi * 5
  say "Circumference: {circumference}"
]

with "Hello, World!" as greeting [
  say greeting
  say greeting.upper()
  say greeting.length()
]
```

---

## Part 17: The `becomes` Shorthand

```
set score to 0

~ Instead of: set score to score + 1 ~
~ Or: change score to 50 ~
~ You can write: ~
score becomes 50
say "Score is now {score}"
```

---

## Part 18: Repeat with a Counter

```
~ Expose the loop counter as a variable ~
repeat 5 times as i [
  say "Round {i} of 5"
]
```

Output:
```
Round 1 of 5
Round 2 of 5
Round 3 of 5
Round 4 of 5
Round 5 of 5
```

---

## Part 19: Iterating Object Fields

```
create Settings [
  has: [
    theme is "dark"
    language is "en"
    notifications is true
  ]
]

for each key, val in Settings [
  say "{key} = {val}"
]
```

Output:
```
theme = dark
language = en
notifications = yes
```

---

## Part 20: Recursion

Abilities can call themselves:

```
create Math [
  can: [
    factorial(n) [
      if n <= 1 then [ return 1 ]
      return n * Math.factorial(n - 1)
    ]
  ]
]

say Math.factorial(5)   ~ 120 ~
say Math.factorial(10)  ~ 3628800 ~
```

HE automatically protects you if a recursion goes too deep — you'll get a clear error message instead of a crash.

---


## Part 21: Compiling Your Program

So far you've been running programs with `he run`. There's a faster way for programs you run a lot.

```
he build myapp.he --o myapp.hbc
```

This compiles your program into a single file (`myapp.hbc`) that starts up instantly — no re-reading or re-checking your code every time.

Run it directly:

```
he run myapp.hbc
```

Or run it through the experimental fast VM:

```
he vm myapp.hbc
```

You can compare speed between the normal way and the compiled way:

```
he bench myapp.he
```

This prints how long each approach takes, so you can see the difference for yourself.

---


## Part 22: Bringing In Other Files

You can split your program across multiple `.he` files and pull them together with `summon`.

### Just using it (no name needed)

If you're not going to spell out the module's name later, leave off `as`:

```
summon "math.he"

say sqrt(16)
say double(21)
```

Everything `math.he` defines becomes usable immediately, by name, with no prefix.

### Wanting more control (give it a name)

If you'd rather be explicit — especially for things like styling or UI, where you want `ui.something()` to read clearly — add `as`:

```
summon "UI.he" as ui

ui.nav("Home")
ui.button("Sign In")
```

You can do **both at once** — the name is never required, just an option:

```
summon "math.he" as maths

say sqrt(16)          ~ still works, no prefix needed ~
say maths.sqrt(16)     ~ also works, if you want to be explicit ~
```

### If two files define the same thing

```
summon "physics.he"
summon "ui.he"

render   ~ ✗ error — both files have something called "render" ~
```

HE won't guess which one you mean. Add names to tell them apart:

```
summon "physics.he" as phys
summon "ui.he" as ui

phys.render()
ui.render()
```

### Your project's starting point: `genesis.he`

Name your entry-point file `genesis.he` — the same role `index.html` plays for a website. Then you can just run:

```
he run
```

No filename needed — HE finds `genesis.he` automatically in the current folder.

---

## Part 23: Packaging Your App

Once your program works, you can bundle it into a single file — the HE equivalent of an `.exe`, `.apk`, or `.ipa`:

```
he build genesis.he --o myapp.hex
```

This produces `myapp.hex` — one file containing your whole program, including anything it `summon`s. Hand that file to someone else, and they can run it without needing any of your original `.he` files:

```
he run myapp.hex
```

---

## Quick Reference Card

```
~ OUTPUT ~
say "Hello"             print to screen
show "Hello"            same thing
print "Hello"           same thing

~ VARIABLES ~
set x to 5              create / update variable
let x be 5              same thing
change x to 10          same thing
grow x by 2             x = x + 2
shrink x by 1           x = x - 1
set x, y to 1, 2        set multiple at once

~ CONDITIONS ~
if x > 5 then [...]
if x is 5 then [...]
if x is not 5 then [...]
if x is between 1 and 10 then [...]
check if x > 0 then [...]

~ LOOPS ~
repeat 3 times [...]
repeat while x < 10 [...]
repeat until x == 10 [...]
for each item in list [...]
each i from 1 to 10 [...]
each n from 0 to 100 step 5 [...]

~ OBJECTS ~
create Name [...body...]
tell Name to action
tell Name to action with value
say Name.field
set Name.field to value
give value to Name.field

~ INPUT / OUTPUT ~
ask "Question?" as answer

~ MODULES ~
summon "math"    as m
summon "text"    as t
summon "list"    as lst
summon "io"      as io
summon "clock"   as clock
summon "net"     as net
summon "mylib.he" as lib

~ MUTATION ~
  score becomes 50         same as: set score to 50

~ MEMBERSHIP ~
  x is one of [a, b, c]    true if x equals any item

~ METHOD CHAINS ~
  word.upper()              on any text value
  items.length()            on any list value
  n.abs()                   on any number value

~ WITH SCOPE ~
  with expr as name [...]   bind expr to name for the block

~ FIELD ITERATION ~
  for each key, val in obj [...]

~ ERROR HANDLING ~
try [...] or (err) [...]
try [...] or if it fails [...]

~ PERSISTENCE ~
remember varName
recall varName
forget varName
```

---

## Common Mistakes

**Forgot `then` after a condition?**
```
~ Wrong ~
if x > 5 [
~ Right ~
if x > 5 then [
```

**Used `=` for comparison?**
```
~ Wrong — = isn't a thing in HE ~
if x = 5
~ Right ~
if x is 5
if x == 5
```

**Forgot to close a block?**
```
~ Every [ needs a matching ] ~
if x > 5 then [
  say "yes"
]            ~ ← don't forget this ~
```

**Tried to do math on text?**
```
~ Wrong ~
set age to "25"
grow age by 1        ~ won't work — "25" is text, not a number ~

~ Right ~
summon "text" as t
set age to t.number("25")   ~ convert first ~
grow age by 1
```

---

*That's everything you need to write real programs in HE.*
*Start small, experiment, and read errors carefully — they tell you exactly what's wrong.*
