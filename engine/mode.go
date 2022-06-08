package engine

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jamestunnell/topdown/resource"
)

//go:generate mockgen -destination=mock_engine/mockmode.go . Mode

type Mode interface {
	Initialize(resource.Manager) error

	Draw(screen *ebiten.Image)
	Update() error
	Layout(w, h int) (int, int)
}
