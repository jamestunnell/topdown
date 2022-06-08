package imageset

const SchemaStr = `{
  "$id": "https://github.com/jamestunnell/prosper/imageset.json",
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "title": "Image set",
  "type": "object",
  "required": ["imageRef", "subImages"],
  "properties": {
	"imageRef": {
		"type": "string",
		"minLength": 1
	},
	"subImages": {
		"type": "array",
		"items": {
		"$ref": "#/$defs/subImage"
		}
	}
  },
  "$defs": {
    "subImage": {
      "title": "Sub-image definition",
	  "type": "object",
	  "required": ["x", "y", "width", "height"],
      "properties": {
		"x": {
			"type": "integer",
			"minimum": 0
		},
		"y": {
			"type": "integer",
			"minimum": 0
		},
		"width": {
			"type": "integer",
			"minimum": 1
		},
		"height": {
			"type": "integer",
			"minimum": 1
		},
		"tags": {
			"type": "array",
			"items": {
			  "type": "string"
			}
		},
		"props": {
			"type": "object"
		}
      }
    }
  }
}`
