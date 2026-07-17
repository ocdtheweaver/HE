package eval

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"hunterlang/lang/types"
)

// persistStore saves HE values to a JSON file in the user's home directory.
// This implements remember / recall / forget.
type persistStore struct {
	path string
	data map[string]jsonValue
}

// jsonValue is a serialisable form of a types.Value
type jsonValue struct {
	Type    string      `json:"t"`
	Number  float64     `json:"n,omitempty"`
	Str     string      `json:"s,omitempty"`
	Boolean bool        `json:"b,omitempty"`
	Array   []jsonValue `json:"a,omitempty"`
}

func newPersistStore() *persistStore {
	home, err := os.UserHomeDir()
	if err != nil {
		home = "."
	}
	path := filepath.Join(home, ".he_memory.json")
	ps := &persistStore{path: path, data: map[string]jsonValue{}}
	_ = ps.load_file() // ignore error if file doesn't exist yet
	return ps
}

func (ps *persistStore) save(key string, v types.Value) error {
	ps.data[key] = toJSON(v)
	return ps.flush()
}

func (ps *persistStore) load(key string) (types.Value, error) {
	jv, ok := ps.data[key]
	if !ok {
		return types.Nil(), fmt.Errorf("nothing remembered as %q", key)
	}
	return fromJSON(jv), nil
}

func (ps *persistStore) remove(key string) error {
	delete(ps.data, key)
	return ps.flush()
}

func (ps *persistStore) flush() error {
	b, err := json.MarshalIndent(ps.data, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(ps.path, b, 0644)
}

func (ps *persistStore) load_file() error {
	b, err := os.ReadFile(ps.path)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, &ps.data)
}

// ── Conversion helpers ────────────────────────────────────────────────────────

func toJSON(v types.Value) jsonValue {
	switch v.Type {
	case types.NumberT:
		return jsonValue{Type: "n", Number: v.Number}
	case types.StringT:
		return jsonValue{Type: "s", Str: v.Str}
	case types.BooleanT:
		b := false
		if v.Boolean {
			b = true
		}
		return jsonValue{Type: "b", Boolean: b}
	case types.ArrayT:
		arr := make([]jsonValue, len(v.Array))
		for i, el := range v.Array {
			arr[i] = toJSON(el)
		}
		return jsonValue{Type: "a", Array: arr}
	default:
		return jsonValue{Type: "nil"}
	}
}

func fromJSON(jv jsonValue) types.Value {
	switch jv.Type {
	case "n":
		return types.FromNumber(jv.Number)
	case "s":
		return types.FromString(jv.Str)
	case "b":
		return types.FromBoolean(jv.Boolean)
	case "a":
		arr := make([]types.Value, len(jv.Array))
		for i, el := range jv.Array {
			arr[i] = fromJSON(el)
		}
		return types.FromArray(arr)
	default:
		return types.Nil()
	}
}

// AllKeys returns all persisted keys (for REPL :state display)
func (ps *persistStore) AllKeys() []string {
	keys := make([]string, 0, len(ps.data))
	for k := range ps.data {
		keys = append(keys, k)
	}
	return keys
}

// AllValues returns key→value pairs as strings for display
func (ps *persistStore) AllValues() map[string]string {
	out := make(map[string]string, len(ps.data))
	for k, jv := range ps.data {
		v := fromJSON(jv)
		switch v.Type {
		case types.NumberT:
			out[k] = strconv.FormatFloat(v.Number, 'f', -1, 64)
		default:
			out[k] = v.String()
		}
	}
	return out
}
