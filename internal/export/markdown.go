package export

import (
	"fmt"
	"io"
	"os"
	"strings"

	"dtree/internal/tree"
)

// ExportToMarkdown exports the tree to Markdown format
func ExportToMarkdown(root *tree.Node, writer io.Writer) error {
	fmt.Fprintf(writer, "# Directory Tree: %s\n\n", root.Name)
	fmt.Fprintf(writer, "```\n")
	
	var sb strings.Builder
	renderMarkdownNode(root, &sb, "", true, true)
	fmt.Fprint(writer, sb.String())
	
	fmt.Fprintf(writer, "```\n")
	return nil
}

func renderMarkdownNode(node *tree.Node, sb *strings.Builder, prefix string, isLast bool, skipRoot bool) {
	if !skipRoot {
		connector := "└── "
		if !isLast {
			connector = "├── "
		}
		sb.WriteString(prefix)
		sb.WriteString(connector)
		sb.WriteString(node.Name)
		sb.WriteString("\n")
	}

	children := node.Children
	for i, child := range children {
		isLastChild := i == len(children)-1
		
		childPrefix := prefix
		if !skipRoot {
			if isLast {
				childPrefix += "    "
			} else {
				childPrefix += "│   "
			}
		}

		renderMarkdownNode(child, sb, childPrefix, isLastChild, false)
	}
}

// ExportToMarkdownFile exports the tree to a Markdown file
func ExportToMarkdownFile(root *tree.Node, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	return ExportToMarkdown(root, file)
}

