package export

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"dtree/internal/tree"
)

// JSONNode represents a node in JSON format
type JSONNode struct {
	Name      string      `json:"name"`
	Path      string      `json:"path"`
	Type      string      `json:"type"`
	Size      int64       `json:"size,omitempty"`
	ModTime   string      `json:"modTime,omitempty"`
	Children  []*JSONNode `json:"children,omitempty"`
}

// ExportToJSON exports the tree to JSON format
func ExportToJSON(root *tree.Node, writer io.Writer) error {
	jsonNode := nodeToJSON(root)
	encoder := json.NewEncoder(writer)
	encoder.SetIndent("", "  ")
	return encoder.Encode(jsonNode)
}

func nodeToJSON(node *tree.Node) *JSONNode {
	jsonNode := &JSONNode{
		Name: node.Name,
		Path: node.Path,
	}

	if node.IsSymlink {
		jsonNode.Type = "symlink"
	} else if node.IsDir {
		jsonNode.Type = "directory"
	} else {
		jsonNode.Type = "file"
		jsonNode.Size = node.Size
		jsonNode.ModTime = node.ModTime.Format(time.RFC3339)
	}

	if len(node.Children) > 0 {
		jsonNode.Children = make([]*JSONNode, 0, len(node.Children))
		for _, child := range node.Children {
			jsonNode.Children = append(jsonNode.Children, nodeToJSON(child))
		}
	}

	return jsonNode
}

// ExportToJSONFile exports the tree to a JSON file
func ExportToJSONFile(root *tree.Node, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	return ExportToJSON(root, file)
}

