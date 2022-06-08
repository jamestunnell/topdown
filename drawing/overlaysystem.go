package drawing

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/igrmk/treemap/v2"
)

//go:generate mockgen -destination=mock_drawing/mockoverlaysystem.go . OverlaySystem

type OverlaySystem interface {
	System

	DrawOverlay()
	Resize(w, h int)
}

type overlaySystem struct {
	surface *ebiten.Image
	layers  *treemap.TreeMap[int, *OverlayLayer]
}

func NewOverlaySystem(w, h int) OverlaySystem {
	return &overlaySystem{
		surface: ebiten.NewImage(w, h),
		layers:  treemap.New[int, *OverlayLayer](),
	}
}

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

func (s *overlaySystem) Remove(id string) {
	for it := s.layers.Iterator(); it.Valid(); it.Next() {
		delete(it.Value().Drawables, id)
	}
}

func (s *overlaySystem) Clear() {
	for it := s.layers.Iterator(); it.Valid(); it.Next() {
		it.Value().Clear()
	}
}

func (s *overlaySystem) DrawOverlay() {
	for it := s.layers.Iterator(); it.Valid(); it.Next() {
		for _, d := range it.Value().Drawables {
			d.OverlayDraw(s.surface)
		}
	}
}

func (s *overlaySystem) Resize(w, h int) {
	s.surface = ebiten.NewImage(w, h)
}

func (s *overlaySystem) Surface() *ebiten.Image {
	return s.surface
}
