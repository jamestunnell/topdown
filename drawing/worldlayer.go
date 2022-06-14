package drawing

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jamestunnell/topdown"
	"github.com/jamestunnell/topdown/sliceutil"
	"golang.org/x/exp/slices"
)

type WorldLayer struct {
	ids       []string
	drawables []WorldDrawable
}

func NewWorldLayer(order int) *WorldLayer {
	return &WorldLayer{
		ids:       []string{},
		drawables: []WorldDrawable{},
	}
}

func (l *WorldLayer) Add(id string, d WorldDrawable) {
	l.ids = append(l.ids, id)
	l.drawables = append(l.drawables, d)
}

func (l *WorldLayer) Remove(id string) bool {
	idx := slices.Index(l.ids, id)
	if idx == -1 {
		return false
	}

	slices.Delete(l.ids, idx, idx+1)
	slices.Delete(l.drawables, idx, idx+1)

	return true
}

func (l *WorldLayer) Clear() {
	l.ids = nil
	l.drawables = nil
}

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
