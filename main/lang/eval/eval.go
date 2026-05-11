package eval

import (
	"hunterlang/lang/ast"
	"hunterlang/lang/types"
)

type Interpreter struct {
	runtime *Runtime

	// Injected "current event" to trigger matching reactions after program execution.
	eventType  string
	eventRight *string
}

func NewInterpreter() *Interpreter {
	return &Interpreter{}
}

func (i *Interpreter) SetEvent(eventType string, right *string) {
	i.eventType = eventType
	i.eventRight = right
}

func (i *Interpreter) Run(prog *ast.Program) error {
	i.runtime = newRuntime()
	if err := i.runtime.execProgram(prog); err != nil {
		return err
	}
	// After executing the program, apply the injected event (if any).
	if i.eventType != "" {
		i.runtime.triggerReactionsAll(i.eventType, i.eventRight)
	}
	return nil
}

// Snapshot is a read-only view of runtime state for debugging/visualization.
type Snapshot struct {
	Platform string
	Env      map[string]types.Value
	Objects  map[string]types.Object
}

// Snapshot returns a best-effort copy of current runtime variables/objects.
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
		if obj == nil {
			continue
		}
		// Shallow copy of Fields/Actions/Reactions is fine for our debug view.
		objCopy[name] = *obj
	}

	return Snapshot{
		Platform: i.runtime.Platform,
		Env:      envCopy,
		Objects:  objCopy,
	}
}
