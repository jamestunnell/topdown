package main

import (
	"fmt"

	"github.com/jamestunnell/topdown"
	"github.com/jamestunnell/topdown/jsonfile"
	"github.com/jamestunnell/topdown/resource"
	"github.com/jamestunnell/topdown/tilegrid"
)

type WorldType struct {
}

type World struct {
	Size          topdown.Size `json:"size"`
	BackgroundRef string       `json:"backgroundRef"`

	*tilegrid.Background
}

func (t *WorldType) Name() string {
	return "world"
}

func (t *WorldType) Load(path string) (resource.Resource, error) {
	return jsonfile.Read[*World](path)
}

func (m *World) Initialize(mgr resource.Manager) error {
	bg, err := resource.GetAs[*tilegrid.Background](mgr, m.BackgroundRef)
	if err != nil {
		return fmt.Errorf("failed to get background: %w", err)
	}

	m.Background = bg

	return nil
}
