package tilegrid_test

import (
	"testing"

	"github.com/jamestunnell/topdown/drawing"
	"github.com/jamestunnell/topdown/tilegrid"
	"github.com/stretchr/testify/assert"
)

func TestBackgroundIsDrawable(t *testing.T) {
	testBackgroundIs[drawing.WorldDrawable](t)
}

func testBackgroundIs[T any](t *testing.T) {
	var x interface{}

	x = &tilegrid.Background{}

	_, ok := x.(T)

	assert.True(t, ok)
}
