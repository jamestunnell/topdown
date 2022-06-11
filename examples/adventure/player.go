package main

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/rs/zerolog/log"
	"github.com/zergon321/cirno"

	"github.com/jamestunnell/topdown"
	"github.com/jamestunnell/topdown/animation"
	"github.com/jamestunnell/topdown/drawing"
	"github.com/jamestunnell/topdown/input"
	"github.com/jamestunnell/topdown/jsonfile"
	"github.com/jamestunnell/topdown/resource"
)

type PlayerType struct {
}

type Player struct {
	Animations   *animation.Animations `json:"animations"`
	ColliderSize topdown.Size[float64] `jons:"colliderSize"`
	Position     topdown.Vector        `json:"position"`

	Collider  cirno.Shape
	Direction topdown.Vector
	Velocity  topdown.Vector
}

const (
	OneOverSqrtTwo = 0.7071067811865475244
	PlayerSpeed    = 65.0
)

func (t *PlayerType) Name() string {
	return "player"
}

func (t *PlayerType) Load(path string) (resource.Resource, error) {
	return jsonfile.Read[*Player](path)
}

func (p *Player) Initialize(mgr resource.Manager) error {
	if err := p.Animations.Initialize(mgr); err != nil {
		return fmt.Errorf("failed to initialize animations: %w", err)
	}

	if !p.Animations.Start("idleDown") {
		return errors.New("failed to start idle animation")
	}

	cCenter := cirno.NewVector(p.Position.X, p.Position.Y)
	colliderRect, err := cirno.NewRectangle(
		cCenter, p.ColliderSize.Width, p.ColliderSize.Height, 0)

	if err != nil {
		return fmt.Errorf("failed to make collider rect: %w", err)
	}

	p.Collider = colliderRect
	p.Velocity = topdown.Vector{}
	p.Direction = topdown.Vec(0, 1)

	return nil
}

func (p *Player) WorldLayer() int {
	return drawing.LayerForeground
}

func (p *Player) WorldDraw(world *ebiten.Image, visible topdown.Rectangle[float64]) {
	img := p.Animations.Controller.CurrentFrameImage()
	w, h := img.Size()
	wFlt := float64(w)
	hFlt := float64(h)

	// the image bottom lines up with the collider bottom
	maxY := p.Position.Y + p.ColliderSize.Height/2.0
	minX := p.Position.X - wFlt/2.0

	worldArea := topdown.Rect(minX, maxY-hFlt, minX+wFlt, maxY)

	if worldArea.Intersect(visible).Empty() {
		return
	}

	// draw the image relative to the char position

	op := &ebiten.DrawImageOptions{}

	op.GeoM.Translate(float64(worldArea.Min.X), float64(worldArea.Min.Y))

	world.DrawImage(img, op)
}

func (p *Player) WatchKeys() []ebiten.Key {
	return []ebiten.Key{
		ebiten.KeyArrowLeft,
		ebiten.KeyArrowRight,
		ebiten.KeyArrowUp,
		ebiten.KeyArrowDown,
	}
}

func (p *Player) Control(deltaSec float64, inputMgr input.Manager) {
	p.controlMovement(deltaSec, inputMgr)
}

func (p *Player) PlanMovement(deltaSec float64) topdown.Vector {
	return p.Velocity.Multiply(deltaSec)
}

func (p *Player) Move(moveDiff topdown.Vector) {
	if !moveDiff.Zero() {
		p.Position = p.Position.Add(moveDiff)
	}
}

func (p *Player) UpdateAnimation(delta time.Duration) {
	p.Animations.Controller.Update(delta)
}

func (p *Player) ColliderShape() cirno.Shape {
	return p.Collider
}

func (p *Player) ResolveCollision(moveDiff cirno.Vector, shapes cirno.Shapes) cirno.Vector {
	shapeIDs := []string{}
	for s := range shapes {
		shapeIDs = append(shapeIDs, s.Data().(string))
	}

	log.Debug().Str("shapeIDs", strings.Join(shapeIDs, ",")).Msg("resolving collision")

	// most basic collision resolution: back up to where there is no collision
	newPos, _, _, err := cirno.Approximate(p.Collider, moveDiff, 0, shapes, 100, false)
	if err != nil {
		log.Warn().Err(err).Msg("failed to resolve collision")

		return cirno.Zero()
	}

	pos := cirno.NewVector(p.Position.X, p.Position.Y)

	return newPos.Subtract(pos)
}

func (p *Player) controlMovement(deltaSec float64, inputMgr input.Manager) {
	l := inputMgr.KeyPressed(ebiten.KeyArrowLeft)
	r := inputMgr.KeyPressed(ebiten.KeyArrowRight)
	u := inputMgr.KeyPressed(ebiten.KeyArrowUp)
	d := inputMgr.KeyPressed(ebiten.KeyArrowDown)

	if l && r {
		l = false
		r = false
	}

	if u && d {
		u = false
		d = false
	}

	dir := topdown.Vector{}

	switch {
	case u && l:
		dir.X = -OneOverSqrtTwo
		dir.Y = -OneOverSqrtTwo
	case u && r:
		dir.X = OneOverSqrtTwo
		dir.Y = -OneOverSqrtTwo
	case d && l:
		dir.X = -OneOverSqrtTwo
		dir.Y = OneOverSqrtTwo
	case d && r:
		dir.X = OneOverSqrtTwo
		dir.Y = OneOverSqrtTwo
	case u:
		dir.Y = -1
	case d:
		dir.Y = 1
	case l:
		dir.X = -1
	case r:
		dir.X = 1
	}

	dirChanged := !dir.Equal(p.Direction)

	if dirChanged {
		log.Debug().Float64("x", dir.X).Float64("y", dir.Y).Msg("direction changed")
	}

	if dir.Zero() && dirChanged {
		p.Velocity = topdown.Vector{}

		// idle direction based on the last direction, which should be moving
		if p.Direction.Y == 0 {
			if p.Direction.X > 0 {
				p.Animations.Start("idleRight")
			} else {
				p.Animations.Start("idleLeft")
			}
		} else {
			if p.Direction.Y > 0 {
				p.Animations.Start("idleDown")
			} else {
				p.Animations.Start("idleUp")
			}
		}

		p.Direction = dir
		p.Velocity = topdown.Vector{}

		return
	}

	if dirChanged {
		switch {
		case u:
			p.Animations.Start("walkUp")
		case d:
			p.Animations.Start("walkDown")
		case l:
			p.Animations.Start("walkLeft")
		case r:
			p.Animations.Start("walkRight")
		}

		p.Direction = dir
		p.Velocity = dir.Multiply(PlayerSpeed)

		return
	}

	return
}
