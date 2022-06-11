package animation

import (
	"fmt"

	"github.com/xeipuuv/gojsonschema"
)

var schema *gojsonschema.Schema

const SchemaStr = `{
  "$id": "https://github.com/jamestunnell/topdown/animations.json",
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "title": "Animations",
  "description": "Organizes tagged image set frames into for animations.",
  "type": "object",
  "required": ["spriteSetRef", "frameDuration"],
  "properties": {
    "spriteSetRef": {
      "type": "string",
	  "minLength": 1
	},
	"frameDuration": {
	  "type": "string",
	  "minLength": 1
	}
  }
}`

func LoadSchema() (*gojsonschema.Schema, error) {
	if schema != nil {
		return schema, nil
	}

	l := gojsonschema.NewStringLoader(SchemaStr)
	newSchema, err := gojsonschema.NewSchema(l)
	if err != nil {
		err = fmt.Errorf("failed to make JSON schema: %w", err)

		return nil, err
	}

	schema = newSchema

	return schema, nil
}
