package books

import (
	"io/fs"
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/assert"
)

func TestScanDir(t *testing.T) {
	fsys := fstest.MapFS{
		"Book Title.epub":                                     {},
		"Another Book Title.epub":                             {},
		"Series I":                                            {Mode: fs.ModeDir},
		"Series I/Series - Vol 01.epub":                       {},
		"Series I/Series - Vol 02.epub":                       {},
		"Series II":                                           {Mode: fs.ModeDir},
		"Series II/Series II - Vol 01.epub":                   {},
		"Series II/Series II - Vol 02.epub":                   {},
		"Series II/Specials":                                  {Mode: fs.ModeDir},
		"Series II/Specials/Series II Specials - Vol 01.epub": {},
	}

	item, err := ScanDir(fsys, ".")
	assert.Nil(t, err)

	assert.Equal(t, item.Name(), ".")
	assert.Equal(t, item.IsDir(), true)

	assert.Len(t, item.Children(), 4)

	child := item.Children()[2]
	assert.False(t, child.(Dir).ChildrenHasDir(), "%s should have no directory/ies as child", child.Name())

	child = item.Children()[3]
	assert.True(t, child.(Dir).ChildrenHasDir(), "%s should have directory/ies as child", child.Name())

	item, err = ScanDir(fsys, "Series I")
	assert.Nil(t, err)

	assert.Equal(t, item.Name(), "Series I")
	assert.Equal(t, item.IsDir(), true)

	assert.Len(t, item.Children(), 2)
}
