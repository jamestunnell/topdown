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

	"github.com/jamestunnell/topdown"
	"github.com/jamestunnell/topdown/animation"
	"github.com/jamestunnell/topdown/jsonfile"
	"github.com/jamestunnell/topdown/resource/restest"
	"github.com/jamestunnell/topdown/sprite"
)

func TestAnimations(t *testing.T) {
	dir, err := ioutil.TempDir("", "testanimation")

	require.NoError(t, err)

	defer os.RemoveAll(dir)

	w := 128
	h := 128

	// Break up the image into quadrants
	sprites := []*sprite.Sprite{
		{
			Start: topdown.Pt(0, 0),
			Size:  topdown.Sz(w/2, h/2),
			Tags:  []string{"idle"},
		},
		{
			Start: topdown.Pt(w/2, 0),
			Size:  topdown.Sz(w/2, h/2),
			Tags:  []string{"idle"},
		},
		{
			Start: topdown.Pt(0, h/2),
			Size:  topdown.Sz(w/2, h/2),
			Tags:  []string{"walk"},
		},
		{
			Start: topdown.Pt(w/2, h/2),
			Size:  topdown.Sz(w/2, h/2),
			Tags:  []string{"walk"},
		},
	}

	// write a test image and test image set
	path := writeTestImageSet(t, dir, w, h, sprites...)

	imageSetRef := filepath.Base(path)
	anims := animation.NewAnimations(imageSetRef, 20*time.Millisecond)

	mgr := restest.SetupManager(t, dir, sprite.Types()...)

	f, err := ioutil.TempFile(dir, "testanimations*.animations")

	require.NoError(t, err)

	require.NoError(t, jsonfile.Write(f.Name(), anims))

	anims2, err := jsonfile.Read[*animation.Animations](f.Name())

	require.NoError(t, err)

	require.NoError(t, anims2.Initialize(mgr))

	assert.Len(t, anims2.TaggedImages, 2)
	assert.Contains(t, anims2.TaggedImages, "idle")
	assert.Contains(t, anims2.TaggedImages, "walk")
	assert.Len(t, anims2.TaggedImages["idle"], 2)
	assert.Len(t, anims2.TaggedImages["walk"], 2)

	// f, err := ioutil.TempFile(dir, "testanimations*.animations")

	// require.NoError(t, err)

	// require.NoError(t, jsonfile.Write(f.Name(), anims))

	// a, err := mgr.Load(filepath.Base(f.Name()))

	// assert.NoError(t, err)
	// assert.NotNil(t, a)
}

func writeTestImageSet(t *testing.T, dir string, w, h int, sprites ...*sprite.Sprite) string {
	imgPath := writeTestPNG(t, dir, w, h)
	imgRef := filepath.Base(imgPath)

	ss := sprite.NewSpriteSet(imgRef, sprites...)

	f, err := os.CreateTemp(dir, "testSprites*.spriteset")

	require.NoError(t, err)

	require.NoError(t, jsonfile.Write(f.Name(), ss))

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
