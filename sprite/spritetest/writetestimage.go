package spritetest

import (
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"os"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/stretchr/testify/require"
)

func WriteTestPNG(t *testing.T, dir string, w, h int) string {
	img := image.NewNRGBA(image.Rect(0, 0, w, h))

	img2 := ebiten.NewImageFromImage(img)
	img2.Clear()
	img2.Fill(color.Black)

	f, err := os.CreateTemp(dir, "testImg*.png")

	require.NoError(t, err)

	require.NoError(t, png.Encode(f, img))

	return f.Name()
}

func WriteTestJPEG(t *testing.T, dir string, w, h int) string {
	return writeTestJPEG(t, dir, w, h, "jpeg")
}

func WriteTestJPG(t *testing.T, dir string, w, h int) string {
	return writeTestJPEG(t, dir, w, h, "jpg")
}

func writeTestJPEG(t *testing.T, dir string, w, h int, fileExt string) string {
	img := image.NewRGBA(image.Rect(0, 0, w, h))

	img2 := ebiten.NewImageFromImage(img)
	img2.Clear()
	img2.Fill(color.Black)

	f, err := os.CreateTemp(dir, "testImg*."+fileExt)

	require.NoError(t, err)

	require.NoError(t, jpeg.Encode(f, img, &jpeg.Options{Quality: 50}))

	return f.Name()
}
