package control

import (
	"github.com/hajimehoshi/ebiten/v2"

	"github.com/jamestunnell/topdown/input"
)

//go:generate mockgen -destination=mock_control/mockcontrollable.go . Controllable

// Controllable is the component used in the control system.
type Controllable interface {
	WatchKeys() []ebiten.Key
	Control(deltaSec float64, inputMgr input.Manager)
}
