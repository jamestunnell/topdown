package topdown

import (
	"fmt"
)

// Point is a 2D point.
type Point[T Number] struct {
	X T `json:"x"`
	Y T `json:"y"`
}

// PointSchemaStr is the JSON schema for a point.
const PointSchemaStr = `{
  "$id": "https://github.com/jamestunnell/topdown/point.json",
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "title": "Point",
  "description": "Coordinates for a Point.",
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

// Pt makes a new point.
func Pt[T Number](x, y T) Point[T] {
	return Point[T]{X: x, Y: y}
}

func ParseIntPoint(s string) (Point[int], error) {
	var x int
	var y int

	_, err := fmt.Sscanf(s, "(%d,%d)", &x, &y)
	if err != nil {
		return Point[int]{}, fmt.Errorf("failed to parse '%s' as int point: %w", s, err)
	}

	return Pt(x, y), nil
}

// String makes a string representing the point.
func (p Point[T]) String() string {
	return fmt.Sprintf("(%v,%v)", p.X, p.Y)
}

// Equal checks for equal X and Y values.
func (p Point[T]) Equal(other Point[T]) bool {
	return (p.X == other.X) && (p.Y == other.Y)
}

// Add makes a new point by adding the given and current points.
func (p Point[T]) Add(other Point[T]) Point[T] {
	return Pt(p.X+other.X, p.Y+other.Y)
}

// Add makes a new point by subtracting the given point from the current point.
func (p Point[T]) Sub(other Point[T]) Point[T] {
	return Pt(p.X-other.X, p.Y-other.Y)
}

// In checks if the current point is in the given rectangle.
func (p Point[T]) In(r Rectangle[T]) bool {
	return r.Min.X <= p.X && p.X < r.Max.X &&
		r.Min.Y <= p.Y && p.Y < r.Max.Y
}
