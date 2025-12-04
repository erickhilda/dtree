package tree

import (
	"os"
	"path/filepath"
)

// WalkerOptions contains options for directory walking
type WalkerOptions struct {
	ShowHidden bool
	MaxDepth   int
	RootPath   string
}

// WalkTree builds a tree structure from the given root path
func WalkTree(rootPath string, options WalkerOptions) (*Node, error) {
	absPath, err := filepath.Abs(rootPath)
	if err != nil {
		return nil, err
	}

	info, err := os.Stat(absPath)
	if err != nil {
		return nil, err
	}

	root := NewNode(filepath.Dir(absPath), info)
	root.Name = filepath.Base(absPath)
	root.Path = filepath.Dir(absPath)

	err = walkDirectory(root, absPath, options, 0)
	if err != nil {
		return nil, err
	}

	return root, nil
}

func walkDirectory(parent *Node, dirPath string, options WalkerOptions, currentDepth int) error {
	// Check depth limit
	if options.MaxDepth > 0 && currentDepth >= options.MaxDepth {
		return nil
	}

	entries, err := os.ReadDir(dirPath)
	if err != nil {
		// Skip directories we can't read (permission denied, etc.)
		return nil
	}

	for _, entry := range entries {
		// Skip hidden files if not showing them
		if !options.ShowHidden && entry.Name()[0] == '.' {
			continue
		}

		info, err := entry.Info()
		if err != nil {
			// Skip files we can't stat
			continue
		}

		fullPath := filepath.Join(dirPath, entry.Name())
		node := NewNode(dirPath, info)

		// Handle symlinks
		if entry.Type()&os.ModeSymlink != 0 {
			node.IsSymlink = true
			// Try to resolve symlink target
			resolved, err := filepath.EvalSymlinks(fullPath)
			if err == nil {
				info, err := os.Stat(resolved)
				if err == nil {
					node.IsDir = info.IsDir()
					node.Size = info.Size()
					node.ModTime = info.ModTime()
				}
			} else {
				// Symlink target doesn't exist, keep original info
				node.IsDir = false
			}
		}

		parent.AddChild(node)

		// Recursively walk subdirectories
		if node.IsDir && !node.IsSymlink {
			err := walkDirectory(node, fullPath, options, currentDepth+1)
			if err != nil {
				// Continue even if subdirectory can't be read
				continue
			}
		} else if node.IsDir && node.IsSymlink {
			// Follow symlink if it points to a directory
			resolved, err := filepath.EvalSymlinks(fullPath)
			if err == nil {
				info, err := os.Stat(resolved)
				if err == nil && info.IsDir() {
					err := walkDirectory(node, resolved, options, currentDepth+1)
					if err != nil {
						continue
					}
				}
			}
		}
	}

	return nil
}

