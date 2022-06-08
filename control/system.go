package control

import (
	"time"

	"github.com/jamestunnell/topdown/input"
	"github.com/rs/zerolog/log"
	"golang.org/x/exp/maps"
)

//go:generate mockgen -destination=mock_control/mocksystem.go . System

type System interface {
	Add(id string, x any)
	Remove(id string)
	Clear()

	Control(deltaSec float64)
}

type system struct {
	controllables map[string]Controllable
	inputMgr      input.Manager
}

func NewSystem() System {
	return &system{
		controllables: map[string]Controllable{},
		inputMgr:      input.NewManager(),
	}
}

func (s *system) Add(id string, x any) {
	if c, ok := x.(Controllable); ok {
		for _, key := range c.WatchKeys() {
			log.Debug().Str("id", id).Stringer("key", key).Msg("watching key")

			s.inputMgr.WatchKey(key)
		}

		log.Debug().Str("id", id).Msg("adding controllable")

		s.controllables[id] = c
	}
}

func (s *system) Remove(id string) {
	c, found := s.controllables[id]
	if !found {
		return
	}

	for _, key := range c.WatchKeys() {
		log.Debug().Str("id", id).Stringer("key", key).Msg("un-watching key")

		s.inputMgr.UnwatchKey(key)
	}

	delete(s.controllables, id)
}

func (s *system) Clear() {
	for id, c := range s.controllables {
		for _, key := range c.WatchKeys() {
			log.Debug().Str("id", id).Stringer("key", key).Msg("un-watching key")

			s.inputMgr.UnwatchKey(key)
		}
	}

	maps.Clear(s.controllables)
}

func (s *system) Control(deltaSec float64) {
	s.inputMgr.UpdateKeys(time.Duration(deltaSec * 1e9))

	for _, c := range s.controllables {
		c.Control(deltaSec, s.inputMgr)
	}
}
