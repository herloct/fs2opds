package books

import (
	"log"
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
