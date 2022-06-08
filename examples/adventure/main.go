package main

import (
	"os"

	"github.com/rs/zerolog/log"

	"github.com/jamestunnell/topdown"
	"github.com/jamestunnell/topdown/engine"
	"github.com/jamestunnell/topdown/resource"
)

func main() {
	resourcesDir, set := os.LookupEnv("RESOURCES_DIR")
	if !set {
		log.Fatal().Msg("RESOURCES_DIR is not set")
	}

	play := &Play{
		PlayerRef: "adventurer.player",
		WorldRef:  "adventure.world",
	}
	types := []resource.Type{
		&PlayerType{},
		&WorldType{},
	}
	cfg := &engine.Config{
		ResourcesDir: resourcesDir,
		StartMode:    play,
		ExtraTypes:   types,
		WindowSize:   topdown.NewSize(800, 600),
		Fullscreen:   false,
	}
	eng := engine.New(cfg)

	if err := eng.Initialize(); err != nil {
		log.Fatal().Err(err).Msg("failed to initialize engine")
	}

	if err := eng.Run(); err != nil {
		log.Fatal().Err(err).Msg("failed to run game engine")
	}
}
