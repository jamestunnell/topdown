package drawing

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jamestunnell/topdown"
	"github.com/rs/zerolog/log"
	"golang.org/x/exp/slices"
)

//go:generate mockgen -destination=mock_drawing/mockworldsystem.go . WorldSystem

// WorldSystem is used to draw a world-sized image.
type WorldSystem interface {
	Add(id string, r interface{})
	Remove(id string) bool
	Clear()

	DrawWorld(worldSurface *ebiten.Image, visible topdown.Rectangle[float64])
}

type worldSystem struct {
	layers []*WorldLayer
}

// NewWorldSystem makes a new world drawing system.
func NewWorldSystem(w, h int) WorldSystem {
	return &worldSystem{
		layers: []*WorldLayer{},
	}
}

// Add will add the given object as a world drawable if it conforms
// to the WorldDrawable interface.
func (s *worldSystem) Add(id string, x interface{}) {
	d, ok := x.(WorldDrawable)
	if !ok {
		return
	}

	order := d.WorldLayer()

	idx := slices.IndexFunc(s.layers, func(l *WorldLayer) bool {
		return l.Order() == order
	})

	if idx == -1 {
		l := NewWorldLayer(order)

		s.layers = append(s.layers, l)

		slices.SortFunc(s.layers, func(a, b *WorldLayer) bool {
			return a.Order() < b.Order()
		})

		idx = slices.IndexFunc(s.layers, func(l *WorldLayer) bool {
			return l.Order() == order
		})
	}

	s.layers[idx].Add(id, d)

	log.Debug().Str("id", id).Msg("added world drawable")
}

// Remove will remove a drawable with the given ID if it is found.
// Returns true if found.
func (s *worldSystem) Remove(id string) bool {
	for _, l := range s.layers {
		if l.Remove(id) {
			return true
		}
	}

	return false
}

// Clear will remove all drawables.
func (s *worldSystem) Clear() {
	for _, l := range s.layers {
		l.Clear()
	}
}

// DrawOverlay will draw all drawables, first by layer and then by sort order.
func (s *worldSystem) DrawWorld(worldSurface *ebiten.Image, visible topdown.Rectangle[float64]) {
	for _, l := range s.layers {
		l.Draw(worldSurface, visible)
	}

	// opts := &ebiten.DrawImageOptions{}

	// translateX := float64(0)
	// translateY := float64(0)

	// if visible.Min.X < 0 {
	// 	translateX = float64(-visible.Min.X)
	// }

	// if visible.Min.Y < 0 {
	// 	translateY = float64(-visible.Min.Y)
	// }

	// if translateX != 0 || translateY != 0 {
	// 	opts.GeoM.Translate(translateX, translateY)
	// }

	// rect := image.Rect(
	// 	int(math.Round(visible.Min.X)), int(math.Round(visible.Min.Y)),
	// 	int(math.Round(visible.Max.X)), int(math.Round(visible.Max.Y)))

	// subImg := s.surface.SubImage(rect).(*ebiten.Image)

	// // draws visible world to camera drawing surface
	// screenSurface.DrawImage(subImg, opts)

}

// // Resize resizes the drawing surface.
// func (s *worldSystem) Resize(w, h int) {
// 	s.surface.Dispose()

// 	s.surface = ebiten.NewImage(w, h)
// }

// // Surface returns the world drawing surface.
// func (s *worldSystem) Surface() *ebiten.Image {
// 	return s.surface
// }
