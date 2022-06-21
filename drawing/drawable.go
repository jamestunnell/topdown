package drawing

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jamestunnell/topdown/camera"
)

//go:generate mockgen -destination=mock_drawing/mockdrawable.go . Drawable

type Drawable interface {
	DrawLayer() int
	DrawSortValue() float64
	Draw(screen *ebiten.Image, cam camera.Camera)
}
