# TextEdit

A lightweight text editor for editing `.env`, `.conf`, and other config files on Windows servers. Available as a GUI desktop app (Wails v2 + Svelte + CodeMirror 6) or a TUI for terminal use (Bubbletea).

## GUI

### Development

```bash
wails dev
```

### Build

```bash
wails build
```

The compiled binary will be at `build/bin/textedit.exe`.

## TUI

### Build

```bash
go build -o ./build/bin/textedit.exe ./cmd/tui/
```

### Usage

```bash
# Open current directory
textedit

# Open a folder
textedit C:\path\to\folder

# Open a file directly
textedit C:\path\to\file.txt
```

### Install

Copy `textedit.exe` to a directory on your PATH (e.g. `C:\Windows\System32` or a custom bin directory).
