package resource

import (
	"fmt"

	"github.com/xeipuuv/gojsonschema"

	"github.com/jamestunnell/topdown/sliceutil"
)

func MakeJSONSchema(schemaStr string, reqdSchemas ...string) (*gojsonschema.Schema, error) {
	reqdLoaders := sliceutil.Map(reqdSchemas, gojsonschema.NewStringLoader)

	sl := gojsonschema.NewSchemaLoader()
	if err := sl.AddSchemas(reqdLoaders...); err != nil {
		return nil, fmt.Errorf("failed to add required schemas: %w", err)
	}

	schema, err := sl.Compile(gojsonschema.NewStringLoader(schemaStr))
	if err != nil {
		err = fmt.Errorf("failed to make JSON schema: %w", err)

		return nil, err
	}

	return schema, nil
}
