// Command he — Hunter's Engine interpreter
//
// Usage:
//
//	he <file.he>          run a HE source file
//	he                    start the interactive REPL
//	he --tokens <file>    dump the token stream (debug)
//	he --ast    <file>    dump the AST (debug)
package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/hunter/he/pkg/lexer"
	"github.com/hunter/he/pkg/parser"
	"github.com/hunter/he/pkg/resolver"
	"github.com/hunter/he/pkg/runtime"
)

const version = "0.1.0"

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		runREPL()
		return
	}

	switch args[0] {
	case "--version", "-v":
		fmt.Printf("he (Hunter's Engine) interpreter v%s\n", version)
	case "--tokens":
		if len(args) < 2 {
			die("--tokens requires a file argument")
		}
		debugTokens(args[1])
	case "--ast":
		if len(args) < 2 {
			die("--ast requires a file argument")
		}
		debugAST(args[1])
	case "--help", "-h":
		printHelp()
	default:
		runFile(args[0])
	}
}

// ─── run a file ───────────────────────────────────────────────────────────────

func runFile(path string) {
	src, err := os.ReadFile(path)
	if err != nil {
		die("cannot open file: %v", err)
	}
	if exitCode := execute(string(src), path); exitCode != 0 {
		os.Exit(exitCode)
	}
}

// execute lexes, parses, resolves, and runs source. Returns 0 on success.
func execute(src, filename string) int {
	// Lex
	l := lexer.New(src)
	tokens := l.Tokenise()
	if errs := l.Errors(); len(errs) > 0 {
		for _, e := range errs {
			fmt.Fprintf(os.Stderr, "[%s] %s\n", filename, e.Error())
		}
		return 1
	}

	// Parse
	p := parser.New(tokens)
	prog := p.Parse()
	if errs := p.Errors(); len(errs) > 0 {
		for _, e := range errs {
			fmt.Fprintf(os.Stderr, "[%s] %s\n", filename, e.Error())
		}
		return 1
	}

	// Resolve
	r := resolver.New()
	r.Resolve(prog)
	if errs := r.Errors(); len(errs) > 0 {
		for _, e := range errs {
			fmt.Fprintf(os.Stderr, "[%s] %s\n", filename, e.Error())
		}
		return 1
	}

	// Run
	interp := runtime.New()
	if err := interp.Run(prog); err != nil {
		fmt.Fprintf(os.Stderr, "[%s] %s\n", filename, err.Error())
		return 1
	}
	return 0
}

// ─── REPL ─────────────────────────────────────────────────────────────────────

func runREPL() {
	fmt.Printf("Hunter's Engine v%s — interactive REPL\n", version)
	fmt.Println("Type 'exit' or 'quit' to leave.  Type 'help' for language tips.")
	fmt.Println()

	scanner := bufio.NewScanner(os.Stdin)
	interp := runtime.New()
	// We share one interpreter across REPL iterations so variables persist.
	_ = interp // used via execute-with-shared-state below

	// For the REPL we accumulate multi-line input until the user submits a blank line
	// or a complete statement.
	var buf []string
	for {
		if len(buf) == 0 {
			fmt.Print("he> ")
		} else {
			fmt.Print("... ")
		}

		if !scanner.Scan() {
			break
		}
		line := scanner.Text()

		switch strings.TrimSpace(line) {
		case "exit", "quit":
			fmt.Println("Goodbye.")
			return
		case "help":
			printLangHelp()
			continue
		}

		buf = append(buf, line)
		src := strings.Join(buf, "\n")

		// Try to parse; if there are only "expected end/else" errors we're
		// likely waiting for more input — keep buffering.
		p := parser.New(lexer.New(src).Tokenise())
		prog := p.Parse()
		errs := p.Errors()

		incomplete := false
		for _, e := range errs {
			if strings.Contains(e.Message, "expected 'end'") ||
				strings.Contains(e.Message, "expected 'else'") {
				incomplete = true
				break
			}
		}

		if incomplete {
			continue
		}

		// Reset buffer.
		buf = nil

		if len(errs) > 0 {
			for _, e := range errs {
				fmt.Fprintln(os.Stderr, e.Error())
			}
			continue
		}

		r := resolver.New()
		r.Resolve(prog)
		for _, e := range r.Errors() {
			fmt.Fprintln(os.Stderr, e.Error())
		}
		if len(r.Errors()) > 0 {
			continue
		}

		if err := interp.Run(prog); err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
		}
	}
}

// ─── debug modes ─────────────────────────────────────────────────────────────

func debugTokens(path string) {
	src, err := os.ReadFile(path)
	if err != nil {
		die("cannot open file: %v", err)
	}
	l := lexer.New(string(src))
	tokens := l.Tokenise()
	for _, tok := range tokens {
		fmt.Println(tok)
	}
}

func debugAST(path string) {
	src, err := os.ReadFile(path)
	if err != nil {
		die("cannot open file: %v", err)
	}
	l := lexer.New(string(src))
	p := parser.New(l.Tokenise())
	prog := p.Parse()
	for _, s := range prog.Statements {
		fmt.Println(s)
	}
}

// ─── helpers ──────────────────────────────────────────────────────────────────

func die(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "he: "+format+"\n", args...)
	os.Exit(1)
}

func printHelp() {
	fmt.Print(`he — Hunter's Engine interpreter

Usage:
  he [file.he]         run a .he source file (omit for REPL)
  he --tokens <file>   dump token stream
  he --ast    <file>   dump AST
  he --version         print version
  he --help            show this message
`)
}

func printLangHelp() {
	fmt.Print(`
Hunter's Engine (HE) — quick reference
───────────────────────────────────────
Variables       x is 5
                set name to "Alice"
Print           print "Hello, " + name
Input           input name as "Enter your name: "
If              if x > 0 then
                    print "positive"
                else
                    print "not positive"
                end
While           while x > 0 do
                    set x to x - 1
                end
For-each        for each item in myList do
                    print item
                end
Classes         class Dog extends Animal
                    property name
                    property breed
                    method bark
                        print this.name + " says woof!"
                    end
                end
Objects         object rex of Dog with name is "Rex", breed is "Lab"
                call bark on rex
New (inline)    set d to new Dog with name is "Buddy"
Lists           set nums to [1, 2, 3]
                append(nums, 4)
Builtins        len, toNumber, toString, typeOf,
                append, remove, contains, sqrt, abs,
                floor, ceil, upper, lower, trim, split, join
Comments        -- this is a comment
`)
}
