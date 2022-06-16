package drawing

import (
	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/exp/slices"
)

// OverlayLayer is used to organize the Overlay into layers.
type OverlayLayer struct {
	order     int
	ids       []string
	drawables []OverlayDrawable
}

// NewOverlayLayer makes a new Overlay layer.
func NewOverlayLayer(order int) *OverlayLayer {
	return &OverlayLayer{
		order:     order,
		ids:       []string{},
		drawables: []OverlayDrawable{},
	}
}

// Order gets the layer order.
func (l *OverlayLayer) Order() int {
	return l.order
}

// Add adds a drawable with ID.
func (l *OverlayLayer) Add(id string, d OverlayDrawable) {
	l.ids = append(l.ids, id)
	l.drawables = append(l.drawables, d)
}

// Remove removes the drawable with the given ID.
// Returns true if removed.
func (l *OverlayLayer) Remove(id string) bool {
	idx := slices.Index(l.ids, id)
	if idx == -1 {
		return false
	}

	slices.Delete(l.ids, idx, idx+1)
	slices.Delete(l.drawables, idx, idx+1)

	return true
}

// Clear removes all drawables.
func (l *OverlayLayer) Clear() {
	l.ids = nil
	l.drawables = nil
}

// Draw draws the layer drawables in sorted order.
func (l *OverlayLayer) Draw(screen *ebiten.Image) {
	for _, d := range l.drawables {
		d.OverlayDraw(screen)
	}
}
