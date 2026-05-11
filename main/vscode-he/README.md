# Hunterlang (HE) VSCode Extension

Provides:
- Syntax configuration for `.he` files
- Commands:
  - **HE: Run file** (`hunterlang.runHe`)
  - **HE: Debug file** (`hunterlang.debugHe`)

The commands currently execute the language via:
`go run ./cmd/hunterlang -file <activeFile>`

## Build

From `vscode-he/`:
- `npm install`
- `npm run compile`
- Output goes to `out/extension.js`
