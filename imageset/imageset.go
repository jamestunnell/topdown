package imageset

import (
	"fmt"
	"image"

	mapset "github.com/deckarep/golang-set/v2"
	"github.com/hajimehoshi/ebiten/v2"

	"github.com/jamestunnell/topdown"
	"github.com/jamestunnell/topdown/imageresource"
	"github.com/jamestunnell/topdown/resource"
)

type ImageSet struct {
	ImageRef  string      `json:"imageRef"`
	SubImages []*SubImage `json:"subImages"`
}

type SubImage struct {
	Properties map[string]any `json:"props,omitempty"`
	X          int            `json:"x"`
	Y          int            `json:"y"`
	Width      int            `json:"width"`
	Height     int            `json:"height"`
	Tags       []string       `json:"tags,omitempty"`
	Image      *ebiten.Image
}

const TypeName = "imageset"

func New(imageRef string, subImages ...*SubImage) *ImageSet {
	return &ImageSet{
		ImageRef:  imageRef,
		SubImages: subImages,
	}
}

func (is *ImageSet) Initialize(mgr resource.Manager) error {
	ir, err := resource.GetAs[*imageresource.ImageResource](mgr, is.ImageRef)
	if err != nil {
		return fmt.Errorf("failed to get '%s' from dependencies: %w", is.ImageRef, err)
	}

	for _, si := range is.SubImages {
		bounds := image.Rect(si.X, si.Y, si.X+si.Width, si.Y+si.Height)

		si.Image = ir.Image.SubImage(bounds).(*ebiten.Image)
	}

	return nil
}

func (is *ImageSet) Tags() []string {
	tags := mapset.NewSet[string]()

	for _, subImage := range is.SubImages {
		for _, tag := range subImage.Tags {
			tags.Add(tag)
		}
	}

	return tags.ToSlice()
}

func (is *ImageSet) SubImage(startPixel topdown.Pixel) (*ebiten.Image, *SubImage, bool) {
	for _, subImage := range is.SubImages {
		if subImage.X == startPixel.X && subImage.Y == startPixel.Y {
			if subImage.Image == nil {
				return nil, nil, false
			}

			return subImage.Image, subImage, true
		}
	}

	return nil, nil, false
}
