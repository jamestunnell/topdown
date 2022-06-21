package drawing

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jamestunnell/topdown/camera"
	"github.com/jamestunnell/topdown/sliceutil"
	"golang.org/x/exp/slices"
)

type Layer struct {
	order     int
	ids       []string
	drawables []Drawable
}

// NewLayer makes a new layer.
func NewLayer(order int) *Layer {
	return &Layer{
		order:     order,
		ids:       []string{},
		drawables: []Drawable{},
	}
}

// Order gets the layer order.
func (l *Layer) Order() int {
	return l.order
}

// Add adds a drawable with ID.
func (l *Layer) Add(id string, d Drawable) {
	l.ids = append(l.ids, id)
	l.drawables = append(l.drawables, d)
}

// Remove removes the drawable with the given ID.
// Returns true if removed.
func (l *Layer) Remove(id string) bool {
	idx := slices.Index(l.ids, id)
	if idx == -1 {
		return false
	}

	slices.Delete(l.ids, idx, idx+1)
	slices.Delete(l.drawables, idx, idx+1)

	return true
}

// Clear removes all drawables.
func (l *Layer) Clear() {
	l.ids = nil
	l.drawables = nil
}

// Draw draws the layer drawables in sorted order.
func (l *Layer) Draw(screen *ebiten.Image, cam camera.Camera) {
	n := len(l.drawables)
	order := sliceutil.Make(n, func(i int) int { return i })

	slices.SortFunc(order, func(a, b int) bool {
		return l.drawables[a].DrawSortValue() < l.drawables[b].DrawSortValue()
	})

	for _, idx := range order {
		l.drawables[idx].Draw(screen, cam)
	}
}
