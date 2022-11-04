package diff

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"regexp"
	"strings"

	"github.com/kylelemons/godebug/diff"
	"github.com/mitchellh/colorstring"
	"gopkg.in/yaml.v3"
)

// Diff type
type Diff struct {
	color *colorstring.Colorize
}

// New func
func New() *Diff {
	color := &colorstring.Colorize{
		Colors:  colorstring.DefaultColors,
		Disable: false,
		Reset:   true,
	}
	return &Diff{color}
}

// Config func
func (d *Diff) Config(noColor bool) {
	d.color.Disable = noColor
}

// PlanChange func
func (d *Diff) PlanChange(r io.Reader, w io.Writer, noColor bool) error {
	s1, s2, err := parseInput(r)
	if err != nil {
		return err
	}
	d.Config(noColor)
	// diffStr := diffString(deepDecode(s1), deepDecode(s2), color)
	diffStr := d.diffParts(toParts(s1), toParts(s2))
	_, err = w.Write([]byte(diffStr))
	return err
}
func parseInput(r io.Reader) (string, string, error) {
	var err error
	// var s1, s2 string
	// _, err := fmt.Fscanf(r, "%s -> %s", &s1, &s2)
	b, _ := ioutil.ReadAll(r)
	parts := strings.Split(string(b), " -> ")
	s1, s2 := parts[0], parts[1]
	if s1 == "" || s2 == "" {
		err = fmt.Errorf("Input does not match format 'a -> b'")
	}
	return strings.Trim(s1, "\""), strings.Trim(s2, "\""), err
}

// deepDecode decodes then gunzip, then decode base64 encoded content in YAML part (if exists)
func deepDecode(s string) string {
	parts := toParts(s)
	partBodies := make([][]byte, len(parts))
	for i, part := range parts {
		partBodies[i] = part.body
		if part.isYAML() {
			partBodies[i] = decodeYAMLPreserveOrder(part.body)
		}
	}
	return string(bytes.Join(partBodies, []byte("--boundary")))
}
func toParts(s string) []*part {
	if s == "" {
		return nil
	}
	data, err := base64DecodeGunzip(s)
	if err != nil {
		data, _ = base64Decode(s)
	}
	_, parts := parse(data)
	return parts
}
func (d *Diff) diffParts(partsA, partsB []*part) string {
	sb := strings.Builder{}
	for i := 0; i < len(partsA); i++ {
		diff := d.diffPart(*partsA[i], *partsB[i])
		if diff == "" {
			continue
		}
		sb.WriteString(diff)
		sb.WriteString("\n")
	}
	return strings.TrimRight(sb.String(), "\n")
}
func (d *Diff) diffPart(partA, partB part) string {
	var diffStr string
	if partA.isYAML() {
		diffStr = d.diffYAML(string(partA.body), string(partB.body))
	} else {
		diffStr = d.diffString(string(partA.body), string(partB.body))
	}
	if diffStr == "" {
		return ""
	}
	sb := strings.Builder{}
	delimitedLine := fmt.Sprintf("Content-Type: %s\n", partA.header.Get("Content-Type"))
	if val := partA.header.Get("Content-Disposition"); val != "" {
		delimitedLine = fmt.Sprintf("Content-Disposition: %s\n", val)
	}
	sb.WriteString(delimitedLine)
	sb.WriteString(diffStr)
	return sb.String()
}
func (d *Diff) diffString(A, B string) string {
	aLines := strings.Split(A, "\n")
	bLines := strings.Split(B, "\n")

	chunks := diff.DiffChunks(aLines, bLines)

	return formatChunks(chunks, d.color, 2)
}
func (d *Diff) diffYAML(s1, s2 string) string {
	m1, m2 := toMapPreserveStyle(s1), toMapPreserveStyle(s2)
	sb := strings.Builder{}
	for _, obj := range diffMap(m1, m2) {
		sb.WriteString(obj.toString(d.color))
	}
	return sb.String()
}
func toMap(s string) map[string]map[string]interface{} {
	data := []byte(s)
	var v interface{}
	err := yaml.Unmarshal(data, &v)
	if err != nil {
		return nil
	}
	schemas := v.(map[string]interface{})["write_files"]
	if schemas == nil {
		return nil
	}
	pathToContent := make(map[string]map[string]interface{})
	for _, schema := range schemas.([]interface{}) {
		obj := schema.(map[string]interface{})
		pathToContent[obj["path"].(string)] = obj
	}
	return pathToContent
}
func toMapPreserveStyle(s string) map[string]map[string]interface{} {
	document := yaml.Node{}
	err := yaml.Unmarshal([]byte(s), &document)
	if err != nil {
		return nil
	}
	pathToObject := make(map[string]map[string]interface{})
	for _, node := range document.Content {
		if node.Kind == yaml.MappingNode {
			seqNode := getNodeByKey(node, "write_files")
			if seqNode != nil && seqNode.Kind == yaml.SequenceNode {
				for i := 0; i < len(seqNode.Content); i++ {
					mappingNode := seqNode.Content[i]
					object := make(map[string]interface{})
					path := ""
					for i := 0; i < len(mappingNode.Content); i += 2 {
						keyNode := mappingNode.Content[i]
						valueNode := mappingNode.Content[i+1]
						key := keyNode.Value
						value := valueWithStyle(valueNode)
						if key == "path" {
							path = value
						} else {
							object[key] = value
						}
					}
					if _, ok := object["content"]; ok {
						object["content"], _ = decode(toString(object["content"]), toString(object["encoding"]))
					}
					pathToObject[path] = object
				}
			}
		}
	}
	return pathToObject
}
func valueWithStyle(node *yaml.Node) string {
	value := node.Value
	switch node.Style {
	case yaml.DoubleQuotedStyle:
		value = "\"" + value + "\""
	case yaml.SingleQuotedStyle:
		value = "'" + value + "'"
	}
	return value
}
func formatChunks(chunks []diff.Chunk, color *colorstring.Colorize, indentSize int) string {
	buf := new(bytes.Buffer)
	indent := strings.Repeat(" ", indentSize)
	oidx, nidx := 1, 1
	padding := 5
	for _, c := range chunks {
		for _, line := range c.Added {
			fmt.Fprintf(buf, fmt.Sprintf("%*s|%-*d ", padding, " ", padding, nidx)+color.Color(diffActionSymbol(Create)+fmt.Sprintf("%s%s\n", indent, line)))
			nidx++
		}
		for _, line := range c.Deleted {
			fmt.Fprintf(buf, fmt.Sprintf("%*d|%*s ", padding, oidx, padding, " ")+color.Color(diffActionSymbol(Delete)+fmt.Sprintf("%s%s\n", indent, line)))
			oidx++
		}
		delimitedLine := indent + " ...\n"
		if len(c.Equal) > 0 {
			fmt.Fprint(buf, delimitedLine)
		}
		oidx += len(c.Equal)
		nidx += len(c.Equal)
	}
	return buf.String()
}
func diffActionSymbol(action Action) string {
	switch action {
	case Create:
		return "[green]" + string(Create)
	case Delete:
		return "[red]" + string(Delete)
	default:
		return " "
	}
}
func skipLine(n, i int, line string) bool {
	appendLines := 5
	r, _ := regexp.Compile("(^[ ]+[a-z]+): (.+)")
	if i < appendLines || i >= n-appendLines || r.MatchString(line) {
		return false
	}
	return true
}
