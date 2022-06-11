package tilegrid

import "github.com/hajimehoshi/ebiten/v2"

type Row struct {
	Images  []*ebiten.Image
	XScales []float64
	YScales []float64
}

func NewRow(nCols int) *Row {
	return &Row{
		Images:  make([]*ebiten.Image, nCols),
		XScales: make([]float64, nCols),
		YScales: make([]float64, nCols),
	}
}
