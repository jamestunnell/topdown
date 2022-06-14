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
	Size          topdown.Size[float64] `json:"size"`
	BackgroundRef string                `json:"background"`
	NPCRefs       []string              `json:"npcs"`

	*tilegrid.Background

	NPCs []*NonPlayer
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

	npcs, err := resource.GetManyAs[*NonPlayer](mgr, m.NPCRefs)
	if err != nil {
		return fmt.Errorf("failed to get NPCs: %w", err)
	}

	m.Background = bg
	m.NPCs = npcs

	return nil
}
