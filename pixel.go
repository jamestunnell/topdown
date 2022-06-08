package topdown

import "fmt"

type Pixel struct {
	X int `json:"x"`
	Y int `json:"y"`
}

const PixelSchemaStr = `{
  "$id": "https://github.com/jamestunnell/prosper/pixel.json",
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "title": "Pixel",
  "description": "Coordinates for a pixel.",
  "type": "object",
  "required": ["x", "y"],
  "properties": {
	"x": {
	  "type": "integer",
	  "minimum": 0
	},
	"y": {
	  "type": "integer",
	  "minimum": 0
	}
  }
}`

func NewPixel(x, y int) Pixel {
	return Pixel{X: x, Y: y}
}

func (p Pixel) String() string {
	return fmt.Sprintf("(%d,%d)", p.X, p.Y)
}
