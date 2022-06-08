package jsonfile

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"reflect"

	"github.com/xeipuuv/gojsonschema"
)

func Read[T any](path string) (T, error) {
	return readAs[T](path, nil)
}

func ReadAndValidate[T any](path string, schema *gojsonschema.Schema) (T, error) {
	return readAs[T](path, schema)
}

func readAs[T any](path string, schema *gojsonschema.Schema) (T, error) {
	var z T

	d, err := ioutil.ReadFile(path)
	if err != nil {
		return z, fmt.Errorf("failed to read file: %w", err)
	}

	if schema != nil {
		loader := gojsonschema.NewBytesLoader(d)

		result, err := schema.Validate(loader)
		if err != nil {
			return z, fmt.Errorf("failed to validate JSON with schema: %w", err)
		}

		if !result.Valid() {
			return z, &ErrValidation{Errors: result.Errors()}
		}
	}

	ptrType := reflect.TypeOf(z)
	baseType := ptrType.Elem()
	val := reflect.New(baseType)

	err = json.Unmarshal(d, val.Interface())
	if err != nil {
		return z, fmt.Errorf("failed to unmarshal: %w", err)
	}

	return val.Interface().(T), nil
}
