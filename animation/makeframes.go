package animation

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/exp/slices"

	"github.com/jamestunnell/topdown"
	"github.com/jamestunnell/topdown/imageset"
)

func FrameImages(tag string, imageSet *imageset.ImageSet) (Images, error) {
	frameImages := Images{}

	for _, subImage := range imageSet.SubImages {
		if slices.Contains(subImage.Tags, tag) {
			img, _, found := imageSet.SubImage(topdown.NewPixel(subImage.X, subImage.Y))
			if !found {
				err := fmt.Errorf("failed to find sub-image at (%d, %d)", subImage.X, subImage.Y)

				return Images{}, err
			}

			frameImages = append(frameImages, img)
		}
	}

	SortFramesTopBottomLeftRight(frameImages)

	return frameImages, nil
}

func SortFramesTopBottomLeftRight(frameImages Images) {
	slices.SortFunc(frameImages, func(a, b *ebiten.Image) bool {
		aMin := a.Bounds().Min
		bMin := b.Bounds().Min

		if aMin.Y < bMin.Y {
			return true
		}

		return aMin.X < bMin.X
	})
}
