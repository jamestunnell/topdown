package camera_test

import (
	"image"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jamestunnell/topdown/camera"
)

type testCase struct {
	Name          string
	Center        image.Point
	Width, Height int
	Zoom          float64
	ExpectedZoom  float64
}

func TestCameraDefaults(t *testing.T) {
	cam, err := camera.New(500, 300)

	require.NoError(t, err)
	assert.Equal(t, image.Pt(0, 0), cam.WorldPosition())
	assert.Equal(t, float64(1), cam.ZoomLevel())
}

func TestCameraScreenTooBig(t *testing.T) {
	cam, err := camera.New(camera.MaxWidth+1, 300)

	assert.Error(t, err)
	assert.Nil(t, cam)

	cam, err = camera.New(300, camera.MaxHeight+1)

	assert.Error(t, err)
	assert.Nil(t, cam)
}

func TestCameraBascis(t *testing.T) {
	testCases := []*testCase{
		{
			Name:         "default zoom at origin",
			Center:       image.Pt(0, 0),
			Width:        400,
			Height:       200,
			Zoom:         1,
			ExpectedZoom: 1,
		},
		{
			Name:         "2X zoom",
			Center:       image.Pt(-10, 150),
			Width:        640,
			Height:       480,
			Zoom:         2,
			ExpectedZoom: 2,
		},
		{
			Name:         "0.5X zoom",
			Center:       image.Pt(250, -170),
			Width:        320,
			Height:       240,
			Zoom:         0.5,
			ExpectedZoom: 0.5,
		},
		{
			Name:         "min zoom",
			Center:       image.Pt(-25, -50),
			Width:        800,
			Height:       600,
			Zoom:         camera.MinZoomLevel(800, 600),
			ExpectedZoom: camera.MinZoomLevel(800, 600),
		},
		{
			Name:         "below min zoom",
			Center:       image.Pt(44, 26),
			Width:        1200,
			Height:       900,
			Zoom:         camera.MinZoomLevel(1200, 900) / 2,
			ExpectedZoom: camera.MinZoomLevel(1200, 900),
		},
	}

	for _, tc := range testCases {
		testCameraBasics(t, tc)
	}
}

func TestCameraResizeTooBig(t *testing.T) {
	cam, err := camera.New(200, 200)

	require.NoError(t, err)

	assert.Error(t, cam.Resize(camera.MaxWidth+1, 200))
	assert.Error(t, cam.Resize(200, camera.MaxHeight+1))
}

func TestCameraResizeRequiresZoomAdjust(t *testing.T) {
	cam, err := camera.New(200, 200)

	require.NoError(t, err)

	minZoom := cam.MinZoomLevel()

	cam.Zoom(minZoom)

	assert.NoError(t, cam.Resize(400, 400))

	assert.Equal(t, minZoom*2, cam.ZoomLevel())
}

func TestCameraConvertCoordsOutOfArea(t *testing.T) {
	cam, err := camera.New(500, 300)

	require.NoError(t, err)

	_, ok := cam.ConvertWorldToScreen(image.Pt(-251, -150))

	assert.False(t, ok)

	_, ok = cam.ConvertWorldToScreen(image.Pt(-250, -151))

	assert.False(t, ok)

	_, ok = cam.ConvertScreenToWorld(image.Pt(-1, 0))

	assert.False(t, ok)

	_, ok = cam.ConvertScreenToWorld(image.Pt(0, -1))

	assert.False(t, ok)
}

func TestCameraConvertCoordsDefaultZoom(t *testing.T) {
	cam, err := camera.New(500, 300)

	require.NoError(t, err)

	testCameraConvertCoordsOK(
		t, "world pos (0,0) - screen center", cam, cam.WorldPosition(), cam.ScreenPosition())
	testCameraConvertCoordsOK(
		t, "world pos (0,0) - screen start", cam, image.Pt(-250, -150), image.Pt(0, 0))
	testCameraConvertCoordsOK(
		t, "world pos (0,0) - screen end", cam, image.Pt(250, 150), image.Pt(500, 300))

	cam.Move(image.Pt(100, 50))

	testCameraConvertCoordsOK(
		t, "world pos (100,50) - screen center", cam, image.Pt(100, 50), cam.ScreenPosition())
	testCameraConvertCoordsOK(
		t, "world pos (100,50) - screen start", cam, image.Pt(-150, -100), image.Pt(0, 0))
	testCameraConvertCoordsOK(
		t, "world pos (100,50) - screen end", cam, image.Pt(350, 200), image.Pt(500, 300))
}

func TestCameraConvertCoordsHalfZoom(t *testing.T) {
	cam, err := camera.New(500, 300)

	require.NoError(t, err)

	// this will make the world area 2X the screen area
	cam.Zoom(0.5)

	testCameraConvertCoordsOK(
		t, "world pos (0,0) - screen center", cam, cam.WorldPosition(), cam.ScreenPosition())
	testCameraConvertCoordsOK(
		t, "world pos (0,0) - screen start", cam, image.Pt(-500, -300), image.Pt(0, 0))
	testCameraConvertCoordsOK(
		t, "world pos (0,0) - screen end", cam, image.Pt(500, 300), image.Pt(500, 300))

	cam.Move(image.Pt(100, 50))

	testCameraConvertCoordsOK(
		t, "world pos (100,50) - screen center", cam, image.Pt(100, 50), cam.ScreenPosition())
	testCameraConvertCoordsOK(
		t, "world pos (100,50) - screen start", cam, image.Pt(-400, -250), image.Pt(0, 0))
	testCameraConvertCoordsOK(
		t, "world pos (100,50) - screen end", cam, image.Pt(600, 350), image.Pt(500, 300))
}

func testCameraConvertCoordsOK(
	t *testing.T,
	name string,
	cam camera.Camera,
	worldPos, screenPos image.Point) {
	t.Run(name, func(t *testing.T) {
		actualScreenPos, ok := cam.ConvertWorldToScreen(worldPos)

		require.True(t, ok)
		assert.Equal(t, screenPos, actualScreenPos)

		actualWorldPos, ok := cam.ConvertScreenToWorld(screenPos)

		require.True(t, ok)
		assert.Equal(t, worldPos, actualWorldPos)
	})
}

func testCameraBasics(t *testing.T, tc *testCase) {
	t.Run(tc.Name, func(t *testing.T) {
		cam, err := camera.New(tc.Width, tc.Height)
		require.NoError(t, err)

		assert.NotNil(t, cam.DrawSurface())

		cam.Move(tc.Center)
		cam.Zoom(tc.Zoom)

		assert.Equal(t, tc.ExpectedZoom, cam.ZoomLevel())
		assert.Equal(t, tc.Center, cam.WorldPosition())

		area := cam.ScreenArea()

		assert.Equal(t, tc.Width, area.Dx())
		assert.Equal(t, tc.Height, area.Dy())

		area = cam.WorldArea()

		zoomMul := 1 / tc.ExpectedZoom
		width := int(float64(tc.Width) * zoomMul)
		height := int(float64(tc.Height) * zoomMul)

		dx := area.Dx()
		dy := area.Dy()

		assert.Equal(t, width, dx)
		assert.Equal(t, height, dy)

		assert.Equal(t, tc.Center.X, (area.Min.X + dx/2))
		assert.Equal(t, tc.Center.Y, (area.Min.Y + dy/2))
	})
}
