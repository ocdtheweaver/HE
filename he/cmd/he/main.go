package main

import (
    "fmt"
    "os"
)

func main() {
    if len(os.Args) < 2 {
        fmt.Println("HE Language v2")
        fmt.Println("Usage:")
        fmt.Println("  he run <file.he>    - Run a HE program")
        fmt.Println("  he build <file.he>  - Build for target platform")
        fmt.Println("  he fmt <file.he>    - Format HE code")
        fmt.Println("  he help             - Show this help")
        os.Exit(1)
    }

    command := os.Args[1]
    switch command {
    case "run":
        if len(os.Args) < 3 {
            fmt.Println("Error: missing file argument for run")
            os.Exit(1)
        }
        runFile(os.Args[2])
    case "build":
        if len(os.Args) < 3 {
            fmt.Println("Error: missing file argument for build")
            os.Exit(1)
        }
        buildFile(os.Args[2])
    case "fmt":
        if len(os.Args) < 3 {
            fmt.Println("Error: missing file argument for fmt")
            os.Exit(1)
        }
        fmtFile(os.Args[2])
    case "help", "-h", "--help":
        printHelp()
    default:
        fmt.Printf("Unknown command: %s\n", command)
        printHelp()
        os.Exit(1)
    }
}

func runFile(filename string) {
    fmt.Printf("Running %s (not implemented yet)\n", filename)
    // TODO: Load, parse, and execute HE file
}

func buildFile(filename string) {
    fmt.Printf("Building %s (not implemented yet)\n", filename)
    // TODO: Compile HE file to target platform
}

func fmtFile(filename string) {
    fmt.Printf("Formatting %s (not implemented yet)\n", filename)
    // TODO: Format HE code
}

func printHelp() {
    fmt.Println(`
HE Language v2 - A friendly programming language for creators

Commands:
  run <file.he>    - Execute a HE program
  build <file.he>  - Build for target platform (web/mobile/desktop)
  fmt <file.he>    - Format HE source code
  help             - Show this help message

Examples:
  he run examples/hello.he
  he build game.he --target=web
  he fmt myscript.he

Learn more: https://github.com/he-lang/he
`)
}
