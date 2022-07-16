package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func executeCommand(t *testing.T, cmd *cobra.Command, args ...string) (string, error) {
	t.Helper()
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetArgs(args)
	cmd.SetErr(buf)
	err := cmd.Execute()
	return strings.TrimSpace(buf.String()), err
}
func TestPrintHelpWhenNoArgs(t *testing.T) {
	output, err := executeCommand(t, rootCmd)
	if err == nil {
		t.Errorf("Expected error")
	}
	checkStringContains(t, output, "[flags]")
}
func TestPrintUsageWhenUnknownSubCommand(t *testing.T) {
	output, err := executeCommand(t, rootCmd, []string{"unknown"}...)
	if err == nil {
		t.Errorf("Expected error")
	}
	checkStringContains(t, output, "[flags]")

}
func TestRequiredFlagMutuallyExclusive(t *testing.T) {
	output, err := executeCommand(t, rootCmd, []string{"--input", "file1", "--json", "file2"}...)
	if err == nil {
		t.Errorf("Expected error")
	}
	checkStringContains(t, output, "[flags]")
}
func checkStringContains(t *testing.T, got, expected string) {
	if !strings.Contains(got, expected) {
		t.Errorf("Expected to contain: \n %v\nGot:\n %v\n", expected, got)
	}
}
