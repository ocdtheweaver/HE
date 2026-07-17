package eval

import (
	"fmt"
	"html"
	"math"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"hunterlang/lang/types"
	"io"
	"net/http"
	"os"
)

// registerModule wires up a named module's builtin methods onto an object.
func registerModule(obj *types.Object, moduleName string) {
	switch strings.ToLower(moduleName) {
	case "ui":
		registerBuiltin(obj, "window", builtinUIWindow)
		registerBuiltin(obj, "navbar", builtinNavbar)
		registerBuiltin(obj, "renderDocs", builtinRenderDocs)
		registerBuiltin(obj, "button", builtinUIButton)
		registerBuiltin(obj, "text", builtinUIText)
	case "physics", "phys":
		registerBuiltin(obj, "gravity", builtinGravity)
		registerBuiltin(obj, "collision", builtinCollision)
	case "math":
		registerBuiltin(obj, "abs", func(args []types.Value) (types.Value, error) {
			n, err := requireNumber(args, 1, "math.abs")
			if err != nil {
				return types.Nil(), err
			}
			return types.FromNumber(math.Abs(n)), nil
		})
		registerBuiltin(obj, "sqrt", func(args []types.Value) (types.Value, error) {
			n, err := requireNumber(args, 1, "math.sqrt")
			if err != nil {
				return types.Nil(), err
			}
			return types.FromNumber(math.Sqrt(n)), nil
		})
		registerBuiltin(obj, "floor", func(args []types.Value) (types.Value, error) {
			n, err := requireNumber(args, 1, "math.floor")
			if err != nil {
				return types.Nil(), err
			}
			return types.FromNumber(math.Floor(n)), nil
		})
		registerBuiltin(obj, "ceil", func(args []types.Value) (types.Value, error) {
			n, err := requireNumber(args, 1, "math.ceil")
			if err != nil {
				return types.Nil(), err
			}
			return types.FromNumber(math.Ceil(n)), nil
		})
		registerBuiltin(obj, "round", func(args []types.Value) (types.Value, error) {
			n, err := requireNumber(args, 1, "math.round")
			if err != nil {
				return types.Nil(), err
			}
			return types.FromNumber(math.Round(n)), nil
		})
		registerBuiltin(obj, "max", func(args []types.Value) (types.Value, error) {
			if len(args) != 2 {
				return types.Nil(), fmt.Errorf("math.max expects 2 numbers")
			}
			a, b := args[0].Number, args[1].Number
			if a > b {
				return types.FromNumber(a), nil
			}
			return types.FromNumber(b), nil
		})
		registerBuiltin(obj, "min", func(args []types.Value) (types.Value, error) {
			if len(args) != 2 {
				return types.Nil(), fmt.Errorf("math.min expects 2 numbers")
			}
			a, b := args[0].Number, args[1].Number
			if a < b {
				return types.FromNumber(a), nil
			}
			return types.FromNumber(b), nil
		})
		registerBuiltin(obj, "random", func(args []types.Value) (types.Value, error) {
			if len(args) == 0 {
				return types.FromNumber(rand.Float64()), nil
			}
			if len(args) == 2 {
				lo, hi := args[0].Number, args[1].Number
				return types.FromNumber(lo + rand.Float64()*(hi-lo)), nil
			}
			return types.Nil(), fmt.Errorf("math.random expects 0 or 2 args")
		})
		registerBuiltin(obj, "pow", func(args []types.Value) (types.Value, error) {
			if len(args) != 2 {
				return types.Nil(), fmt.Errorf("math.pow expects 2 numbers")
			}
			return types.FromNumber(math.Pow(args[0].Number, args[1].Number)), nil
		})
		registerBuiltin(obj, "sin", func(args []types.Value) (types.Value, error) {
			n, err := requireNumber(args, 1, "math.sin")
			if err != nil {
				return types.Nil(), err
			}
			return types.FromNumber(math.Sin(n)), nil
		})
		registerBuiltin(obj, "cos", func(args []types.Value) (types.Value, error) {
			n, err := requireNumber(args, 1, "math.cos")
			if err != nil {
				return types.Nil(), err
			}
			return types.FromNumber(math.Cos(n)), nil
		})
		registerBuiltin(obj, "pi", func(args []types.Value) (types.Value, error) {
			return types.FromNumber(math.Pi), nil
		})

	case "text":
		registerBuiltin(obj, "join", func(args []types.Value) (types.Value, error) {
			// text.join(list, separator)
			if len(args) < 1 || args[0].Type != types.ArrayT {
				return types.Nil(), fmt.Errorf("text.join expects a list as first argument")
			}
			sep := ""
			if len(args) >= 2 {
				sep = args[1].Str
			}
			parts := make([]string, len(args[0].Array))
			for i, v := range args[0].Array {
				parts[i] = v.String()
			}
			return types.FromString(strings.Join(parts, sep)), nil
		})
		registerBuiltin(obj, "split", func(args []types.Value) (types.Value, error) {
			if len(args) < 2 {
				return types.Nil(), fmt.Errorf("text.split expects (text, separator)")
			}
			parts := strings.Split(args[0].Str, args[1].Str)
			vals := make([]types.Value, len(parts))
			for i, p := range parts {
				vals[i] = types.FromString(p)
			}
			return types.FromArray(vals), nil
		})
		registerBuiltin(obj, "upper", func(args []types.Value) (types.Value, error) {
			if len(args) != 1 {
				return types.Nil(), fmt.Errorf("text.upper expects 1 argument")
			}
			return types.FromString(strings.ToUpper(args[0].Str)), nil
		})
		registerBuiltin(obj, "lower", func(args []types.Value) (types.Value, error) {
			if len(args) != 1 {
				return types.Nil(), fmt.Errorf("text.lower expects 1 argument")
			}
			return types.FromString(strings.ToLower(args[0].Str)), nil
		})
		registerBuiltin(obj, "length", func(args []types.Value) (types.Value, error) {
			if len(args) != 1 {
				return types.Nil(), fmt.Errorf("text.length expects 1 argument")
			}
			return types.FromNumber(float64(len(args[0].Str))), nil
		})
		registerBuiltin(obj, "contains", func(args []types.Value) (types.Value, error) {
			if len(args) != 2 {
				return types.Nil(), fmt.Errorf("text.contains expects (text, substr)")
			}
			return types.FromBoolean(strings.Contains(args[0].Str, args[1].Str)), nil
		})
		registerBuiltin(obj, "starts", func(args []types.Value) (types.Value, error) {
			if len(args) != 2 {
				return types.Nil(), fmt.Errorf("text.starts expects (text, prefix)")
			}
			return types.FromBoolean(strings.HasPrefix(args[0].Str, args[1].Str)), nil
		})
		registerBuiltin(obj, "ends", func(args []types.Value) (types.Value, error) {
			if len(args) != 2 {
				return types.Nil(), fmt.Errorf("text.ends expects (text, suffix)")
			}
			return types.FromBoolean(strings.HasSuffix(args[0].Str, args[1].Str)), nil
		})
		registerBuiltin(obj, "replace", func(args []types.Value) (types.Value, error) {
			if len(args) != 3 {
				return types.Nil(), fmt.Errorf("text.replace expects (text, old, new)")
			}
			return types.FromString(strings.ReplaceAll(args[0].Str, args[1].Str, args[2].Str)), nil
		})
		registerBuiltin(obj, "trim", func(args []types.Value) (types.Value, error) {
			if len(args) != 1 {
				return types.Nil(), fmt.Errorf("text.trim expects 1 argument")
			}
			return types.FromString(strings.TrimSpace(args[0].Str)), nil
		})
		registerBuiltin(obj, "number", func(args []types.Value) (types.Value, error) {
			if len(args) != 1 {
				return types.Nil(), fmt.Errorf("text.number expects 1 argument")
			}
			n, err := strconv.ParseFloat(args[0].Str, 64)
			if err != nil {
				return types.Nil(), fmt.Errorf("can't convert %q to a number", args[0].Str)
			}
			return types.FromNumber(n), nil
		})
		registerBuiltin(obj, "from", func(args []types.Value) (types.Value, error) {
			if len(args) != 1 {
				return types.Nil(), fmt.Errorf("text.from expects 1 argument")
			}
			return types.FromString(args[0].String()), nil
		})



	case "net":
		registerBuiltin(obj, "get", func(args []types.Value) (types.Value, error) {
			if len(args) < 1 {
				return types.Nil(), fmt.Errorf("net.get expects a URL")
			}
			resp, err := http.Get(args[0].Str)
			if err != nil {
				return types.Nil(), fmt.Errorf("net.get failed: %v", err)
			}
			defer resp.Body.Close()
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				return types.Nil(), fmt.Errorf("net.get read error: %v", err)
			}
			// Return object with body, status, ok fields
			resultObj := &types.Object{
				Name:   "response",
				Fields: map[string]types.Value{
					"body":   types.FromString(string(body)),
					"status": types.FromNumber(float64(resp.StatusCode)),
					"ok":    types.FromBoolean(resp.StatusCode >= 200 && resp.StatusCode < 300),
				},
				Actions:  map[string]*types.Action{},
				Builtins: map[string]types.BuiltinFn{},
			}
			return types.FromObject(resultObj), nil
		})
		registerBuiltin(obj, "post", func(args []types.Value) (types.Value, error) {
			if len(args) < 2 {
				return types.Nil(), fmt.Errorf("net.post expects (url, body)")
			}
			body := strings.NewReader(args[1].Str)
			resp, err := http.Post(args[0].Str, "application/json", body)
			if err != nil {
				return types.Nil(), fmt.Errorf("net.post failed: %v", err)
			}
			defer resp.Body.Close()
			respBody, _ := io.ReadAll(resp.Body)
			resultObj := &types.Object{
				Name:   "response",
				Fields: map[string]types.Value{
					"body":   types.FromString(string(respBody)),
					"status": types.FromNumber(float64(resp.StatusCode)),
					"ok":    types.FromBoolean(resp.StatusCode >= 200 && resp.StatusCode < 300),
				},
				Actions:  map[string]*types.Action{},
				Builtins: map[string]types.BuiltinFn{},
			}
			return types.FromObject(resultObj), nil
		})
		registerBuiltin(obj, "status", func(args []types.Value) (types.Value, error) {
			if len(args) != 1 || args[0].Type != types.ObjectT {
				return types.Nil(), fmt.Errorf("net.status expects a response object")
			}
			if v, ok := args[0].Object.Fields["status"]; ok {
				return v, nil
			}
			return types.Nil(), nil
		})

	case "wolfhead", "os":
		// WolfHead OS bindings — workspace, context, notifications
		registerBuiltin(obj, "workspace", func(args []types.Value) (types.Value, error) {
			if len(args) < 1 {
				// Return current workspace info
				w := &types.Object{
					Name: "workspace",
					Fields: map[string]types.Value{
						"id":     types.FromNumber(0),
						"name":   types.FromString("Genesis"),
						"active": types.FromBoolean(true),
					},
					Actions: map[string]*types.Action{},
					Builtins: map[string]types.BuiltinFn{},
				}
				return types.FromObject(w), nil
			}
			id := int(args[0].Number)
			fmt.Printf("[WolfHead] switching to workspace %d\n", id)
			w := &types.Object{
				Name: "workspace",
				Fields: map[string]types.Value{
					"id":     types.FromNumber(float64(id)),
					"name":   types.FromString(fmt.Sprintf("Workspace %d", id)),
					"active": types.FromBoolean(true),
				},
				Actions: map[string]*types.Action{},
				Builtins: map[string]types.BuiltinFn{},
			}
			return types.FromObject(w), nil
		})
		registerBuiltin(obj, "context", func(args []types.Value) (types.Value, error) {
			contexts := []string{"Work", "Social", "Chill"}
			if len(args) >= 1 {
				name := args[0].Str
				fmt.Printf("[WolfHead] switching context to %q\n", name)
				return types.FromString(name), nil
			}
			// Return all available contexts
			vals := make([]types.Value, len(contexts))
			for i, c := range contexts {
				vals[i] = types.FromString(c)
			}
			return types.FromArray(vals), nil
		})
		registerBuiltin(obj, "notify", func(args []types.Value) (types.Value, error) {
			if len(args) < 1 {
				return types.Nil(), fmt.Errorf("os.notify expects (message) or (title, message)")
			}
			if len(args) == 1 {
				fmt.Printf("[WolfHead] 🔔 %s\n", args[0].Str)
			} else {
				fmt.Printf("[WolfHead] 🔔 %s: %s\n", args[0].Str, args[1].Str)
			}
			return types.Nil(), nil
		})
		registerBuiltin(obj, "launch", func(args []types.Value) (types.Value, error) {
			if len(args) < 1 {
				return types.Nil(), fmt.Errorf("os.launch expects an app name")
			}
			fmt.Printf("[WolfHead] launching %q\n", args[0].Str)
			return types.Nil(), nil
		})
		registerBuiltin(obj, "gesture", func(args []types.Value) (types.Value, error) {
			if len(args) < 1 {
				return types.Nil(), fmt.Errorf("os.gesture expects a gesture name or number")
			}
			fmt.Printf("[WolfHead] gesture triggered: %s\n", args[0].String())
			return types.Nil(), nil
		})
		registerBuiltin(obj, "platform", func(args []types.Value) (types.Value, error) {
			return types.FromString("WolfHead/Linux"), nil
		})

	case "io":
		registerBuiltin(obj, "read", func(args []types.Value) (types.Value, error) {
			if len(args) != 1 || args[0].Type != types.StringT {
				return types.Nil(), fmt.Errorf("io.read expects a file path")
			}
			data, err := os.ReadFile(args[0].Str)
			if err != nil {
				return types.Nil(), fmt.Errorf("io.read: %v", err)
			}
			return types.FromString(string(data)), nil
		})
		registerBuiltin(obj, "write", func(args []types.Value) (types.Value, error) {
			if len(args) != 2 {
				return types.Nil(), fmt.Errorf("io.write expects (path, content)")
			}
			err := os.WriteFile(args[0].Str, []byte(args[1].Str), 0644)
			if err != nil {
				return types.Nil(), fmt.Errorf("io.write: %v", err)
			}
			return types.Nil(), nil
		})
		registerBuiltin(obj, "append", func(args []types.Value) (types.Value, error) {
			if len(args) != 2 {
				return types.Nil(), fmt.Errorf("io.append expects (path, content)")
			}
			f, err := os.OpenFile(args[0].Str, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				return types.Nil(), fmt.Errorf("io.append: %v", err)
			}
			defer f.Close()
			if _, err := f.WriteString(args[1].Str); err != nil {
				return types.Nil(), fmt.Errorf("io.append: %v", err)
			}
			return types.Nil(), nil
		})
		registerBuiltin(obj, "exists", func(args []types.Value) (types.Value, error) {
			if len(args) != 1 {
				return types.Nil(), fmt.Errorf("io.exists expects a path")
			}
			_, err := os.Stat(args[0].Str)
			return types.FromBoolean(!os.IsNotExist(err)), nil
		})
		registerBuiltin(obj, "delete", func(args []types.Value) (types.Value, error) {
			if len(args) != 1 {
				return types.Nil(), fmt.Errorf("io.delete expects a path")
			}
			return types.Nil(), os.Remove(args[0].Str)
		})
		registerBuiltin(obj, "say", func(args []types.Value) (types.Value, error) {
			parts := make([]string, len(args))
			for i, a := range args {
				parts[i] = a.String()
			}
			fmt.Println(strings.Join(parts, " "))
			return types.Nil(), nil
		})

	case "time", "clock":
		registerBuiltin(obj, "now", func(args []types.Value) (types.Value, error) {
			now := time.Now()
			t := &types.Object{
				Name: "time",
				Fields: map[string]types.Value{
					"year":    types.FromNumber(float64(now.Year())),
					"month":   types.FromNumber(float64(now.Month())),
					"day":     types.FromNumber(float64(now.Day())),
					"hour":    types.FromNumber(float64(now.Hour())),
					"minute":  types.FromNumber(float64(now.Minute())),
					"second":  types.FromNumber(float64(now.Second())),
					"unix":    types.FromNumber(float64(now.Unix())),
					"weekday": types.FromString(now.Weekday().String()),
				},
				Actions:  map[string]*types.Action{},
				Builtins: map[string]types.BuiltinFn{},
			}
			return types.FromObject(t), nil
		})
		registerBuiltin(obj, "today", func(args []types.Value) (types.Value, error) {
			return types.FromString(time.Now().Format("2006-01-02")), nil
		})
		registerBuiltin(obj, "format", func(args []types.Value) (types.Value, error) {
			t := time.Now()
			if len(args) >= 1 && args[0].Type == types.ObjectT && args[0].Object != nil {
				if u, ok := args[0].Object.Fields["unix"]; ok {
					t = time.Unix(int64(u.Number), 0)
				}
			}
			layout := "2006-01-02 15:04:05"
			for _, a := range args {
				if a.Type == types.StringT {
					layout = a.Str
					layout = strings.ReplaceAll(layout, "YYYY", "2006")
					layout = strings.ReplaceAll(layout, "MM", "01")
					layout = strings.ReplaceAll(layout, "DD", "02")
					layout = strings.ReplaceAll(layout, "HH", "15")
					layout = strings.ReplaceAll(layout, "mm", "04")
					layout = strings.ReplaceAll(layout, "ss", "05")
				}
			}
			return types.FromString(t.Format(layout)), nil
		})
		registerBuiltin(obj, "since", func(args []types.Value) (types.Value, error) {
			if len(args) < 1 || args[0].Type != types.ObjectT {
				return types.Nil(), fmt.Errorf("clock.since expects a time object")
			}
			var t time.Time
			if u, ok := args[0].Object.Fields["unix"]; ok {
				t = time.Unix(int64(u.Number), 0)
			}
			return types.FromNumber(time.Since(t).Seconds()), nil
		})
		registerBuiltin(obj, "sleep", func(args []types.Value) (types.Value, error) {
			if len(args) < 1 {
				return types.Nil(), fmt.Errorf("clock.sleep expects seconds")
			}
			time.Sleep(time.Duration(float64(time.Second) * args[0].Number))
			return types.Nil(), nil
		})
		registerBuiltin(obj, "timestamp", func(args []types.Value) (types.Value, error) {
			return types.FromNumber(float64(time.Now().Unix())), nil
		})
	case "list":
		registerBuiltin(obj, "length", func(args []types.Value) (types.Value, error) {
			if len(args) != 1 || args[0].Type != types.ArrayT {
				return types.Nil(), fmt.Errorf("list.length expects a list")
			}
			return types.FromNumber(float64(len(args[0].Array))), nil
		})
		registerBuiltin(obj, "get", func(args []types.Value) (types.Value, error) {
			if len(args) != 2 || args[0].Type != types.ArrayT {
				return types.Nil(), fmt.Errorf("list.get expects (list, index)")
			}
			idx := int(args[1].Number)
			if idx < 0 || idx >= len(args[0].Array) {
				return types.Nil(), fmt.Errorf("list index %d out of range", idx)
			}
			return args[0].Array[idx], nil
		})
		registerBuiltin(obj, "add", func(args []types.Value) (types.Value, error) {
			if len(args) != 2 || args[0].Type != types.ArrayT {
				return types.Nil(), fmt.Errorf("list.add expects (list, item)")
			}
			newArr := append(append([]types.Value{}, args[0].Array...), args[1])
			return types.FromArray(newArr), nil
		})
		registerBuiltin(obj, "remove", func(args []types.Value) (types.Value, error) {
			if len(args) != 2 || args[0].Type != types.ArrayT {
				return types.Nil(), fmt.Errorf("list.remove expects (list, index)")
			}
			idx := int(args[1].Number)
			arr := args[0].Array
			if idx < 0 || idx >= len(arr) {
				return types.Nil(), fmt.Errorf("list index %d out of range", idx)
			}
			newArr := append(append([]types.Value{}, arr[:idx]...), arr[idx+1:]...)
			return types.FromArray(newArr), nil
		})
		registerBuiltin(obj, "contains", func(args []types.Value) (types.Value, error) {
			if len(args) != 2 || args[0].Type != types.ArrayT {
				return types.Nil(), fmt.Errorf("list.contains expects (list, item)")
			}
			needle := args[1].String()
			for _, v := range args[0].Array {
				if v.String() == needle {
					return types.FromBoolean(true), nil
				}
			}
			return types.FromBoolean(false), nil
		})
		registerBuiltin(obj, "first", func(args []types.Value) (types.Value, error) {
			if len(args) != 1 || args[0].Type != types.ArrayT || len(args[0].Array) == 0 {
				return types.Nil(), fmt.Errorf("list.first expects a non-empty list")
			}
			return args[0].Array[0], nil
		})
		registerBuiltin(obj, "last", func(args []types.Value) (types.Value, error) {
			if len(args) != 1 || args[0].Type != types.ArrayT || len(args[0].Array) == 0 {
				return types.Nil(), fmt.Errorf("list.last expects a non-empty list")
			}
			return args[0].Array[len(args[0].Array)-1], nil
		})
	}
}

// ── Helpers ────────────────────────────────────────────────────────────────────

func registerBuiltin(obj *types.Object, name string, fn types.BuiltinFn) {
	if obj.Builtins == nil {
		obj.Builtins = map[string]types.BuiltinFn{}
	}
	if obj.Actions == nil {
		obj.Actions = map[string]*types.Action{}
	}
	obj.Builtins[name] = fn
	obj.Actions[name] = &types.Action{Name: name}
}

func requireNumber(args []types.Value, n int, name string) (float64, error) {
	if len(args) != n || args[0].Type != types.NumberT {
		return 0, fmt.Errorf("%s expects %d number argument(s)", name, n)
	}
	return args[0].Number, nil
}

// ── UI builtins ───────────────────────────────────────────────────────────────

func builtinNavbar(args []types.Value) (types.Value, error) {
	if len(args) != 1 || args[0].Type != types.ArrayT {
		return types.Nil(), fmt.Errorf("ui.navbar expects one list of items")
	}
	fmt.Println("UI.navbar:", args[0].String())
	return types.Nil(), nil
}

func builtinUIText(args []types.Value) (types.Value, error) {
	if len(args) != 1 {
		return types.Nil(), fmt.Errorf("ui.text expects one argument")
	}
	return types.FromArray([]types.Value{
		types.FromString("text"),
		types.FromString(args[0].String()),
	}), nil
}

func builtinUIButton(args []types.Value) (types.Value, error) {
	if len(args) < 1 {
		return types.Nil(), fmt.Errorf("ui.button expects (label) or (label, clickLabel)")
	}
	label := args[0].String()
	clickLabel := label
	if len(args) >= 2 {
		clickLabel = args[1].String()
	}
	return types.FromArray([]types.Value{
		types.FromString("button"),
		types.FromString(label),
		types.FromString(clickLabel),
	}), nil
}

func builtinUIWindow(args []types.Value) (types.Value, error) {
	var title string
	var children []types.Value

	// Support named arg block: [title is "...", children: [...]]
	if len(args) == 1 && args[0].Type == types.ObjectT && args[0].Object != nil {
		obj := args[0].Object
		if v, ok := obj.Fields["title"]; ok {
			title = v.String()
		}
		if v, ok := obj.Fields["children"]; ok && v.Type == types.ArrayT {
			children = v.Array
		}
	} else if len(args) == 1 && args[0].Type == types.ArrayT {
		// Positional array: [title, childrenArray]
		arr := args[0].Array
		if len(arr) >= 1 {
			title = arr[0].String()
		}
		if len(arr) >= 2 && arr[1].Type == types.ArrayT {
			children = arr[1].Array
		}
	} else {
		return types.Nil(), fmt.Errorf("ui.window expects [title is \"...\", children: [...]]")
	}

	output := buildWindowHTML(title, children)
	fmt.Println(output)
	return types.Nil(), nil
}

func buildWindowHTML(title string, children []types.Value) string {
	return fmt.Sprintf(`<!doctype html>
<html>
<head>
  <meta charset="utf-8"/>
  <meta name="viewport" content="width=device-width, initial-scale=1"/>
  <title>%s</title>
  <style>
    body{font-family:ui-monospace,SFMono-Regular,Menlo,Monaco,Consolas,monospace;margin:16px;background:#fff;color:#000;}
    .card{border:4px solid #000;background:#f8f8f8;box-shadow:6px 6px 0 #000;padding:14px;margin-top:12px;}
    button{border:4px solid #000;background:#fff;box-shadow:4px 4px 0 #000;padding:10px 14px;font-weight:900;cursor:pointer;}
    button:active{transform:translate(2px,2px);box-shadow:2px 2px 0 #000;}
  </style>
</head>
<body>
  <h1>%s</h1>
  <div class="card">%s</div>
<script>function heClick(label){console.log("HE click:", label);}</script>
</body>
</html>`, title, title, renderChildrenHTML(children))
}

func renderChildrenHTML(children []types.Value) string {
	var out strings.Builder
	for _, ch := range children {
		if ch.Type != types.ArrayT || len(ch.Array) < 2 {
			continue
		}
		nodeType := ch.Array[0].String()
		switch nodeType {
		case "text":
			out.WriteString(fmt.Sprintf(`<div>%s</div>`, escapeHTML(ch.Array[1].String())))
		case "button":
			label := ch.Array[1].String()
			clickLabel := label
			if len(ch.Array) >= 3 {
				clickLabel = ch.Array[2].String()
			}
			out.WriteString(fmt.Sprintf(
				`<div style="margin-top:10px;"><button onclick="heClick('%s')">%s</button></div>`,
				escapeHTML(clickLabel), escapeHTML(label),
			))
		}
	}
	return out.String()
}

func escapeHTML(s string) string {
	escaped := html.EscapeString(s)
	escaped = strings.ReplaceAll(escaped, "\\", "\\\\")
	escaped = strings.ReplaceAll(escaped, "'", "\\'")
	return escaped
}

func builtinRenderDocs(args []types.Value) (types.Value, error) {
	if len(args) != 1 || args[0].Type != types.ArrayT {
		return types.Nil(), fmt.Errorf("ui.renderDocs expects [title, tab1, tab2, ...]")
	}
	arr := args[0].Array
	if len(arr) < 2 {
		return types.Nil(), fmt.Errorf("ui.renderDocs needs at least [title, tab1]")
	}
	title := arr[0].String()
	tabs := make([]string, len(arr)-1)
	for i, v := range arr[1:] {
		tabs[i] = v.String()
	}

	var tabBtns, sections strings.Builder
	for i, t := range tabs {
		sel := "false"
		hidden := "true"
		if i == 0 {
			sel = "true"
			hidden = "false"
		}
		tabBtns.WriteString(fmt.Sprintf(
			`<div class="tab" role="tab" tabindex="0" aria-selected='%s' data-tab='%s'>%s</div>`+"\n",
			sel, t, t,
		))
		sections.WriteString(fmt.Sprintf(
			`<div class="card section" data-tab='%s' aria-hidden='%s'><h3>%s</h3></div>`+"\n",
			t, hidden, t,
		))
	}

	fmt.Printf(`<!doctype html>
<html>
<head>
  <meta charset="utf-8"/>
  <title>%s</title>
  <style>
    body{font-family:ui-monospace,monospace;margin:0;background:#fff;color:#000;}
    header{position:sticky;top:0;background:#fff;border-bottom:4px solid #000;padding:12px 16px;}
    .tabs{margin-top:10px;display:flex;gap:10px;flex-wrap:wrap;}
    .tab{border:3px solid #000;padding:8px 12px;box-shadow:4px 4px 0 #000;cursor:pointer;font-weight:900;}
    .tab[aria-selected='true']{background:#f8f8f8;}
    main{padding:16px;max-width:980px;margin:0 auto;}
    .card{border:4px solid #000;background:#f8f8f8;box-shadow:6px 6px 0 #000;padding:14px;margin-top:12px;}
    .section{display:none;}
    .section[aria-hidden='false']{display:block;}
  </style>
</head>
<body>
  <header>
    <div style="font-size:18px;font-weight:900;">%s</div>
    <div class="tabs">%s</div>
  </header>
  <main>%s</main>
<script>
  function setTab(n){
    document.querySelectorAll('.tab').forEach(t=>t.setAttribute('aria-selected',t.dataset.tab===n?'true':'false'));
    document.querySelectorAll('.section').forEach(s=>s.setAttribute('aria-hidden',s.dataset.tab===n?'false':'true'));
  }
  document.querySelectorAll('.tab').forEach(t=>t.addEventListener('click',()=>setTab(t.dataset.tab)));
</script>
</body>
</html>`, title, title, tabBtns.String(), sections.String())
	return types.Nil(), nil
}

// ── Physics builtins ──────────────────────────────────────────────────────────

func builtinGravity(args []types.Value) (types.Value, error) {
	if len(args) != 1 || args[0].Type != types.NumberT {
		return types.Nil(), fmt.Errorf("phys.gravity expects one number (e.g. -9.8)")
	}
	fmt.Printf("Physics: gravity set to %s\n", args[0].String())
	return types.Nil(), nil
}

func builtinCollision(args []types.Value) (types.Value, error) {
	if len(args) != 2 {
		return types.Nil(), fmt.Errorf("phys.collision expects two object names")
	}
	if args[0].String() == args[1].String() {
		return types.FromBoolean(true), nil
	}
	return types.FromBoolean(false), nil
}
