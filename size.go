package topdown

import "fmt"

// Size is a 2D size.
type Size[T Number] struct {
	Width  T `json:"w"`
	Height T `json:"h"`
}

// SizeSchemaStr is the JSON schema for a size.
const SizeSchemaStr = `{
  "$id": "https://github.com/jamestunnell/topdown/size.json",
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

// Sz makes a new size.
func Sz[T Number](w, h T) Size[T] {
	return Size[T]{Width: w, Height: h}
}

// Equal checks for equal width and height values.
func (s Size[T]) Equal(other Size[T]) bool {
	return (s.Width == other.Width) && (s.Height == other.Height)
}

// Center makes a point that would be at the center of a rectangle with the current size, starting at (0,0).
func (s Size[T]) Center() Point[T] {
	return Pt(s.Width/T(2), s.Height/T(2))
}

// String makes a string representing the size.
func (s Size[T]) String() string {
	return fmt.Sprintf("%v x %v", s.Width, s.Height)
}
