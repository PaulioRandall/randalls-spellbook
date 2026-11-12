package business

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"

	"github.com/PaulioRandall/randalls-spellbook/pkg/sourcery"
)

func randomEntityId() string {
	return uuid.New().String()
}

type FilePicker struct{}

func init() {
	sourcery.Append(
		sourcery.OptionEntity(&FilePicker{}),
	)
}

type DirEntries struct {
	AbsPath   string
	AbsParent string
	Entries   []DirEntry
}

type DirEntry struct {
	IsDir   bool
	AbsPath string
	Name    string
}

// DoesFileExist returns true if the passed path is a file
// that exists.
func (fp *FilePicker) DoesFileExist(path string) (bool, error) {
	_, e := os.Stat(path)

	if errors.Is(e, os.ErrNotExist) {
		return false, nil
	}

	if e != nil {
		return false, e
	}

	return true, nil
}

// ListFilesInDir lists directory contents of the path.
//
// AI generated.
//
// TODO: Clean.
func (fp *FilePicker) ListFilesInDir(path string) (DirEntries, error) {
	var zero DirEntries

	if strings.TrimSpace(path) == "" {
		return zero, errors.New("path must not be empty")
	}

	if strings.ContainsRune(path, 0) {
		return zero, errors.New("path contains a null byte")
	}

	absDir, err := filepath.Abs(filepath.Clean(path))
	if err != nil {
		return zero, fmt.Errorf("resolving path %q: %w", path, err)
	}

	info, err := os.Stat(absDir)
	if err != nil {
		return zero, fmt.Errorf("accessing %q: %w", absDir, err)
	}

	if !info.IsDir() {
		return zero, fmt.Errorf("%q is not a directory", absDir)
	}

	entries, err := os.ReadDir(absDir)
	if err != nil {
		return zero, fmt.Errorf("reading directory %q: %w", absDir, err)
	}

	files := make([]DirEntry, 0, len(entries))

	for _, entry := range entries {
		// Ignore hidden files.
		if strings.HasPrefix(entry.Name(), ".") {
			continue
		}

		entryPath := filepath.Join(absDir, entry.Name())

		isDir := entry.IsDir()
		// Symlinks report IsDir() == false, so follow them to
		// see whether they point at a directory. A broken link
		// is simply treated as a non-directory entry.
		if entry.Type()&os.ModeSymlink != 0 {
			if target, err := os.Stat(entryPath); err == nil {
				isDir = target.IsDir()
			}
		}

		files = append(files, DirEntry{
			IsDir:   isDir,
			AbsPath: entryPath,
			Name:    entry.Name(),
		})
	}

	dirContent := DirEntries{
		AbsPath:   absDir,
		AbsParent: filepath.Dir(absDir),
		Entries:   files,
	}

	return dirContent, nil
}
