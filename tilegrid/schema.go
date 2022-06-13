package tilegrid

const TileGridSchemaStr = `{
  "$id": "https://example.com/person.schema.json",
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "title": "Image set",
  "type": "object",
  "required": ["origin", "tileSize", "tileLinks", "tileRows"],
  "properties": {
	"origin": { "$ref": "https://github.com/jamestunnell/topdown/vector.json" },
	"tileSize": { "$ref": "https://github.com/jamestunnell/topdown/size.json" },
	"tileLinks": {
		"type": "object",
		"patternProperties" :{
			".*": {"type": "string", "minLength": 1}
		}
	},
	"tileRows": {
		"type": "array",
		"items": {"type": "string", "minLength": 1},
		"minLength": 1
	}
  }
}`
