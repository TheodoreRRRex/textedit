package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"textedit/internal/editor"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// toSlashEntries converts DirEntry paths to forward slashes for the frontend.
func toSlashEntries(entries []editor.DirEntry) []editor.DirEntry {
	for i := range entries {
		entries[i].Path = filepath.ToSlash(entries[i].Path)
	}
	return entries
}

type App struct {
	ctx   context.Context
	dirty bool
}

type FileResult struct {
	Path     string `json:"path"`
	Name     string `json:"name"`
	Content  string `json:"content"`
	Size     int64  `json:"size"`
	ReadOnly bool   `json:"readOnly"`
	HasCRLF  bool   `json:"hasCRLF"`
}

func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) OpenFile(path string) (*FileResult, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	if editor.IsBinary(data) {
		return nil, fmt.Errorf("cannot open binary file")
	}

	hasCRLF := strings.Contains(string(data), "\r\n")

	// Check if file is read-only by attempting to open for writing
	readOnly := false
	f, err := os.OpenFile(path, os.O_WRONLY, 0)
	if err != nil {
		readOnly = true
	} else {
		f.Close()
	}

	return &FileResult{
		Path:     filepath.ToSlash(path),
		Name:     filepath.Base(path),
		Content:  string(data),
		Size:     info.Size(),
		ReadOnly: readOnly,
		HasCRLF:  hasCRLF,
	}, nil
}

func (a *App) SaveFile(path string, content string) error {
	return os.WriteFile(path, []byte(content), 0644)
}

func (a *App) ListDirectory(path string) ([]editor.DirEntry, error) {
	entries, err := editor.ListDirectory(path)
	if err != nil {
		return nil, err
	}
	return toSlashEntries(entries), nil
}

func (a *App) PickFileDialog() (string, error) {
	path, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Open File",
		Filters: []runtime.FileFilter{
			{DisplayName: "Config Files", Pattern: "*.env;*.conf;*.cfg;*.ini;*.yaml;*.yml;*.json;*.toml;*.txt;*.properties"},
			{DisplayName: "All Files", Pattern: "*.*"},
		},
	})
	return path, err
}

func (a *App) SaveFileDialog(defaultFilename string) (string, error) {
	path, err := runtime.SaveFileDialog(a.ctx, runtime.SaveDialogOptions{
		Title:           "Save As",
		DefaultFilename: defaultFilename,
		Filters: []runtime.FileFilter{
			{DisplayName: "All Files", Pattern: "*.*"},
		},
	})
	return path, err
}

func (a *App) PickDirectoryDialog() (string, error) {
	path, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Open Folder",
	})
	return path, err
}

func (a *App) CreateFile(dirPath string, name string) (string, error) {
	fullPath := filepath.Join(dirPath, name)
	f, err := os.Create(fullPath)
	if err != nil {
		return "", err
	}
	f.Close()
	return filepath.ToSlash(fullPath), nil
}

func (a *App) SetDirty(dirty bool) {
	a.dirty = dirty
}

func (a *App) GetDefaultDirectory() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return "C:\\"
	}
	return filepath.ToSlash(home)
}

func (a *App) SetWindowTitle(title string) {
	runtime.WindowSetTitle(a.ctx, title)
}
