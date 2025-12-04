package tree

import (
	"fmt"
	"io"
	"strings"
)

const (
	// Box-drawing characters for tree structure
	TreeBranch  = "├── "
	TreeLast    = "└── "
	TreePipe    = "│   "
	TreeSpace   = "    "
)

// Renderer handles rendering the tree structure
type Renderer struct {
	writer io.Writer
}

// NewRenderer creates a new renderer
func NewRenderer(writer io.Writer) *Renderer {
	return &Renderer{writer: writer}
}

// RenderTree renders the entire tree structure
func (r *Renderer) RenderTree(root *Node, showRoot bool) error {
	if showRoot {
		fmt.Fprintf(r.writer, "%s\n", root.Name)
	}
	return r.renderNode(root, "", true, showRoot)
}

func (r *Renderer) renderNode(node *Node, prefix string, isLast bool, skipRoot bool) error {
	if !skipRoot {
		// Print the node itself
		connector := TreeLast
		if !isLast {
			connector = TreeBranch
		}
		fmt.Fprintf(r.writer, "%s%s%s\n", prefix, connector, node.Name)
	}

	// Process children
	children := node.Children
	for i, child := range children {
		isLastChild := i == len(children)-1
		
		// Determine prefix for child
		childPrefix := prefix
		if !skipRoot {
			if isLast {
				childPrefix += TreeSpace
			} else {
				childPrefix += TreePipe
			}
		}

		// Render child
		err := r.renderNode(child, childPrefix, isLastChild, false)
		if err != nil {
			return err
		}
	}

	return nil
}

// RenderPlain renders the tree without box-drawing characters
func (r *Renderer) RenderPlain(root *Node, showRoot bool) error {
	if showRoot {
		fmt.Fprintf(r.writer, "%s\n", root.Name)
	}
	return r.renderPlainNode(root, "", showRoot)
}

func (r *Renderer) renderPlainNode(node *Node, prefix string, skipRoot bool) error {
	if !skipRoot {
		fmt.Fprintf(r.writer, "%s%s\n", prefix, node.Name)
	}

	// Process children
	for _, child := range node.Children {
		childPrefix := prefix + "  "
		err := r.renderPlainNode(child, childPrefix, false)
		if err != nil {
			return err
		}
	}

	return nil
}

// GetTreeString returns the tree as a string
func GetTreeString(root *Node, showRoot bool) string {
	var sb strings.Builder
	renderer := NewRenderer(&sb)
	renderer.RenderTree(root, showRoot)
	return sb.String()
}

