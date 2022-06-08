package animation_test

import (
	"image"
	"image/color"
	"image/png"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jamestunnell/topdown/animation"
	"github.com/jamestunnell/topdown/imageresource"
	"github.com/jamestunnell/topdown/imageset"
	"github.com/jamestunnell/topdown/jsonfile"
	"github.com/jamestunnell/topdown/resource/restest"
)

func TestAnimations(t *testing.T) {
	dir, err := ioutil.TempDir("", "testanimation")

	require.NoError(t, err)

	defer os.RemoveAll(dir)

	w := 128
	h := 128

	// Break up the image into quadrants
	subImages := []*imageset.SubImage{
		{X: 0, Y: 0, Width: w / 2, Height: h / 2, Tags: []string{"idle"}},
		{X: w / 2, Y: 0, Width: w / 2, Height: h / 2, Tags: []string{"idle"}},
		{X: 0, Y: h / 2, Width: w / 2, Height: h / 2, Tags: []string{"walk"}},
		{X: w / 2, Y: h / 2, Width: w / 2, Height: h / 2, Tags: []string{"walk"}},
	}

	// write a test image and test image set
	path := writeTestImageSet(t, dir, w, h, subImages...)

	imageSetRef := filepath.Base(path)
	anims := animation.NewAnimations(imageSetRef, 20*time.Millisecond)

	isType, err := imageset.NewType()

	require.NoError(t, err)

	types := append(imageresource.Types(), isType)
	mgr := restest.SetupManager(t, dir, types...)

	f, err := ioutil.TempFile(dir, "testanimations*.animations")

	require.NoError(t, err)

	require.NoError(t, jsonfile.Write(f.Name(), anims))

	anims2, err := jsonfile.Read[*animation.Animations](f.Name())

	require.NoError(t, err)

	require.NoError(t, anims2.Initialize(mgr))

	assert.Len(t, anims2.TaggedFrames, 2)
	assert.Contains(t, anims2.TaggedFrames, "idle")
	assert.Contains(t, anims2.TaggedFrames, "walk")
	assert.Len(t, anims2.TaggedFrames["idle"], 2)
	assert.Len(t, anims2.TaggedFrames["walk"], 2)

	// f, err := ioutil.TempFile(dir, "testanimations*.animations")

	// require.NoError(t, err)

	// require.NoError(t, jsonfile.Write(f.Name(), anims))

	// a, err := mgr.Load(filepath.Base(f.Name()))

	// assert.NoError(t, err)
	// assert.NotNil(t, a)
}

func writeTestImageSet(t *testing.T, dir string, w, h int, subImages ...*imageset.SubImage) string {
	imgPath := writeTestPNG(t, dir, w, h)
	imgRef := filepath.Base(imgPath)

	is := imageset.New(imgRef, subImages...)

	f, err := os.CreateTemp(dir, "testImgSet*.imageset")

	require.NoError(t, err)

	require.NoError(t, jsonfile.Write(f.Name(), is))

	return f.Name()
}

func writeTestPNG(t *testing.T, dir string, w, h int) string {
	img := image.NewNRGBA(image.Rect(0, 0, w, h))

	img2 := ebiten.NewImageFromImage(img)
	img2.Clear()
	img2.Fill(color.Black)

	f, err := os.CreateTemp(dir, "testImg*.png")

	require.NoError(t, err)

	require.NoError(t, png.Encode(f, img))

	return f.Name()
}
