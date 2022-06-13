package sprite_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jamestunnell/topdown"
	"github.com/jamestunnell/topdown/jsonfile"
	"github.com/jamestunnell/topdown/resource"
	"github.com/jamestunnell/topdown/resource/restest"
	"github.com/jamestunnell/topdown/sprite"
	"github.com/jamestunnell/topdown/sprite/spritetest"
)

func TestImageSet(t *testing.T) {
	dir, err := ioutil.TempDir("", "testspriteset")

	require.NoError(t, err)

	defer os.RemoveAll(dir)

	subImages := []*sprite.Sprite{
		{ID: "a", Origin: topdown.Pt(0, 0), Size: topdown.Sz(16, 16)},
		{ID: "b", Origin: topdown.Pt(0, 16), Size: topdown.Sz(16, 16)},
	}
	ss := sprite.NewSheet("missing.png", subImages...)
	partialPath := "bad.spritesheet"

	require.NoError(t, jsonfile.Write(filepath.Join(dir, partialPath), ss))

	mgr := restest.SetupManager(t, dir, sprite.Types()...)

	// fails due to missing image
	r, err := resource.GetAs[*sprite.Sheet](mgr, partialPath)

	assert.Error(t, err)
	assert.Nil(t, r)

	imagePath := spritetest.WriteTestPNG(t, dir, 128, 128)

	ss.ImageRef = filepath.Base(imagePath)
	partialPath = "good.spritesheet"

	require.NoError(t, jsonfile.Write(filepath.Join(dir, partialPath), ss))

	// should not fail this time with a good source image
	r, err = resource.GetAs[*sprite.Sheet](mgr, partialPath)

	require.NoError(t, err)

	ss2, err := resource.As[*sprite.Sheet](r)

	require.NoError(t, err)

	sprt, found := ss2.FindSpriteByOrigin(topdown.Pt(0, 0))

	assert.True(t, found)
	assert.NotNil(t, sprt)

	// also look for one that we don't expect to find
	_, found = ss2.FindSpriteByOrigin(topdown.Pt(555, 60))

	assert.False(t, found)
}
