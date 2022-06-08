package input

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

type Manager interface {
	WatchKey(key ebiten.Key)
	UnwatchKey(key ebiten.Key)
	UpdateKeys(delta time.Duration)

	KeyPressed(key ebiten.Key) bool
	KeyJustPressed(key ebiten.Key) bool
}

type manager struct {
	KeyStates map[ebiten.Key]*KeyState
}

type KeyState struct {
	Pressed     bool
	Duration    time.Duration
	numWatchers int
}

func NewManager() Manager {
	return &manager{
		KeyStates: map[ebiten.Key]*KeyState{},
	}
}

func (m *manager) WatchKey(key ebiten.Key) {
	state, found := m.KeyStates[key]

	if found {
		state.numWatchers++
	} else {
		m.KeyStates[key] = &KeyState{
			Pressed:     false,
			Duration:    0,
			numWatchers: 1,
		}
	}
}

func (m *manager) UnwatchKey(key ebiten.Key) {
	if state, found := m.KeyStates[key]; found {
		state.numWatchers -= 1

		if state.numWatchers < 1 {
			delete(m.KeyStates, key)
		}
	}
}

func (m *manager) UpdateKeys(delta time.Duration) {
	for key, state := range m.KeyStates {
		pressed := ebiten.IsKeyPressed(key)

		// if pressed hasn't changed, increase duration
		if pressed == state.Pressed {
			state.Duration += delta
		} else {
			state.Pressed = pressed
			state.Duration = 0
		}
	}
}

func (m *manager) KeyPressed(key ebiten.Key) bool {
	state, found := m.KeyStates[key]
	if !found {
		return false
	}

	return state.Pressed
}

func (m *manager) KeyJustPressed(key ebiten.Key) bool {
	state, found := m.KeyStates[key]
	if !found {
		return false
	}

	return state.Pressed && state.Duration == 0
}
