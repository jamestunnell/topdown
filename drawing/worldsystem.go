package drawing

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/igrmk/treemap/v2"
	"github.com/jamestunnell/topdown"
	"github.com/rs/zerolog/log"
)

//go:generate mockgen -destination=mock_drawing/mockworldsystem.go . WorldSystem

type WorldSystem interface {
	System

	DrawWorld(visible topdown.Rectangle[float64])
}

type worldSystem struct {
	surface *ebiten.Image
	layers  *treemap.TreeMap[int, *WorldLayer]
}

func NewWorldSystem(w, h int) WorldSystem {
	return &worldSystem{
		surface: ebiten.NewImage(w, h),
		layers:  treemap.New[int, *WorldLayer](),
	}
}

func (s *worldSystem) Add(id string, x interface{}) {
	d, ok := x.(WorldDrawable)
	if !ok {
		return
	}

	log.Debug().Str("id", id).Msg("adding world drawable")

	order := d.WorldLayer()

	l, found := s.layers.Get(order)
	if !found {
		l = NewWorldLayer(order)

		s.layers.Set(order, l)
	}

	l.Add(id, d)
}

func (s *worldSystem) Remove(id string) bool {
	for it := s.layers.Iterator(); it.Valid(); it.Next() {
		if it.Value().Remove(id) {
			return true
		}
	}

	return false
}

func (s *worldSystem) Clear() {
	for it := s.layers.Iterator(); it.Valid(); it.Next() {
		it.Value().Clear()
	}
}

func (s *worldSystem) DrawWorld(visible topdown.Rectangle[float64]) {
	for it := s.layers.Iterator(); it.Valid(); it.Next() {
		it.Value().Draw(s.surface, visible)
	}
}

func (s *worldSystem) Surface() *ebiten.Image {
	return s.surface
}
