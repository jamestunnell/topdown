package topdown

import "math"

// Vector is a 2D vector.
type Vector struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// VectorSchemaStr is the JSON schema for a vector.
const VectorSchemaStr = `{
  "$id": "https://github.com/jamestunnell/topdown/vector.json",
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "title": "Vector",
  "description": "Components for a 2D vector.",
  "type": "object",
  "required": ["x", "y"],
  "properties": {
    "x": { "type": "number" },
	  "y": { "type": "number" }
  }
}`

// Vec makes a new vector value.
func Vec(x, y float64) Vector {
	return Vector{X: x, Y: y}
}

// AsPoint make a point from the vector.
func (v Vector) AsPoint() Point[float64] {
	return Pt(v.X, v.Y)
}

// Pair gets the vector X-Y pair.
func (v Vector) Pair() (float64, float64) {
	return v.X, v.Y
}

// Equal checks for equal X and Y values.
func (v Vector) Equal(w Vector) bool {
	return v.X == w.X && v.Y == w.Y
}

// Unit makes a new vector resized to a magnitude of 1.
func (v Vector) Unit() Vector {
	return v.Resize(1)
}

// Resize makes a new vector resized to the given magnitude.
func (v Vector) Resize(newMag float64) Vector {
	return v.Multiply(newMag / v.Magnitude())
}

// Multiply makes a new vector with scaled X and Y.
func (v Vector) Multiply(scale float64) Vector {
	return Vec(v.X*scale, v.Y*scale)
}

// Add makes a new vector with the vectors added.
func (v Vector) Add(w Vector) Vector {
	return Vec(v.X+w.X, v.Y+w.Y)
}

// Zero returns true if X and Y are 0.
func (v Vector) Zero() bool {
	return v.X == 0 && v.Y == 0
}

// Magnitude returns the vector magnitude.
func (v Vector) Magnitude() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}
