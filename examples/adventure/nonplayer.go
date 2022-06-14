package main

import (
	"github.com/jamestunnell/topdown/jsonfile"
	"github.com/jamestunnell/topdown/resource"
)

type NonPlayerType struct {
}

type NonPlayer struct {
	*Character
}

func (t *NonPlayerType) Name() string {
	return "nonplayer"
}

func (t *NonPlayerType) Load(path string) (resource.Resource, error) {
	return jsonfile.Read[*NonPlayer](path)
}
