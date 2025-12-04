package tree

import (
	"os"
	"path/filepath"
	"time"
)

// Node represents a file or directory in the tree
type Node struct {
	Name     string
	Path     string
	IsDir    bool
	IsSymlink bool
	Children []*Node
	Parent   *Node
	Size     int64
	ModTime  time.Time
	Mode     os.FileMode
}

// NewNode creates a new Node from file info
func NewNode(path string, info os.FileInfo) *Node {
	return &Node{
		Name:      info.Name(),
		Path:      path,
		IsDir:     info.IsDir(),
		IsSymlink: info.Mode()&os.ModeSymlink != 0,
		Size:      info.Size(),
		ModTime:   info.ModTime(),
		Mode:      info.Mode(),
		Children:  make([]*Node, 0),
	}
}

// AddChild adds a child node to this node
func (n *Node) AddChild(child *Node) {
	child.Parent = n
	n.Children = append(n.Children, child)
}

// IsHidden returns true if the file/directory is hidden
func (n *Node) IsHidden() bool {
	return len(n.Name) > 0 && n.Name[0] == '.'
}

// GetDepth returns the depth of this node in the tree
func (n *Node) GetDepth() int {
	depth := 0
	parent := n.Parent
	for parent != nil {
		depth++
		parent = parent.Parent
	}
	return depth
}

// GetFullPath returns the full path of the node
func (n *Node) GetFullPath() string {
	return filepath.Join(n.Path, n.Name)
}

