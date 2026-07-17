package compiler

import (
	"fmt"
	"math"
	"os"
	"strings"
	"time"

	"hunterlang/protect"
)

// ── VM Value ──────────────────────────────────────────────────────────────────

type VMValue struct {
	tag  byte // 'n' number, 's' string, 'b' bool, 'a' array, 'o' object, '0' nil
	num  float64
	str  string
	flag bool
	arr  []VMValue
	obj  *VMObject
}

func vmObj(o *VMObject) VMValue { return VMValue{tag: 'o', obj: o} }
func vmNum(n float64) VMValue   { return VMValue{tag: 'n', num: n} }
func vmStr(s string) VMValue    { return VMValue{tag: 's', str: s} }
func vmBool(b bool) VMValue     { return VMValue{tag: 'b', flag: b} }
func vmNil() VMValue            { return VMValue{tag: '0'} }
func vmArr(a []VMValue) VMValue { return VMValue{tag: 'a', arr: a} }

func (v VMValue) isTruthy() bool {
	switch v.tag {
	case 'b':
		return v.flag
	case 'n':
		return v.num != 0
	case 's':
		return v.str != ""
	case 'a':
		return len(v.arr) > 0
	case '0':
		return false
	default:
		return true
	}
}

func (v VMValue) String() string {
	switch v.tag {
	case 'n':
		if v.num == float64(int64(v.num)) {
			return fmt.Sprintf("%d", int64(v.num))
		}
		return fmt.Sprintf("%g", v.num)
	case 's':
		return v.str
	case 'b':
		if v.flag {
			return "yes"
		}
		return "no"
	case 'a':
		parts := make([]string, len(v.arr))
		for i, el := range v.arr {
			parts[i] = el.String()
		}
		return "[" + strings.Join(parts, ", ") + "]"
	case 'o':
		if v.obj != nil {
			return "[" + v.obj.Name + "]"
		}
		return "[object]"
	default:
		return "nothing"
	}
}

func (v VMValue) equals(other VMValue) bool {
	return v.String() == other.String()
}

// ── VM Object ─────────────────────────────────────────────────────────────────

type VMObject struct {
	Name       string
	Fields     map[string]VMValue
	Methods    map[string]*MethodChunk
	Builtin    map[string]func(args []VMValue) (VMValue, error)
	ProtectTag string // "" = unprotected; "protected", "protected1", ...
}

func newVMObject(name string) *VMObject {
	return &VMObject{
		Name:    name,
		Fields:  map[string]VMValue{},
		Methods: map[string]*MethodChunk{},
		Builtin: map[string]func(args []VMValue) (VMValue, error){},
	}
}

// ── VM ────────────────────────────────────────────────────────────────────────

type VM struct {
	chunk       *Chunk
	ip          int
	stack       []VMValue
	vars        map[string]VMValue
	persist     map[string]VMValue
	output      []string
	maxOps      int
	objects     map[string]*VMObject
	callDepth   int
	tryHandlers []tryHandler

	// nameOrigin/ambiguous mirror the interpreter's flat-summon collision
	// tracking: which summoned file last contributed a flat-merged name,
	// and which names are ambiguous because two different files define them.
	nameOrigin map[string]string
	ambiguous  map[string]bool

	// entitlement answers "is this #protected[N] tag currently granted?"
	// Defaults to AlwaysDeny (fail closed).
	entitlement protect.Checker
}

type tryHandler struct {
	target    int
	errVarIdx int
}

func NewVM(chunk *Chunk) *VM {
	return &VM{
		chunk:       chunk,
		ip:          0,
		stack:       make([]VMValue, 0, 64),
		vars:        map[string]VMValue{"pi": vmNum(math.Pi), "nothing": vmNil()},
		persist:     map[string]VMValue{},
		objects:     map[string]*VMObject{},
		maxOps:      1_000_000,
		nameOrigin:  map[string]string{},
		ambiguous:   map[string]bool{},
		entitlement: protect.AlwaysDeny{},
	}
}

// NewVMWithChecker creates a VM with a configured entitlement checker.
// Used by the CLI when a real or stub checker is available.
func NewVMWithChecker(chunk *Chunk, checker protect.Checker) *VM {
	vm := NewVM(chunk)
	if checker != nil {
		vm.entitlement = checker
	}
	return vm
}

func (vm *VM) checkEntitlement(tag string) error {
	granted, err := vm.entitlement.Check(tag)
	if err != nil {
		return fmt.Errorf("couldn't verify access for %q: %v", tag, err)
	}
	if !granted {
		return fmt.Errorf("access denied — %q is not currently available to you", tag)
	}
	return nil
}

func (vm *VM) Run() error {
	ops := 0
	for vm.ip < len(vm.chunk.Instructions) {
		ops++
		if ops > vm.maxOps {
			return fmt.Errorf("VM: execution limit reached (%d ops)", vm.maxOps)
		}
		ins := vm.chunk.Instructions[vm.ip]
		vm.ip++
		if err := vm.exec(ins); err != nil {
			if vm.tryRecover(err) {
				continue
			}
			return err
		}
	}
	return nil
}

func (vm *VM) tryRecover(err error) bool {
	if len(vm.tryHandlers) == 0 {
		return false
	}
	h := vm.tryHandlers[len(vm.tryHandlers)-1]
	vm.tryHandlers = vm.tryHandlers[:len(vm.tryHandlers)-1]
	if h.errVarIdx >= 0 && h.errVarIdx < len(vm.chunk.Names) {
		vm.vars[vm.chunk.Names[h.errVarIdx]] = vmStr(err.Error())
	}
	vm.stack = vm.stack[:0]
	vm.ip = h.target
	return true
}

func (vm *VM) Output() []string { return vm.output }

func (vm *VM) push(v VMValue) { vm.stack = append(vm.stack, v) }

func (vm *VM) pop() (VMValue, error) {
	if len(vm.stack) == 0 {
		return vmNil(), fmt.Errorf("VM: stack underflow at instruction %d", vm.ip)
	}
	v := vm.stack[len(vm.stack)-1]
	vm.stack = vm.stack[:len(vm.stack)-1]
	return v, nil
}

func nameIdx(operand interface{}) int {
	switch v := operand.(type) {
	case int:
		return v
	case [2]int:
		return v[0]
	}
	return -1
}

// ── Opcode execution ──────────────────────────────────────────────────────────

func (vm *VM) exec(ins Instruction) error {
	switch ins.Op {

	case OP_PUSH_NUM:
		vm.push(vmNum(ins.Operand.(float64)))
	case OP_PUSH_STR:
		vm.push(vmStr(ins.Operand.(string)))
	case OP_PUSH_BOOL:
		vm.push(vmBool(ins.Operand.(bool)))
	case OP_PUSH_NIL:
		vm.push(vmNil())
	case OP_POP:
		_, err := vm.pop()
		return err

	case OP_LOAD:
		idx := nameIdx(ins.Operand)
		if idx < 0 || idx >= len(vm.chunk.Names) {
			return fmt.Errorf("LOAD: bad index %v", ins.Operand)
		}
		name := vm.chunk.Names[idx]
		if vm.ambiguous[name] {
			return fmt.Errorf(
				"%q is defined in more than one summoned file — use the file's alias to choose which one",
				name,
			)
		}
		v, ok := vm.vars[name]
		if !ok {
			v = vmNil()
		}
		vm.push(v)

	case OP_STORE:
		idx := nameIdx(ins.Operand)
		if idx < 0 || idx >= len(vm.chunk.Names) {
			return fmt.Errorf("STORE: bad index %v", ins.Operand)
		}
		name := vm.chunk.Names[idx]
		v, err := vm.pop()
		if err != nil {
			return err
		}
		vm.vars[name] = v
		if v.tag == 'o' && v.obj != nil {
			vm.objects[name] = v.obj
		}

	case OP_ADD:
		b, err := vm.pop()
		if err != nil {
			return err
		}
		a, err := vm.pop()
		if err != nil {
			return err
		}
		if a.tag == 'n' && b.tag == 'n' {
			vm.push(vmNum(a.num + b.num))
		} else {
			vm.push(vmStr(a.String() + b.String()))
		}

	case OP_SUB:
		b, err := vm.pop()
		if err != nil {
			return err
		}
		a, err := vm.pop()
		if err != nil {
			return err
		}
		vm.push(vmNum(a.num - b.num))

	case OP_MUL:
		b, err := vm.pop()
		if err != nil {
			return err
		}
		a, err := vm.pop()
		if err != nil {
			return err
		}
		vm.push(vmNum(a.num * b.num))

	case OP_DIV:
		b, err := vm.pop()
		if err != nil {
			return err
		}
		a, err := vm.pop()
		if err != nil {
			return err
		}
		if b.num == 0 {
			return fmt.Errorf("can't divide by zero")
		}
		vm.push(vmNum(a.num / b.num))

	case OP_POW:
		b, err := vm.pop()
		if err != nil {
			return err
		}
		a, err := vm.pop()
		if err != nil {
			return err
		}
		vm.push(vmNum(math.Pow(a.num, b.num)))

	case OP_NEG:
		a, err := vm.pop()
		if err != nil {
			return err
		}
		vm.push(vmNum(-a.num))

	case OP_CONCAT:
		b, err := vm.pop()
		if err != nil {
			return err
		}
		a, err := vm.pop()
		if err != nil {
			return err
		}
		vm.push(vmStr(a.String() + b.String()))

	case OP_EQ:
		b, err := vm.pop()
		if err != nil {
			return err
		}
		a, err := vm.pop()
		if err != nil {
			return err
		}
		vm.push(vmBool(a.equals(b)))

	case OP_NEQ:
		b, err := vm.pop()
		if err != nil {
			return err
		}
		a, err := vm.pop()
		if err != nil {
			return err
		}
		vm.push(vmBool(!a.equals(b)))

	case OP_GT:
		b, err := vm.pop()
		if err != nil {
			return err
		}
		a, err := vm.pop()
		if err != nil {
			return err
		}
		vm.push(vmBool(a.num > b.num))

	case OP_LT:
		b, err := vm.pop()
		if err != nil {
			return err
		}
		a, err := vm.pop()
		if err != nil {
			return err
		}
		vm.push(vmBool(a.num < b.num))

	case OP_GTE:
		b, err := vm.pop()
		if err != nil {
			return err
		}
		a, err := vm.pop()
		if err != nil {
			return err
		}
		vm.push(vmBool(a.num >= b.num))

	case OP_LTE:
		b, err := vm.pop()
		if err != nil {
			return err
		}
		a, err := vm.pop()
		if err != nil {
			return err
		}
		vm.push(vmBool(a.num <= b.num))

	case OP_BETWEEN:
		high, err := vm.pop()
		if err != nil {
			return err
		}
		low, err := vm.pop()
		if err != nil {
			return err
		}
		val, err := vm.pop()
		if err != nil {
			return err
		}
		vm.push(vmBool(val.num >= low.num && val.num <= high.num))

	case OP_AND:
		b, err := vm.pop()
		if err != nil {
			return err
		}
		a, err := vm.pop()
		if err != nil {
			return err
		}
		vm.push(vmBool(a.isTruthy() && b.isTruthy()))

	case OP_OR:
		b, err := vm.pop()
		if err != nil {
			return err
		}
		a, err := vm.pop()
		if err != nil {
			return err
		}
		vm.push(vmBool(a.isTruthy() || b.isTruthy()))

	case OP_NOT:
		a, err := vm.pop()
		if err != nil {
			return err
		}
		vm.push(vmBool(!a.isTruthy()))

	case OP_JUMP:
		vm.ip = ins.Operand.(int)

	case OP_JUMP_IF_NOT:
		cond, err := vm.pop()
		if err != nil {
			return err
		}
		if !cond.isTruthy() {
			vm.ip = ins.Operand.(int)
		}

	case OP_JUMP_IF:
		cond, err := vm.pop()
		if err != nil {
			return err
		}
		if cond.isTruthy() {
			vm.ip = ins.Operand.(int)
		}

	case OP_SAY:
		v, err := vm.pop()
		if err != nil {
			return err
		}
		line := v.String()
		fmt.Println(line)
		vm.output = append(vm.output, line)

	case OP_NEW_ARRAY:
		count := ins.Operand.(int)
		arr := make([]VMValue, count)
		for i := count - 1; i >= 0; i-- {
			v, err := vm.pop()
			if err != nil {
				return err
			}
			arr[i] = v
		}
		vm.push(vmArr(arr))

	case OP_ARRAY_LEN:
		v, err := vm.pop()
		if err != nil {
			return err
		}
		if v.tag != 'a' {
			return fmt.Errorf("ARRAY_LEN: expected list, got %c", v.tag)
		}
		vm.push(vmNum(float64(len(v.arr))))

	case OP_ARRAY_GET:
		idx, err := vm.pop()
		if err != nil {
			return err
		}
		arr, err := vm.pop()
		if err != nil {
			return err
		}
		i := int(idx.num)
		if i < 0 || i >= len(arr.arr) {
			return fmt.Errorf("list index %d out of range", i)
		}
		vm.push(arr.arr[i])

	case OP_REMEMBER:
		idx := nameIdx(ins.Operand)
		key := vm.chunk.Names[idx]
		v, err := vm.pop()
		if err != nil {
			return err
		}
		vm.persist[key] = v

	case OP_RECALL:
		idx := nameIdx(ins.Operand)
		key := vm.chunk.Names[idx]
		v, ok := vm.persist[key]
		if !ok {
			v = vmNil()
		}
		vm.push(v)

	case OP_FORGET:
		idx := nameIdx(ins.Operand)
		delete(vm.persist, vm.chunk.Names[idx])

	case OP_NEW_OBJECT:
		idx := nameIdx(ins.Operand)
		name := vm.chunk.Names[idx]
		vm.push(vmObj(newVMObject(name)))

	case OP_NEW_MODULE:
		arr, ok := ins.Operand.([2]int)
		if !ok {
			return fmt.Errorf("NEW_MODULE: bad operand")
		}
		alias := vm.chunk.Names[arr[0]]
		modName := vm.chunk.Names[arr[1]]
		obj := newVMObject(alias)
		registerVMModule(obj, modName)
		vm.vars[alias] = vmObj(obj)
		vm.objects[alias] = obj

	case OP_INHERIT:
		idx := nameIdx(ins.Operand)
		parentName := vm.chunk.Names[idx]
		childVal, err := vm.pop()
		if err != nil {
			return err
		}
		if childVal.tag != 'o' {
			return fmt.Errorf("INHERIT: top of stack is not an object")
		}
		if parent, ok := vm.objects[parentName]; ok && parent != nil {
			for k, v := range parent.Fields {
				childVal.obj.Fields[k] = v
			}
			for k, m := range parent.Methods {
				childVal.obj.Methods[k] = m
			}
		}
		vm.push(childVal)

	case OP_INIT_FIELD:
		idx := nameIdx(ins.Operand)
		fieldName := vm.chunk.Names[idx]
		val, err := vm.pop()
		if err != nil {
			return err
		}
		objVal, err := vm.pop()
		if err != nil {
			return err
		}
		if objVal.tag != 'o' {
			return fmt.Errorf("INIT_FIELD: expected object on stack")
		}
		objVal.obj.Fields[fieldName] = val
		vm.push(objVal)

	case OP_DEF_METHOD:
		arr, ok := ins.Operand.([2]int)
		if !ok {
			return fmt.Errorf("DEF_METHOD: bad operand")
		}
		methodName := vm.chunk.Names[arr[0]]
		chunkIdx := arr[1]
		objVal, err := vm.pop()
		if err != nil {
			return err
		}
		if objVal.tag != 'o' {
			return fmt.Errorf("DEF_METHOD: expected object on stack")
		}
		if chunkIdx < 0 || chunkIdx >= len(vm.chunk.SubChunks) {
			return fmt.Errorf("DEF_METHOD: bad chunk index %d", chunkIdx)
		}
		objVal.obj.Methods[methodName] = vm.chunk.SubChunks[chunkIdx]
		vm.push(objVal)

	case OP_LOAD_FIELD:
		idx := nameIdx(ins.Operand)
		fieldName := vm.chunk.Names[idx]
		objVal, err := vm.pop()
		if err != nil {
			return err
		}
		if objVal.tag != 'o' || objVal.obj == nil {
			vm.push(vmNil())
			return nil
		}
		if v, ok := objVal.obj.Fields[fieldName]; ok {
			vm.push(v)
		} else {
			vm.push(vmNil())
		}

	case OP_STORE_FIELD:
		idx := nameIdx(ins.Operand)
		fieldName := vm.chunk.Names[idx]
		val, err := vm.pop()
		if err != nil {
			return err
		}
		objVal, err := vm.pop()
		if err != nil {
			return err
		}
		if objVal.tag != 'o' || objVal.obj == nil {
			return fmt.Errorf("STORE_FIELD: not an object")
		}
		objVal.obj.Fields[fieldName] = val

	case OP_CALL_VALUE_METHOD:
		arr, ok := ins.Operand.([2]int)
		if !ok {
			return fmt.Errorf("CALL_VALUE_METHOD: bad operand")
		}
		methodName := vm.chunk.Names[arr[0]]
		argc := arr[1]
		args := make([]VMValue, argc)
		for i := argc - 1; i >= 0; i-- {
			v, err := vm.pop()
			if err != nil {
				return err
			}
			args[i] = v
		}
		recv, err := vm.pop()
		if err != nil {
			return err
		}
		result, err := vm.dispatchValueMethod(recv, methodName, args)
		if err != nil {
			return err
		}
		vm.push(result)

	case OP_CALL_METHOD:
		arr, ok := ins.Operand.([2]int)
		if !ok {
			return fmt.Errorf("CALL_METHOD: bad operand")
		}
		methodName := vm.chunk.Names[arr[0]]
		argc := arr[1]
		args := make([]VMValue, argc)
		for i := argc - 1; i >= 0; i-- {
			v, err := vm.pop()
			if err != nil {
				return err
			}
			args[i] = v
		}
		objVal, err := vm.pop()
		if err != nil {
			return err
		}
		if objVal.tag != 'o' || objVal.obj == nil {
			return fmt.Errorf("can't call .%s — not an object", methodName)
		}
		result, err := vm.invokeMethod(objVal.obj, methodName, args)
		if err != nil {
			return err
		}
		vm.push(result)

	case OP_CALL:
		arr, ok := ins.Operand.([2]int)
		if !ok {
			return fmt.Errorf("CALL: bad operand")
		}
		calleeName := vm.chunk.Names[arr[0]]
		if vm.ambiguous[calleeName] {
			return fmt.Errorf(
				"%q is defined in more than one summoned file — use the file's alias to choose which one",
				calleeName,
			)
		}
		argc := arr[1]
		args := make([]VMValue, argc)
		for i := argc - 1; i >= 0; i-- {
			v, err := vm.pop()
			if err != nil {
				return err
			}
			args[i] = v
		}
		if v, ok := vm.vars[calleeName]; ok && v.tag == 'o' && v.obj != nil {
			methodName := calleeName
			if _, hasCall := v.obj.Methods["__call__"]; hasCall {
				methodName = "__call__"
			}
			result, err := vm.invokeMethod(v.obj, methodName, args)
			if err != nil {
				return err
			}
			vm.push(result)
			return nil
		}
		return fmt.Errorf("I don't know what %q is", calleeName)

	case OP_FIELD_PAIRS:
		objVal, err := vm.pop()
		if err != nil {
			return err
		}
		if objVal.tag != 'o' || objVal.obj == nil {
			return fmt.Errorf("'for each key, val in' expects an object")
		}
		pairs := make([]VMValue, 0, len(objVal.obj.Fields))
		for k, v := range objVal.obj.Fields {
			pairs = append(pairs, vmArr([]VMValue{vmStr(k), v}))
		}
		vm.push(vmArr(pairs))

	case OP_CONTAINS:
		list, err := vm.pop()
		if err != nil {
			return err
		}
		val, err := vm.pop()
		if err != nil {
			return err
		}
		if list.tag != 'a' {
			return fmt.Errorf("'is one of' expects a list")
		}
		found := false
		needle := val.String()
		for _, item := range list.arr {
			if item.String() == needle {
				found = true
				break
			}
		}
		vm.push(vmBool(found))

	case OP_TRY_START:
		arr, ok := ins.Operand.([2]int)
		if !ok {
			return fmt.Errorf("TRY_START: bad operand")
		}
		vm.tryHandlers = append(vm.tryHandlers, tryHandler{target: arr[0], errVarIdx: arr[1]})

	case OP_TRY_END:
		if len(vm.tryHandlers) > 0 {
			vm.tryHandlers = vm.tryHandlers[:len(vm.tryHandlers)-1]
		}

	case OP_RETURN:
		vm.ip = len(vm.chunk.Instructions)

	case OP_RETURN_VOID:
		// no-op at top level

	case OP_RUN_MODULE_CHUNK:
		arr, ok := ins.Operand.([3]int)
		if !ok {
			return fmt.Errorf("RUN_MODULE_CHUNK: bad operand")
		}
		pathIdx, aliasIdx, chunkIdx := arr[0], arr[1], arr[2]
		filePath := vm.chunk.Names[pathIdx]
		hasAlias := aliasIdx >= 0
		var alias string
		if hasAlias {
			alias = vm.chunk.Names[aliasIdx]
		}
		if chunkIdx < 0 || chunkIdx >= len(vm.chunk.SubChunks) {
			return fmt.Errorf("RUN_MODULE_CHUNK: bad chunk index %d", chunkIdx)
		}
		mc := vm.chunk.SubChunks[chunkIdx]

		sub := &VM{
			chunk: &Chunk{
				Name:         mc.Name,
				Instructions: mc.Instructions,
				Constants:    vm.chunk.Constants,
				Names:        vm.chunk.Names,
				SubChunks:    vm.chunk.SubChunks,
			},
			vars:       map[string]VMValue{"pi": vmNum(math.Pi), "nothing": vmNil()},
			persist:    vm.persist,
			objects:    map[string]*VMObject{},
			maxOps:     vm.maxOps,
			callDepth:  vm.callDepth,
			nameOrigin: map[string]string{},
			ambiguous:  map[string]bool{},
		}
		if _, err := sub.runMethodBody(); err != nil {
			return fmt.Errorf("summon %q: %v", filePath, err)
		}

		// Build flatNames: raw top-level vars/objects PLUS, for every
		// exported object, a directly-callable wrapper per builtin/method
		// — mirrors the interpreter so summon "math.he" → sqrt(16) works
		// identically whether run from source or from a compiled .hex.
		flatNames := map[string]VMValue{}
		for k, v := range sub.vars {
			if k == "pi" || k == "nothing" {
				continue
			}
			flatNames[k] = v
		}
		for _, val := range flatNames {
			if val.tag != 'o' || val.obj == nil {
				continue
			}
			for methodName, mchunk := range val.obj.Methods {
				wrapper := newVMObject(methodName)
				wrapper.Fields = val.obj.Fields
				wrapper.Methods[methodName] = mchunk
				flatNames[methodName] = vmObj(wrapper)
			}
			for builtinName, bfn := range val.obj.Builtin {
				if _, already := flatNames[builtinName]; already {
					continue
				}
				wrapper := newVMObject(builtinName)
				wrapper.Builtin[builtinName] = bfn
				flatNames[builtinName] = vmObj(wrapper)
			}
		}

		// Flat-merge into the importing VM, tracking origin for collisions.
		for name, val := range flatNames {
			if origin, exists := vm.nameOrigin[name]; exists && origin != filePath {
				vm.ambiguous[name] = true
				continue
			}
			vm.nameOrigin[name] = filePath
			if !vm.ambiguous[name] {
				vm.vars[name] = val
				if val.tag == 'o' && val.obj != nil {
					vm.objects[name] = val.obj
				}
			}
		}

		// If an alias was given, also expose every flat name as a
		// qualified object field — additive, works even when ambiguous.
		if hasAlias {
			mod := newVMObject(alias)
			for k, v := range flatNames {
				mod.Fields[k] = v
			}
			vm.vars[alias] = vmObj(mod)
			vm.objects[alias] = mod
		}

	case OP_SET_PROTECT_TAG:
		idx := nameIdx(ins.Operand)
		tag := vm.chunk.Names[idx]
		objVal, err := vm.pop()
		if err != nil {
			return err
		}
		if objVal.tag != 'o' || objVal.obj == nil {
			return fmt.Errorf("SET_PROTECT_TAG: expected object on stack")
		}
		objVal.obj.ProtectTag = tag
		vm.push(objVal)

	case OP_ASK:
		promptVal, err := vm.pop()
		if err != nil {
			return err
		}
		fmt.Print(promptVal.String() + " ")
		var input string
		fmt.Scanln(&input)
		vm.push(vmStr(input))

	case OP_NOP:
		// no-op

	case OP_HALT:
		vm.ip = len(vm.chunk.Instructions)

	default:
		return fmt.Errorf("VM: unknown opcode %s at instruction %d", ins.Op, vm.ip-1)
	}
	return nil
}

// ── Virtual method dispatch ────────────────────────────────────────────────────

func (vm *VM) dispatchValueMethod(recv VMValue, method string, args []VMValue) (VMValue, error) {
	switch recv.tag {
	case 'o':
		if recv.obj == nil {
			return vmNil(), fmt.Errorf("can't call .%s — not an object", method)
		}
		return vm.invokeMethod(recv.obj, method, args)

	case 's':
		s := recv.str
		switch method {
		case "upper":
			return vmStr(strings.ToUpper(s)), nil
		case "lower":
			return vmStr(strings.ToLower(s)), nil
		case "length", "len":
			return vmNum(float64(len(s))), nil
		case "trim":
			return vmStr(strings.TrimSpace(s)), nil
		case "split":
			sep := ","
			if len(args) > 0 {
				sep = args[0].str
			}
			parts := strings.Split(s, sep)
			arr := make([]VMValue, len(parts))
			for i, p := range parts {
				arr[i] = vmStr(p)
			}
			return vmArr(arr), nil
		case "contains":
			if len(args) == 0 {
				return vmNil(), fmt.Errorf(".contains() expects an argument")
			}
			return vmBool(strings.Contains(s, args[0].str)), nil
		case "starts", "startsWith":
			return vmBool(strings.HasPrefix(s, args[0].str)), nil
		case "ends", "endsWith":
			return vmBool(strings.HasSuffix(s, args[0].str)), nil
		case "replace":
			if len(args) < 2 {
				return vmNil(), fmt.Errorf(".replace() expects (old, new)")
			}
			return vmStr(strings.ReplaceAll(s, args[0].str, args[1].str)), nil
		case "isEmpty", "empty":
			return vmBool(len(s) == 0), nil
		case "repeat":
			if len(args) == 0 {
				return vmNil(), fmt.Errorf(".repeat() expects a count")
			}
			return vmStr(strings.Repeat(s, int(args[0].num))), nil
		}
		return vmNil(), fmt.Errorf("text doesn't have a .%s() method", method)

	case 'n':
		n := recv.num
		switch method {
		case "abs":
			if n < 0 {
				return vmNum(-n), nil
			}
			return vmNum(n), nil
		case "floor":
			return vmNum(math.Floor(n)), nil
		case "ceil":
			return vmNum(math.Ceil(n)), nil
		case "round":
			return vmNum(math.Round(n)), nil
		case "sqrt":
			return vmNum(math.Sqrt(n)), nil
		case "text", "toString":
			return vmStr(vmNum(n).String()), nil
		case "isPositive":
			return vmBool(n > 0), nil
		case "isNegative":
			return vmBool(n < 0), nil
		case "isZero":
			return vmBool(n == 0), nil
		}
		return vmNil(), fmt.Errorf("number doesn't have a .%s() method", method)

	case 'a':
		arr := recv.arr
		switch method {
		case "length", "len":
			return vmNum(float64(len(arr))), nil
		case "first":
			if len(arr) == 0 {
				return vmNil(), fmt.Errorf("list is empty")
			}
			return arr[0], nil
		case "last":
			if len(arr) == 0 {
				return vmNil(), fmt.Errorf("list is empty")
			}
			return arr[len(arr)-1], nil
		case "add", "append", "push":
			if len(args) == 0 {
				return vmNil(), fmt.Errorf(".add() expects an item")
			}
			return vmArr(append(append([]VMValue{}, arr...), args[0])), nil
		case "contains", "has", "includes":
			if len(args) == 0 {
				return vmNil(), fmt.Errorf(".contains() expects an item")
			}
			needle := args[0].String()
			for _, item := range arr {
				if item.String() == needle {
					return vmBool(true), nil
				}
			}
			return vmBool(false), nil
		case "reverse":
			out := make([]VMValue, len(arr))
			for i, item := range arr {
				out[len(arr)-1-i] = item
			}
			return vmArr(out), nil
		case "join":
			sep := ""
			if len(args) > 0 {
				sep = args[0].str
			}
			parts := make([]string, len(arr))
			for i, item := range arr {
				parts[i] = item.String()
			}
			return vmStr(strings.Join(parts, sep)), nil
		case "isEmpty", "empty":
			return vmBool(len(arr) == 0), nil
		}
		return vmNil(), fmt.Errorf("list doesn't have a .%s() method", method)
	}
	return vmNil(), fmt.Errorf("can't call .%s on this value", method)
}

// ── Method invocation ─────────────────────────────────────────────────────────

func (vm *VM) invokeMethod(obj *VMObject, methodName string, args []VMValue) (VMValue, error) {
	// Stage A, Step 3: enforce whole-object protection tag first.
	// If the object itself is tagged, every method call on it is gated.
	if obj.ProtectTag != "" {
		if err := vm.checkEntitlement(obj.ProtectTag); err != nil {
			return vmNil(), err
		}
	}

	if fn, ok := obj.Builtin[methodName]; ok {
		return fn(args)
	}

	mc, ok := obj.Methods[methodName]
	if !ok {
		// Field-as-callable fallback: summon "math.he" as maths wraps
		// each ability as a single-purpose object stored under
		// maths.Fields[method] — maths.sqrt(16) means "call the wrapper
		// object stored in field sqrt".
		if fv, fok := obj.Fields[methodName]; fok && fv.tag == 'o' && fv.obj != nil {
			return vm.invokeMethod(fv.obj, methodName, args)
		}
		return vmNil(), fmt.Errorf("%q doesn't know how to %q", obj.Name, methodName)
	}

	// Stage A, Step 3: also enforce per-ability protection tag.
	if mc.ProtectTag != "" {
		if err := vm.checkEntitlement(mc.ProtectTag); err != nil {
			return vmNil(), err
		}
	}

	vm.callDepth++
	if vm.callDepth > 500 {
		vm.callDepth--
		return vmNil(), fmt.Errorf("maximum recursion depth reached in %q", methodName)
	}
	defer func() { vm.callDepth-- }()

	sub := &VM{
		chunk: &Chunk{
			Name:         mc.Name,
			Instructions: mc.Instructions,
			Constants:    vm.chunk.Constants,
			Names:        vm.chunk.Names,
			SubChunks:    vm.chunk.SubChunks,
		},
		vars:      map[string]VMValue{},
		persist:   vm.persist,
		objects:   vm.objects,
		maxOps:    vm.maxOps,
		callDepth: vm.callDepth,
	}

	for k, v := range vm.vars {
		sub.vars[k] = v
	}
	for k, v := range obj.Fields {
		sub.vars[k] = v
	}

	paramBackup := map[string]VMValue{}
	hadBackup := map[string]bool{}
	for i, p := range mc.Params {
		if prev, had := obj.Fields[p]; had {
			paramBackup[p] = prev
			hadBackup[p] = true
		}
		if i < len(args) {
			sub.vars[p] = args[i]
		} else {
			sub.vars[p] = vmNil()
		}
	}

	result, err := sub.runMethodBody()

	for k, v := range sub.vars {
		if _, isParam := paramIndex(mc.Params, k); !isParam {
			obj.Fields[k] = v
		}
	}
	for p, had := range hadBackup {
		if had {
			obj.Fields[p] = paramBackup[p]
		} else {
			delete(obj.Fields, p)
		}
	}

	return result, err
}

func paramIndex(params []string, name string) (int, bool) {
	for i, p := range params {
		if p == name {
			return i, true
		}
	}
	return -1, false
}

func (vm *VM) runMethodBody() (VMValue, error) {
	ops := 0
	for vm.ip < len(vm.chunk.Instructions) {
		ops++
		if ops > vm.maxOps {
			return vmNil(), fmt.Errorf("VM: execution limit reached in method body")
		}
		ins := vm.chunk.Instructions[vm.ip]
		vm.ip++

		if ins.Op == OP_RETURN {
			v, err := vm.pop()
			if err != nil {
				return vmNil(), nil
			}
			return v, nil
		}
		if ins.Op == OP_RETURN_VOID {
			return vmNil(), nil
		}
		if err := vm.exec(ins); err != nil {
			if vm.tryRecover(err) {
				continue
			}
			return vmNil(), err
		}
	}
	return vmNil(), nil
}

// ── Stdlib modules for the VM ──────────────────────────────────────────────────

func registerVMModule(obj *VMObject, moduleName string) {
	switch moduleName {
	case "math":
		obj.Builtin["abs"] = func(args []VMValue) (VMValue, error) {
			n := args[0].num
			if n < 0 {
				n = -n
			}
			return vmNum(n), nil
		}
		obj.Builtin["sqrt"] = func(args []VMValue) (VMValue, error) { return vmNum(math.Sqrt(args[0].num)), nil }
		obj.Builtin["round"] = func(args []VMValue) (VMValue, error) { return vmNum(math.Round(args[0].num)), nil }
		obj.Builtin["floor"] = func(args []VMValue) (VMValue, error) { return vmNum(math.Floor(args[0].num)), nil }
		obj.Builtin["ceil"] = func(args []VMValue) (VMValue, error) { return vmNum(math.Ceil(args[0].num)), nil }
		obj.Builtin["max"] = func(args []VMValue) (VMValue, error) {
			if args[0].num > args[1].num {
				return args[0], nil
			}
			return args[1], nil
		}
		obj.Builtin["min"] = func(args []VMValue) (VMValue, error) {
			if args[0].num < args[1].num {
				return args[0], nil
			}
			return args[1], nil
		}
		obj.Builtin["pow"] = func(args []VMValue) (VMValue, error) { return vmNum(math.Pow(args[0].num, args[1].num)), nil }
		obj.Builtin["pi"] = func(args []VMValue) (VMValue, error) { return vmNum(math.Pi), nil }
		obj.Builtin["sin"] = func(args []VMValue) (VMValue, error) { return vmNum(math.Sin(args[0].num)), nil }
		obj.Builtin["cos"] = func(args []VMValue) (VMValue, error) { return vmNum(math.Cos(args[0].num)), nil }
		obj.Builtin["random"] = func(args []VMValue) (VMValue, error) {
			if len(args) == 0 {
				return vmNum(0.5), nil
			}
			if len(args) == 2 {
				return vmNum(args[0].num + 0.5*(args[1].num-args[0].num)), nil
			}
			return vmNil(), fmt.Errorf("math.random expects 0 or 2 args")
		}

	case "text":
		obj.Builtin["upper"] = func(args []VMValue) (VMValue, error) { return vmStr(strings.ToUpper(args[0].str)), nil }
		obj.Builtin["lower"] = func(args []VMValue) (VMValue, error) { return vmStr(strings.ToLower(args[0].str)), nil }
		obj.Builtin["length"] = func(args []VMValue) (VMValue, error) { return vmNum(float64(len(args[0].str))), nil }
		obj.Builtin["trim"] = func(args []VMValue) (VMValue, error) { return vmStr(strings.TrimSpace(args[0].str)), nil }
		obj.Builtin["contains"] = func(args []VMValue) (VMValue, error) {
			return vmBool(strings.Contains(args[0].str, args[1].str)), nil
		}
		obj.Builtin["replace"] = func(args []VMValue) (VMValue, error) {
			return vmStr(strings.ReplaceAll(args[0].str, args[1].str, args[2].str)), nil
		}

	case "list":
		obj.Builtin["length"] = func(args []VMValue) (VMValue, error) { return vmNum(float64(len(args[0].arr))), nil }
		obj.Builtin["first"] = func(args []VMValue) (VMValue, error) {
			if len(args[0].arr) == 0 {
				return vmNil(), fmt.Errorf("list is empty")
			}
			return args[0].arr[0], nil
		}
		obj.Builtin["last"] = func(args []VMValue) (VMValue, error) {
			if len(args[0].arr) == 0 {
				return vmNil(), fmt.Errorf("list is empty")
			}
			return args[0].arr[len(args[0].arr)-1], nil
		}
		obj.Builtin["add"] = func(args []VMValue) (VMValue, error) {
			return vmArr(append(append([]VMValue{}, args[0].arr...), args[1])), nil
		}

	case "io":
		obj.Builtin["read"] = func(args []VMValue) (VMValue, error) {
			data, err := os.ReadFile(args[0].str)
			if err != nil {
				return vmNil(), fmt.Errorf("io.read: %v", err)
			}
			return vmStr(string(data)), nil
		}
		obj.Builtin["write"] = func(args []VMValue) (VMValue, error) {
			if err := os.WriteFile(args[0].str, []byte(args[1].str), 0644); err != nil {
				return vmNil(), fmt.Errorf("io.write: %v", err)
			}
			return vmNil(), nil
		}
		obj.Builtin["exists"] = func(args []VMValue) (VMValue, error) {
			_, err := os.Stat(args[0].str)
			return vmBool(!os.IsNotExist(err)), nil
		}

	case "clock", "time":
		obj.Builtin["today"] = func(args []VMValue) (VMValue, error) { return vmStr(time.Now().Format("2006-01-02")), nil }
		obj.Builtin["timestamp"] = func(args []VMValue) (VMValue, error) { return vmNum(float64(time.Now().Unix())), nil }
		obj.Builtin["now"] = func(args []VMValue) (VMValue, error) {
			now := time.Now()
			t := newVMObject("time")
			t.Fields["year"] = vmNum(float64(now.Year()))
			t.Fields["month"] = vmNum(float64(now.Month()))
			t.Fields["day"] = vmNum(float64(now.Day()))
			t.Fields["hour"] = vmNum(float64(now.Hour()))
			t.Fields["minute"] = vmNum(float64(now.Minute()))
			t.Fields["second"] = vmNum(float64(now.Second()))
			t.Fields["unix"] = vmNum(float64(now.Unix()))
			t.Fields["weekday"] = vmStr(now.Weekday().String())
			return vmObj(t), nil
		}
		obj.Builtin["format"] = func(args []VMValue) (VMValue, error) {
			t := time.Now()
			if len(args) >= 1 && args[0].tag == 'o' && args[0].obj != nil {
				if u, ok := args[0].obj.Fields["unix"]; ok {
					t = time.Unix(int64(u.num), 0)
				}
			}
			layout := "2006-01-02 15:04:05"
			for _, a := range args {
				if a.tag == 's' {
					layout = a.str
					layout = strings.ReplaceAll(layout, "YYYY", "2006")
					layout = strings.ReplaceAll(layout, "MM", "01")
					layout = strings.ReplaceAll(layout, "DD", "02")
					layout = strings.ReplaceAll(layout, "HH", "15")
					layout = strings.ReplaceAll(layout, "mm", "04")
					layout = strings.ReplaceAll(layout, "ss", "05")
				}
			}
			return vmStr(t.Format(layout)), nil
		}
		obj.Builtin["since"] = func(args []VMValue) (VMValue, error) {
			if len(args) < 1 || args[0].tag != 'o' || args[0].obj == nil {
				return vmNil(), fmt.Errorf("clock.since expects a time object")
			}
			var t time.Time
			if u, ok := args[0].obj.Fields["unix"]; ok {
				t = time.Unix(int64(u.num), 0)
			}
			return vmNum(time.Since(t).Seconds()), nil
		}

	case "wolfhead", "os":
		obj.Builtin["notify"] = func(args []VMValue) (VMValue, error) {
			if len(args) == 1 {
				fmt.Printf("[WolfHead] 🔔 %s\n", args[0].str)
			} else if len(args) >= 2 {
				fmt.Printf("[WolfHead] 🔔 %s: %s\n", args[0].str, args[1].str)
			}
			return vmNil(), nil
		}
		obj.Builtin["platform"] = func(args []VMValue) (VMValue, error) { return vmStr("WolfHead/Linux"), nil }
		obj.Builtin["context"] = func(args []VMValue) (VMValue, error) {
			if len(args) >= 1 {
				fmt.Printf("[WolfHead] switching context to %q\n", args[0].str)
				return vmStr(args[0].str), nil
			}
			return vmArr([]VMValue{vmStr("Work"), vmStr("Social"), vmStr("Chill")}), nil
		}
		obj.Builtin["workspace"] = func(args []VMValue) (VMValue, error) {
			if len(args) == 0 {
				w := newVMObject("workspace")
				w.Fields["id"] = vmNum(0)
				w.Fields["name"] = vmStr("Genesis")
				w.Fields["active"] = vmBool(true)
				return vmObj(w), nil
			}
			id := int(args[0].num)
			fmt.Printf("[WolfHead] switching to workspace %d\n", id)
			w := newVMObject("workspace")
			w.Fields["id"] = vmNum(float64(id))
			w.Fields["name"] = vmStr(fmt.Sprintf("Workspace %d", id))
			w.Fields["active"] = vmBool(true)
			return vmObj(w), nil
		}
		obj.Builtin["launch"] = func(args []VMValue) (VMValue, error) {
			if len(args) >= 1 {
				fmt.Printf("[WolfHead] launching %q\n", args[0].str)
			}
			return vmNil(), nil
		}
		obj.Builtin["gesture"] = func(args []VMValue) (VMValue, error) {
			if len(args) >= 1 {
				fmt.Printf("[WolfHead] gesture triggered: %s\n", args[0].String())
			}
			return vmNil(), nil
		}
	}
}

// ── Timing helper ─────────────────────────────────────────────────────────────

func RunTimed(chunk *Chunk) ([]string, time.Duration, error) {
	vm := NewVM(chunk)
	start := time.Now()
	err := vm.Run()
	return vm.Output(), time.Since(start), err
}
