package eval

import (
	"hunterlang/lang/ast"
	"hunterlang/lang/types"
)

type Interpreter struct {
	runtime     *Runtime
	eventType   string
	eventRight  *string
	baseDir     string
	entitlement EntitlementChecker
}

func NewInterpreter() *Interpreter {
	return &Interpreter{}
}

func (i *Interpreter) SetEvent(eventType string, right *string) {
	i.eventType = eventType
	i.eventRight = right
}

// SetBaseDir sets the directory used to resolve relative summon paths
// (e.g. summon "lib/helpers.he"). Call before Run.
func (i *Interpreter) SetBaseDir(dir string) {
	i.baseDir = dir
}

// SetEntitlementChecker configures how #protected[N] tags will be
// resolved once enforcement (Stage A, Step 3) is built. Call before Run.
// If never called, the runtime defaults to AlwaysDeny (fail closed).
func (i *Interpreter) SetEntitlementChecker(c EntitlementChecker) {
	i.entitlement = c
}

func (i *Interpreter) Run(prog *ast.Program) error {
	i.runtime = newRuntime()
	i.runtime.baseDir = i.baseDir
	if i.entitlement != nil {
		i.runtime.entitlement = i.entitlement
	}
	if prog != nil {
		if err := i.runtime.execProgram(prog); err != nil {
			return err
		}
	}
	if i.eventType != "" {
		i.runtime.triggerReactionsAll(i.eventType, i.eventRight)
	}
	return nil
}


// RunProg executes an additional parsed program on the current runtime state.
// Used by the REPL to maintain variables and objects between lines.
func (i *Interpreter) RunProg(prog *ast.Program) error {
	if i.runtime == nil {
		i.runtime = newRuntime()
	}
	return i.runtime.execProgram(prog)
}

// Snapshot is a read-only view of runtime state for debugging.
type Snapshot struct {
	Platform string
	Env      map[string]types.Value
	Objects  map[string]types.Object
}

func (i *Interpreter) Snapshot() Snapshot {
	if i.runtime == nil {
		return Snapshot{
			Platform: "unknown",
			Env:      map[string]types.Value{},
			Objects:  map[string]types.Object{},
		}
	}
	envCopy := make(map[string]types.Value, len(i.runtime.Env))
	for k, v := range i.runtime.Env {
		envCopy[k] = v
	}
	objCopy := make(map[string]types.Object, len(i.runtime.Objects))
	for name, obj := range i.runtime.Objects {
		if obj != nil {
			objCopy[name] = *obj
		}
	}
	return Snapshot{
		Platform: i.runtime.Platform,
		Env:      envCopy,
		Objects:  objCopy,
	}
}
