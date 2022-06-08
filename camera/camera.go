package camera

import (
	"fmt"
	"image"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/igrmk/treemap/v2"
)

type Camera interface {
	// ScreenPosition returns the camera position in screen space (post-zoom).
	// Should always be screen center
	ScreenPosition() image.Point
	// ScreenArea returns the camera area in screen space (post-zoom).
	ScreenArea() image.Rectangle
	// WorldPosition gets the camera position in world space (pre-zoom).
	WorldPosition() image.Point
	// WorldArea returns the camera area in world space (pre-zoom).
	WorldArea() image.Rectangle
	// ZoomLevel returns the zoom level (e.g. 0.5 -> 50% zoom).
	ZoomLevel() float64
	// MinZoomLevel returns the minimum zoom level allowed given the current camera size.
	// The value is chosen to prevent the camera surface from exceeding MaxWidth or
	// MaxHeight.
	MinZoomLevel() float64
	// ConvertWorldToScreen converts the given world position to screen position.
	// Returns false if the world position is outside of the camera world area.
	ConvertWorldToScreen(worldPos image.Point) (image.Point, bool)
	// ConvertScreenToWorld converts the given screen position to world position.
	// Returns false if the screen position is outside of the camera display area.
	ConvertScreenToWorld(screenPos image.Point) (image.Point, bool)

	// Resize resizes the camera.
	Resize(width, height int) error
	// Move moves the camera position in world space.
	Move(image.Point)
	// Zoom sets the camera zoom level.
	Zoom(float64)

	// DrawSurface returns the camera drawing surface (pre-zoom).
	DrawSurface() *ebiten.Image
	// Translation produces a translation operation given a position
	// in world space. This translation can be used to draw an image
	// on the camera draw surface.
	Translation(image.Point) *ebiten.DrawImageOptions
	// Blit draws the camera surface to the screen and applies zoom.
	Blit(screen *ebiten.Image)
}

type camera struct {
	position              image.Point
	zoom                  float64
	width, height         int
	worldArea, screenArea image.Rectangle
	surface               *ebiten.Image
	debugInfo             *treemap.TreeMap[string, string]
}

func New(width, height int) (Camera, error) {
	cam := &camera{
		position:  image.Pt(0, 0),
		zoom:      1,
		debugInfo: treemap.New[string, string](),
	}

	if err := cam.Resize(width, height); err != nil {
		return nil, err
	}

	return cam, nil
}

func (c *camera) ScreenPosition() image.Point {
	return image.Pt(c.width/2, c.height/2)
}

func (c *camera) ScreenArea() image.Rectangle {
	return c.screenArea
}

func (c *camera) WorldPosition() image.Point {
	return c.position
}

func WorldArea(zoom float64, width, height int, position image.Point) image.Rectangle {
	scale := 1 / zoom
	width = int(scale * float64(width))
	height = int(scale * float64(height))

	origin := image.Pt(position.X-width/2, position.Y-height/2)
	end := image.Pt(origin.X+width, origin.Y+height)

	return image.Rect(origin.X, origin.Y, end.X, end.Y)
}

func (c *camera) WorldArea() image.Rectangle {
	return c.worldArea
}

func (c *camera) ZoomLevel() float64 {
	return c.zoom
}

func (c *camera) MinZoomLevel() float64 {
	return MinZoomLevel(c.width, c.height)
}

func (c *camera) ConvertWorldToScreen(worldPos image.Point) (image.Point, bool) {
	if !(worldPos.In(c.worldArea) || worldPos.Eq(c.worldArea.Max)) {
		return image.Point{}, false
	}

	fracX := float64(worldPos.X-c.worldArea.Min.X) / float64(c.worldArea.Dx())
	fracY := float64(worldPos.Y-c.worldArea.Min.Y) / float64(c.worldArea.Dy())
	x := c.screenArea.Min.X + int(float64(c.screenArea.Dx())*fracX)
	y := c.screenArea.Min.Y + int(float64(c.screenArea.Dy())*fracY)

	return image.Pt(x, y), true
}

func (c *camera) ConvertScreenToWorld(screenPos image.Point) (image.Point, bool) {
	screenArea := c.ScreenArea()
	if !(screenPos.In(screenArea) || screenPos.Eq(screenArea.Max)) {
		return image.Point{}, false
	}

	worldArea := c.worldArea
	fracX := float64(screenPos.X-screenArea.Min.X) / float64(screenArea.Dx())
	fracY := float64(screenPos.Y-screenArea.Min.Y) / float64(screenArea.Dy())
	x := worldArea.Min.X + int(float64(worldArea.Dx())*fracX)
	y := worldArea.Min.Y + int(float64(worldArea.Dy())*fracY)

	return image.Pt(x, y), true
}

func (c *camera) Resize(width, height int) error {
	if width > MaxWidth || height > MaxHeight {
		return fmt.Errorf("camera size %d x %d is too big", width, height)
	}

	c.width = width
	c.height = height
	c.screenArea = image.Rect(0, 0, c.width, c.height)

	minZoom := c.MinZoomLevel()
	if c.zoom < minZoom {
		c.zoom = minZoom
	}

	c.worldArea = WorldArea(c.zoom, width, height, c.position)

	if c.surface != nil {
		c.surface.Dispose()
	}

	zoomedWidth := float64(width) / c.zoom
	zoomedHeight := float64(height) / c.zoom

	c.surface = ebiten.NewImage(int(zoomedWidth), int(zoomedHeight))

	return nil
}

func (c *camera) Move(position image.Point) {
	c.position = position
	c.worldArea = WorldArea(c.zoom, c.width, c.height, position)
}

func (c *camera) Zoom(newZoom float64) {
	minZoom := c.MinZoomLevel()
	if newZoom < minZoom {
		newZoom = minZoom
	}

	c.zoom = newZoom

	_ = c.Resize(c.width, c.height)
}

func (c *camera) DrawSurface() *ebiten.Image { return c.surface }

func (c *camera) Translation(p image.Point) *ebiten.DrawImageOptions {
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

func (c *camera) DebugInfo() *treemap.TreeMap[string, string] {
	c.debugInfo.Set("x", strconv.Itoa(c.position.X))
	c.debugInfo.Set("y", strconv.Itoa(c.position.Y))
	c.debugInfo.Set("zoom", strconv.FormatFloat(c.zoom, 'f', 2, 64))
	c.debugInfo.Set("w", strconv.Itoa(c.worldArea.Dx()))
	c.debugInfo.Set("h", strconv.Itoa(c.worldArea.Dy()))

	return c.debugInfo
}
