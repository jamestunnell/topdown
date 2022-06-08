package topdown

type Size struct {
	Width  float64 `json:"w"`
	Height float64 `json:"h"`
}

const SizeSchemaStr = `{
  "$id": "https://github.com/jamestunnell/prosper/size.json",
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "title": "Size",
  "description": "Dimensions for a 2D size.",
  "type": "object",
  "required": ["w", "h"],
  "properties": {
    "w": {
      "type": "number",
	    "exclusiveMinimum": 0
	  },
	  "h": {
	    "type": "number",
	    "exclusiveMinimum": 0
	  }
  }
}`

func NewSize(w, h float64) Size {
	return Size{Width: w, Height: h}
}
