package imageset_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jamestunnell/topdown"
	"github.com/jamestunnell/topdown/imageresource"
	"github.com/jamestunnell/topdown/imageresource/imagetest"
	"github.com/jamestunnell/topdown/imageset"
	"github.com/jamestunnell/topdown/jsonfile"
	"github.com/jamestunnell/topdown/resource"
	"github.com/jamestunnell/topdown/resource/restest"
)

func TestImageSet(t *testing.T) {
	dir, err := ioutil.TempDir("", "testimageset")

	require.NoError(t, err)

	defer os.RemoveAll(dir)

	subImages := []*imageset.SubImage{
		{X: 0, Y: 0, Width: 16, Height: 16},
		{X: 0, Y: 16, Width: 16, Height: 16},
	}
	is := imageset.New("missing.png", subImages...)
	partialPath := "bad.imageset"

	require.NoError(t, jsonfile.Write(filepath.Join(dir, partialPath), is))

	isType, err := imageset.NewType()

	require.NoError(t, err)

	types := append(imageresource.Types(), isType)
	mgr := restest.SetupManager(t, dir, types...)

	// fails due to missing image
	r, err := resource.GetAs[*imageset.ImageSet](mgr, partialPath)

	assert.Error(t, err)
	assert.Nil(t, r)

	imagePath := imagetest.WriteTestPNG(t, dir, 128, 128)

	is.ImageRef = filepath.Base(imagePath)
	partialPath = "good.imageset"

	require.NoError(t, jsonfile.Write(filepath.Join(dir, partialPath), is))

	// should not fail this time with a good source image
	r, err = resource.GetAs[*imageset.ImageSet](mgr, partialPath)

	require.NoError(t, err)

	is2, err := resource.As[*imageset.ImageSet](r)

	require.NoError(t, err)

	subImg, subImgInfo, found := is2.SubImage(topdown.NewPixel(0, 0))

	assert.True(t, found)
	assert.NotNil(t, subImg)
	assert.NotNil(t, subImgInfo)

	// also look for one that we don't expect to find
	_, _, found = is2.SubImage(topdown.NewPixel(555, 60))

	assert.False(t, found)
}
