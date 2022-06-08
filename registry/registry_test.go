package registry_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/jamestunnell/topdown/registry"
)

type named struct {
	name string
}

func TestRegistry(t *testing.T) {
	reg := registry.New[*named]("*named")

	assert.Empty(t, reg.Names())

	n1 := &named{name: "fido"}
	n2 := &named{name: "fifi"}

	reg.Add(n1)

	assert.Equal(t, []string{n1.name}, reg.Names())

	reg.Add(n1)

	assert.Equal(t, []string{n1.name}, reg.Names())

	reg.Add(n2)

	assert.ElementsMatch(t, []string{n1.name, n2.name}, reg.Names())

	n, found := reg.Get(n1.name)

	assert.True(t, found)
	assert.Equal(t, n1, n)

	n, found = reg.Get("unknown")

	assert.False(t, found)
	assert.Nil(t, n)
}

func (n *named) Name() string {
	return n.name
}
