package tree

import (
	"fmt"
	"io"
	"sort"
	"strings"

	"dtree/internal/stats"
)

// RendererStats extends Renderer with statistics support
type RendererStats struct {
	*Renderer
	showSize   bool
	showDate   bool
	showLong   bool
	sortBy     string
	collectStats bool
}

// NewRendererStats creates a new stats-enabled renderer
func NewRendererStats(writer io.Writer, showSize, showDate, showLong bool, sortBy string, collectStats bool) *RendererStats {
	return &RendererStats{
		Renderer:     NewRenderer(writer),
		showSize:     showSize || showLong,
		showDate:     showDate || showLong,
		showLong:     showLong,
		sortBy:       sortBy,
		collectStats: collectStats,
	}
}

// RenderTreeWithStats renders the tree with statistics
func (r *RendererStats) RenderTreeWithStats(root *Node, showRoot bool) error {
	// Collect statistics if needed
	var fileStats *stats.FileStats
	if r.collectStats {
		fileStats = stats.NewStats()
		r.collectNodeStats(root, fileStats, showRoot)
	}

	// Sort children if needed
	r.sortNode(root)

	// Render header with stats if available
	if fileStats != nil && showRoot {
		totalItems := fileStats.TotalFiles + fileStats.TotalDirs
		fmt.Fprintf(r.writer, "%s (%d items, %s)\n", root.Name, totalItems, stats.FormatSize(fileStats.TotalSize))
	} else if showRoot {
		fmt.Fprintf(r.writer, "%s\n", root.Name)
	}

	// Render tree
	err := r.renderNodeWithStats(root, "", true, showRoot, fileStats)
	if err != nil {
		return err
	}

	// Render footer with summary if stats collected
	if fileStats != nil && r.collectStats {
		fmt.Fprintf(r.writer, "\n")
		fmt.Fprintf(r.writer, "Total: %d files, %d directories, %s\n", 
			fileStats.TotalFiles, fileStats.TotalDirs, stats.FormatSize(fileStats.TotalSize))
	}

	return nil
}

func (r *RendererStats) renderNodeWithStats(node *Node, prefix string, isLast bool, skipRoot bool, fileStats *stats.FileStats) error {
	if !skipRoot {
		connector := TreeLast
		if !isLast {
			connector = TreeBranch
		}

		// Build the line with optional stats
		line := fmt.Sprintf("%s%s%s", prefix, connector, node.Name)

		// Add size and/or date if requested
		if r.showSize || r.showDate {
			var parts []string
			if r.showSize && !node.IsDir {
				parts = append(parts, stats.FormatSizeCompact(node.Size))
			} else if r.showSize && node.IsDir {
				// Calculate directory size
				dirSize := r.calculateDirSize(node)
				parts = append(parts, fmt.Sprintf("[%s]", stats.FormatSizeCompact(dirSize)))
			}
			if r.showDate {
				dateStr := stats.FormatDate(node.ModTime)
				if r.showLong {
					dateStr = stats.FormatDateLong(node.ModTime)
				}
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

	// Process children
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

		err := r.renderNodeWithStats(child, childPrefix, isLastChild, false, fileStats)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *RendererStats) collectNodeStats(node *Node, fileStats *stats.FileStats, skipRoot bool) {
	if !skipRoot {
		if node.IsDir {
			fileStats.AddDir()
		} else {
			fileStats.AddFile(node.Size, node.ModTime)
		}
	}

	for _, child := range node.Children {
		r.collectNodeStats(child, fileStats, false)
	}
}

func (r *RendererStats) calculateDirSize(node *Node) int64 {
	var size int64
	for _, child := range node.Children {
		if child.IsDir {
			size += r.calculateDirSize(child)
		} else {
			size += child.Size
		}
	}
	return size
}

func (r *RendererStats) sortNode(node *Node) {
	if r.sortBy == "" {
		return
	}

	sort.Slice(node.Children, func(i, j int) bool {
		switch r.sortBy {
		case "size":
			return node.Children[i].Size > node.Children[j].Size
		case "date":
			return node.Children[i].ModTime.After(node.Children[j].ModTime)
		case "name":
			fallthrough
		default:
			return strings.ToLower(node.Children[i].Name) < strings.ToLower(node.Children[j].Name)
		}
	})

	// Recursively sort children
	for _, child := range node.Children {
		r.sortNode(child)
	}
}

