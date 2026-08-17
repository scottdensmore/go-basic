package corpus

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Variant identifies one byte-distinct original BASIC program.
type Variant struct {
	Path      string
	Alternate bool
}

// Discover returns every main-tree program plus byte-different alternates.
func Discover(root string) ([]Variant, error) {
	var mains []Variant
	mainContents := make(map[string][]byte)
	entries, err := os.ReadDir(root)
	if err != nil {
		return nil, fmt.Errorf("read corpus root: %w", err)
	}
	for _, entry := range entries {
		if !entry.IsDir() || entry.Name() == "00_Alternate_Languages" ||
			!gameDirectory.MatchString(entry.Name()) {
			continue
		}
		files, err := immediateBASICFiles(filepath.Join(root, entry.Name()))
		if err != nil {
			return nil, err
		}
		for _, name := range files {
			path := filepath.ToSlash(filepath.Join(entry.Name(), name))
			contents, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(path)))
			if err != nil {
				return nil, fmt.Errorf("read corpus variant %s: %w", path, err)
			}
			mains = append(mains, Variant{Path: path})
			mainContents[path] = contents
		}
	}
	sort.Slice(mains, func(left, right int) bool {
		return mains[left].Path < mains[right].Path
	})

	var alternates []Variant
	alternateRoot := filepath.Join(root, "00_Alternate_Languages")
	alternateEntries, err := os.ReadDir(alternateRoot)
	if os.IsNotExist(err) {
		return mains, nil
	}
	if err != nil {
		return nil, fmt.Errorf("read alternate corpus root: %w", err)
	}
	for _, entry := range alternateEntries {
		if !entry.IsDir() || !gameDirectory.MatchString(entry.Name()) {
			continue
		}
		files, err := immediateBASICFiles(filepath.Join(alternateRoot, entry.Name()))
		if err != nil {
			return nil, err
		}
		for _, name := range files {
			mainPath := filepath.ToSlash(filepath.Join(entry.Name(), name))
			alternatePath := filepath.ToSlash(filepath.Join("00_Alternate_Languages", entry.Name(), name))
			contents, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(alternatePath)))
			if err != nil {
				return nil, fmt.Errorf("read corpus variant %s: %w", alternatePath, err)
			}
			if main, exists := mainContents[mainPath]; exists && bytes.Equal(main, contents) {
				continue
			}
			alternates = append(alternates, Variant{Path: alternatePath, Alternate: true})
		}
	}
	sort.Slice(alternates, func(left, right int) bool {
		return alternates[left].Path < alternates[right].Path
	})
	return append(mains, alternates...), nil
}

func immediateBASICFiles(directory string) ([]string, error) {
	entries, err := os.ReadDir(directory)
	if err != nil {
		return nil, fmt.Errorf("read corpus directory %s: %w", directory, err)
	}
	var files []string
	for _, entry := range entries {
		if !entry.Type().IsRegular() || !strings.EqualFold(filepath.Ext(entry.Name()), ".bas") {
			continue
		}
		files = append(files, entry.Name())
	}
	sort.Strings(files)
	return files, nil
}
