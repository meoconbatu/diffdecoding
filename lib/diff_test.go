package diff

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDiffYAML_OnePath(t *testing.T) {
	tests := []struct {
		name   string
		m1, m2 string
		expect string
	}{
		{"no change", `owner: root:root`, `owner: root:root`, ""},
		{"no change and diff order", "owner: root:root\n   encoding: b64", "encoding: b64\n   owner: root:root", ""},
		{"change one-line field", "owner: root:root", "owner: root:user", "-  owner: root:root\n+  owner: root:user"},
		{"delete one-line field", "owner: root:root", "", "-  owner: root:root"},
		{"delete one-line field", "owner: root:root\n  encoding: b64", "encoding: b64", "-  owner: root:root"},
		{"add one-line field", "", "owner: root:user", "+  owner: root:user"},
		{"no change content field", "encoding: text/plain\n  content: |\n    line1\n    line2", "encoding: text/plain\n  content: |\n     line1\n     line2", ""},
		{"change content field", "encoding: text/plain\n  content: |\n     line1\n     line2", "encoding: text/plain\n  content: |\n     line3", "   content:\n-    line1\n-    line2\n+    line3"},
		{"add content field", "", "encoding: text/plain\n  content: |\n     line1\n     line2", "+  content:\n+    line1\n+    line2\n+  encoding: text/plain"},
		{"delete content field", "encoding: text/plain\n  content: |\n     line1\n     line2", "", "-  content:\n-    line1\n-    line2\n-  encoding: text/plain"},
		{"change and diff order", "owner: root:root\n  encoding: b64", "encoding: base64\n  owner: root:user", "-  encoding: b64\n+  encoding: base64\n-  owner: root:root\n+  owner: root:user"},
		{"preserve single quoted style", ``, `permissions: '0644'`, "+  permissions: '0644'"},
		{"preserve double quoted style", `permissions: '0644'`, `permissions: "0644"`, "-  permissions: '0644'\n+  permissions: \"0644\""},
	}
	d := New()
	d.Config(true)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := d.diffYAML(buildYAML(tt.m1), buildYAML(tt.m2))
			if !assert.Equal(t, buildDiff(tt.expect), actual) {
				fmt.Println(actual)
				fmt.Println(buildDiff(tt.expect))
			}
		})
	}
}

func TestDiffYAML_OnePath_Encoding(t *testing.T) {
	tests := []struct {
		name   string
		m1, m2 string
		expect string
	}{
		{"add", ``, "content: bGluZTEKbGluZTIK\n  encoding: b64", "+  content:\n+    line1\n+    line2\n+  encoding: b64"},
		{"delete", "content: bGluZTEKbGluZTIK\n  encoding: b64", ``, "-  content:\n-    line1\n-    line2\n-  encoding: b64"},
		{"no change", "content: bGluZTEKbGluZTIK\n  encoding: b64", "content: bGluZTEKbGluZTIK\n  encoding: b64", ""},
	}
	d := New()
	d.Config(true)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := d.diffYAML(buildYAML(tt.m1), buildYAML(tt.m2))
			if !assert.Equal(t, buildDiff(tt.expect), actual) {
				fmt.Println(actual)
				fmt.Println()
				fmt.Println(buildDiff(tt.expect))
			}
		})
	}
}
func TestDiffYAML_MultiPath(t *testing.T) {
	tests := []struct {
		name   string
		m1, m2 string
		expect string
	}{
		{"no change", `
write_files:
- owner: root:root
  path: /etc/sysconfig/selinux
- owner: root:root
  path: /etc/sysconfig/selinux2
`, `
write_files:
- owner: root:root
  path: /etc/sysconfig/selinux2
- owner: root:root
  path: /etc/sysconfig/selinux
`, ""},
		{"delete 1 path and add 1 path", `
write_files:
- owner: root:root
  path: /etc/sysconfig/selinux2
- owner: root:root
  path: /etc/sysconfig/selinux
`, `
write_files:
- owner: root:root
  path: /etc/sysconfig/selinux
- owner: root:root
  path: /etc/sysconfig/selinux3
  encoding: text/plain
  content: |
    line1
    line2
`, `-- path: /etc/sysconfig/selinux2
-  owner: root:root
+- path: /etc/sysconfig/selinux3
+  content:
+    line1
+    line2
+  encoding: text/plain
+  owner: root:root`},
	}

	d := New()
	d.Config(true)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := d.diffYAML((tt.m1), (tt.m2))
			if !assert.Equal(t, (tt.expect), actual) {
				fmt.Println(actual)
				fmt.Println()
				fmt.Println(tt.expect)
			}
		})
	}
}
func buildYAML(a string) string {
	return fmt.Sprintf(`
write_files:
- %s
  path: /etc/sysconfig/selinux
`, a)
}
func buildDiff(a string) string {
	if a == "" {
		return ""
	}
	return fmt.Sprintf(` - path: /etc/sysconfig/selinux
%s`, a)
}
