package sprite

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/jamestunnell/topdown"
)

type Sprite struct {
	ID     string             `json:"id"`
	Origin topdown.Point[int] `json:"origin"`
	Size   topdown.Size[int]  `json:"size"`
	Tags   []string           `json:"tags,omitempty"`

	Image *ebiten.Image
}

func (c *Sprite) Initialize(parent *ebiten.Image) {
	r := image.Rect(c.Origin.X, c.Origin.Y, c.Origin.X+c.Size.Width, c.Origin.Y+c.Size.Height)

	c.Image = parent.SubImage(r).(*ebiten.Image)
}
