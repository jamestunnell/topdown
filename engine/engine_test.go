package engine_test

import (
	"errors"
	"io/ioutil"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jamestunnell/topdown"
	"github.com/jamestunnell/topdown/engine"
	"github.com/jamestunnell/topdown/engine/mock_engine"
	"github.com/jamestunnell/topdown/resource"
)

func TestEngineResourceManagerFailsToInit(t *testing.T) {
	ctrl := gomock.NewController(t)
	mode := mock_engine.NewMockMode(ctrl)

	cfg := &engine.Config{
		ResourcesDir: "unknown",
		StartMode:    mode,
		ExtraTypes:   []resource.Type{},
	}

	eng := engine.New(cfg)

	assert.Error(t, eng.Initialize())
}

func TestEngineStartModeFailsToInit(t *testing.T) {
	dir, err := ioutil.TempDir("", "enginetest")

	require.NoError(t, err)

	defer os.RemoveAll(dir)

	ctrl := gomock.NewController(t)
	mode := mock_engine.NewMockMode(ctrl)

	cfg := &engine.Config{
		ResourcesDir: dir,
		StartMode:    mode,
		ExtraTypes:   []resource.Type{},
		WindowSize:   topdown.NewSize(100, 100),
	}

	eng := engine.New(cfg)

	mode.EXPECT().Initialize(cfg.WindowSize, gomock.Any()).Return(errors.New(""))

	assert.Error(t, eng.Initialize())
}

func TestEngine(t *testing.T) {
	dir, err := ioutil.TempDir("", "enginetest")

	require.NoError(t, err)

	defer os.RemoveAll(dir)

	ctrl := gomock.NewController(t)
	mode := mock_engine.NewMockMode(ctrl)

	cfg := &engine.Config{
		ResourcesDir: dir,
		StartMode:    mode,
		ExtraTypes:   []resource.Type{},
		WindowSize:   topdown.NewSize(200, 200),
	}

	eng := engine.New(cfg)

	mode.EXPECT().Initialize(cfg.WindowSize, gomock.Any()).Return(nil)

	assert.NoError(t, eng.Initialize())
}
