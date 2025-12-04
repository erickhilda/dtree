# dtree

A directory tree visualization tool written in Go, inspired by the Windows `tree` command.

## Features

- **Basic Tree Visualization**: Display directory structures with ASCII box-drawing characters
- **Colorized Output**: Color-coded files and directories (directories in blue, executables in green, etc.)
- **File Statistics**: Show file sizes, modification dates, and summary statistics
- **Export Formats**: Export to JSON, Markdown, or plain text

## Installation

```bash
go build -o dtree .
```

Or install globally:

```bash
go install .
```

## Usage

### Basic Usage

```bash
# Show tree of current directory
dtree

# Show tree of specific directory
dtree /path/to/directory
```

### Options

- `-a, --all`: Show hidden files and directories
- `-d, --depth N`: Limit traversal depth (0 = unlimited)
- `--no-color`: Disable color output
- `-s, --size`: Show file sizes
- `-t, --date`: Show modification dates
- `-l, --long`: Show detailed information (size and date)
- `--sort TYPE`: Sort by name, size, or date
- `--json`: Export as JSON
- `--md`: Export as Markdown
- `--plain`: Export as plain text (no box characters)
- `-o, --output FILE`: Output to file

### Examples

```bash
# Show tree with colors
dtree .

# Show tree with file sizes
dtree --size /usr/local

# Show tree with dates
dtree --date ~/projects

# Show detailed tree
dtree --long /var/log

# Export to JSON
dtree --json -o tree.json .

# Export to Markdown
dtree --md -o tree.md .
```

## Project Structure

```bash
dtree/
├── cmd/
│   └── root.go           # CLI entry, cobra commands
├── internal/
│   ├── tree/
│   │   ├── walker.go     # Directory traversal
│   │   ├── node.go       # Tree node structure
│   │   ├── renderer.go   # ASCII/Unicode output
│   │   ├── renderer_color.go  # Color rendering
│   │   └── renderer_stats.go  # Statistics rendering
│   ├── color/
│   │   └── theme.go      # Color schemes
│   ├── stats/
│   │   └── stats.go      # File statistics
│   └── export/
│       ├── json.go
│       ├── markdown.go
│       └── plain.go
├── main.go
├── go.mod
└── README.md
```
