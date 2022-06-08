package engine

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/jamestunnell/topdown"
	"github.com/jamestunnell/topdown/imageresource"
	"github.com/jamestunnell/topdown/imageset"
	"github.com/jamestunnell/topdown/registry"
	"github.com/jamestunnell/topdown/resource"
	"github.com/jamestunnell/topdown/tilegrid"
)

//go:generate mockgen -destination=mock_engine/mockengine.go . Engine

type Engine interface {
	Initialize() error
	Run() error

	ebiten.Game
}

type engine struct {
	config          *Config
	typeRegistry    registry.Registry[resource.Type]
	resourceManager resource.Manager
	// serviceRegistry service.Registry
	mode       Mode
	windowSize topdown.Size
}

type Config struct {
	WindowSize   topdown.Size
	Fullscreen   bool
	ResourcesDir string
	ExtraTypes   []resource.Type
	StartMode    Mode
}

type MakeTypeFunc func() (resource.Type, error)

func New(cfg *Config) Engine {
	tr := registry.New[resource.Type]("resource.Type")

	return &engine{
		config:          cfg,
		typeRegistry:    tr,
		resourceManager: resource.NewManager(cfg.ResourcesDir, tr),
		// serviceRegistry: service.NewRegistry(),
	}
}

func (eng *engine) Initialize() error {
	eng.typeRegistry.Add(imageresource.Types()...)

	isType, err := imageset.NewType()
	if err != nil {
		return fmt.Errorf("failed to make imageset type: %w", err)
	}
	eng.typeRegistry.Add(isType)

	bgType, err := tilegrid.NewBackgroundType()
	if err != nil {
		return fmt.Errorf("failed to make background type: %w", err)
	}
	eng.typeRegistry.Add(bgType)

	for _, t := range eng.config.ExtraTypes {
		eng.typeRegistry.Add(t)
	}

	if err := eng.resourceManager.Initialize(); err != nil {
		return fmt.Errorf("failed to initialize resource manager: %w", err)
	}

	if err := eng.config.StartMode.Initialize(eng.config.WindowSize, eng.resourceManager); err != nil {
		return fmt.Errorf("failed to initialize start mode: %w", err)
	}

	eng.windowSize = eng.config.WindowSize
	eng.mode = eng.config.StartMode

	return nil
}

func (eng *engine) Run() error {
	ebiten.SetFullscreen(eng.config.Fullscreen)
	ebiten.SetWindowSize(int(eng.config.WindowSize.Width), int(eng.config.WindowSize.Height))

	return ebiten.RunGame(eng)
}

func (eng *engine) Update() error {
	newMode, err := eng.mode.Update()
	if err != nil {
		return fmt.Errorf("failure during update: %w", err)
	}

	if newMode != nil {
		eng.resourceManager.Clear()

		newMode.Initialize(eng.windowSize, eng.resourceManager)
	}

	return nil
}

func (eng *engine) Draw(screen *ebiten.Image) {
	eng.mode.Draw(screen)
}

func (eng *engine) Layout(w, h int) (int, int) {
	w, h = eng.mode.Layout(w, h)
	if w != int(eng.windowSize.Width) || h != int(eng.windowSize.Height) {
		eng.windowSize = topdown.NewSize(float64(w), float64(h))
	}

	return w, h
}
