package diff

import (
	"bytes"
	"fmt"
	"sort"
	"strings"

	"github.com/kylelemons/godebug/diff"
	"github.com/mitchellh/colorstring"
)

// Action is a action type for a resource change.
type Action rune

const (
	// NoOp denotes a no-op operation.
	NoOp Action = 0
	// Create denotes a create operation.
	Create Action = '+'
	// Update denotes a update operation.
	Update Action = '~'
	// Delete denotes a delete operation.
	Delete Action = '-'
)

type object struct {
	key       string
	diffType  Action
	keychunks []keychunk
}

func (obj object) toString(color *colorstring.Colorize) string {
	buf := new(bytes.Buffer)
	fmt.Fprintf(buf, color.Color(fmt.Sprintf(diffActionSymbol(obj.diffType)+"- path: %s\n", obj.key)))
	for _, chunks := range obj.keychunks {
		fmt.Fprintf(buf, chunks.toString(color, 2))
	}
	return buf.String()
}

type keychunk struct {
	key          string
	chunks       []diff.Chunk
	isBlockStyle bool
	diffType     Action
}

func (kc keychunk) toString(color *colorstring.Colorize, indentSize int) string {
	sb := strings.Builder{}
	indent := strings.Repeat(" ", indentSize)
	if !kc.isBlockStyle {
		for _, c := range kc.chunks {
			for _, line := range c.Added {
				sb.WriteString(color.Color(diffActionSymbol(Create) + fmt.Sprintf("%s%s: %s\n", indent, kc.key, line)))
			}
			for _, line := range c.Deleted {
				sb.WriteString(color.Color(diffActionSymbol(Delete) + fmt.Sprintf("%s%s: %s\n", indent, kc.key, line)))
			}
		}
	} else {
		sb.WriteString(color.Color(diffActionSymbol(kc.diffType) + fmt.Sprintf("%s%s:\n", indent, kc.key)))
		sb.WriteString(formatChunks(kc.chunks, color, indentSize+2))
	}
	return sb.String()
}
func diffMapToChunks(m1, m2 map[string]interface{}) []keychunk {
	diffs := make([]keychunk, 0)
	for k1, v1 := range m1 {
		// at this time, treat every value as string to compare
		a := strings.Split(strings.TrimRight(toString(v1), "\n"), "\n")
		isBlockStyle := len(a) > 1
		if v2, ok := m2[k1]; ok {
			b := strings.Split(strings.TrimRight(toString(v2), "\n"), "\n")
			chunks := diff.DiffChunks(a, b)
			if len(chunks) > 0 {
				diffs = append(diffs, keychunk{k1, chunks, isBlockStyle || len(b) > 1, Update})
			}
		} else {
			chunks := diff.DiffChunks(a, nil)
			diffs = append(diffs, keychunk{k1, chunks, isBlockStyle, Delete})
		}
	}
	for k2, v2 := range m2 {
		if _, ok := m1[k2]; !ok {
			b := strings.Split(strings.TrimRight(v2.(string), "\n"), "\n")
			chunks := diff.DiffChunks(nil, b)
			diffs = append(diffs, keychunk{k2, chunks, len(b) > 1, Create})
		}
	}
	sort.SliceStable(diffs, func(i, j int) bool {
		return diffs[i].key < diffs[j].key
	})
	return diffs
}

func toString(in interface{}) string {
	return fmt.Sprint(in)
}
func diffMap(m1, m2 map[string]map[string]interface{}) []object {
	objs := make([]object, 0)
	for k1, v1 := range m1 {
		if v2, ok := m2[k1]; ok {
			chunks := diffMapToChunks(v1, v2)
			if len(chunks) > 0 {
				objs = append(objs, object{k1, Update, chunks})
			}
		} else {
			objs = append(objs, object{k1, Delete, diffMapToChunks(v1, nil)})
		}
	}
	for k2, v2 := range m2 {
		if _, ok := m1[k2]; !ok {
			objs = append(objs, object{k2, Create, diffMapToChunks(nil, v2)})
		}
	}
	sort.SliceStable(objs, func(i, j int) bool {
		return objs[i].key < objs[j].key
	})
	return objs
}
