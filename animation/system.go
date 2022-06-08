package animation

import (
	"time"

	"golang.org/x/exp/maps"
)

type System interface {
	Add(id string, x interface{})
	Remove(id string)
	Clear()

	Animate(deltaSec float64)
}

type system struct {
	animatables map[string]Animatable
}

func NewSystem() System {
	return &system{
		animatables: map[string]Animatable{},
	}
}

func (s *system) Add(id string, x interface{}) {
	if a, ok := x.(Animatable); ok {
		s.animatables[id] = a
	}
}

func (s *system) Remove(id string) {
	delete(s.animatables, id)
}

func (s *system) Clear() {
	maps.Clear(s.animatables)
}

func (s *system) Animate(deltaSec float64) {
	delta := time.Duration(deltaSec * 1e9)

	for _, a := range s.animatables {
		a.UpdateAnimation(delta)
	}
}
