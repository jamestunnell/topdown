package drawing

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jamestunnell/topdown/camera"
	"github.com/rs/zerolog/log"
	"golang.org/x/exp/slices"
)

//go:generate mockgen -destination=mock_drawing/mocksystem.go . System

// System is used to draw all of the game drawable.
type System interface {
	Add(id string, r interface{})
	Remove(id string)
	Clear()

	Draw(screen *ebiten.Image)
}

type system struct {
	cam               camera.Camera
	layers            []*Layer
	debugPrintables   []DebugPrintable
	debugPrintableIDs []string
}

// NewSystem makes a new overlay drawing system.
func NewSystem(cam camera.Camera) System {
	return &system{
		cam:               cam,
		layers:            []*Layer{},
		debugPrintables:   []DebugPrintable{},
		debugPrintableIDs: []string{},
	}
}

// Add will add the given object as an overlay drawable if it conforms
// to the Drawable interface.
func (s *system) Add(id string, x interface{}) {
	d, ok := x.(Drawable)
	if ok {
		order := d.DrawLayer()

		idx := slices.IndexFunc(s.layers, func(l *Layer) bool {
			return l.Order() == order
		})

		if idx == -1 {
			l := NewLayer(order)

			s.layers = append(s.layers, l)

			slices.SortFunc(s.layers, func(a, b *Layer) bool {
				return a.Order() < b.Order()
			})

			idx = slices.IndexFunc(s.layers, func(l *Layer) bool {
				return l.Order() == order
			})
		}

		s.layers[idx].Add(id, d)

		log.Debug().Str("id", id).Msg("added drawable")
	}

	dp, ok := x.(DebugPrintable)
	if ok {
		s.debugPrintables = append(s.debugPrintables, dp)
		s.debugPrintableIDs = append(s.debugPrintableIDs, id)

		log.Debug().Str("id", id).Msg("added debug printable")
	}
}

// Remove will remove a drawable or debug printable with the given
// ID if it is found.
func (s *system) Remove(id string) {
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
func (s *system) Clear() {
	for _, l := range s.layers {
		l.Clear()
	}

	s.debugPrintableIDs = []string{}
	s.debugPrintables = []DebugPrintable{}
}

// Draw will draw all drawables by layer order.
func (s *system) Draw(screen *ebiten.Image) {
	for _, l := range s.layers {
		l.Draw(screen, s.cam)
	}

	DebugPrint(screen, s.debugPrintableIDs, s.debugPrintables)

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
