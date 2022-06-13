package engine

import (
	"fmt"

	"github.com/jamestunnell/topdown/registry"
	"github.com/jamestunnell/topdown/resource"
	"github.com/jamestunnell/topdown/sprite"
	"github.com/jamestunnell/topdown/tilegrid"
)

func SetupTypes(reg registry.Registry[resource.Type], extraTypes ...resource.Type) error {
	reg.Add(sprite.Types()...)

	bgType, err := tilegrid.NewBackgroundType()
	if err != nil {
		return fmt.Errorf("failed to make background type: %w", err)
	}

	reg.Add(bgType)

	for _, t := range extraTypes {
		reg.Add(t)
	}

	return nil
}
