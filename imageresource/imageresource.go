package imageresource

import (
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/jamestunnell/topdown/resource"
)

type ImageResource struct {
	Image *ebiten.Image
}

type Type struct {
	name string
}

const (
	ImageTypePNG  = "png"
	ImageTypeJPEG = "jpeg"
	ImageTypeJPG  = "jpg"
)

func Load(path string) (*ImageResource, error) {
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

	a := &ImageResource{
		Image: ebiten.NewImageFromImage(img),
	}

	return a, nil
}

func (t *Type) Name() string {
	return t.name
}

func (t *Type) Load(path string) (resource.Resource, error) {
	return Load(path)
}

func Types() []resource.Type {
	return []resource.Type{
		&Type{name: ImageTypePNG},
		&Type{name: ImageTypeJPG},
		&Type{name: ImageTypeJPEG},
	}
}

func (img *ImageResource) Initialize(mgr resource.Manager) error {
	return nil
}
