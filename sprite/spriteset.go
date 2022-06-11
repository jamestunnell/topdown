package sprite

import (
	"fmt"

	mapset "github.com/deckarep/golang-set/v2"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/xeipuuv/gojsonschema"

	"github.com/jamestunnell/topdown"
	"github.com/jamestunnell/topdown/jsonfile"
	"github.com/jamestunnell/topdown/resource"
)

type SpriteSet struct {
	ImageRef string    `json:"imageRef"`
	Sprites  []*Sprite `json:"sprites"`
}

type SpritesType struct {
	schema *gojsonschema.Schema
}

func NewSpriteSet(imageRef string, sprites ...*Sprite) *SpriteSet {
	return &SpriteSet{
		ImageRef: imageRef,
		Sprites:  sprites,
	}
}

func NewSpritesType() (resource.Type, error) {
	schema, err := resource.MakeJSONSchema(
		SpritesSchemaStr, topdown.SizeSchemaStr, topdown.PointSchemaStr)
	if err != nil {
		return nil, fmt.Errorf("failed to make JSON schema: %w", err)
	}

	return &SpritesType{schema: schema}, nil
}

func (s *SpriteSet) Initialize(mgr resource.Manager) error {
	ir, err := resource.GetAs[*Image](mgr, s.ImageRef)
	if err != nil {
		return fmt.Errorf("failed to get image resource: %w", err)
	}

	for _, sprite := range s.Sprites {
		sprite.Initialize(ir.Image)
	}

	return nil
}

func (s *SpriteSet) Tags() []string {
	tags := mapset.NewSet[string]()

	for _, sprite := range s.Sprites {
		for _, tag := range sprite.Tags {
			tags.Add(tag)
		}
	}

	return tags.ToSlice()
}

func (s *SpriteSet) FindSprite(start topdown.Point[int]) (*ebiten.Image, *Sprite, bool) {
	for _, sprite := range s.Sprites {
		if start.Equal(sprite.Start) {
			if sprite.Image == nil {
				return nil, nil, false
			}

			return sprite.Image, sprite, true
		}
	}

	return nil, nil, false
}

func (pt *SpritesType) Name() string {
	return "spriteset"
}

func (pt *SpritesType) Load(path string) (resource.Resource, error) {
	return jsonfile.ReadAndValidate[*SpriteSet](path, pt.schema)
}
