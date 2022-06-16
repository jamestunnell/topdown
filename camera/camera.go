package camera

import (
	"fmt"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/jamestunnell/topdown"
	"github.com/jamestunnell/topdown/debug"
)

type Camera interface {
	// ScreenPosition returns the camera position in screen space (post-zoom).
	// Should always be screen center
	ScreenPosition() topdown.Point[float64]
	// ScreenArea returns the camera area in screen space (post-zoom).
	ScreenArea() topdown.Rectangle[float64]
	// WorldPosition gets the camera position in world space (pre-zoom).
	WorldPosition() topdown.Point[float64]
	// WorldArea returns the camera area in world space (pre-zoom).
	WorldArea() topdown.Rectangle[float64]
	// ZoomLevel returns the zoom level (e.g. 0.5 -> 50% zoom).
	ZoomLevel() float64
	// MinZoomLevel returns the minimum zoom level allowed given the current camera size.
	// The value is chosen to prevent the camera surface from exceeding MaxWidth or
	// MaxHeight.
	MinZoomLevel() float64
	// ConvertWorldToScreen converts the given world position to screen position.
	// Returns false if the world position is outside of the camera world area.
	ConvertWorldToScreen(worldPos topdown.Point[float64]) (topdown.Point[float64], bool)
	// ConvertScreenToWorld converts the given screen position to world position.
	// Returns false if the screen position is outside of the camera display area.
	ConvertScreenToWorld(screenPos topdown.Point[float64]) (topdown.Point[float64], bool)

	// Resize resizes the camera.
	Resize(size topdown.Size[int]) error
	// Move moves the camera position in world space.
	Move(topdown.Point[float64])
	// Zoom sets the camera zoom level.
	Zoom(float64)

	// DrawSurface returns the camera drawing surface (pre-zoom).
	DrawSurface() *ebiten.Image
	// Translation produces a translation operation given a position
	// in world space. This translation can be used to draw an image
	// on the camera draw surface.
	Translation(topdown.Point[float64]) *ebiten.DrawImageOptions
	// Blit draws the camera surface to the screen and applies zoom.
	Blit(screen *ebiten.Image)
}

type camera struct {
	position   topdown.Point[float64]
	zoom       float64
	size       topdown.Size[int]
	worldArea  topdown.Rectangle[float64]
	screenArea topdown.Rectangle[float64]
	surface    *ebiten.Image
	debugData  *debug.Dataset
}

func New(size topdown.Size[int]) (Camera, error) {
	cam := &camera{
		position:  topdown.Pt[float64](0, 0),
		zoom:      1,
		debugData: debug.NewDataset(),
	}

	if err := cam.Resize(size); err != nil {
		return nil, err
	}

	return cam, nil
}

func (c *camera) ScreenPosition() topdown.Point[float64] {
	p := c.size.Center()

	return topdown.Pt(float64(p.X), float64(p.Y))
}

func (c *camera) ScreenArea() topdown.Rectangle[float64] {
	return c.screenArea
}

func (c *camera) WorldPosition() topdown.Point[float64] {
	return c.position
}

func WorldArea(zoom float64, screenSize topdown.Size[int], position topdown.Point[float64]) topdown.Rectangle[float64] {
	scale := 1 / zoom
	width := scale * float64(screenSize.Width)
	height := scale * float64(screenSize.Height)

	origin := topdown.Pt(position.X-width/2.0, position.Y-height/2.0)

	return topdown.Rect(
		origin.X, origin.Y, origin.X+width, origin.Y+height)
}

func (c *camera) WorldArea() topdown.Rectangle[float64] {
	return c.worldArea
}

func (c *camera) ZoomLevel() float64 {
	return c.zoom
}

func (c *camera) MinZoomLevel() float64 {
	return MinZoomLevel(c.size)
}

func (c *camera) ConvertWorldToScreen(worldPos topdown.Point[float64]) (topdown.Point[float64], bool) {
	if !(worldPos.In(c.worldArea) || worldPos.Equal(c.worldArea.Max)) {
		return topdown.Pt(0.0, 0.0), false
	}

	fracX := (worldPos.X - c.worldArea.Min.X) / c.worldArea.Dx()
	fracY := (worldPos.Y - c.worldArea.Min.Y) / c.worldArea.Dy()
	x := c.screenArea.Dx() * fracX
	y := c.screenArea.Dy() * fracY

	return topdown.Pt(x, y), true
}

func (c *camera) ConvertScreenToWorld(screenPos topdown.Point[float64]) (topdown.Point[float64], bool) {
	if !(screenPos.In(c.screenArea) || screenPos.Equal(c.screenArea.Max)) {
		return topdown.Pt[float64](0, 0), false
	}

	worldArea := c.worldArea
	fracX := screenPos.X / c.screenArea.Dx()
	fracY := screenPos.Y / c.screenArea.Dy()
	x := worldArea.Min.X + worldArea.Dx()*fracX
	y := worldArea.Min.Y + worldArea.Dy()*fracY

	return topdown.Pt(x, y), true
}

func (c *camera) Resize(size topdown.Size[int]) error {
	if size.Width > MaxWidth || size.Height > MaxHeight {
		return fmt.Errorf("camera size %v is too big", size)
	}

	c.size = size
	c.screenArea = topdown.Rect(0.0, 0.0, float64(size.Width), float64(size.Height))

	minZoom := c.MinZoomLevel()
	if c.zoom < minZoom {
		c.zoom = minZoom
	}

	c.worldArea = WorldArea(c.zoom, size, c.position)

	if c.surface != nil {
		c.surface.Dispose()
	}

	c.surface = ebiten.NewImage(size.Width, size.Height)

	return nil
}

func (c *camera) Move(position topdown.Point[float64]) {
	c.position = position
	c.worldArea = WorldArea(c.zoom, c.size, position)
}

func (c *camera) Zoom(newZoom float64) {
	c.zoom = newZoom

	_ = c.Resize(c.size)
}

func (c *camera) DrawSurface() *ebiten.Image { return c.surface }

func (c *camera) Translation(p topdown.Point[float64]) *ebiten.DrawImageOptions {
	w, h := c.surface.Size()
	p2 := p.Sub(c.position)
	op := &ebiten.DrawImageOptions{}

	op.GeoM.Translate(float64(w)/2, float64(h)/2)
	op.GeoM.Translate(float64(p2.X), float64(p2.Y))

	return op
}

func (c *camera) Blit(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	w, h := c.surface.Size()
	cx := float64(w) / 2.0
	cy := float64(h) / 2.0

	op.GeoM.Translate(-cx, -cy)
	op.GeoM.Scale(c.zoom, c.zoom)
	op.GeoM.Translate(cx*c.zoom, cy*c.zoom)

	screen.DrawImage(c.surface, op)
}

func (c *camera) DebugData() *debug.Dataset {
	c.debugData.Set("x", strconv.FormatFloat(c.position.X, 'f', 2, 64))
	c.debugData.Set("y", strconv.FormatFloat(c.position.Y, 'f', 2, 64))
	c.debugData.Set("zoom", strconv.FormatFloat(c.zoom, 'f', 2, 64))
	c.debugData.Set("w", strconv.FormatFloat(c.worldArea.Dx(), 'f', 2, 64))
	c.debugData.Set("h", strconv.FormatFloat(c.worldArea.Dy(), 'f', 2, 64))

	return c.debugData
}
