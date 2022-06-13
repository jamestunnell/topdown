package engine

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/jamestunnell/topdown"
	"github.com/jamestunnell/topdown/registry"
	"github.com/jamestunnell/topdown/resource"
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
	windowSize topdown.Size[int]
}

type Config struct {
	WindowSize   topdown.Size[int]
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
	err := SetupTypes(eng.typeRegistry, eng.config.ExtraTypes...)
	if err != nil {
		return fmt.Errorf("failed to set up types: %w", err)
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
	if w != eng.windowSize.Width || h != eng.windowSize.Height {
		eng.windowSize = topdown.Sz[int](w, h)
	}

	return w, h
}
