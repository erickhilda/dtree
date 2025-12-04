package tree

import (
	"fmt"
	"io"
	"strings"
	"time"

	"dtree/internal/color"
)

// RendererColor extends Renderer with color support
type RendererColor struct {
	*Renderer
	theme *color.Theme
}

// NewRendererColor creates a new color-enabled renderer
func NewRendererColor(writer io.Writer, theme *color.Theme) *RendererColor {
	return &RendererColor{
		Renderer: NewRenderer(writer),
		theme:    theme,
	}
}

// RenderTree renders the tree with colors
func (r *RendererColor) RenderTree(root *Node, showRoot bool) error {
	if showRoot {
		coloredName := r.theme.Colorize(root.Name, root.IsDir, root.IsSymlink, root.Mode)
		fmt.Fprintf(r.writer, "%s\n", coloredName)
	}
	return r.renderNode(root, "", true, showRoot)
}

func (r *RendererColor) renderNode(node *Node, prefix string, isLast bool, skipRoot bool) error {
	if !skipRoot {
		connector := TreeLast
		if !isLast {
			connector = TreeBranch
		}
		
		coloredName := r.theme.Colorize(node.Name, node.IsDir, node.IsSymlink, node.Mode)
		fmt.Fprintf(r.writer, "%s%s%s\n", prefix, connector, coloredName)
	}

	children := node.Children
	for i, child := range children {
		isLastChild := i == len(children)-1

		childPrefix := prefix
		if !skipRoot {
			if isLast {
				childPrefix += TreeSpace
			} else {
				childPrefix += TreePipe
			}
		}

		err := r.renderNode(child, childPrefix, isLastChild, false)
		if err != nil {
			return err
		}
	}

	return nil
}

// RenderTreeWithStats renders the tree with colors and statistics
func (r *RendererColor) RenderTreeWithStats(root *Node, showRoot bool, showSize, showDate, showLong bool, fileStats interface{}) error {
	// This is a simplified version - full stats rendering would need more integration
	if showRoot {
		coloredName := r.theme.Colorize(root.Name, root.IsDir, root.IsSymlink, root.Mode)
		fmt.Fprintf(r.writer, "%s\n", coloredName)
	}
	return r.renderNodeWithStats(root, "", true, showRoot, showSize, showDate, showLong)
}

func (r *RendererColor) renderNodeWithStats(node *Node, prefix string, isLast bool, skipRoot bool, showSize, showDate, showLong bool) error {
	if !skipRoot {
		connector := TreeLast
		if !isLast {
			connector = TreeBranch
		}

		coloredName := r.theme.Colorize(node.Name, node.IsDir, node.IsSymlink, node.Mode)
		line := fmt.Sprintf("%s%s%s", prefix, connector, coloredName)

		if showSize || showDate {
			var parts []string
			if showSize && !node.IsDir {
				parts = append(parts, formatSizeCompact(node.Size))
			}
			if showDate {
				dateStr := formatDate(node.ModTime, showLong)
				if dateStr != "" {
					parts = append(parts, dateStr)
				}
			}
			if len(parts) > 0 {
				line += "  " + strings.Join(parts, "  ")
			}
		}

		fmt.Fprintf(r.writer, "%s\n", line)
	}

	children := node.Children
	for i, child := range children {
		isLastChild := i == len(children)-1

		childPrefix := prefix
		if !skipRoot {
			if isLast {
				childPrefix += TreeSpace
			} else {
				childPrefix += TreePipe
			}
		}

		err := r.renderNodeWithStats(child, childPrefix, isLastChild, false, showSize, showDate, showLong)
		if err != nil {
			return err
		}
	}

	return nil
}

func formatSizeCompact(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%dB", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f%cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

func formatDate(t time.Time, long bool) string {
	if t.IsZero() {
		return ""
	}
	if long {
		return t.Format("Jan 02 15:04")
	}
	return t.Format("Jan 02")
}

