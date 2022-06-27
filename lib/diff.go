package diff

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/kylelemons/godebug/diff"
	"github.com/mitchellh/colorstring"
)

func Diff(r io.Reader, w io.Writer, noColor bool) error {
	s1, s2, err := parseInput(r)
	if err != nil {
		return err
	}
	rs1, err := processData(s1)
	if err != nil {
		return err
	}
	rs2, err := processData(s2)
	if err != nil {
		return err
	}
	diffStr := colorDiff(rs1, rs2, noColor)
	_, err = w.Write([]byte(diffStr))
	return err
}
func colorDiff(A, B string, noColor bool) string {
	color := &colorstring.Colorize{
		Colors:  colorstring.DefaultColors,
		Disable: noColor,
		Reset:   false,
	}
	aLines := strings.Split(A, "\n")
	bLines := strings.Split(B, "\n")

	chunks := diff.DiffChunks(aLines, bLines)

	buf := new(bytes.Buffer)
	for _, c := range chunks {
		for _, line := range c.Added {
			fmt.Fprintf(buf, color.Color(fmt.Sprintf("[green]+[reset] %s\n", line)))
		}
		for _, line := range c.Deleted {
			fmt.Fprintf(buf, color.Color(fmt.Sprintf("[red]-[reset] %s\n", line)))
		}
		for _, line := range c.Equal {
			fmt.Fprintf(buf, " %s\n", line)
		}
	}
	return strings.TrimRight(buf.String(), "\n")
}
func parseInput(r io.Reader) (string, string, error) {
	var s1, s2 string
	_, err := fmt.Fscanf(r, "%s -> %s", &s1, &s2)
	if err != nil {
		err = fmt.Errorf("Input does not match format 'a -> b'")
	}
	return strings.Trim(s1, "\""), strings.Trim(s2, "\""), err
}

// processData decodes then gunzip, then decode base64 encoded content in YAML part (if exists)
func processData(s string) (string, error) {
	data, err := base64DecodeGunzip(s)
	if err != nil {
		return "", fmt.Errorf("'%s': %w", s, err)
	}
	sep := extractSeperator(data)

	parts := split(data, sep)
	for i, part := range parts {
		parts[i] = decodeYAMLPreserveOrder(part)
	}

	return string(bytes.Join(parts, sep)), nil
}

// split slices s into all subslices separated by sep and returns a slice of the subslices
// between those separators. If sep is empty, do nothing
func split(s, sep []byte) [][]byte {
	if sep == nil {
		return [][]byte{s}
	}
	return bytes.Split(s, sep)
}

// extractSeperator extracts boundary separator from first line of s
func extractSeperator(s []byte) []byte {
	buf := bytes.NewBuffer(s)
	reader := bufio.NewReader(buf)

	firstLine, err := reader.ReadString('\n')
	if err != nil {
		return nil
	}
	re, _ := regexp.Compile("Content-Type: multipart/mixed; boundary=\"(.+)\"")
	match := re.FindStringSubmatch(firstLine)
	if len(match) == 0 {
		return nil
	}
	return []byte(fmt.Sprintf("--%s\r\n", match[1]))
}
