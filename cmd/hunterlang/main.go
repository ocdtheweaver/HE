package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"hunterlang/compiler"
	"hunterlang/lang/ast"
	"hunterlang/lang/eval"
	"hunterlang/lang/lexer"
	"hunterlang/lang/parser"
	"hunterlang/repl"
)

const version = "5.0.0"

const helpText = `HE — a near-English programming language

Usage:
  he <command> [file] [flags]

Commands:
  run     <file>   Run a .he program
  check   <file>   Show type inference report
  build   <file>   Compile to bytecode and show disassembly
                   --o <out.hex>  save compiled binary (standard format)
                   --o <out.hbc>  save compiled binary (alias of .hex)
  vm      <file>   Run via bytecode VM (accepts .he, .hex, or .hbc)
  bench   <file>   Benchmark interpreter vs VM
  repl             Start the interactive shell
  new     <name>   Scaffold a new HE project
  version          Show version information
  help             Show this help

Examples:
  he run game.he
  he build game.he --o game.hex
  he run game.hex
  he vm game.hex
  he check myapp.he
  he build myapp.he
  he vm myapp.he
  he bench myapp.he
  he repl
  he new myproject

Flags (for run, check, build, vm, bench):
  --event <name>   Trigger a named reaction after running
`

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		// Bare `he` with no arguments — try genesis.he in the current
		// directory before falling back to the help text.
		if _, err := os.Stat("genesis.he"); err == nil {
			cmdRun([]string{"genesis.he"})
			return
		}
		fmt.Print(helpText)
		os.Exit(0)
	}

	cmd := args[0]
	rest := args[1:]

	switch cmd {
	case "run":
		cmdRun(rest)
	case "check":
		cmdCheck(rest)
	case "build":
		cmdBuild(rest)
	case "vm":
		cmdVM(rest)
	case "bench":
		cmdBench(rest)
	case "repl":
		cmdRepl()
	case "new":
		cmdNew(rest)
	case "version", "--version", "-v":
		cmdVersion()
	case "help", "--help", "-h":
		if len(rest) > 0 {
			cmdHelpFor(rest[0])
		} else {
			fmt.Print(helpText)
		}
	default:
		// If the argument looks like a .he file, treat as implicit "run"
		if strings.HasSuffix(cmd, ".he") || strings.HasSuffix(cmd, ".hex") || strings.HasSuffix(cmd, ".hbc") {
			cmdRun(args)
			return
		}
		fmt.Fprintf(os.Stderr, "  ✗ unknown command %q\n\nRun 'he help' for usage.\n", cmd)
		os.Exit(1)
	}
}

// ── Commands ──────────────────────────────────────────────────────────────────

func cmdRun(args []string) {
	file, flags := parseFileArgs(args, "run")
	baseDir := filepath.Dir(file)

	if strings.HasSuffix(file, ".hex") || strings.HasSuffix(file, ".hbc") {
		chunk := mustLoadHBC(file)
		_, dur, err := compiler.RunTimed(chunk)
		if err != nil {
			die(err.Error())
		}
		fmt.Fprintf(os.Stderr, "\n  ran %s in %v\n", file, dur)
		return
	}

	prog := mustParse(file)

	if event := flags["event"]; event != "" {
		interp := eval.NewInterpreter()
		interp.SetBaseDir(baseDir)
		interp.SetEvent(event, nil)
		if err := interp.Run(prog); err != nil {
			die(err.Error())
		}
		return
	}

	interp := eval.NewInterpreter()
	interp.SetBaseDir(baseDir)
	if err := interp.Run(prog); err != nil {
		die(err.Error())
	}
}

func cmdCheck(args []string) {
	file, _ := parseFileArgs(args, "check")
	prog := mustParse(file)

	inf := compiler.NewInferencer()
	_, errs := inf.InferProgram(prog)
	fmt.Print(inf.Report())
	if len(errs) > 0 {
		os.Exit(1)
	}
}

func cmdBuild(args []string) {
	file, flags := parseFileArgs(args, "build")
	prog := mustParse(file)

	bc := compiler.NewBytecodeCompiler()
	bc.BaseDir = filepath.Dir(file)
	chunk, errs := bc.Compile(prog)

	if outPath := flags["o"]; outPath != "" {
		data, err := compiler.EncodeChunk(chunk)
		if err != nil {
			die(fmt.Sprintf("can't encode bytecode: %v", err))
		}
		if err := os.WriteFile(outPath, data, 0644); err != nil {
			die(fmt.Sprintf("can't write %q: %v", outPath, err))
		}
		fmt.Printf("  ✓ wrote %s (%d bytes)\n", outPath, len(data))
		if len(errs) > 0 {
			fmt.Fprintln(os.Stderr)
			for _, e := range errs {
				fmt.Fprintf(os.Stderr, "  ⚠ %s\n", e)
			}
		}
		return
	}

	fmt.Print(chunk.Disassemble())
	if len(errs) > 0 {
		fmt.Fprintln(os.Stderr)
		for _, e := range errs {
			fmt.Fprintf(os.Stderr, "  ✗ %s\n", e)
		}
		os.Exit(1)
	}
}

func cmdVM(args []string) {
	file, _ := parseFileArgs(args, "vm")

	if strings.HasSuffix(file, ".hex") || strings.HasSuffix(file, ".hbc") {
		chunk := mustLoadHBC(file)
		_, dur, err := compiler.RunTimed(chunk)
		if err != nil {
			fmt.Fprintf(os.Stderr, "\n  ✗ VM error: %s\n\n", err)
			os.Exit(1)
		}
		fmt.Fprintf(os.Stderr, "\n  VM completed in %v\n", dur)
		return
	}

	prog := mustParse(file)

	bc := compiler.NewBytecodeCompiler()
	bc.BaseDir = filepath.Dir(file)
	chunk, compErrs := bc.Compile(prog)
	if len(compErrs) > 0 {
		fmt.Fprintln(os.Stderr, "  Compile warnings:")
		for _, e := range compErrs {
			fmt.Fprintf(os.Stderr, "    %s\n", e)
		}
	}
	_, dur, err := compiler.RunTimed(chunk)
	if err != nil {
		fmt.Fprintf(os.Stderr, "\n  ✗ VM error: %s\n\n", err)
		os.Exit(1)
	}
	fmt.Fprintf(os.Stderr, "\n  VM completed in %v\n", dur)
}

func mustLoadHBC(path string) *compiler.Chunk {
	data, err := os.ReadFile(path)
	if err != nil {
		die(fmt.Sprintf("can't read %q: %v", path, err))
	}
	chunk, err := compiler.DecodeChunk(data)
	if err != nil {
		die(fmt.Sprintf("can't load %q: %v", path, err))
	}
	return chunk
}

func cmdBench(args []string) {
	file, _ := parseFileArgs(args, "bench")
	prog := mustParse(file)

	fmt.Printf("  Benchmarking: %s\n\n", file)

	// Interpreter
	start := time.Now()
	interp := eval.NewInterpreter()
	_ = interp.Run(prog)
	interpDur := time.Since(start)

	// VM
	bc := compiler.NewBytecodeCompiler()
	bc.BaseDir = filepath.Dir(file)
	chunk, _ := bc.Compile(prog)
	_, vmDur, _ := compiler.RunTimed(chunk)

	fmt.Printf("  Interpreter : %v\n", interpDur)
	fmt.Printf("  VM          : %v\n", vmDur)
	if vmDur > 0 && interpDur > 0 {
		ratio := float64(interpDur) / float64(vmDur)
		if ratio >= 1 {
			fmt.Printf("  VM is %.1fx faster\n", ratio)
		} else {
			fmt.Printf("  Interpreter is %.1fx faster\n", 1/ratio)
		}
	}
}

func cmdRepl() {
	repl.Start()
}

func cmdNew(args []string) {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "  Usage: he new <project-name>")
		os.Exit(1)
	}
	name := args[0]
	if err := scaffold(name); err != nil {
		die(err.Error())
	}
	fmt.Printf("  ✓ Created project %q\n\n", name)
	fmt.Printf("    %s/\n", name)
	fmt.Printf("    ├── main.he\n")
	fmt.Printf("    ├── lib/\n")
	fmt.Printf("    │   └── helpers.he\n")
	fmt.Printf("    └── README.md\n\n")
	fmt.Printf("  Run it:  he run %s/main.he\n\n", name)
}

func cmdVersion() {
	fmt.Printf("HE language  v%s\n", version)
	fmt.Printf("Build:       Pass 14 — protection enforcement\n")
	fmt.Printf("Runtime:     tree-walk interpreter + bytecode VM\n")
	fmt.Printf("Compiler:    type inferencer + bytecode IR + .hex bundles\n")
	fmt.Printf("Protection:  #protected[N] tags enforced on all execution paths\n")
}

func cmdHelpFor(cmd string) {
	helps := map[string]string{
		"run":     "he run <file.he> [--event <name>]\n\n  Runs a HE program file.\n  Use --event to trigger a named reaction after the program executes.",
		"check":   "he check <file.he>\n\n  Runs type inference on the program and prints a report.\n  Shows inferred types for all variables and object fields.",
		"build":   "he build <file.he>\n\n  Compiles the program to bytecode and prints the disassembly.\n  Useful for understanding how HE programs compile.",
		"vm":      "he vm <file.he>\n\n  Runs the program through the bytecode VM instead of the interpreter.\n  Experimental — object method calls fall back to nil.",
		"bench":   "he bench <file.he>\n\n  Runs the program through both the interpreter and VM,\n  then prints timing and relative speedup.",
		"repl":    "he repl\n\n  Starts the interactive HE shell.\n  Type :help inside the REPL for guidance.",
		"new":     "he new <project-name>\n\n  Scaffolds a new HE project with a main.he, a lib/ folder,\n  and a README.md.",
		"version": "he version\n\n  Prints version information.",
	}
	if h, ok := helps[cmd]; ok {
		fmt.Println(h)
	} else {
		fmt.Printf("  No help available for %q\n\nRun 'he help' for all commands.\n", cmd)
	}
}

// ── Project scaffolding ───────────────────────────────────────────────────────

func scaffold(name string) error {
	dirs := []string{
		name,
		filepath.Join(name, "lib"),
	}
	for _, d := range dirs {
		if err := os.MkdirAll(d, 0755); err != nil {
			return err
		}
	}

	mainHE := fmt.Sprintf(`~ %s — main entry point ~

~ Import your helpers ~
summon "lib/helpers.he" as helpers

~ Your program starts here ~
say "Hello from {helpers.appName}!"
say "Version {helpers.version}"

~ Write your program below ~
`, name)

	helpersHE := fmt.Sprintf(`~ %s/lib/helpers.he — shared utilities ~

set appName to "%s"
set version to "0.1.0"

create Utils [
  can: [
    greet(name) [
      return "Hello, {name}!"
    ]
  ]
]
`, name, name)

	readme := fmt.Sprintf(`# %s

A HE language project.

## Running

`+"```"+`
he run main.he
`+"```"+`

## Structure

- `+"`main.he`"+` — entry point
- `+"`lib/helpers.he`"+` — shared utilities

## Learn More

- [HE for Dummies](docs/HE_FOR_DUMMIES.md)
- [HE Professional Reference](docs/PROFESSIONAL.md)
`, name)

	files := map[string]string{
		filepath.Join(name, "main.he"):          mainHE,
		filepath.Join(name, "lib", "helpers.he"): helpersHE,
		filepath.Join(name, "README.md"):         readme,
	}
	for path, content := range files {
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			return err
		}
	}
	return nil
}

// ── Helpers ───────────────────────────────────────────────────────────────────

func parseFileArgs(args []string, cmd string) (string, map[string]string) {
	file := ""
	flags := map[string]string{}

	for i := 0; i < len(args); i++ {
		a := args[i]
		if strings.HasPrefix(a, "--") {
			key := strings.TrimPrefix(a, "--")
			if i+1 < len(args) && !strings.HasPrefix(args[i+1], "--") {
				flags[key] = args[i+1]
				i++
			} else {
				flags[key] = "true"
			}
		} else if file == "" {
			file = a
		}
	}

	// No file given — fall back to genesis.he in the current directory,
	// the same role index.html / App.tsx / main() plays elsewhere.
	if file == "" {
		if _, err := os.Stat("genesis.he"); err == nil {
			fmt.Printf("  → no file given, using genesis.he\n")
			return "genesis.he", flags
		}
		fmt.Fprintf(os.Stderr, "  Usage: he %s <file.he>\n", cmd)
		fmt.Fprintf(os.Stderr, "  (or create a genesis.he in this directory to omit the filename)\n")
		os.Exit(1)
	}

	return file, flags
}

func mustParse(file string) *ast.Program {
	b, err := os.ReadFile(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "  ✗ Can't read %q: %v\n", file, err)
		os.Exit(1)
	}
	if len(b) >= 3 && b[0] == 0xEF && b[1] == 0xBB && b[2] == 0xBF {
		b = b[3:]
	}
	lx := lexer.New(string(b))
	p := parser.New(lx)
	prog, err := p.ParseProgram()
	if err != nil {
		fmt.Fprintf(os.Stderr, "\n  ✗ %s\n\n", friendly(err.Error()))
		os.Exit(1)
	}
	return prog
}

var _ = (*ast.Program)(nil) // keep import

func die(msg string) {
	fmt.Fprintf(os.Stderr, "\n  ✗ %s\n\n", friendly(msg))
	os.Exit(1)
}

func friendly(msg string) string {
	return strings.ReplaceAll(msg, "token.TokenType", "")
}
