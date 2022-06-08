package camera

import "math"

const (
	MaxWidth  = 2048
	MaxHeight = 2048
)

func MinZoomLevel(unzoomedWidth, unzoomedHeight int) float64 {
	minZoom1 := float64(unzoomedWidth) / float64(MaxWidth)
	minZoom2 := float64(unzoomedHeight) / float64(MaxHeight)

	return math.Max(minZoom1, minZoom2)
}
