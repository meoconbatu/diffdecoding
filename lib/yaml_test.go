package diff

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func TestDecodeContent_DecodeBase64(t *testing.T) {
	tests := []struct {
		name           string
		encoding       string
		content        string
		encodedContent string
	}{
		{"b64", "b64", "test b64 content", base64Encode([]byte("test b64 content"))},
		{"base64", "base64", "test base64 content", base64Encode([]byte("test base64 content"))},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			documentNode := createDocumentNode(tt.encoding, tt.encodedContent)
			decodeContent(documentNode)

			expectedDocumentNode := createDocumentNode(tt.encoding, tt.content)
			assert.Equal(t, expectedDocumentNode, documentNode)
		})
	}
}

func TestDecodeContent_Gunzip(t *testing.T) {
	tests := []struct {
		name     string
		encoding string
		content  string
	}{
		{"gzip", "gzip", "test gzip content"},
		{"gz", "gz", "test gzip content"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			compressedContent, _ := gzipData([]byte(tt.content))
			documentNode := createDocumentNode(tt.encoding, string(compressedContent))
			decodeContent(documentNode)

			expectedDocumentNode := createDocumentNode(tt.encoding, tt.content)

			assert.Equal(t, expectedDocumentNode, documentNode)
		})
	}
}
func createDocumentNode(encoding, content string) *yaml.Node {
	return &yaml.Node{Kind: yaml.DocumentNode, Content: []*yaml.Node{
		{Kind: yaml.MappingNode, Content: []*yaml.Node{
			{Kind: yaml.ScalarNode, Value: "write_files"},
			{Kind: yaml.SequenceNode, Content: []*yaml.Node{
				{Kind: yaml.MappingNode, Content: []*yaml.Node{
					{Kind: yaml.ScalarNode, Value: "encoding"},
					{Kind: yaml.ScalarNode, Value: encoding},
					{Kind: yaml.ScalarNode, Value: "content"},
					{Kind: yaml.ScalarNode, Value: content},
				}},
				{Kind: yaml.MappingNode, Content: []*yaml.Node{
					{Kind: yaml.ScalarNode, Value: "encoding"},
					{Kind: yaml.ScalarNode, Value: "text/plain"},
					{Kind: yaml.ScalarNode, Value: "content"},
					{Kind: yaml.ScalarNode, Value: "test content"},
				}},
				{Kind: yaml.MappingNode, Content: []*yaml.Node{
					{Kind: yaml.ScalarNode, Value: "encoding"},
					{Kind: yaml.ScalarNode, Value: encoding},
					{Kind: yaml.ScalarNode, Value: "content"},
					{Kind: yaml.ScalarNode, Value: content},
				}},
				{Kind: yaml.MappingNode, Content: []*yaml.Node{
					{Kind: yaml.ScalarNode, Value: "content"},
					{Kind: yaml.ScalarNode, Value: "test content"},
				}},
			}},
		}},
	}}
}

func TestAccYAML_DecodeContent(t *testing.T) {
	tests := []struct {
		name            string
		encoding        string
		content         string
		expectedContent string
	}{
		{"b64", "b64", "testDecodeYAML", base64Encode([]byte("testDecodeYAML"))},
		{"no content decoding", "text/plain", "testDecodeYAML", "testDecodeYAML"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := fmt.Sprintf(testDecodeYAML, tt.encoding, tt.expectedContent)
			actual := decodeYAMLPreserveOrder([]byte(in))

			out := fmt.Sprintf(testDecodeYAML, tt.encoding, tt.content)
			expect := formatYAMLDocument([]byte(out))
			assert.Equal(t, string(expect), string(actual))
		})
	}
}

func TestAccYAML_DecodeMultiContents(t *testing.T) {
	tests := []struct {
		name            string
		encoding        []string
		content         []string
		expectedContent []string
	}{
		{"b64", []string{"b64", "b64"}, []string{"testDecodeYAML", "testDecodeYAML"}, []string{base64Encode([]byte("testDecodeYAML")), base64Encode([]byte("testDecodeYAML"))}},
		{"b64", []string{"text/plain", "text/plain"}, []string{"testDecodeYAML", "testDecodeYAML"}, []string{("testDecodeYAML"), ("testDecodeYAML")}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := fmt.Sprintf(testDecodeYAMLMultiFiles, tt.encoding[0], tt.expectedContent[0], tt.encoding[1], tt.expectedContent[1])
			actual := decodeYAMLPreserveOrder([]byte(in))

			out := fmt.Sprintf(testDecodeYAMLMultiFiles, tt.encoding[0], tt.content[0], tt.encoding[1], tt.content[1])
			expect := formatYAMLDocument([]byte(out))
			assert.Equal(t, expect, actual)
		})
	}
}

const testDecodeYAML = `
write_files:
- encoding: %s
  content: %s
  owner: root:root
  path: /etc/sysconfig/selinux
  permissions: '0644'
- content: |
    15 * * * * root ship_logs
  path: /etc/crontab
  append: true
`
const testDecodeYAMLMultiFiles = `
write_files:
- encoding: %s
  content: %s
  owner: root:root
  path: /etc/sysconfig/selinux
  permissions: '0644'
- content: |
    15 * * * * root ship_logs
  path: /etc/crontab
  append: true
- encoding: %s
  content: %s
  owner: root:root
  path: /etc/sysconfig/selinux
  permissions: '0644'
`

func formatYAMLDocument(data []byte) []byte {
	node := yaml.Node{}
	err := yaml.Unmarshal(data, &node)
	if err != nil {
		return data
	}
	out, _ := yaml.Marshal(&node)
	return out
}
