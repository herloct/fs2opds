package books

import (
	"fmt"
	"io/fs"
	"path/filepath"
)

func ScanDir(fsys fs.FS, path string) (*Dir, error) {
	info, err := fs.Stat(fsys, path)
	if err != nil {
		return nil, fmt.Errorf("failed to read path: %w", err)
	}

	bookDir := NewDir(info.Name(), path, info.ModTime())

	files, err := fs.ReadDir(fsys, path)
	if err != nil {
		return nil, fmt.Errorf("failed to read path children: %w", err)
	}

	for _, f := range files {
		i, err := f.Info()
		if err != nil {
			continue
		}

		var bookItem ItemLike
		if f.IsDir() {
			bookItem = NewDir(f.Name(), filepath.Join(path, f.Name()), i.ModTime())

			// use glob as faster way than using recursive scan, to chech whether the child has directory/ies
			bookItem, _ = globChildDir(fsys, bookItem.(Dir))
		} else {
			bookItem = NewFile(f.Name(), filepath.Join(path, f.Name()), i.ModTime(), i.Size())
		}

		bookDir.AddChild(bookItem)
	}

	return &bookDir, nil
}

func globChildDir(fsys fs.FS, childDir Dir) (Dir, error) {
	matches, err := fs.Glob(fsys, filepath.Join(childDir.Path(), "**/*"))
	if err != nil {
		return childDir, fmt.Errorf("failed to read path children: %w", err)
	}

	if len(matches) > 0 {
		childDir.MarkChildrenHasDir()
	}

	return childDir, nil
}
