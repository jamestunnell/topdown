package sprite

import (
	"fmt"

	mapset "github.com/deckarep/golang-set/v2"
	"github.com/xeipuuv/gojsonschema"

	"github.com/jamestunnell/topdown"
	"github.com/jamestunnell/topdown/jsonfile"
	"github.com/jamestunnell/topdown/resource"
)

type Sheet struct {
	ImageRef string    `json:"imageRef"`
	Sprites  []*Sprite `json:"sprites"`
}

type SheetType struct {
	schema *gojsonschema.Schema
}

const SheetSchemaStr = `{
  "$id": "https://github.com/jamestunnell/topdown/spritesheet.json",
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "title": "Sprite sheet",
  "description": "Sprites from a sprite sheet image",
  "type": "object",
  "required": ["imageRef", "sprites"],
  "properties": {
	"imageRef": {
		"type": "string",
		"minLength": 1
	},
	"sprites": {
		"type": "array",
		"items": {
		  "$ref": "#/$defs/sprite"
		}
	}
  },
  "$defs": {
    "sprite": {
      "title": "Sprite",
	  "type": "object",
	  "required": ["id", "origin", "size"],
      "properties": {
		"id": {"type": "string", "minLength": 1},
		"origin": { "$ref": "https://github.com/jamestunnell/topdown/point.json"},
		"size": { "$ref": "https://github.com/jamestunnell/topdown/size.json"},
		"tags": {
			"type": "array",
			"items": {"type": "string"}
		}
      }
    }
  }
}`

func NewSheet(imageRef string, sprites ...*Sprite) *Sheet {
	return &Sheet{
		ImageRef: imageRef,
		Sprites:  sprites,
	}
}

func NewSheetType() (resource.Type, error) {
	schema, err := resource.MakeJSONSchema(
		SheetSchemaStr, topdown.SizeSchemaStr, topdown.PointSchemaStr)
	if err != nil {
		return nil, fmt.Errorf("failed to make JSON schema: %w", err)
	}

	return &SheetType{schema: schema}, nil
}

func (s *Sheet) Initialize(mgr resource.Manager) error {
	ir, err := resource.GetAs[*Image](mgr, s.ImageRef)
	if err != nil {
		return fmt.Errorf("failed to get image resource: %w", err)
	}

	for _, sprite := range s.Sprites {
		sprite.Initialize(ir.Image)
	}

	return nil
}

func (s *Sheet) Tags() []string {
	tags := mapset.NewSet[string]()

	for _, sprite := range s.Sprites {
		for _, tag := range sprite.Tags {
			tags.Add(tag)
		}
	}

	return tags.ToSlice()
}

func (s *Sheet) FindSpriteByID(id string) (*Sprite, bool) {
	for _, sprite := range s.Sprites {
		if id == sprite.ID {
			return sprite, true
		}
	}

	return nil, false
}

func (s *Sheet) FindSpriteByOrigin(origin topdown.Point[int]) (*Sprite, bool) {
	for _, sprite := range s.Sprites {
		if origin.Equal(sprite.Origin) {
			return sprite, true
		}
	}

	return nil, false
}

func (pt *SheetType) Name() string {
	return "spritesheet"
}

func (pt *SheetType) Load(path string) (resource.Resource, error) {
	return jsonfile.ReadAndValidate[*Sheet](path, pt.schema)
}
