package topdown_test

import (
	"testing"

	"github.com/jamestunnell/topdown"
	"github.com/stretchr/testify/assert"
)

func TestRectangleCenter(t *testing.T) {
	ctr := topdown.Rect(0, 0, 50, 50).Center()

	assert.True(t, ctr.Equal(topdown.Point[int]{25, 25}))

	ctr = topdown.Rect(-50, -50, 0, 0).Center()

	assert.True(t, ctr.Equal(topdown.Point[int]{-25, -25}))

	ctr = topdown.Rect(-50, -50, 50, 50).Center()

	assert.True(t, ctr.Equal(topdown.Point[int]{0, 0}))
}
