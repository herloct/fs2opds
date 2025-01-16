package books

import (
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewFile(t *testing.T) {
	name := "Some Book"
	ext := ".epub"
	nameWithExt := strings.Join([]string{name, ext}, "")
	path := filepath.Join("books", nameWithExt)
	modTime := time.Now()
	size := int64(1024 * 4)

	file := NewFile(nameWithExt, path, modTime, size)

	assert.Equal(t, name, file.Name())
	assert.Equal(t, path, file.Path())
	assert.Equal(t, modTime, file.ModTime())
	assert.Equal(t, false, file.IsDir())
	assert.Equal(t, size, file.Size())
	assert.Equal(t, ext, file.Ext())

	assert.Implements(t, (*ItemLike)(nil), file)
	assert.Implements(t, (*FileLike)(nil), file)
	assert.NotImplements(t, (*DirLike)(nil), file)
}

func TestFileFormattedSize(t *testing.T) {
	var tests = []struct {
		message       string
		size          int64
		formattedSize string
	}{
		{"Kilobytes", 768, "768 KB"},
		{"Megabytes", 5643, "5.5 MB"},
		{"Gigabytes", 9986877, "9.5 GB"},
	}

	for _, tt := range tests {
		t.Run(tt.message, func(t *testing.T) {
			file := File{size: tt.size}
			if got := file.FormattedSize(); got != tt.formattedSize {
				t.Errorf("File.FormattedSize() with %d should be %v, got %v", tt.size, tt.formattedSize, got)
			}
		})
	}
}

func TestFileExtName(t *testing.T) {
	var tests = []struct {
		message string
		ext     string
		extName string
	}{
		{"EPUB", ".epub", "EPUB"},
		{"MOBI", ".mobi", "MOBI"},
		{"PDF", ".pdf", "PDF"},
		{"AZW", ".azw", "AZW"},
		{"AZW3", ".azw3", "AZW3"},
		{"Unknown", ".mkv", "Unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.message, func(t *testing.T) {
			file := File{ext: tt.ext}
			if got := file.ExtName(); got != tt.extName {
				t.Errorf("File.ExtName() with %s file should be %v, got %v", tt.message, tt.extName, got)
			}
		})
	}
}

func TestFileMimeType(t *testing.T) {
	var tests = []struct {
		message string
		ext     string
		mime    string
	}{
		{"EPUB", ".epub", "application/epub+zip"},
		{"MOBI", ".mobi", "application/x-mobipocket-ebook"},
		{"PDF", ".pdf", "application/pdf"},
		{"AZW", ".azw", "application/vnd.amazon.ebook"},
		{"AZW3", ".azw3", "application/vnd.amazon.mobi8-ebook"},
		{"Unknown", ".mkv", "application/octet-stream"},
	}

	for _, tt := range tests {
		t.Run(tt.message, func(t *testing.T) {
			file := File{ext: tt.ext}
			if got := file.MimeType(); got != tt.mime {
				t.Errorf("File.MimeType() with %s file should be %v, got %v", tt.message, tt.mime, got)
			}
		})
	}
}

func TestNewDir(t *testing.T) {
	dirName := "Some Directory"
	path := filepath.Join("books", dirName)
	modTime := time.Now()

	dir := NewDir(dirName, path, modTime)

	assert.Equal(t, dirName, dir.Name())
	assert.Equal(t, path, dir.Path())
	assert.Equal(t, modTime, dir.ModTime())
	assert.Equal(t, true, dir.IsDir())

	assert.Implements(t, (*ItemLike)(nil), dir)
	assert.Implements(t, (*DirLike)(nil), dir)
	assert.NotImplements(t, (*FileLike)(nil), dir)

	dir.AddChild(NewFile("Some File 1.ext", filepath.Join(path, "Some File 1.ext"), time.Now(), int64(1024*5)))
	dir.AddChild(NewFile("Some File 2.ext", filepath.Join(path, "Some File 2.ext"), time.Now(), int64(1024*5)))
	dir.AddChild(NewFile("Some File 3.ext", filepath.Join(path, "Some File 3.ext"), time.Now(), int64(1024*5)))

	assert.Len(t, dir.Children(), 3)
	assert.Equal(t, false, dir.ChildrenHasDir())

	dir.AddChild(NewDir("Another Directory", filepath.Join(path, "Another Directory"), time.Now()))

	assert.Len(t, dir.Children(), 4)
	assert.Equal(t, true, dir.ChildrenHasDir())

	for _, c := range dir.Children() {
		assert.Implements(t, (*ItemLike)(nil), c)

		if c.IsDir() {
			assert.Implements(t, (*DirLike)(nil), c.(DirLike))
		} else {
			assert.Implements(t, (*FileLike)(nil), c.(FileLike))
		}
	}
}
