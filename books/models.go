package books

import (
	"fmt"
	"math"
	"path/filepath"
	"strings"
	"time"
)

var baseSymbols = []string{"KB", "MB", "GB"}

var extNameMap = map[string]string{
	".epub": "EPUB",
	".mobi": "MOBI",
	".pdf":  "PDF",
	".azw":  "AZW",
	".azw3": "AZW3",
}

var extMimeMap = map[string]string{
	".epub": "application/epub+zip",
	".mobi": "application/x-mobipocket-ebook",
	".pdf":  "application/pdf",
	".azw":  "application/vnd.amazon.ebook",
	".azw3": "application/vnd.amazon.mobi8-ebook",
}

type ItemLike interface {
	Name() string
	Path() string
	ModTime() time.Time
	IsDir() bool
}

type FileLike interface {
	ItemLike
	Size() int64
	Ext() string
	FormattedSize() string
	ExtName() string
	MimeType() string
}

type DirLike interface {
	ItemLike
	Children() []ItemLike
	ChildrenHasDir() bool
}

type Item struct {
	name    string
	path    string
	modTime time.Time
	isDir   bool
}

func (i Item) Name() string {
	return i.name
}

func (i Item) Path() string {
	return i.path
}

func (i Item) ModTime() time.Time {
	return i.modTime
}

func (i Item) IsDir() bool {
	return i.isDir
}

type File struct {
	Item
	size int64
	ext  string
}

func (f File) Size() int64 {
	return f.size
}

func (f File) Ext() string {
	return f.ext
}

func (f File) FormattedSize() string {
	raw := float64(f.size)
	base := math.Floor(math.Log(float64(raw)) / math.Log(float64(1024)))
	if base == 0 {
		return fmt.Sprintf("%d %s", f.size, baseSymbols[int64(base)])
	}

	value := float64(raw) / math.Pow(float64(1024), base)

	return fmt.Sprintf("%.1f %s", value, baseSymbols[int64(base)])
}

func (f File) ExtName() string {
	name, ok := extNameMap[f.ext]
	if !ok {
		return "Unknown"
	}

	return name
}

func (f File) MimeType() string {
	mime, ok := extMimeMap[f.ext]
	if !ok {
		return "application/octet-stream"
	}

	return mime
}

func NewFile(name string, path string, modTime time.Time, size int64) File {
	base := filepath.Base(name)
	ext := filepath.Ext(base)
	nameOnly := strings.TrimSuffix(base, ext)

	return File{
		Item: Item{
			name:    nameOnly,
			path:    path,
			modTime: modTime,
			isDir:   false,
		},
		size: size,
		ext:  ext,
	}
}

type Dir struct {
	Item
	children       []ItemLike
	childrenHasDir bool
}

func (d Dir) Children() []ItemLike {
	return d.children
}

func (d Dir) ChildrenHasDir() bool {
	return d.childrenHasDir
}

func (d *Dir) AddChild(child ItemLike) {
	d.children = append(d.children, child)

	if child.IsDir() {
		d.MarkChildrenHasDir()
	}
}

func (d *Dir) MarkChildrenHasDir() {
	d.childrenHasDir = true
}

func NewDir(name string, path string, modTime time.Time) Dir {
	base := filepath.Base(name)

	return Dir{
		Item: Item{
			name:    base,
			path:    path,
			modTime: modTime,
			isDir:   true,
		},
		childrenHasDir: false,
	}
}
