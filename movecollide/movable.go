package movecollide

import "github.com/jamestunnell/topdown"

//go:generate mockgen -destination=mock_movecollide/mockmovable.go . Movable

// Movable is a component used in the move-collide system.
type Movable interface {
	PlanMovement(deltaSec float64) topdown.Vector
	Move(topdown.Vector)
}
