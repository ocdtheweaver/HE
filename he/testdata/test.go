package main

import (
"fmt"
"io/ioutil"

"github.com/user/he/pkg/lexer"
"github.com/user/he/pkg/parser"
"github.com/user/he/pkg/runtime"
)

func main() {
// Read example file
data, err := ioutil.ReadFile("examples/game.he")
if err != nil {
fmt.Printf("Error reading file: %v\n", err)
return
}

source := string(data)

// Lexing
fmt.Println("=== Lexing ===")
lx := lexer.New(source)
for {
tok := lx.Next()
if tok.Type == lexer.EOF {
break
}
fmt.Printf("%v: %q\n", tok.Type, tok.Lit)
}

// Parsing
fmt.Println("\n=== Parsing ===")
lx2 := lexer.New(source)
parser := parser.New(lx2)
program := parser.ParseProgram()

if len(parser.Errors()) > 0 {
fmt.Println("Parser errors:")
for _, err := range parser.Errors() {
fmt.Println("  ", err)
}
} else {
fmt.Println("Parsed successfully!")
fmt.Printf("Program structure: %v\n", program)
}

// Runtime
fmt.Println("\n=== Runtime ===")
rt := runtime.New()
rt.SetPlatform("mobile") // Simulate mobile platform

// Execute the program
err = rt.ExecuteProgram(program)
if err != nil {
fmt.Printf("Runtime error: %v\n", err)
}

fmt.Println("\n=== Program Complete ===")
}
