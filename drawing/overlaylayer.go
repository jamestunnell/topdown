package drawing

import "golang.org/x/exp/maps"

// OverlayLayer is used to organize the overlay into layers.
type OverlayLayer struct {
	Drawables map[string]OverlayDrawable
}

// NewOverlayLayer makes a new overlay layer.
func NewOverlayLayer(order int) *OverlayLayer {
	return &OverlayLayer{
		Drawables: map[string]OverlayDrawable{},
	}
}

// Clear removes all drawables.
func (l *OverlayLayer) Clear() {
	maps.Clear(l.Drawables)
}
