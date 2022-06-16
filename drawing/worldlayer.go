package drawing

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jamestunnell/topdown"
	"github.com/jamestunnell/topdown/sliceutil"
	"golang.org/x/exp/slices"
)

// WorldLayer is used to organize the world into layers.
type WorldLayer struct {
	order     int
	ids       []string
	drawables []WorldDrawable
}

// NewWorldLayer makes a new world layer.
func NewWorldLayer(order int) *WorldLayer {
	return &WorldLayer{
		order:     order,
		ids:       []string{},
		drawables: []WorldDrawable{},
	}
}

// Order gets the layer order.
func (l *WorldLayer) Order() int {
	return l.order
}

// Add adds a drawable with ID.
func (l *WorldLayer) Add(id string, d WorldDrawable) {
	l.ids = append(l.ids, id)
	l.drawables = append(l.drawables, d)
}

// Remove removes the drawable with the given ID.
// Returns true if removed.
func (l *WorldLayer) Remove(id string) bool {
	idx := slices.Index(l.ids, id)
	if idx == -1 {
		return false
	}

	slices.Delete(l.ids, idx, idx+1)
	slices.Delete(l.drawables, idx, idx+1)

	return true
}

// Clear removes all drawables.
func (l *WorldLayer) Clear() {
	l.ids = nil
	l.drawables = nil
}

// Draw draws the layer drawables in sorted order.
func (l *WorldLayer) Draw(surface *ebiten.Image, visible topdown.Rectangle[float64]) {
	n := len(l.drawables)
	order := sliceutil.Make(n, func(i int) int { return i })

	slices.SortFunc(order, func(a, b int) bool {
		return l.drawables[a].WorldSortValue() < l.drawables[b].WorldSortValue()
	})

	for _, idx := range order {
		l.drawables[idx].WorldDraw(surface, visible)
	}
}
