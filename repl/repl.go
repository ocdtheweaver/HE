// HE REPL — interactive shell for the HE language.
// Run: go run ./cmd/repl
//
// Maintains state (variables, objects) across inputs.
// Detects open blocks and waits for them to close before executing.
package repl

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"hunterlang/lang/eval"
	"hunterlang/lang/lexer"
	"hunterlang/lang/parser"
)

const banner = `
  ██╗  ██╗███████╗
  ██║  ██║██╔════╝
  ███████║█████╗
  ██╔══██║██╔══╝
  ██║  ██║███████╗
  ╚═╝  ╚═╝╚══════╝  interactive

  Type HE code and press Enter.  :help for guidance.  :quit to exit.
`

const helpText = `
── HE Quick Reference ─────────────────────────────────────

  Output         say "Hello, world"
  Variables      set name to "Hunter"   /   let x be 42
                 change x to 100        /   x becomes 99
  Grow / Shrink  grow score by 10       /   shrink hp by 5
  Conditions     if x > 5 then [ ... ] else [ ... ]
                 check if done is not true then [ ... ]
  Loops          repeat 3 times [ ... ]
                 repeat while x < 10 [ grow x by 1 ]
                 for each item in list [ say item ]
  Objects        create Dog [ has: [name is "Rex"] can: [bark [say "Woof!"]] ]
  Call           tell Dog to bark
  Dot access     say Dog.name
  Dot assign     set Dog.name to "Max"   /   give 100 to Dog.score
  Grow field     grow Dog.score by 10
  Inherit        create Puppy like Dog [ has: [age is 1] ]
  Input          ask "Your name?" as myName
  Modules        summon "math" as m      →  m.sqrt(16), m.abs(-3), m.round(3.7)
                 summon "text" as t      →  t.upper("hi"), t.length("hello")
                 summon "list" as lst    →  lst.length(x), lst.add(x, item)
                 summon "io"   as io     →  io.read("file.txt"), io.write("f","content")
  Comments       ~ anything between tildes is a comment ~

── REPL Commands ──────────────────────────────────────────

  :help           show this guide
  :clear          reset all state (variables, objects)
  :state          show current variables and objects
  :run <file>     run a .he file in current state
  :quit  :exit    leave the REPL
───────────────────────────────────────────────────────────
`

func Start() {
	fmt.Print(banner)

	interp := eval.NewInterpreter()
	// Boot with empty program so runtime is initialised
	if err := interp.Run(nil); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	scanner := bufio.NewScanner(os.Stdin)
	var buf []string // multiline accumulator

	for {
		if len(buf) == 0 {
			fmt.Print("\n  → ")
		} else {
			fmt.Printf("  %s… ", strings.Repeat("  ", depthLevel(buf)))
		}

		if !scanner.Scan() {
			fmt.Println("\n  Goodbye.")
			break
		}

		line := scanner.Text()
		trimmed := strings.TrimSpace(line)

		// REPL commands (only at fresh prompt)
		if len(buf) == 0 && strings.HasPrefix(trimmed, ":") {
			switch {
			case trimmed == ":quit" || trimmed == ":exit":
				fmt.Println("  Goodbye.")
				return

			case trimmed == ":help":
				fmt.Print(helpText)
				continue

			case trimmed == ":clear":
				interp = eval.NewInterpreter()
				if err := interp.Run(nil); err == nil {
					fmt.Println("  State cleared.")
				}
				continue

			case trimmed == ":state":
				printState(interp)
				continue

			case strings.HasPrefix(trimmed, ":run "):
				path := strings.TrimSpace(strings.TrimPrefix(trimmed, ":run "))
				if err := runFile(interp, path); err != nil {
					fmt.Printf("  ✗ %s\n", err)
				} else {
					fmt.Printf("  ✓ ran %s\n", path)
				}
				continue

			default:
				fmt.Printf("  Unknown command %q — type :help\n", trimmed)
				continue
			}
		}

		// Skip blank lines at fresh prompt
		if len(buf) == 0 && trimmed == "" {
			continue
		}

		buf = append(buf, line)
		src := strings.Join(buf, "\n")

		// Check bracket balance — wait for open blocks to close
		if openBrackets(src) > 0 {
			continue
		}

		// Execute
		if err := execSnippet(interp, src); err != nil {
			fmt.Printf("  ✗ %s\n", friendlyErr(err.Error()))
		}
		buf = nil
	}
}

// ── Helpers ───────────────────────────────────────────────────────────────────

func execSnippet(interp *eval.Interpreter, src string) error {
	lx := lexer.New(src)
	p := parser.New(lx)
	prog, err := p.ParseProgram()
	if err != nil {
		return err
	}
	return interp.RunProg(prog)
}

func runFile(interp *eval.Interpreter, path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("can't read %q: %v", path, err)
	}
	return execSnippet(interp, string(data))
}

func printState(interp *eval.Interpreter) {
	snap := interp.Snapshot()
	skip := map[string]bool{"platform": true, "nothing": true, "pi": true}

	vars := []string{}
	for k, v := range snap.Env {
		if !skip[k] {
			vars = append(vars, fmt.Sprintf("    %-16s = %s", k, v.String()))
		}
	}

	fmt.Println("\n  ── Variables ──")
	if len(vars) == 0 {
		fmt.Println("    (none)")
	} else {
		for _, v := range vars {
			fmt.Println(v)
		}
	}

	fmt.Println("\n  ── Objects ──")
	if len(snap.Objects) == 0 {
		fmt.Println("    (none)")
	} else {
		for name, obj := range snap.Objects {
			fmt.Printf("    %s\n", name)
			for k, v := range obj.Fields {
				fmt.Printf("      .%-14s = %s\n", k, v.String())
			}
		}
	}
	fmt.Println()
}

func openBrackets(src string) int {
	// Count unmatched [ — ignoring those inside strings
	open := 0
	inStr := false
	for _, ch := range src {
		switch {
		case ch == '"':
			inStr = !inStr
		case inStr:
			continue
		case ch == '[':
			open++
		case ch == ']':
			open--
		}
	}
	if open < 0 {
		return 0
	}
	return open
}

func depthLevel(buf []string) int {
	src := strings.Join(buf, "\n")
	return openBrackets(src)
}

func friendlyErr(msg string) string {
	return strings.ReplaceAll(msg, "token.TokenType", "")
}
