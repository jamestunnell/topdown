package sprite

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/jamestunnell/topdown"
)

type Sprite struct {
	Start topdown.Point[int] `json:"start"`
	Size  topdown.Size[int]  `json:"size"`
	Tags  []string           `json:"tags,omitempty"`

	Image *ebiten.Image
}

func (c *Sprite) Initialize(parent *ebiten.Image) {
	r := image.Rect(c.Start.X, c.Start.Y, c.Start.X+c.Size.Width, c.Start.Y+c.Size.Height)

	c.Image = parent.SubImage(r).(*ebiten.Image)
}
