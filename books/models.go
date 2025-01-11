package books

import (
	"path/filepath"
	"strings"
	"time"
)

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
