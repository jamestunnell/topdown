package imageresource_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jamestunnell/topdown/imageresource"
	"github.com/jamestunnell/topdown/imageresource/imagetest"
	"github.com/jamestunnell/topdown/resource"
	"github.com/jamestunnell/topdown/resource/restest"
)

func TestLoadImageResourceMissingFile(t *testing.T) {
	_, err := imageresource.Load("missing.png")

	assert.Error(t, err)
}

func TestPNGHappyPath(t *testing.T) {
	testImageResource(t, func(dir string, mgr resource.Manager) {
		path := imagetest.WriteTestPNG(t, dir, 60, 60)

		verifyLoadOK(t, mgr, path)
	})
}

func TestJPGHappyPath(t *testing.T) {
	testImageResource(t, func(dir string, mgr resource.Manager) {
		path := imagetest.WriteTestJPG(t, dir, 60, 60)

		verifyLoadOK(t, mgr, path)
	})
}

func TestJPEGHappyPath(t *testing.T) {
	testImageResource(t, func(dir string, mgr resource.Manager) {
		path := imagetest.WriteTestJPEG(t, dir, 60, 60)

		verifyLoadOK(t, mgr, path)
	})
}

func verifyLoadOK(t *testing.T, mgr resource.Manager, path string) {
	r, err := mgr.Get(filepath.Base(path))

	require.NoError(t, err)
	require.NotNil(t, r)

	_, ok := r.(*imageresource.ImageResource)

	require.True(t, ok)
}

func testImageResource(t *testing.T, test func(dir string, mgr resource.Manager)) {
	dir, err := ioutil.TempDir("", "testimageresource")

	require.NoError(t, err)

	defer os.RemoveAll(dir)

	mgr := restest.SetupManager(t, dir, imageresource.Types()...)

	test(dir, mgr)
}
