package drawing

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

//go:generate mockgen -destination=mock_drawing/mockdrawable.go . WorldDrawable,OverlayDrawable

// WorldDrawable is the component used in the drawing system to draw
// on the world surface.
type WorldDrawable interface {
	WorldLayer() int
	WorldDraw(world *ebiten.Image, visible image.Rectangle)
}

// OverlayDrawable is the component used in the drawing system to draw
// on the screen overlay.
type OverlayDrawable interface {
	OverlayLayer() int
	OverlayDraw(screen *ebiten.Image)
}
