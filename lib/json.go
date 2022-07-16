package diff

import (
	"fmt"
	"io"
	"io/ioutil"
	"strings"

	tfjson "github.com/hashicorp/terraform-json"
)

// supported arg for resource type
var supportedResourceTypeArgs = map[string]string{
	"aws_instance": "user_data_base64", "aws_launch_template": "user_data", "local_file": "content_base64"}

// PlanJSON func
// reads input from r, extracts supported resource change data,
// decodes the content, then compares and writes diff result to w.
// List supported args for resource type is declared in supportedResourceTypeArgs.
//
// r contains the plan format output by "terraform show -json" command.
func (d *Diff) PlanJSON(r io.Reader, w io.Writer, noColor bool) error {
	var planSchema tfjson.Plan
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	err = planSchema.UnmarshalJSON(b)
	if err != nil {
		return err
	}
	d.Config(noColor)
	var sb strings.Builder
	for _, resourceChange := range planSchema.ResourceChanges {
		arg, ok := supportedResourceTypeArgs[resourceChange.Type]
		if !ok || resourceChange.Change.Actions.NoOp() {
			continue
		}
		before := getArgValue(resourceChange.Change.Before, arg)
		after := getArgValue(resourceChange.Change.After, arg)

		if diffStr := d.diffParts(toParts(before), toParts(after)); diffStr != "" {
			sb.WriteString(d.color.Color(fmt.Sprintf("[cyan]@@ %s[reset]\n", resourceChange.Address)))
			sb.WriteString(diffStr)
			sb.WriteString("\n")
		}
	}
	w.Write([]byte(strings.TrimRight(sb.String(), "\n")))
	return nil
}
func getArgValue(obj interface{}, name string) string {
	if obj == nil || obj.(map[string]interface{})[name] == nil {
		return ""
	}
	return obj.(map[string]interface{})[name].(string)
}
