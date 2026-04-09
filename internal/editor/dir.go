package editor

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type DirEntry struct {
	Name  string `json:"name"`
	Path  string `json:"path"`
	IsDir bool   `json:"isDir"`
	Ext   string `json:"ext"`
	Size  int64  `json:"size"`
}

// ListDirectory reads a directory and returns entries sorted with directories
// first, then alphabetically. Hidden directories (starting with .) are skipped.
func ListDirectory(path string) ([]DirEntry, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	var result []DirEntry
	for _, entry := range entries {
		name := entry.Name()

		if strings.HasPrefix(name, ".") && entry.IsDir() {
			continue
		}

		info, err := entry.Info()
		if err != nil {
			continue
		}

		fullPath := filepath.Join(path, name)

		result = append(result, DirEntry{
			Name:  name,
			Path:  fullPath,
			IsDir: entry.IsDir(),
			Ext:   strings.TrimPrefix(filepath.Ext(name), "."),
			Size:  info.Size(),
		})
	}

	sort.Slice(result, func(i, j int) bool {
		if result[i].IsDir != result[j].IsDir {
			return result[i].IsDir
		}
		return strings.ToLower(result[i].Name) < strings.ToLower(result[j].Name)
	})

	return result, nil
}
