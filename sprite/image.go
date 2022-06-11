package sprite

import (
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/jamestunnell/topdown/resource"
)

type Image struct {
	Image *ebiten.Image
}

type ImageType struct {
	name string
}

const (
	ImageTypePNG  = "png"
	ImageTypeJPEG = "jpeg"
	ImageTypeJPG  = "jpg"
)

func LoadImage(path string) (*Image, error) {
	f, err := os.Open(path)
	if err != nil {
		err = fmt.Errorf("failed to open: %w", err)

		return nil, err
	}

	img, _, err := image.Decode(f)
	if err != nil {
		err = fmt.Errorf("failed to decode: %w", err)

		return nil, err
	}

	a := &Image{
		Image: ebiten.NewImageFromImage(img),
	}

	return a, nil
}

func (img *Image) Initialize(mgr resource.Manager) error {
	return nil
}

func (t *ImageType) Name() string {
	return t.name
}

func (t *ImageType) Load(path string) (resource.Resource, error) {
	return LoadImage(path)
}
