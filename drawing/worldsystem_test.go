package drawing_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/hajimehoshi/ebiten/v2"

	"github.com/jamestunnell/topdown"
	"github.com/jamestunnell/topdown/drawing"
	"github.com/jamestunnell/topdown/drawing/mock_drawing"
)

func TestDrawingSystemNoDrawables(t *testing.T) {
	sys := drawing.NewWorldSystem(200, 200)
	surf := ebiten.NewImage(100, 100)

	sys.DrawWorld(surf, topdown.Rect[float64](0, 0, 100, 100))
}

func TestDrawingSystem(t *testing.T) {
	ctrl := gomock.NewController(t)

	s := drawing.NewWorldSystem(300, 3000)
	d1 := mock_drawing.NewMockWorldDrawable(ctrl)
	d2 := mock_drawing.NewMockWorldDrawable(ctrl)

	d1.EXPECT().WorldLayer().Return(0)
	d1.EXPECT().WorldSortValue().AnyTimes().Return(5.0)

	d2.EXPECT().WorldLayer().Return(0)
	d2.EXPECT().WorldSortValue().AnyTimes().Return(10.0)

	s.Add("a", d1)
	s.Add("b", d2)

	img := ebiten.NewImage(100, 00)
	visible := topdown.Rect[float64](0, 0, 100, 100)

	gomock.InOrder(
		d1.EXPECT().WorldDraw(img, visible),
		d2.EXPECT().WorldDraw(img, visible),
	)

	// one layer
	s.DrawWorld(img, visible)

	d3 := mock_drawing.NewMockWorldDrawable(ctrl)

	d3.EXPECT().WorldLayer().Return(1)
	d3.EXPECT().WorldSortValue().AnyTimes().Return(8.0)

	s.Add("c", d3)

	gomock.InOrder(
		d1.EXPECT().WorldDraw(img, visible),
		d2.EXPECT().WorldDraw(img, visible),
		d3.EXPECT().WorldDraw(img, visible),
	)

	// two layers
	s.DrawWorld(img, visible)

	s.Clear()

	// nothing to draw
	s.DrawWorld(img, visible)
}
