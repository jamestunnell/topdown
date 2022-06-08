package animation

import "time"

//go:generate mockgen -destination=mock_animation/mockanimatable.go . Animatable

// Animatable is the component used in the animation system.
type Animatable interface {
	UpdateAnimation(delta time.Duration)
}
