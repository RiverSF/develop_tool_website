package json2struct

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"unicode"
)

// Json2struct converts JSON text into a Go struct definition.
func Json2struct(input string) string {
	var data interface{}
	if err := json.Unmarshal([]byte(input), &data); err != nil {
		return fmt.Sprintf("// parse error: %v", err)
	}

	gen := newGenerator()
	gen.defineStruct("Request", data)
	return gen.render()
}

type fieldDef struct {
	jsonKey string
	goName  string
	goType  string
}

type generator struct {
	fields map[string][]fieldDef
	order  []string
	seen   map[string]bool
	seq    int
}

func newGenerator() *generator {
	return &generator{
		fields: map[string][]fieldDef{},
		seen:   map[string]bool{},
	}
}

func (g *generator) defineStruct(name string, data interface{}) {
	obj, ok := data.(map[string]interface{})
	if !ok {
		return
	}
	if g.seen[name] {
		return
	}
	g.seen[name] = true
	g.order = append(g.order, name)

	fields := make([]fieldDef, 0, len(obj))
	for key, val := range obj {
		fields = append(fields, fieldDef{
			jsonKey: key,
			goName:  exportField(key),
			goType:  g.inferType(val, exportField(key)),
		})
	}
	sort.Slice(fields, func(i, j int) bool { return fields[i].jsonKey < fields[j].jsonKey })
	g.fields[name] = fields
}

func (g *generator) inferType(val interface{}, prefix string) string {
	switch v := val.(type) {
	case nil:
		return "interface{}"
	case bool:
		return "bool"
	case string:
		return "string"
	case float64:
		if v == float64(int64(v)) {
			return "int"
		}
		return "float64"
	case []interface{}:
		if len(v) == 0 {
			return "[]interface{}"
		}
		return "[]" + g.inferType(v[0], prefix+"Item")
	case map[string]interface{}:
		name := g.uniqueStructName(prefix)
		g.defineStruct(name, v)
		return name
	default:
		return "interface{}"
	}
}

func (g *generator) uniqueStructName(prefix string) string {
	name := exportField(prefix)
	if name == "" {
		g.seq++
		name = fmt.Sprintf("Nested%d", g.seq)
	}
	if g.seen[name] {
		g.seq++
		name = fmt.Sprintf("%s%d", name, g.seq)
	}
	return name
}

func (g *generator) render() string {
	var b strings.Builder
	b.WriteString("package main\n\n")
	for _, name := range g.order {
		b.WriteString(fmt.Sprintf("type %s struct {\n", name))
		for _, f := range g.fields[name] {
			fmt.Fprintf(&b, "\t%s %s `json:\"%s\"`\n", f.goName, f.goType, f.jsonKey)
		}
		b.WriteString("}\n\n")
	}
	return strings.TrimRight(b.String(), "\n") + "\n"
}

func exportField(name string) string {
	name = strings.TrimSpace(name)
	if name == "" {
		return "Field"
	}
	parts := strings.FieldsFunc(name, func(r rune) bool {
		return r == '_' || r == '-' || r == ' '
	})
	for i, part := range parts {
		if part == "" {
			continue
		}
		runes := []rune(part)
		runes[0] = unicode.ToUpper(runes[0])
		parts[i] = string(runes)
	}
	out := strings.Join(parts, "")
	if out == "" {
		return "Field"
	}
	return out
}
