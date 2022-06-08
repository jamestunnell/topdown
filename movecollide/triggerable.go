package movecollide

import "github.com/zergon321/cirno"

//go:generate mockgen -destination=mock_movecollide/mocktriggerable.go . Triggerable

// Triggerable is a component used in the move-collide system.
type Triggerable interface {
	TriggerShape() cirno.Shape
	TriggerEnter(cirno.Shape)
	TriggerRemain(cirno.Shape)
	TriggerExit(cirno.Shape)
}
