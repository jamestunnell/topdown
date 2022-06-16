package main

import (
	"strconv"
	"time"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/jamestunnell/topdown"
	"github.com/jamestunnell/topdown/debug"
	"github.com/jamestunnell/topdown/input"
	"github.com/jamestunnell/topdown/jsonfile"
	"github.com/jamestunnell/topdown/resource"
)

type PlayerType struct {
}

type Player struct {
	*Character

	debugData *debug.Dataset
}

const (
	OneOverSqrtTwo = 0.7071067811865475244
	PlayerSpeed    = 25.0
)

func (t *PlayerType) Name() string {
	return "player"
}

func (t *PlayerType) Load(path string) (resource.Resource, error) {
	return jsonfile.Read[*Player](path)
}

func (p *Player) Initialize(mgr resource.Manager) error {
	p.debugData = debug.NewDataset()

	return p.Character.Initialize(mgr)
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

func (p *Player) DebugData() *debug.Dataset {
	p.debugData.Set("x", strconv.FormatFloat(p.Position.X, 'f', 2, 64))
	p.debugData.Set("y", strconv.FormatFloat(p.Position.Y, 'f', 2, 64))

	return p.debugData
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

	// if dirChanged {
	// 	log.Debug().Float64("x", dir.X).Float64("y", dir.Y).Msg("direction changed")
	// }

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
