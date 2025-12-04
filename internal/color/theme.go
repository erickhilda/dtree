package color

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
)

var (
	// Color functions
	DirColor      = color.New(color.FgBlue, color.Bold)
	ExecColor     = color.New(color.FgGreen)
	SymlinkColor  = color.New(color.FgCyan)
	ImageColor    = color.New(color.FgMagenta)
	ArchiveColor  = color.New(color.FgYellow)
	CodeColor     = color.New(color.FgCyan)
	DocColor      = color.New(color.FgYellow)
	DefaultColor  = color.New(color.Reset)
)

// Theme manages color output
type Theme struct {
	enabled bool
}

// NewTheme creates a new theme with color enabled/disabled
func NewTheme(enabled bool) *Theme {
	return &Theme{enabled: enabled}
}

// IsEnabled returns whether colors are enabled
func (t *Theme) IsEnabled() bool {
	return t.enabled
}

// Colorize applies appropriate color to a filename based on its type
func (t *Theme) Colorize(name string, isDir bool, isSymlink bool, mode os.FileMode) string {
	if !t.enabled {
		return name
	}

	if isSymlink {
		return SymlinkColor.Sprint(name)
	}

	if isDir {
		return DirColor.Sprint(name)
	}

	// Check if executable
	if mode&0111 != 0 {
		return ExecColor.Sprint(name)
	}

	// Color by extension
	ext := strings.ToLower(filepath.Ext(name))
	
	switch ext {
	case ".jpg", ".jpeg", ".png", ".gif", ".bmp", ".svg", ".webp", ".ico":
		return ImageColor.Sprint(name)
	case ".zip", ".tar", ".gz", ".bz2", ".xz", ".rar", ".7z":
		return ArchiveColor.Sprint(name)
	case ".go", ".js", ".ts", ".py", ".java", ".cpp", ".c", ".h", ".rs", ".rb", ".php", ".swift", ".kt":
		return CodeColor.Sprint(name)
	case ".md", ".txt", ".doc", ".docx", ".pdf", ".rtf":
		return DocColor.Sprint(name)
	default:
		return name
	}
}

// DisableColors disables color output
func (t *Theme) DisableColors() {
	t.enabled = false
	color.NoColor = true
}

// EnableColors enables color output
func (t *Theme) EnableColors() {
	t.enabled = true
	color.NoColor = false
}

// IsTTY checks if stdout is a terminal
func IsTTY() bool {
	fileInfo, err := os.Stdout.Stat()
	if err != nil {
		return false
	}
	return (fileInfo.Mode() & os.ModeCharDevice) != 0
}

