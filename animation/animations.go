package animation

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/jamestunnell/topdown/resource"
	"github.com/jamestunnell/topdown/sprite"
)

// Animations organizes tagged frames from an image set into animations.
type Animations struct {
	SpriteSetRef  string
	FrameDuration time.Duration

	TaggedImages map[string]Images
	Controller   *Controller
}

type Images = []*ebiten.Image

type AnimationsJSON struct {
	SpriteSetRef  string `json:"spriteSetRef"`
	FrameDuration string `json:"frameDuration"`
}

type Frame struct {
	Image    *ebiten.Image
	Duration time.Duration
}

const TypeName = "animations"

func NewAnimations(
	spriteSetRef string, frameDur time.Duration) *Animations {
	return &Animations{
		SpriteSetRef:  spriteSetRef,
		FrameDuration: frameDur,
		TaggedImages:  map[string]Images{},
	}
}

func (anims *Animations) MarshalJSON() ([]byte, error) {
	animsJSON := &AnimationsJSON{
		SpriteSetRef:  anims.SpriteSetRef,
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

	frameDur, err = time.ParseDuration(animsJSON.FrameDuration)
	if err != nil {
		return fmt.Errorf("failed to parse frame duration '%s': %w", anims.FrameDuration, err)
	}

	anims.SpriteSetRef = animsJSON.SpriteSetRef
	anims.FrameDuration = frameDur

	return nil
}

func (anims *Animations) Initialize(mgr resource.Manager) error {
	sh, err := resource.GetAs[*sprite.Sheet](mgr, anims.SpriteSetRef)
	if err != nil {
		return fmt.Errorf("failed to get '%s' from dependencies: %w", anims.SpriteSetRef, err)
	}

	taggedImages := map[string]Images{}
	tags := sh.Tags()

	for _, tag := range tags {
		frames, err := FrameImages(tag, sh)
		if err != nil {
			return fmt.Errorf("failed to make frames for tag '%s': %w", tag, err)
		}

		taggedImages[tag] = frames
	}

	anims.TaggedImages = taggedImages

	anims.Controller = NewController()

	return nil
}

func (anims *Animations) Start(tag string) bool {
	images, found := anims.TaggedImages[tag]
	if !found {
		return false
	}

	return anims.Controller.StartAnimation(tag, images, anims.FrameDuration)
}
