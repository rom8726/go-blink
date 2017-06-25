package str

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"
	"sync"
)

// Format replaces ${param} placeholders with map values.
func Format(s string, params map[string]interface{}) string {
	return NewFormatter(s).Format(params)
}

// Format replaces ${param} placeholders with struct fields.
func FormatStruct(s string, src interface{}) string {
	return NewFormatter(s).FormatStruct(src)
}

// Formatter parses a string with ${param} placeholders,
// and can later replace them with map values or struct fields.
type Formatter struct {
	nodes []formatNode
}

// NewFormatter parses a string with ${param} placeholders and returns a formatter.
func NewFormatter(s string) *Formatter {
	return &Formatter{
		nodes: parseFormatNodes(s),
	}
}

// Format replaces ${param} placeholders with map values.
func (f *Formatter) Format(params map[string]interface{}) string {
	buf := getFormatBuffer()

	for _, node := range f.nodes {
		if !node.isParam {
			buf.WriteString(node.s)
			continue
		}

		value, ok := params[node.s]
		if !ok {
			continue
		}

		v, ok := value.(string)
		if !ok {
			if stringer, ok := value.(fmt.Stringer); ok {
				v = stringer.String()
			} else {
				v = fmt.Sprintf("%v", value)
			}
		}
		buf.WriteString(v)
	}

	s := buf.String()
	releaseFormatBuffer(buf)
	return s
}

// Format replaces ${param} placeholders with struct fields.
func (f *Formatter) FormatStruct(struct0 interface{}) string {
	val := reflect.Indirect(reflect.ValueOf(struct0))
	if val.Kind() != reflect.Struct {
		panic("strs: argument to FormatStruct must be a struct or a pointer to a struct")
	}

	buf := getFormatBuffer()
	for _, node := range f.nodes {
		if !node.isParam {
			buf.WriteString(node.s)
			continue
		}

		field := val.FieldByName(node.s)
		if !field.IsValid() {
			continue
		}

		value := field.Interface()
		v, ok := value.(string)
		if !ok {
			if stringer, ok := value.(fmt.Stringer); ok {
				v = stringer.String()
			} else {
				v = fmt.Sprintf("%v", value)
			}
		}
		buf.WriteString(v)
	}

	s := buf.String()
	releaseFormatBuffer(buf)
	return s
}

// Params returns param names in a format string.
func (f *Formatter) Params() []string {
	params := make([]string, 0, len(f.nodes))
	for _, node := range f.nodes {
		if node.isParam {
			params = append(params, node.s)
		}
	}
	return params
}

type formatNode struct {
	isParam bool
	s       string
}

func parseFormatNodes(s string) []formatNode {
	nodes := []formatNode{}
	for len(s) > 0 {
		start := strings.Index(s, "${")
		if start < 0 {
			nodes = append(nodes, formatNode{false, s})
			break
		}

		end := strings.Index(s, "}")
		if end < 0 {
			nodes = append(nodes, formatNode{false, s})
			break
		}

		before := s[:start]
		if len(before) > 0 {
			nodes = append(nodes, formatNode{false, before})
		}

		param := s[start+2 : end]
		nodes = append(nodes, formatNode{true, param})
		s = s[end+1:]
	}

	return nodes
}

var formatBuffers = sync.Pool{}

func getFormatBuffer() *bytes.Buffer {
	buf := formatBuffers.Get()
	if buf != nil {
		return buf.(*bytes.Buffer)
	}
	return &bytes.Buffer{}
}

func releaseFormatBuffer(buf *bytes.Buffer) {
	buf.Reset()
	formatBuffers.Put(buf)
}
