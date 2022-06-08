package topdown

import "math"

type Vector struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

const VectorSchemaStr = `{
  "$id": "https://github.com/jamestunnell/prosper/vector.json",
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

func NewVector(x, y float64) Vector {
	return Vector{X: x, Y: y}
}

func (v Vector) Pair() (float64, float64) {
	return v.X, v.Y
}

func (v Vector) Equal(w Vector) bool {
	return v.X == w.X && v.Y == w.Y
}

func (v Vector) Unit() Vector {
	return v.Resize(1)
}

func (v Vector) Resize(newMag float64) Vector {
	return v.Multiply(newMag / v.Magnitude())
}

func (v Vector) Multiply(scale float64) Vector {
	return NewVector(v.X*scale, v.Y*scale)
}

func (v Vector) Add(w Vector) Vector {
	return NewVector(v.X+w.X, v.Y+w.Y)
}

func (v Vector) Zero() bool {
	return v.X == 0 && v.Y == 0
}

func (v Vector) Magnitude() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}
