package restest

import (
	"testing"

	"github.com/jamestunnell/topdown/registry"
	"github.com/jamestunnell/topdown/resource"
	"github.com/stretchr/testify/require"
)

func SetupManager(t *testing.T, dir string, types ...resource.Type) resource.Manager {
	reg := registry.New[resource.Type]("type")

	reg.Add(types...)

	mgr := resource.NewManager(dir, reg)

	require.NoError(t, mgr.Initialize())

	return mgr
}
