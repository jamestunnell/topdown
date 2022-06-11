package animation

import (
	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/exp/slices"

	"github.com/jamestunnell/topdown/sprite"
)

func FrameImages(tag string, spriteSet *sprite.SpriteSet) (Images, error) {
	frameImages := Images{}

	for _, sprite := range spriteSet.Sprites {
		if slices.Contains(sprite.Tags, tag) {
			frameImages = append(frameImages, sprite.Image)
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
