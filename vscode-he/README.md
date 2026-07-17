# HE Language — VS Code Extension

**HE** is a near-English programming language. If you can write a sentence, you can write HE.

```he
~ This is a comment ~
set name to "Hunter"
say "Hello, {name}!"

create Dog [
  has: [ name is "Rex" ]
  can: [
    bark [ say "{name} says Woof!" ]
  ]
]

tell Dog to bark
```

---

## Features

### Syntax Highlighting
Full highlighting for every HE construct — tilde comments, `[`/`]` blocks, string interpolation `{expr}`, `#protected[N]` tags, keywords, and operators.

### Snippets
Type a prefix and press Tab to expand:

| Prefix | Expands to |
|--------|-----------|
| `create` | Full object with `has:` and `can:` |
| `createp` | Protected object with `#protectedN` tag |
| `ability` | Named ability with parameters |
| `ab` | Anonymous ability (closure) |
| `summon` | `summon "file.he" as alias` |
| `summonf` | Flat summon (no alias) |
| `try` | `try [...] or (err) [...]` |
| `foreach` | `for each item in list [...]` |
| `each` | `each i from 1 to 10 [...]` |
| `if` | `if condition then [...]` |
| `ife` | `if/then/else` |
| `repeati` | `repeat N times as i [...]` |
| `tell` | `tell Object to action` |
| `genesis` | Full `genesis.he` project template |
| `policyblock` | Policy block for `protected.he` |

### Run Commands

| Command | Keyboard | What it does |
|---------|----------|-------------|
| **HE: Run File** | `Ctrl+F5` / `Cmd+F5` | Runs the current `.he` file |
| **HE: Check File** | `Ctrl+Shift+K` / `Cmd+Shift+K` | Shows type inference and protection tag report |
| **HE: Build to .hex** | Right-click menu | Compiles to a self-contained `.hex` bundle |
| **HE: Open REPL** | Command palette | Opens the interactive HE shell |
| **HE: Run genesis.he** | Command palette | Runs the project entry point (no filename needed) |

Right-click any `.he` file in the Explorer or editor for run/check/build options.

### Check on Save
When you save a `.he` file, `he check` runs automatically and any type warnings appear in VS Code's Problems panel. Disable in settings if you prefer manual checking.

### Bracket Matching
`[`, `]`, `(`, `)` are all matched. Tilde comments (`~ ... ~`) are recognized as a comment pair.

### Auto-Close
- Type `[` → get `[]`
- Type `(` → get `()`
- Type `"` → get `""`
- Type `~` → get `~ ~` (comment template)

---

## Requirements

The `he` CLI must be installed and on your `PATH`.

**Build from source:**
```bash
cd HE
go build -o he ./cmd/hunterlang
# Move to PATH:
mv he /usr/local/bin/       # Linux/Mac
# or add the folder to PATH on Windows
```

---

## Extension Settings

| Setting | Default | Description |
|---------|---------|-------------|
| `he.executablePath` | `"he"` | Path to the `he` binary |
| `he.checkOnSave` | `true` | Run `he check` automatically on save |
| `he.buildOutputDir` | `"."` | Where `he build` puts the `.hex` file |

---

## HE Language Basics

```he
~ Variables ~
set score to 0
grow score by 10
set name to "Hunter"

~ Conditions ~
if score > 100 then [
  say "High score!"
] else [
  say "Keep going."
]

~ Loops ~
each i from 1 to 10 [
  say "Step {i}"
]

~ Objects ~
create Cat like Dog [
  has: [ name is "Whiskers" ]
  can: [
    speak [ say "{name} says Meow!" ]
  ]
]
tell Cat to speak

~ Error handling ~
try [
  set data to io.read("config.txt")
] or (err) [
  say "Couldn't load config: {err}"
]

~ Protection tags ~
create Premium #protected1 [
  can: [
    watch4K(title) [ say "Streaming {title} in 4K" ]
  ]
]
try [
  tell Premium to watch4K with "The Last Signal"
] or (err) [
  say "Upgrade to Premium: {err}"
]
```

---

## Protection Tags

Mark features as requiring entitlement with `#protected[N]` tags. Define what each tag means in `protected.he`:

```he
create protected1 [
  has: [
    planName is "Premium"
    mode is "cached"
    endpoint is "https://api.myapp.com/entitlement/premium"
    graceSeconds is 604800
  ]
]
```

`he check` reports all tagged declarations. Enforcement is handled at runtime.

---

## More Information

- [HE Docs](https://github.com/ocdtheweaver/HE/tree/main/docs)
- [HE for Beginners](https://github.com/ocdtheweaver/HE/blob/main/docs/HE_FOR_DUMMIES.md)
- [Report Issues](https://github.com/ocdtheweaver/HE/issues)

---

*Built for WolfHead OS and everywhere else.*

---

## Adding your icon

To add the HE icon before publishing:

1. Place your `icon.png` (128×128px, `#DA7756` orange) in `vscode-he/images/icon.png`
2. Add this back to `package.json`:
   ```json
   "icon": "images/icon.png"
   ```
3. Run `npx vsce package` to rebuild the `.vsix`
