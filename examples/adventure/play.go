package main

import (
	"fmt"
	"image/color"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/rs/zerolog/log"

	"github.com/jamestunnell/topdown"
	"github.com/jamestunnell/topdown/animation"
	"github.com/jamestunnell/topdown/camera"
	"github.com/jamestunnell/topdown/control"
	"github.com/jamestunnell/topdown/drawing"
	"github.com/jamestunnell/topdown/engine"
	"github.com/jamestunnell/topdown/movecollide"
	"github.com/jamestunnell/topdown/resource"
)

type Play struct {
	PlayerRef, WorldRef string

	player      *Player
	world       *World
	cam         camera.Camera
	drawing     drawing.System
	animation   animation.System
	control     control.System
	moveCollide movecollide.System

	screenSize topdown.Size[int]
}

func (p *Play) Initialize(screenSize topdown.Size[int], mgr resource.Manager) error {
	player, err := resource.GetAs[*Player](mgr, p.PlayerRef)
	if err != nil {
		return fmt.Errorf("failed to get player: %w", err)
	}

	world, err := resource.GetAs[*World](mgr, p.WorldRef)
	if err != nil {
		return fmt.Errorf("failed to get world: %w", err)
	}

	cam, err := camera.New(screenSize)
	if err != nil {
		return fmt.Errorf("failed to make camera: %w", err)
	}

	moveCollide, err := movecollide.NewSystem(world.Size.Width, world.Size.Height)
	if err != nil {
		return fmt.Errorf("failed to make move-collide system: %w", err)
	}

	p.player = player
	p.world = world
	p.cam = cam

	p.drawing = drawing.NewSystem(cam)
	p.moveCollide = moveCollide
	p.control = control.NewSystem()
	p.animation = animation.NewSystem()
	p.screenSize = screenSize

	objs := map[string]any{
		"camera": cam,
		"player": p.player,
		"world":  p.world,
	}

	for i, npc := range p.world.NPCs {
		objs[p.world.NPCRefs[i]] = npc
	}

	for id, obj := range objs {
		p.animation.Add(id, obj)
		p.control.Add(id, obj)
		p.drawing.Add(id, obj)
		p.moveCollide.Add(id, obj)
	}

	return nil
}

func (p *Play) Update() (engine.Mode, error) {
	dt := (16667 * time.Microsecond).Seconds()

	p.control.Control(dt)

	p.moveCollide.MoveCollide(dt)

	p.animation.Animate(dt)

	// Update the camera position and zoom
	p.cam.Move(p.player.Position.AsPoint())

	_, scrollAmount := ebiten.Wheel()
	if scrollAmount > 0 {
		p.cam.Zoom(p.cam.ZoomLevel() + 0.1)
	} else if scrollAmount < 0 {
		p.cam.Zoom(p.cam.ZoomLevel() - 0.1)
	}

	return nil, nil
}

func (p *Play) Draw(screen *ebiten.Image) {
	screen.Clear()
	screen.Fill(color.Black)

	p.drawing.Draw(screen)
}

func (p *Play) Layout(w, h int) (int, int) {
	sz := topdown.Sz(w, h)

	if !sz.Equal(p.screenSize) {
		err := p.cam.Resize(sz)
		if err != nil {
			log.Warn().
				Err(err).
				Int("width", w).
				Int("height", h).
				Msg("failed to resize camera")

			// stay at the last size that worked
			return p.screenSize.Width, p.screenSize.Height
		}

		log.Debug().Int("w", w).Int("h", h).Msg("resized window")

		p.screenSize = sz
	}

	return w, h
}
