package imageset

import (
	"fmt"

	"github.com/xeipuuv/gojsonschema"

	"github.com/jamestunnell/topdown/jsonfile"
	"github.com/jamestunnell/topdown/resource"
)

type Type struct {
	schema *gojsonschema.Schema
}

func NewType() (resource.Type, error) {
	schema, err := resource.MakeJSONSchema(SchemaStr)
	if err != nil {
		return nil, fmt.Errorf("failed to make JSON schema: %w", err)
	}

	return &Type{schema: schema}, nil
}

func (pt *Type) Name() string {
	return "imageset"
}

func (pt *Type) Load(path string) (resource.Resource, error) {
	return jsonfile.ReadAndValidate[*ImageSet](path, pt.schema)
}
