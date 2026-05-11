package runtime

import (
"fmt"
"time"

"github.com/user/he/pkg/ast"
)

// Value represents a runtime value.
type Value struct {
Type  string
Num   float64
Str   string
Bool  bool
Func  *Function
Obj   *Object
Array []Value
}

// Function represents a callable function.
type Function struct {
Name   string
Params []string
Body   []ast.Statement
Env    *Environment
}

// Object represents a runtime object.
type Object struct {
Name       string
Like       string
Properties map[string]Value
Abilities  map[string]*Function
Memories   map[string]Value
}

// Environment represents a scope for variable storage.
type Environment struct {
Parent *Environment
Vars   map[string]Value
}

// NewEnvironment creates a new environment.
func NewEnvironment(parent *Environment) *Environment {
return &Environment{
Parent: parent,
Vars:   make(map[string]Value),
}
}

// Get retrieves a variable from the environment.
func (e *Environment) Get(name string) (Value, bool) {
val, ok := e.Vars[name]
if !ok && e.Parent != nil {
return e.Parent.Get(name)
}
return val, ok
}

// Set sets a variable in the current environment.
func (e *Environment) Set(name string, val Value) {
e.Vars[name] = val
}

// Module represents a loaded module.
type Module struct {
Name      string
Functions map[string]*Function
Submodules map[string]*Module
}

// Runtime executes HE programs.
type Runtime struct {
GlobalEnv  *Environment
Modules    map[string]*Module
Objects    map[string]*Object
Platform   string // "web", "mobile", "desktop"
Debug      bool
}

// New creates a new runtime.
func New() *Runtime {
return &Runtime{
GlobalEnv: NewEnvironment(nil),
Modules:   make(map[string]*Module),
Objects:   make(map[string]*Object),
Platform:  "web", // default platform
}
}

// SetPlatform sets the target platform.
func (rt *Runtime) SetPlatform(platform string) {
rt.Platform = platform
}

// LoadModule loads a module into the runtime.
func (rt *Runtime) LoadModule(name string, module *Module) {
rt.Modules[name] = module
}

// ExecuteProgram executes a complete HE program.
func (rt *Runtime) ExecuteProgram(program *ast.Program) error {
// Load summoned modules
for _, summon := range program.Summons {
fmt.Printf("Loading module: %s as %s\n", summon.Path, summon.As)
// TODO: Implement module loading from files
}

// Register modules from AST
for _, module := range program.Modules {
rt.registerModule(module)
}

// Create objects
for _, obj := range program.Objects {
rt.createObject(obj)
}

// Execute global statements
for _, stmt := range program.Statements {
if err := rt.executeStatement(stmt); err != nil {
return err
}
}

return nil
}

// registerModule registers an AST module in the runtime.
func (rt *Runtime) registerModule(module *ast.ModuleDecl) *Module {
rtModule := &Module{
Name:       module.Name,
Functions:  make(map[string]*Function),
Submodules: make(map[string]*Module),
}

// Register functions
for _, fn := range module.Functions {
rtFunc := &Function{
Name:   fn.Name,
Params: make([]string, len(fn.Params)),
Body:   fn.Body,
Env:    NewEnvironment(rt.GlobalEnv),
}
for i, param := range fn.Params {
rtFunc.Params[i] = param.Name
}
rtModule.Functions[fn.Name] = rtFunc
}

// Register submodules
for _, sub := range module.Submodules {
rtSub := rt.registerModule(sub)
rtModule.Submodules[sub.Name] = rtSub
}

rt.Modules[module.Name] = rtModule
return rtModule
}

// createObject creates an object from AST.
func (rt *Runtime) createObject(obj *ast.ObjectDecl) *Object {
rtObj := &Object{
Name:       obj.Name,
Like:       obj.Like,
Properties: make(map[string]Value),
Abilities:  make(map[string]*Function),
Memories:   make(map[string]Value),
}

// Initialize properties
for _, prop := range obj.Properties {
val := rt.evaluateExpression(prop.Value)
rtObj.Properties[prop.Name] = val
}

rt.Objects[obj.Name] = rtObj
return rtObj
}

// executeStatement executes a single statement.
func (rt *Runtime) executeStatement(stmt ast.Statement) error {
switch s := stmt.(type) {
case *ast.PrintStmt:
return rt.executePrint(s)
case *ast.WaitStmt:
return rt.executeWait(s)
case *ast.ReturnStmt:
return rt.executeReturn(s)
default:
return fmt.Errorf("unknown statement type: %T", stmt)
}
}

// executePrint executes a print statement.
func (rt *Runtime) executePrint(stmt *ast.PrintStmt) error {
val := rt.evaluateExpression(stmt.Expr)

switch val.Type {
case "string":
fmt.Println(val.Str)
case "number":
fmt.Println(val.Num)
case "boolean":
fmt.Println(val.Bool)
default:
fmt.Println(val)
}

return nil
}

// executeWait executes a wait statement.
func (rt *Runtime) executeWait(stmt *ast.WaitStmt) error {
val := rt.evaluateExpression(stmt.Seconds)
if val.Type != "number" {
return fmt.Errorf("wait expects a number, got %s", val.Type)
}

seconds := time.Duration(val.Num * float64(time.Second))
if rt.Debug {
fmt.Printf("[DEBUG] Waiting for %v seconds\n", val.Num)
}
time.Sleep(seconds)
return nil
}

// executeReturn executes a return statement (placeholder).
func (rt *Runtime) executeReturn(stmt *ast.ReturnStmt) error {
// Return statements are handled in function execution
if rt.Debug {
fmt.Println("[DEBUG] Return statement")
}
return nil
}

// evaluateExpression evaluates an expression to a value.
func (rt *Runtime) evaluateExpression(expr ast.Expression) Value {
switch e := expr.(type) {
case *ast.StringExpr:
return Value{Type: "string", Str: e.Value}
case *ast.NumberExpr:
return Value{Type: "number", Num: e.Value}
case *ast.BoolExpr:
return Value{Type: "boolean", Bool: e.Value}
case *ast.IdentExpr:
val, ok := rt.GlobalEnv.Get(e.Name)
if !ok {
// Variable not found, return default
return Value{Type: "string", Str: ""}
}
return val
default:
// For now, return a default value
return Value{Type: "string", Str: ""}
}
}

// ResolveFunction resolves a function call with platform-aware routing.
func (rt *Runtime) ResolveFunction(modulePath []string, funcName string) (*Function, error) {
// Example: ui.navbar -> ["ui"], "navbar"
// Platform-aware resolution:
// 1. Try exact path: ui.navbar
// 2. Try platform-specific: ui.mobile.navbar (if platform is mobile)
// 3. Try fallback platforms

if len(modulePath) == 0 {
return nil, fmt.Errorf("empty module path")
}

// Start from root module
current := rt.Modules[modulePath[0]]
if current == nil {
return nil, fmt.Errorf("module not found: %s", modulePath[0])
}

// Navigate through submodules
for i := 1; i < len(modulePath); i++ {
sub := current.Submodules[modulePath[i]]
if sub == nil {
return nil, fmt.Errorf("submodule not found: %s", modulePath[i])
}
current = sub
}

// Try to find function in current module
if fn, ok := current.Functions[funcName]; ok {
return fn, nil
}

// Platform-aware fallback
// If we're looking for ui.navbar, try ui.<platform>.navbar
if len(modulePath) == 1 {
// Single module like "ui", try platform submodule
platformModule := current.Submodules[rt.Platform]
if platformModule != nil {
if fn, ok := platformModule.Functions[funcName]; ok {
return fn, nil
}
}

// Try other common platforms as fallback
fallbacks := []string{"web", "mobile", "desktop"}
for _, fb := range fallbacks {
if fb == rt.Platform {
continue // already tried
}
platformModule := current.Submodules[fb]
if platformModule != nil {
if fn, ok := platformModule.Functions[funcName]; ok {
return fn, nil
}
}
}
}

return nil, fmt.Errorf("function not found: %s.%s", 
strings.Join(modulePath, "."), funcName)
}
