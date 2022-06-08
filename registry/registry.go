package registry

import (
	"github.com/rs/zerolog/log"
	"golang.org/x/exp/maps"

	"github.com/jamestunnell/topdown"
)

type Registry[T topdown.Named] interface {
	Add(vals ...T)
	Names() []string
	Get(name string) (T, bool)
}

type registry[T topdown.Named] struct {
	typeName string
	entries  map[string]T
}

func New[T topdown.Named](typeName string) Registry[T] {
	return &registry[T]{
		typeName: typeName,
		entries:  map[string]T{},
	}
}

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

func (reg *registry[T]) Names() []string {
	return maps.Keys(reg.entries)
}

func (reg *registry[T]) Get(name string) (T, bool) {
	val, found := reg.entries[name]

	return val, found
}
