package diff

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractSeperator(t *testing.T) {
	tests := []struct {
		name   string
		data   []byte
		expect []byte
	}{
		{"default", []byte("Content-Type: multipart/mixed; boundary=\"MIMEBOUNDARY\"\nMIME-Version: 1.0"), []byte("--MIMEBOUNDARY\r\n")},
		{"ASCII", []byte("Content-Type: multipart/mixed; boundary=\"abcABC123+-=*\\\"\nMIME-Version: 1.0"), []byte("--abcABC123+-=*\\\r\n")},
		{"nil", []byte("#cloud-config\n# vim: syntax=yaml\n#"), nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := extractSeperator(tt.data)
			assert.Equal(t, tt.expect, actual)
		})
	}
}
