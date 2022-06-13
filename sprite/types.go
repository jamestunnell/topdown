package sprite

import (
	"github.com/rs/zerolog/log"

	"github.com/jamestunnell/topdown/resource"
)

var spritesType resource.Type

func init() {
	t, err := NewSheetType()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to make sprite sheet type")
	}

	spritesType = t
}

func Types() []resource.Type {
	return []resource.Type{
		spritesType,
		&ImageType{name: ImageTypePNG},
		&ImageType{name: ImageTypeJPG},
		&ImageType{name: ImageTypeJPEG},
	}
}
