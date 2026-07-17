package compiler

import (
	"os"
	"path/filepath"

	"hunterlang/lang/ast"
	"hunterlang/lang/lexer"
	"hunterlang/lang/parser"
)

// These thin wrappers exist so bytecode.go can call short, readable names
// without importing os/path/filepath directly into its own import block,
// keeping that file focused on compilation logic.

func osReadFile(path string) (string, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func osStat(path string) (os.FileInfo, error) {
	return os.Stat(path)
}

func filepathIsAbs(path string) bool  { return filepath.IsAbs(path) }
func filepathJoin(a, b string) string { return filepath.Join(a, b) }
func filepathDir(path string) string  { return filepath.Dir(path) }

// parseModuleSource lexes and parses HE source text into an AST program,
// used when compiling file modules (summon "lib.he" as alias) into bytecode.
func parseModuleSource(src string) (*ast.Program, error) {
	lx := lexer.New(src)
	p := parser.New(lx)
	return p.ParseProgram()
}
