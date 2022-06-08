package mathutil_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/jamestunnell/topdown/mathutil"
)

func TestClamp(t *testing.T) {
	testClamp(t, 0, 2)
}

func testClamp[T mathutil.Number](t *testing.T, min, max T) {
	inbetween := (min + max) / 2

	assert.Equal(t, min, mathutil.Clamp(min, min, max))
	assert.Equal(t, inbetween, mathutil.Clamp(inbetween, min, max))
	assert.Equal(t, max, mathutil.Clamp(max, min, max))
	assert.Equal(t, min, mathutil.Clamp(min-1, min, max))
	assert.Equal(t, max, mathutil.Clamp(max+1, min, max))
}
