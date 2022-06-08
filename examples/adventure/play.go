package main

import (
	"fmt"
	"image"
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

	player       *Player
	world        *World
	cam          camera.Camera
	worldDrawing drawing.WorldSystem
	animation    animation.System
	control      control.System
	moveCollide  movecollide.System

	windowWidth, windowHeight int
}

func (p *Play) Initialize(windowSize topdown.Size, mgr resource.Manager) error {
	player, err := resource.GetAs[*Player](mgr, p.PlayerRef)
	if err != nil {
		return fmt.Errorf("failed to get player: %w", err)
	}

	world, err := resource.GetAs[*World](mgr, p.WorldRef)
	if err != nil {
		return fmt.Errorf("failed to get world: %w", err)
	}

	cam, err := camera.New(int(windowSize.Width), int(windowSize.Height))
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

	p.worldDrawing = drawing.NewWorldSystem(int(world.Size.Width), int(world.Size.Height))
	p.moveCollide = moveCollide
	p.control = control.NewSystem()
	p.animation = animation.NewSystem()
	p.windowWidth = int(windowSize.Width)
	p.windowHeight = int(windowSize.Height)

	objs := map[string]any{
		"player": p.player,
		"world":  p.world,
	}

	for id, obj := range objs {
		p.animation.Add(id, obj)
		p.control.Add(id, obj)
		p.worldDrawing.Add(id, obj)
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
	p.cam.Move(image.Pt(int(p.player.Position.X), int(p.player.Position.Y)))

	_, scrollAmount := ebiten.Wheel()
	if scrollAmount > 0 {
		p.cam.Zoom(p.cam.ZoomLevel() + 0.1)
	} else if scrollAmount < 0 {
		p.cam.Zoom(p.cam.ZoomLevel() - 0.1)
	}

	return nil, nil
}

func (p *Play) Draw(screen *ebiten.Image) {
	surface := p.cam.DrawSurface()

	// Clear camera surface
	surface.Clear()
	surface.Fill(color.Black)

	visible := p.cam.WorldArea()

	p.worldDrawing.DrawWorld(visible)

	opts := &ebiten.DrawImageOptions{}

	translateX := float64(0)
	translateY := float64(0)

	if visible.Min.X < 0 {
		translateX = float64(-visible.Min.X)
	}

	if visible.Min.Y < 0 {
		translateY = float64(-visible.Min.Y)
	}

	if translateX != 0 || translateY != 0 {
		opts.GeoM.Translate(translateX, translateY)
	}

	subImg := p.worldDrawing.Surface().SubImage(visible).(*ebiten.Image)

	// draws visible world to camera drawing surface
	surface.DrawImage(subImg, opts)

	// Draw camera to screen and zoom
	p.cam.Blit(screen)
}

func (p *Play) Layout(w, h int) (int, int) {
	if p.windowWidth != w || p.windowHeight != h {
		err := p.cam.Resize(w, h)
		if err != nil {
			log.Warn().
				Err(err).
				Int("width", w).
				Int("height", h).
				Msg("failed to resize camera")

			// stay at the last size that worked
			return p.windowWidth, p.windowHeight
		}

		log.Debug().Int("w", w).Int("h", h).Msg("resized window")

		p.windowWidth = w
		p.windowHeight = h
	}

	return w, h
}
