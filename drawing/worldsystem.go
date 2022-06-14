package drawing

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/igrmk/treemap/v2"
	"github.com/jamestunnell/topdown"
	"github.com/rs/zerolog/log"
)

//go:generate mockgen -destination=mock_drawing/mockworldsystem.go . WorldSystem

// WorldSystem is used to draw a world-sized image.
type WorldSystem interface {
	System

	DrawWorld(visible topdown.Rectangle[float64])
}

type worldSystem struct {
	surface *ebiten.Image
	layers  *treemap.TreeMap[int, *WorldLayer]
}

// NewWorldSystem makes a new world drawing system.
func NewWorldSystem(w, h int) WorldSystem {
	return &worldSystem{
		surface: ebiten.NewImage(w, h),
		layers:  treemap.New[int, *WorldLayer](),
	}
}

// Add will add the given object as a world drawable if it conforms
// to the WorldDrawable interface.
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

// Remove will remove a drawable with the given ID if it is found.
// Returns true if found.
func (s *worldSystem) Remove(id string) bool {
	for it := s.layers.Iterator(); it.Valid(); it.Next() {
		if it.Value().Remove(id) {
			return true
		}
	}

	return false
}

// Clear will remove all drawables.
func (s *worldSystem) Clear() {
	for it := s.layers.Iterator(); it.Valid(); it.Next() {
		it.Value().Clear()
	}
}

// DrawOverlay will draw all drawables, first by layer and then by sort order.
func (s *worldSystem) DrawWorld(visible topdown.Rectangle[float64]) {
	for it := s.layers.Iterator(); it.Valid(); it.Next() {
		it.Value().Draw(s.surface, visible)
	}
}

// Resize resizes the drawing surface.
func (s *worldSystem) Resize(w, h int) {
	s.surface = ebiten.NewImage(w, h)
}

// Surface returns the world drawing surface.
func (s *worldSystem) Surface() *ebiten.Image {
	return s.surface
}
