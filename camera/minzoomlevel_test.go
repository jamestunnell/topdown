package camera_test

import (
	"fmt"
	"testing"

	"github.com/jamestunnell/topdown/camera"
	"github.com/stretchr/testify/assert"
)

func TestMinZoomLevel(t *testing.T) {
	testMinZoomLevel(t, camera.MaxWidth, camera.MaxHeight, 1)
	testMinZoomLevel(t, 800, 600, float64(800)/float64(camera.MaxWidth))
	testMinZoomLevel(t, 400, 600, float64(600)/float64(camera.MaxHeight))
}

func testMinZoomLevel(t *testing.T, w, h int, expected float64) {
	t.Run(fmt.Sprintf("w=%d,h=%d", w, h), func(t *testing.T) {
		actual := camera.MinZoomLevel(w, h)

		assert.Equal(t, expected, actual)
	})
}
