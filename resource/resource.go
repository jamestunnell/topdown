package resource

import (
	"fmt"
	"reflect"

	"github.com/jamestunnell/topdown"
)

//go:generate mockgen -destination=mock_resource/mockresource.go . Resource

type Resource interface {
	Initialize(Manager) error
}

func As[T Resource](r Resource) (T, error) {
	var val T
	var ok bool

	if val, ok = r.(T); !ok {
		expectedType := reflect.TypeOf(val)
		err := topdown.NewErrUnexpectedType(r, expectedType.Name())

		return val, err
	}

	return val, nil
}

func GetAs[T Resource](mgr Manager, ref string) (T, error) {
	var val T

	r, err := mgr.Get(ref)
	if err != nil {
		return val, err
	}

	return As[T](r)
}

func GetManyAs[T Resource](mgr Manager, refs []string) ([]T, error) {
	vals := make([]T, len(refs))

	for i, ref := range refs {
		val, err := GetAs[T](mgr, ref)
		if err != nil {
			return []T{}, fmt.Errorf("failed to get resource '%s': %w", ref, err)
		}

		vals[i] = val
	}

	return vals, nil
}
