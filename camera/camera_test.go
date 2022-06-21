package camera_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jamestunnell/topdown"
	"github.com/jamestunnell/topdown/camera"
)

type testCase struct {
	Name          string
	Position      topdown.Point[float64]
	Width, Height int
	Zoom          float64
	ExpectedZoom  float64
}

func TestCameraDefaults(t *testing.T) {
	cam, err := camera.New(topdown.Sz(500, 300))

	require.NoError(t, err)
	assert.Equal(t, topdown.Pt[float64](0, 0), cam.WorldPosition())
	assert.Equal(t, float64(1), cam.ZoomLevel())
}

func TestCameraScreenTooBig(t *testing.T) {
	cam, err := camera.New(topdown.Sz(camera.MaxWidth+1, 300))

	assert.Error(t, err)
	assert.Nil(t, cam)

	cam, err = camera.New(topdown.Sz(300, camera.MaxHeight+1))

	assert.Error(t, err)
	assert.Nil(t, cam)
}

func TestCameraBascis(t *testing.T) {
	testCases := []*testCase{
		{
			Name:         "default zoom at origin",
			Position:     topdown.Pt[float64](0, 0),
			Width:        400,
			Height:       200,
			Zoom:         1,
			ExpectedZoom: 1,
		},
		{
			Name:         "2X zoom",
			Position:     topdown.Pt[float64](-10, 150),
			Width:        640,
			Height:       480,
			Zoom:         2,
			ExpectedZoom: 2,
		},
		{
			Name:         "0.5X zoom",
			Position:     topdown.Pt[float64](250, -170),
			Width:        320,
			Height:       240,
			Zoom:         0.5,
			ExpectedZoom: 0.5,
		},
		{
			Name:         "min zoom",
			Position:     topdown.Pt[float64](-25, -50),
			Width:        800,
			Height:       600,
			Zoom:         camera.MinZoomLevel(topdown.Sz(800, 600)),
			ExpectedZoom: camera.MinZoomLevel(topdown.Sz(800, 600)),
		},
		{
			Name:         "below min zoom",
			Position:     topdown.Pt[float64](44, 26),
			Width:        1200,
			Height:       900,
			Zoom:         camera.MinZoomLevel(topdown.Sz(1200, 900)) / 2,
			ExpectedZoom: camera.MinZoomLevel(topdown.Sz(1200, 900)),
		},
	}

	for _, tc := range testCases {
		testCameraBasics(t, tc)
	}
}

func TestCameraResizeTooBig(t *testing.T) {
	cam, err := camera.New(topdown.Sz(200, 200))

	require.NoError(t, err)

	assert.Error(t, cam.Resize(topdown.Sz(camera.MaxWidth+1, 200)))
	assert.Error(t, cam.Resize(topdown.Sz(200, camera.MaxHeight+1)))
}

func TestCameraResizeRequiresZoomAdjust(t *testing.T) {
	cam, err := camera.New(topdown.Sz(200, 200))

	require.NoError(t, err)

	minZoom := cam.MinZoomLevel()

	cam.Zoom(minZoom)

	assert.NoError(t, cam.Resize(topdown.Sz(400, 400)))

	assert.Equal(t, minZoom*2, cam.ZoomLevel())
}

func TestCameraConvertCoordsOutOfArea(t *testing.T) {
	cam, err := camera.New(topdown.Sz(500, 300))

	require.NoError(t, err)

	_, ok := cam.ConvertWorldToScreen(topdown.Pt[float64](-251, -150))

	assert.False(t, ok)

	_, ok = cam.ConvertWorldToScreen(topdown.Pt[float64](-250, -151))

	assert.False(t, ok)

	_, ok = cam.ConvertScreenToWorld(topdown.Pt[float64](-1, 0))

	assert.False(t, ok)

	_, ok = cam.ConvertScreenToWorld(topdown.Pt[float64](0, -1))

	assert.False(t, ok)
}

func TestCameraConvertCoordsDefaultZoom(t *testing.T) {
	cam, err := camera.New(topdown.Sz(500, 300))

	require.NoError(t, err)

	testCameraConvertCoordsOK(
		t, "world pos (0,0) - screen center", cam, cam.WorldPosition(), cam.ScreenPosition())
	testCameraConvertCoordsOK(
		t, "world pos (0,0) - screen start", cam, topdown.Pt[float64](-250, -150), topdown.Pt[float64](0, 0))
	testCameraConvertCoordsOK(
		t, "world pos (0,0) - screen end", cam, topdown.Pt[float64](250, 150), topdown.Pt[float64](500, 300))

	cam.Move(topdown.Pt[float64](100, 50))

	testCameraConvertCoordsOK(
		t, "world pos (100,50) - screen center", cam, topdown.Pt[float64](100, 50), cam.ScreenPosition())
	testCameraConvertCoordsOK(
		t, "world pos (100,50) - screen start", cam, topdown.Pt[float64](-150, -100), topdown.Pt[float64](0, 0))
	testCameraConvertCoordsOK(
		t, "world pos (100,50) - screen end", cam, topdown.Pt[float64](350, 200), topdown.Pt[float64](500, 300))
}

func TestCameraConvertCoordsHalfZoom(t *testing.T) {
	cam, err := camera.New(topdown.Sz(500, 300))

	require.NoError(t, err)

	// this will make the world area 2X the screen area
	cam.Zoom(0.5)

	testCameraConvertCoordsOK(
		t, "world pos (0,0) - screen center", cam, cam.WorldPosition(), cam.ScreenPosition())
	testCameraConvertCoordsOK(
		t, "world pos (0,0) - screen start", cam, topdown.Pt[float64](-500, -300), topdown.Pt[float64](0, 0))
	testCameraConvertCoordsOK(
		t, "world pos (0,0) - screen end", cam, topdown.Pt[float64](500, 300), topdown.Pt[float64](500, 300))

	cam.Move(topdown.Pt[float64](100, 50))

	testCameraConvertCoordsOK(
		t, "world pos (100,50) - screen center", cam, topdown.Pt[float64](100, 50), cam.ScreenPosition())
	testCameraConvertCoordsOK(
		t, "world pos (100,50) - screen start", cam, topdown.Pt[float64](-400, -250), topdown.Pt[float64](0, 0))
	testCameraConvertCoordsOK(
		t, "world pos (100,50) - screen end", cam, topdown.Pt[float64](600, 350), topdown.Pt[float64](500, 300))
}

func testCameraConvertCoordsOK(
	t *testing.T,
	name string,
	cam camera.Camera,
	worldPos topdown.Point[float64],
	screenPos topdown.Point[float64]) {
	t.Run(name, func(t *testing.T) {
		actualScreenPos, _ := cam.ConvertWorldToScreen(worldPos)

		assert.Equal(t, screenPos, actualScreenPos)

		actualWorldPos, _ := cam.ConvertScreenToWorld(screenPos)

		assert.Equal(t, worldPos, actualWorldPos)
	})
}

func testCameraBasics(t *testing.T, tc *testCase) {
	t.Run(tc.Name, func(t *testing.T) {
		cam, err := camera.New(topdown.Sz(tc.Width, tc.Height))
		require.NoError(t, err)

		cam.Move(tc.Position)
		cam.Zoom(tc.Zoom)

		assert.Equal(t, tc.ExpectedZoom, cam.ZoomLevel())
		assert.Equal(t, tc.Position, cam.WorldPosition())

		screenArea := cam.ScreenArea()

		assert.Equal(t, float64(tc.Width), screenArea.Dx())
		assert.Equal(t, float64(tc.Height), screenArea.Dy())

		area := cam.WorldArea()

		zoomMul := 1 / tc.ExpectedZoom
		width := float64(tc.Width) * zoomMul
		height := float64(tc.Height) * zoomMul

		dx := area.Dx()
		dy := area.Dy()

		assert.Equal(t, width, dx)
		assert.Equal(t, height, dy)

		assert.Equal(t, tc.Position.X, (area.Min.X + dx/2))
		assert.Equal(t, tc.Position.Y, (area.Min.Y + dy/2))
	})
}
