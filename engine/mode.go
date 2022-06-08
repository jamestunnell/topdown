package engine

import (
	"github.com/hajimehoshi/ebiten/v2"

	"github.com/jamestunnell/topdown"
	"github.com/jamestunnell/topdown/resource"
)

//go:generate mockgen -destination=mock_engine/mockmode.go . Mode

type Mode interface {
	Initialize(screenSize topdown.Size, mgr resource.Manager) error

	Update() (Mode, error)
	Draw(screen *ebiten.Image)
	Layout(w, h int) (int, int)
}
