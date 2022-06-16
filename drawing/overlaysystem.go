package drawing

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/rs/zerolog/log"
	"golang.org/x/exp/slices"
)

//go:generate mockgen -destination=mock_drawing/mockoverlaysystem.go . OverlaySystem

// OverlaySystem is used to draw a screen-sized HUD/UI overlay image.
type OverlaySystem interface {
	Add(id string, r interface{})
	Remove(id string)
	Clear()

	DrawOverlay(screen *ebiten.Image)
}

type overlaySystem struct {
	layers            []*OverlayLayer
	debugPrintables   []DebugPrintable
	debugPrintableIDs []string
}

// NewOverlaySystem makes a new overlay drawing system.
func NewOverlaySystem() OverlaySystem {
	return &overlaySystem{
		layers:            []*OverlayLayer{},
		debugPrintables:   []DebugPrintable{},
		debugPrintableIDs: []string{},
	}
}

// Add will add the given object as an overlay drawable if it conforms
// to the OverlayDrawable interface.
func (s *overlaySystem) Add(id string, x interface{}) {
	d, ok := x.(OverlayDrawable)
	if ok {
		order := d.OverlayLayer()

		idx := slices.IndexFunc(s.layers, func(l *OverlayLayer) bool {
			return l.Order() == order
		})

		if idx == -1 {
			l := NewOverlayLayer(order)

			s.layers = append(s.layers, l)

			slices.SortFunc(s.layers, func(a, b *OverlayLayer) bool {
				return a.Order() < b.Order()
			})

			idx = slices.IndexFunc(s.layers, func(l *OverlayLayer) bool {
				return l.Order() == order
			})
		}

		s.layers[idx].Add(id, d)

		log.Debug().Str("id", id).Msg("added overlay drawable")
	}

	dp, ok := x.(DebugPrintable)
	if ok {
		s.debugPrintables = append(s.debugPrintables, dp)
		s.debugPrintableIDs = append(s.debugPrintableIDs, id)

		log.Debug().Str("id", id).Msg("added debug printables")
	}
}

// Remove will remove a drawable or debug printable with the given
// ID if it is found.
func (s *overlaySystem) Remove(id string) {
	for _, l := range s.layers {
		if l.Remove(id) {
			break
		}
	}

	if idx := slices.Index(s.debugPrintableIDs, id); idx != -1 {
		slices.Delete(s.debugPrintableIDs, idx, idx+1)
		slices.Delete(s.debugPrintables, idx, idx+1)
	}
}

// Clear will remove all drawables.
func (s *overlaySystem) Clear() {
	for _, l := range s.layers {
		l.Clear()
	}

	s.debugPrintableIDs = []string{}
	s.debugPrintables = []DebugPrintable{}
}

// DrawOverlay will draw all drawables by layer order.
func (s *overlaySystem) DrawOverlay(screen *ebiten.Image) {
	for _, l := range s.layers {
		l.Draw(screen)
	}

	DebugPrint(screen, s.debugPrintableIDs, s.debugPrintables)
}
