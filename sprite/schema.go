package sprite

const SpritesSchemaStr = `{
  "$id": "https://github.com/jamestunnell/topdown/sprites.json",
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "title": "Sprites",
  "description": "Set of sprites from a sprite sheet image",
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
	  "required": ["start", "size"],
      "properties": {
		"start": { "$ref": "https://github.com/jamestunnell/topdown/point.json"},
		"size": { "$ref": "https://github.com/jamestunnell/topdown/size.json"},
		"tags": {
			"type": "array",
			"items": {"type": "string"}
		}
      }
    }
  }
}`
