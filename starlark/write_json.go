package starlark

import (
	"bytes"
	"encoding/json"
	"fmt"

	"go.starlark.net/starlark"
	realyaml "gopkg.in/yaml.v2"
)

func ValToYaml(val starlark.Value) (*bytes.Buffer, error) {
	bufs, err := starlarkValToJSON(val)
	if err != nil {
		return nil, err
	}
	return dumpToYaml(bufs)
}

func dumpToYaml(jsonBufs []*bytes.Buffer) (*bytes.Buffer, error) {
	out := bytes.NewBuffer(nil)
	for i, buf := range jsonBufs {
		dict := make(map[string]interface{})
		err := realyaml.Unmarshal(buf.Bytes(), &dict)
		if err != nil {
			return nil, err
		}
		enc := realyaml.NewEncoder(out)
		err = enc.Encode(dict)
		if err != nil {
			return nil, err
		}
		enc.Close()
		if i != len(jsonBufs)-1 {
			out.WriteString("\n----\n")
		}
	}
	return out, nil
}

// Returns a list of buffers of JSON-encoded bytes
func starlarkValToJSON(input starlark.Value) ([]*bytes.Buffer, error) {
	var out []*bytes.Buffer

	switch v := input.(type) {
	case *starlark.List:
		for i := 0; i < v.Len(); i++ {
			buf := bytes.NewBuffer(nil)
			err := writeJSON(buf, v.Index(i))
			if err != nil {
				return nil, err
			}
			out = append(out, buf)
		}
	case *starlark.Dict:
		buf := bytes.NewBuffer(nil)
		err := writeJSON(buf, v)
		if err != nil {
			return nil, err
		}
		out = append(out, buf)
	default:
		return nil, fmt.Errorf("invalid starlarkValToJSON type (got a %s)", input.Type())
	}
	return out, nil
}

// Adapted from skycfg:
// https://github.com/stripe/skycfg/blob/eaa524101c2a0807c13ed5d2e52576fefc146ec3/internal/go/skycfg/json_write.go#L45
func writeJSON(out *bytes.Buffer, v starlark.Value) error {
	if marshaler, ok := v.(json.Marshaler); ok {
		jsonData, err := marshaler.MarshalJSON()
		if err != nil {
			return err
		}
		out.Write(jsonData)
		return nil
	}

	switch v := v.(type) {
	case starlark.NoneType:
		out.WriteString("null")
	case starlark.Bool:
		fmt.Fprintf(out, "%t", v)
	case starlark.Int:
		out.WriteString(v.String())
	case starlark.Float:
		fmt.Fprintf(out, "%g", v)
	case starlark.String:
		s := string(v)
		if goQuoteIsSafe(s) {
			fmt.Fprintf(out, "%q", s)
		} else {
			// vanishingly rare for text strings
			data, _ := json.Marshal(s)
			out.Write(data)
		}
	case starlark.Indexable: // Tuple, List
		out.WriteByte('[')
		for i, n := 0, starlark.Len(v); i < n; i++ {
			if i > 0 {
				out.WriteString(", ")
			}
			if err := writeJSON(out, v.Index(i)); err != nil {
				return err
			}
		}
		out.WriteByte(']')
	case *starlark.Dict:
		out.WriteByte('{')
		for i, itemPair := range v.Items() {
			key := itemPair[0]
			value := itemPair[1]
			if i > 0 {
				out.WriteString(", ")
			}
			if err := writeJSON(out, key); err != nil {
				return err
			}
			out.WriteString(": ")
			if err := writeJSON(out, value); err != nil {
				return err
			}
		}
		out.WriteByte('}')
	default:
		return fmt.Errorf("TypeError: value %s (type `%s') can't be converted to JSON.", v.String(), v.Type())
	}
	return nil
}

func goQuoteIsSafe(s string) bool {
	for _, r := range s {
		// JSON doesn't like Go's \xHH escapes for ASCII control codes,
		// nor its \UHHHHHHHH escapes for runes >16 bits.
		if r < 0x20 || r >= 0x10000 {
			return false
		}
	}
	return true
}
