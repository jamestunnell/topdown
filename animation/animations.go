package animation

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jamestunnell/topdown/imageset"
	"github.com/jamestunnell/topdown/resource"
)

// Animations organizes tagged frames from an image set into animations.
type Animations struct {
	ImageSetRef   string
	FrameDuration time.Duration

	TaggedFrames map[string]Frames
}

type Frames = []*ebiten.Image

type AnimationsJSON struct {
	ImageSetRef   string `json:"imageSetRef"`
	FrameDuration string `json:"frameDuration,omitempty"`
}

type Frame struct {
	Image    *ebiten.Image
	Duration time.Duration
}

const TypeName = "animations"

func NewAnimations(
	imageSetRef string, frameDur time.Duration) *Animations {
	return &Animations{
		ImageSetRef:   imageSetRef,
		FrameDuration: frameDur,
		TaggedFrames:  map[string]Frames{},
	}
}

func (anims *Animations) MarshalJSON() ([]byte, error) {
	animsJSON := &AnimationsJSON{
		ImageSetRef:   anims.ImageSetRef,
		FrameDuration: anims.FrameDuration.String(),
	}

	d, err := json.Marshal(animsJSON)
	if err != nil {
		return []byte{}, err
	}

	return d, nil
}

func (anims *Animations) UnmarshalJSON(d []byte) error {
	var animsJSON AnimationsJSON

	err := json.Unmarshal(d, &animsJSON)
	if err != nil {
		return err
	}

	var frameDur time.Duration

	if animsJSON.FrameDuration != "" {
		frameDur, err = time.ParseDuration(animsJSON.FrameDuration)
		if err != nil {
			return fmt.Errorf("failed to parse frame duration '%s': %w", anims.FrameDuration, err)
		}
	}

	anims.ImageSetRef = animsJSON.ImageSetRef
	anims.FrameDuration = frameDur

	return nil
}

func (anims *Animations) Initialize(mgr resource.Manager) error {
	is, err := resource.GetAs[*imageset.ImageSet](mgr, anims.ImageSetRef)
	if err != nil {
		return fmt.Errorf("failed to get '%s' from dependencies: %w", anims.ImageSetRef, err)
	}

	taggedFrames := map[string]Frames{}
	tags := is.Tags()

	for _, tag := range tags {
		frames, err := MakeFrames(tag, anims.FrameDuration, is)
		if err != nil {
			return fmt.Errorf("failed to make frames for tag '%s': %w", tag, err)
		}

		taggedFrames[tag] = frames
	}

	anims.TaggedFrames = taggedFrames

	return nil
}
