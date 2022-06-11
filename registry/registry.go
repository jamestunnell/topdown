package registry

import (
	"github.com/rs/zerolog/log"
	"golang.org/x/exp/maps"
)

// Named can be anything that supplies a name.
type Named interface {
	Name() string
}

// Registry stores and gets named items.
type Registry[T Named] interface {
	Add(vals ...T)
	Names() []string
	Get(name string) (T, bool)
}

type registry[T Named] struct {
	typeName string
	entries  map[string]T
}

// New makes a new registry.
func New[T Named](typeName string) Registry[T] {
	return &registry[T]{
		typeName: typeName,
		entries:  map[string]T{},
	}
}

// Add registers the given items.
func (reg *registry[T]) Add(vals ...T) {
	for _, val := range vals {
		name := val.Name()

		if _, found := reg.entries[name]; found {
			log.Debug().Str("name", name).Msgf("re-registering %s", reg.typeName)
		} else {
			log.Debug().Str("name", name).Msgf("registering %s", reg.typeName)

			reg.entries[name] = val
		}
	}
}

// Names gets all of the registered item names.
func (reg *registry[T]) Names() []string {
	return maps.Keys(reg.entries)
}

// Get gets a registered item by name.
func (reg *registry[T]) Get(name string) (T, bool) {
	val, found := reg.entries[name]

	return val, found
}
