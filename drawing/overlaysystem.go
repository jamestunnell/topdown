package drawing

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/igrmk/treemap/v2"
)

//go:generate mockgen -destination=mock_drawing/mockoverlaysystem.go . OverlaySystem

// OverlaySystem is used to draw a screen-sized HUD/UI overlay image.
type OverlaySystem interface {
	System

	DrawOverlay()
	Resize(w, h int)
}

type overlaySystem struct {
	surface *ebiten.Image
	layers  *treemap.TreeMap[int, *OverlayLayer]
}

// NewOverlaySystem makes a new overlay drawing system.
func NewOverlaySystem(w, h int) OverlaySystem {
	return &overlaySystem{
		surface: ebiten.NewImage(w, h),
		layers:  treemap.New[int, *OverlayLayer](),
	}
}

// Add will add the given object as an overlay drawable if it conforms
// to the OverlayDrawable interface.
func (s *overlaySystem) Add(id string, resource interface{}) {
	d, ok := resource.(OverlayDrawable)
	if !ok {
		return
	}

	order := d.OverlayLayer()

	l, found := s.layers.Get(order)
	if !found {
		l = NewOverlayLayer(order)
	}

	l.Drawables[id] = d
}

// Remove will remove a drawable with the given ID if it is found.
// Returns true if found.
func (s *overlaySystem) Remove(id string) bool {
	for it := s.layers.Iterator(); it.Valid(); it.Next() {
		if _, found := it.Value().Drawables[id]; found {
			delete(it.Value().Drawables, id)

			return true
		}
	}

	return false
}

// Clear will remove all drawables.
func (s *overlaySystem) Clear() {
	for it := s.layers.Iterator(); it.Valid(); it.Next() {
		it.Value().Clear()
	}
}

// DrawOverlay will draw all drawables by layer order.
func (s *overlaySystem) DrawOverlay() {
	for it := s.layers.Iterator(); it.Valid(); it.Next() {
		for _, d := range it.Value().Drawables {
			d.OverlayDraw(s.surface)
		}
	}
}

// Resize resizes the drawing surface.
func (s *overlaySystem) Resize(w, h int) {
	s.surface = ebiten.NewImage(w, h)
}

// Surface returns the overlay drawing surface.
func (s *overlaySystem) Surface() *ebiten.Image {
	return s.surface
}
