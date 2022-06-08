package tilegrid

import (
	"fmt"

	"github.com/xeipuuv/gojsonschema"

	"github.com/jamestunnell/topdown"
	"github.com/jamestunnell/topdown/jsonfile"
	"github.com/jamestunnell/topdown/resource"
)

type BackgroundType struct {
	schema *gojsonschema.Schema
}

func NewBackgroundType() (resource.Type, error) {
	schema, err := resource.MakeJSONSchema(
		TileGridSchemaStr,
		topdown.VectorSchemaStr,
		topdown.SizeSchemaStr,
		topdown.PixelSchemaStr)
	if err != nil {
		return nil, fmt.Errorf("failed to make JSON schema: %w", err)
	}

	return &BackgroundType{schema: schema}, nil
}

func (pt *BackgroundType) Name() string {
	return "background"
}

func (pt *BackgroundType) Load(path string) (resource.Resource, error) {
	return jsonfile.ReadAndValidate[*Background](path, pt.schema)
}
