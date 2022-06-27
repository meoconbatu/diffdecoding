package diff

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

// decode encoded content key in write_files config
func decodeYAML(data []byte) []byte {
	var v interface{}
	err := yaml.Unmarshal(data, &v)
	if err != nil {
		return data
	}
	schemas := v.(map[string]interface{})["write_files"]
	if schemas == nil {
		return data
	}
	for _, schema := range schemas.([]interface{}) {
		obj := schema.(map[string]interface{})
		obj["content"], _ = decode(obj["content"].(string), obj["encoding"].(string))
	}
	out, err := yaml.Marshal(v)
	if err != nil {
		fmt.Println(err)
	}
	return out
}

// decodeYAMLPreserveOrder same as decodeYAML, and preserve comments and order of keys in YAML doc
func decodeYAMLPreserveOrder(data []byte) []byte {
	node := yaml.Node{}
	err := yaml.Unmarshal(data, &node)
	if err != nil {
		return data
	}

	decodeContent(&node)

	out, _ := yaml.Marshal(&node)

	return out
}
func decodeContent(document *yaml.Node) {
	if document.Kind != yaml.DocumentNode {
		return
	}
	for _, node := range document.Content {
		if node.Kind == yaml.MappingNode {
			seqNode := getNodeByKey(node, "write_files")
			if seqNode != nil && seqNode.Kind == yaml.SequenceNode {
				for i := 0; i < len(seqNode.Content); i++ {
					encodingNode := getNodeByKey(seqNode.Content[i], "encoding")
					contentNode := getNodeByKey(seqNode.Content[i], "content")
					if encodingNode != nil && encodingNode.Kind == yaml.ScalarNode {
						contentNode.Value, _ = decode(contentNode.Value, encodingNode.Value)
					}
				}
			}
		}
	}
}
func getNodeByKey(node *yaml.Node, key string) *yaml.Node {
	if node.Kind != yaml.MappingNode {
		return nil
	}
	for i := 0; i < len(node.Content); i += 2 {
		keyNode := node.Content[i]
		if keyNode.Value == key && keyNode.Kind == yaml.ScalarNode {
			return node.Content[i+1]
		}
	}
	return nil
}
