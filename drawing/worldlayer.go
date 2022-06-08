package drawing

import "golang.org/x/exp/maps"

type WorldLayer struct {
	Drawables map[string]WorldDrawable
}

func NewWorldLayer(order int) *WorldLayer {
	return &WorldLayer{
		Drawables: map[string]WorldDrawable{},
	}
}

func (l *WorldLayer) Clear() {
	maps.Clear(l.Drawables)
}
