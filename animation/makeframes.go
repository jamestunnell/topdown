package animation

import (
	"fmt"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/exp/slices"

	"github.com/jamestunnell/topdown"
	"github.com/jamestunnell/topdown/imageset"
)

const PropertyNameFrameDurationSec = "frameDurationSec"

func MakeFrames(frameTag string, frameDur time.Duration, imageSet *imageset.ImageSet) (Frames, error) {
	frames := Frames{}

	for _, subImage := range imageSet.SubImages {
		if slices.Contains(subImage.Tags, frameTag) {
			img, _, found := imageSet.SubImage(topdown.NewPixel(subImage.X, subImage.Y))
			if !found {
				err := fmt.Errorf("failed to find sub-image at (%d, %d)", subImage.X, subImage.Y)

				return Frames{}, err
			}

			frames = append(frames, img)
		}
	}

	SortFramesTopBottomLeftRight(frames)

	return frames, nil
}

func SortFramesTopBottomLeftRight(frames Frames) {
	slices.SortFunc(frames, func(a, b *ebiten.Image) bool {
		aMin := a.Bounds().Min
		bMin := b.Bounds().Min

		if aMin.Y < bMin.Y {
			return true
		}

		return aMin.X < bMin.X
	})
}
