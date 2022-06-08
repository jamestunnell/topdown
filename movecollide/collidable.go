package movecollide

import "github.com/zergon321/cirno"

//go:generate mockgen -destination=mock_movecollide/mockcollidable.go . Collidable

// Collidable is a component used in the move-collide system.
type Collidable interface {
	ColliderShape() cirno.Shape
	ResolveCollision(cirno.Vector, cirno.Shapes) cirno.Vector
}
