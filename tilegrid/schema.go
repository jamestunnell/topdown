package tilegrid

const TileGridSchemaStr = `{
  "$id": "https://example.com/person.schema.json",
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "title": "Image set",
  "type": "object",
  "required": ["minPosition", "tileSize", "tileDefs", "tileRows"],
  "properties": {
	"minPosition": { "$ref": "https://github.com/jamestunnell/prosper/vector.json" },
	"tileSize": { "$ref": "https://github.com/jamestunnell/prosper/size.json" },
	"tileDefs": {
		"type": "object",
		"patternProperties" :{
			".*": {"$ref": "#/defs/tileDef"}
		}
	},
	"tileRows": {
		"type": "array",
		"items": {
			"type": "string",
			"minLength": 1
		},
		"minLength": 1
	}
  },
  "defs": {
	"tileDef": {
		"required": ["imageSetRef","startPixel"],
		"properties": {
			"imageSetRef": {
				"type": "string",
				"minLength": 1
			},
			"startPixel": { "$ref": "https://github.com/jamestunnell/prosper/pixel.json" }
		}
	}
  }
}`
