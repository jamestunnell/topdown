package drawing

import "golang.org/x/exp/maps"

type OverlayLayer struct {
	Drawables map[string]OverlayDrawable
}

func NewOverlayLayer(order int) *OverlayLayer {
	return &OverlayLayer{
		Drawables: map[string]OverlayDrawable{},
	}
}

func (l *OverlayLayer) Clear() {
	maps.Clear(l.Drawables)
}
