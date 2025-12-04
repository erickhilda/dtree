package export

import (
	"io"
	"strings"

	"dtree/internal/tree"
)

// ExportToPlain exports the tree in plain text format (no box-drawing characters)
func ExportToPlain(root *tree.Node, writer io.Writer) error {
	var sb strings.Builder
	renderPlainNode(root, &sb, "", true)
	_, err := writer.Write([]byte(sb.String()))
	return err
}

func renderPlainNode(node *tree.Node, sb *strings.Builder, prefix string, skipRoot bool) {
	if !skipRoot {
		sb.WriteString(prefix)
		sb.WriteString(node.Name)
		sb.WriteString("\n")
	}

	for _, child := range node.Children {
		childPrefix := prefix + "  "
		renderPlainNode(child, sb, childPrefix, false)
	}
}

