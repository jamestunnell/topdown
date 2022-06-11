package sprite_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jamestunnell/topdown/resource"
	"github.com/jamestunnell/topdown/resource/restest"
	"github.com/jamestunnell/topdown/sprite"
	"github.com/jamestunnell/topdown/sprite/spritetest"
)

func TestLoadImageResourceMissingFile(t *testing.T) {
	_, err := sprite.LoadImage("missing.png")

	assert.Error(t, err)
}

func TestPNGHappyPath(t *testing.T) {
	testImageResource(t, func(dir string, mgr resource.Manager) {
		path := spritetest.WriteTestPNG(t, dir, 60, 60)

		verifyLoadOK(t, mgr, path)
	})
}

func TestJPGHappyPath(t *testing.T) {
	testImageResource(t, func(dir string, mgr resource.Manager) {
		path := spritetest.WriteTestJPG(t, dir, 60, 60)

		verifyLoadOK(t, mgr, path)
	})
}

func TestJPEGHappyPath(t *testing.T) {
	testImageResource(t, func(dir string, mgr resource.Manager) {
		path := spritetest.WriteTestJPEG(t, dir, 60, 60)

		verifyLoadOK(t, mgr, path)
	})
}

func verifyLoadOK(t *testing.T, mgr resource.Manager, path string) {
	r, err := mgr.Get(filepath.Base(path))

	require.NoError(t, err)
	require.NotNil(t, r)

	_, ok := r.(*sprite.Image)

	require.True(t, ok)
}

func testImageResource(t *testing.T, test func(dir string, mgr resource.Manager)) {
	dir, err := ioutil.TempDir("", "testspriteimage")

	require.NoError(t, err)

	defer os.RemoveAll(dir)

	mgr := restest.SetupManager(t, dir, sprite.Types()...)

	test(dir, mgr)
}
