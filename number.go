package topdown

import "golang.org/x/exp/constraints"

// Number constrains type to float or integer.
type Number interface {
	constraints.Float | constraints.Integer
}
