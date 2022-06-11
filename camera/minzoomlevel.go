package camera

import (
	"math"

	"github.com/jamestunnell/topdown"
)

const (
	MaxWidth  = 2048
	MaxHeight = 2048
)

func MinZoomLevel(unzoomedSize topdown.Size[int]) float64 {
	minZoom1 := float64(unzoomedSize.Width) / float64(MaxWidth)
	minZoom2 := float64(unzoomedSize.Height) / float64(MaxHeight)

	return math.Max(minZoom1, minZoom2)
}
