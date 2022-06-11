package drawing_test

import (
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/jamestunnell/topdown"
	"github.com/jamestunnell/topdown/drawing"
	"github.com/jamestunnell/topdown/drawing/mock_drawing"
)

func TestDrawingSystemNoDrawables(t *testing.T) {
	s := drawing.NewWorldSystem(100, 100)

	s.DrawWorld(topdown.Rect[float64](0, 0, 100, 100))
}

func TestDrawingSystem(t *testing.T) {
	ctrl := gomock.NewController(t)

	s := drawing.NewWorldSystem(100, 100)
	d1 := mock_drawing.NewMockWorldDrawable(ctrl)
	d2 := mock_drawing.NewMockWorldDrawable(ctrl)

	d1.EXPECT().WorldLayer().Return(0)
	d2.EXPECT().WorldLayer().Return(0)

	s.Add("a", d1)
	s.Add("b", d2)

	img := s.Surface()
	visible := topdown.Rect[float64](0, 0, 100, 100)

	d1.EXPECT().WorldDraw(img, visible)
	d2.EXPECT().WorldDraw(img, visible)

	// one layer
	s.DrawWorld(visible)

	d3 := mock_drawing.NewMockWorldDrawable(ctrl)

	d3.EXPECT().WorldLayer().Return(1)

	s.Add("c", d3)

	d1.EXPECT().WorldDraw(img, visible)
	d2.EXPECT().WorldDraw(img, visible)
	d3.EXPECT().WorldDraw(img, visible)

	// two layers
	s.DrawWorld(visible)

	s.Clear()

	// nothing to draw
	s.DrawWorld(visible)
}
