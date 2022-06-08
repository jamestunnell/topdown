package fileindex_test

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jamestunnell/topdown/fileindex"
	"github.com/jamestunnell/topdown/sliceutil"
)

func TestIndexMissingRoot(t *testing.T) {
	idx := fileindex.New("/missing/dir")

	assert.Error(t, idx.Scan())
}

func TestIndexNonDirRoot(t *testing.T) {
	path := writeTempFile(t, "", "", "")
	idx := fileindex.New(path)

	assert.Error(t, idx.Scan())
}

func TestIndex(t *testing.T) {
	testIndex(t, "empty root", 0, func(root string, subDirs []string) {
		idx := fileindex.New(root)

		require.NoError(t, idx.Scan())

		assert.Less(t, time.Since(idx.LastScanComplete()), time.Millisecond)

		assert.Empty(t, idx.Paths())
		assert.Empty(t, idx.PathDirs())
		assert.Empty(t, idx.PathExts())
	})

	testIndex(t, "root contains one .json file", 0, func(root string, subDirs []string) {
		path := writeTempFile(t, root, ".json", `{"abc":123}`)

		idx := fileindex.New(root)

		require.NoError(t, idx.Scan())

		assert.Equal(t, []string{path}, idx.Paths())
		assert.Equal(t, []string{".json"}, idx.PathExts())
		assert.Equal(t, []string{idx.RootDir()}, idx.PathDirs())
	})

	testIndex(t, "files in subdurs", 3, func(root string, subDirs []string) {
		exts := []string{".json", ".sprite", ".level"}
		paths := make([]string, 3)
		for i := 0; i < 3; i++ {
			content := fmt.Sprintf(`{"id":"%d"}`, i)
			paths[i] = writeTempFile(t, subDirs[i], exts[i], content)
		}

		idx := fileindex.New(root)

		require.NoError(t, idx.Scan())

		assert.ElementsMatch(t, subDirs, idx.PathDirs())
		assert.ElementsMatch(t, exts, idx.PathExts())
		assert.ElementsMatch(t, paths, idx.Paths())

		assert.Equal(t, []string{paths[0]}, idx.Paths(fileindex.FilterByDir(subDirs[0])))
		assert.Equal(t, []string{paths[1]}, idx.Paths(fileindex.FilterByDir(subDirs[1])))
		assert.Equal(t, []string{paths[2]}, idx.Paths(fileindex.FilterByDir(subDirs[2])))

		assert.Equal(t, []string{paths[0]}, idx.Paths(fileindex.FilterByExt(exts[0])))
		assert.Equal(t, []string{paths[1]}, idx.Paths(fileindex.FilterByExt(exts[1])))
		assert.Equal(t, []string{paths[2]}, idx.Paths(fileindex.FilterByExt(exts[2])))

		pathsFoundCompound := idx.Paths(
			fileindex.FilterByExt(subDirs[0]),
			fileindex.FilterByDir(exts[2]),
		)

		assert.Empty(t, pathsFoundCompound)
	})
}

func testIndex(t *testing.T, name string, nSubDirs int, test func(root string, subDirs []string)) {
	rootDir, err := ioutil.TempDir("", "testindex")

	require.NoError(t, err)

	subDirs := sliceutil.Make(nSubDirs, func(i int) string {
		dir, err := ioutil.TempDir(rootDir, "testindex")

		require.NoError(t, err)

		return dir
	})

	for _, dir := range subDirs {
		defer os.RemoveAll(dir)
	}

	defer os.RemoveAll(rootDir)

	t.Run(name, func(t *testing.T) {
		test(rootDir, subDirs)
	})
}

func writeTempFile(t *testing.T, dir, ext, content string) string {
	f, err := ioutil.TempFile(dir, fmt.Sprintf("indextest*%s", ext))

	require.NoError(t, err)

	require.NoError(t, ioutil.WriteFile(f.Name(), []byte(content), fs.ModePerm))

	return f.Name()
}
