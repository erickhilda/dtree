package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"dtree/internal/color"
	"dtree/internal/export"
	"dtree/internal/tree"
)

var (
	rootPath    string
	showHidden  bool
	maxDepth    int
	noColor     bool
	showSize    bool
	showDate    bool
	showLong    bool
	sortBy      string
	exportJSON  bool
	exportMD    bool
	exportPlain bool
	outputFile  string
)

var rootCmd = &cobra.Command{
	Use:   "dtree [path]",
	Short: "Display directory tree structure",
	Long: `dtree is a directory tree visualization tool that displays
the structure of directories in a tree format.`,
	Args: cobra.MaximumNArgs(1),
	RunE: runTree,
}

func init() {
	rootCmd.Flags().BoolVarP(&showHidden, "all", "a", false, "Show hidden files and directories")
	rootCmd.Flags().IntVarP(&maxDepth, "depth", "d", 0, "Maximum depth to traverse (0 = unlimited)")
	rootCmd.Flags().BoolVar(&noColor, "no-color", false, "Disable color output")
	rootCmd.Flags().BoolVarP(&showSize, "size", "s", false, "Show file sizes")
	rootCmd.Flags().BoolVarP(&showDate, "date", "t", false, "Show modification dates")
	rootCmd.Flags().BoolVarP(&showLong, "long", "l", false, "Show detailed information (size and date)")
	rootCmd.Flags().StringVar(&sortBy, "sort", "", "Sort by: name, size, or date")
	rootCmd.Flags().BoolVar(&exportJSON, "json", false, "Export as JSON")
	rootCmd.Flags().BoolVar(&exportMD, "md", false, "Export as Markdown")
	rootCmd.Flags().BoolVar(&exportPlain, "plain", false, "Export as plain text (no box characters)")
	rootCmd.Flags().StringVarP(&outputFile, "output", "o", "", "Output to file")
}

func runTree(cmd *cobra.Command, args []string) error {
	// Determine root path
	if len(args) > 0 {
		rootPath = args[0]
	} else {
		var err error
		rootPath, err = os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get current directory: %w", err)
		}
	}

	// Resolve absolute path
	absPath, err := filepath.Abs(rootPath)
	if err != nil {
		return fmt.Errorf("invalid path: %w", err)
	}

	// Build tree
	options := tree.WalkerOptions{
		ShowHidden: showHidden,
		MaxDepth:   maxDepth,
		RootPath:   absPath,
	}

	root, err := tree.WalkTree(absPath, options)
	if err != nil {
		return fmt.Errorf("failed to build tree: %w", err)
	}

	// Setup theme
	themeEnabled := !noColor && color.IsTTY()
	theme := color.NewTheme(themeEnabled)

	// Determine output writer
	var writer *os.File
	if outputFile != "" {
		file, err := os.Create(outputFile)
		if err != nil {
			return fmt.Errorf("failed to create output file: %w", err)
		}
		defer file.Close()
		writer = file
		// Disable colors when writing to file
		theme.DisableColors()
	} else {
		writer = os.Stdout
	}

	// Export formats
	if exportJSON {
		if outputFile != "" {
			return export.ExportToJSONFile(root, outputFile)
		}
		return export.ExportToJSON(root, writer)
	}

	if exportMD {
		if outputFile != "" {
			return export.ExportToMarkdownFile(root, outputFile)
		}
		return export.ExportToMarkdown(root, writer)
	}

	if exportPlain {
		return export.ExportToPlain(root, writer)
	}

	// Render tree
	if showSize || showDate || showLong || sortBy != "" {
		// Use stats renderer
		renderer := tree.NewRendererStats(writer, showSize, showDate, showLong, sortBy, true)
		return renderer.RenderTreeWithStats(root, true)
	}

	// Use color renderer if colors enabled
	if theme.IsEnabled() {
		renderer := tree.NewRendererColor(writer, theme)
		return renderer.RenderTree(root, true)
	}

	// Use basic renderer
	renderer := tree.NewRenderer(writer)
	return renderer.RenderTree(root, true)
}

// Execute runs the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

