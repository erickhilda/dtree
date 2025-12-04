package stats

import (
	"fmt"
	"time"
)

// FileStats contains statistics about files
type FileStats struct {
	TotalFiles      int64
	TotalDirs       int64
	TotalSize       int64
	LargestFile     int64
	OldestFile      time.Time
	NewestFile      time.Time
}

// NewStats creates a new empty stats object
func NewStats() *FileStats {
	return &FileStats{
		OldestFile: time.Now(),
		NewestFile: time.Time{},
	}
}

// AddFile adds file statistics
func (s *FileStats) AddFile(size int64, modTime time.Time) {
	s.TotalFiles++
	s.TotalSize += size
	if size > s.LargestFile {
		s.LargestFile = size
	}
	if modTime.Before(s.OldestFile) {
		s.OldestFile = modTime
	}
	if modTime.After(s.NewestFile) {
		s.NewestFile = modTime
	}
}

// AddDir adds directory statistics
func (s *FileStats) AddDir() {
	s.TotalDirs++
}

// FormatSize formats bytes into human-readable format
func FormatSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

// FormatDate formats a date in a compact format
func FormatDate(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("Jan 02")
}

// FormatDateLong formats a date in a longer format
func FormatDateLong(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("Jan 02 15:04")
}

// FormatSizeCompact formats size in a compact way for inline display
func FormatSizeCompact(bytes int64) string {
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

