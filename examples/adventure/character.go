package main

import (
	"errors"
	"fmt"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/rs/zerolog/log"
	"github.com/zergon321/cirno"

	"github.com/jamestunnell/topdown"
	"github.com/jamestunnell/topdown/animation"
	"github.com/jamestunnell/topdown/camera"
	"github.com/jamestunnell/topdown/drawing"
	"github.com/jamestunnell/topdown/resource"
)

type Character struct {
	Animations   *animation.Animations `json:"animations"`
	ColliderSize topdown.Size[float64] `jons:"colliderSize"`
	Position     topdown.Vector        `json:"position"`

	Collider  cirno.Shape
	Direction topdown.Vector
	Velocity  topdown.Vector
}

func (ch *Character) Initialize(mgr resource.Manager) error {
	if err := ch.Animations.Initialize(mgr); err != nil {
		return fmt.Errorf("failed to initialize animations: %w", err)
	}

	if !ch.Animations.Start("idleDown") {
		return errors.New("failed to start idle animation")
	}

	cCenter := cirno.NewVector(ch.Position.X, ch.Position.Y)
	colliderRect, err := cirno.NewRectangle(
		cCenter, ch.ColliderSize.Width, ch.ColliderSize.Height, 0)

	if err != nil {
		return fmt.Errorf("failed to make collider rect: %w", err)
	}

	ch.Collider = colliderRect
	ch.Velocity = topdown.Vector{}
	ch.Direction = topdown.Vec(0, 1)

	return nil
}

func (ch *Character) DrawLayer() int {
	return drawing.LayerWorldForeground
}

func (ch *Character) DrawSortValue() float64 {
	return ch.maxY()
}

func (ch *Character) Draw(screen *ebiten.Image, cam camera.Camera) {
	img := ch.Animations.Controller.CurrentFrameImage()
	w, h := img.Size()
	wFlt := float64(w)
	hFlt := float64(h)

	// the image bottom lines up with the collider bottom
	maxY := ch.maxY()
	minX := ch.Position.X - wFlt/2.0

	rect := topdown.Rect(minX, maxY-hFlt, minX+wFlt, maxY)
	visible := cam.WorldArea()

	if rect.Intersect(visible).Empty() {
		return
	}

	opts := &ebiten.DrawImageOptions{}
	zoom := cam.ZoomLevel()

	if zoom != 1 {
		opts.GeoM.Scale(zoom, zoom)
	}

	min, _ := cam.ConvertWorldToScreen(rect.Min)

	opts.GeoM.Translate(min.X, min.Y)

	screen.DrawImage(img, opts)
}

func (ch *Character) UpdateAnimation(delta time.Duration) {
	ch.Animations.Controller.Update(delta)
}

func (ch *Character) ColliderShape() cirno.Shape {
	return ch.Collider
}

func (ch *Character) ResolveCollision(moveDiff cirno.Vector, shapes cirno.Shapes) cirno.Vector {
	// shapeIDs := []string{}
	// for s := range shapes {
	// 	shapeIDs = append(shapeIDs, s.Data().(string))
	// }

	// log.Debug().Str("shapeIDs", strings.Join(shapeIDs, ",")).Msg("resolving collision")

	// most basic collision resolution: back up to where there is no collision
	newPos, _, _, err := cirno.Approximate(ch.Collider, moveDiff, 0, shapes, 100, false)
	if err != nil {
		log.Warn().Err(err).Msg("failed to resolve collision")

		return cirno.Zero()
	}

	pos := cirno.NewVector(ch.Position.X, ch.Position.Y)

	return newPos.Subtract(pos)
}

func (ch *Character) maxY() float64 {
	return ch.Position.Y + ch.ColliderSize.Height/2.0
}
