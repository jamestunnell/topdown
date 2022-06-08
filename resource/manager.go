package resource

import (
	"fmt"
	"path/filepath"
	"strings"

	"golang.org/x/exp/maps"

	"github.com/jamestunnell/topdown/fileindex"
	"github.com/jamestunnell/topdown/registry"
	"github.com/jamestunnell/topdown/sliceutil"
)

//go:generate mockgen -destination=mock_resource/mockmanager.go . Manager

type Manager interface {
	Initialize() error

	TypeNames() []string
	FilePartialPaths(typeName string) []string

	Add(string, Resource)
	Get(string) (Resource, error)
	Clear()
}

type TypeRegistry = registry.Registry[Type]

type manager struct {
	resources    map[string]Resource
	typeRegistry TypeRegistry
	fileIndex    fileindex.FileIndex
	dotReplacer  *strings.Replacer
}

func NewManager(rootDir string, reg TypeRegistry) Manager {
	return &manager{
		dotReplacer:  strings.NewReplacer(".", ""),
		fileIndex:    fileindex.New(rootDir),
		typeRegistry: reg,
		resources:    map[string]Resource{},
	}
}

func (mgr *manager) Initialize() error {
	if err := mgr.fileIndex.Scan(); err != nil {
		return fmt.Errorf("failed to scan root dir: %w", err)
	}

	return nil
}

func (mgr *manager) TypeNames() []string {
	return mgr.typeRegistry.Names()
}

func (mgr *manager) FilePartialPaths(typeName string) []string {
	ext := "." + typeName
	dirAndSlash := mgr.fileIndex.RootDir() + string([]rune{filepath.Separator})
	dirReplacer := strings.NewReplacer(dirAndSlash, "")
	paths := mgr.fileIndex.Paths(fileindex.FilterByExt(ext))

	return sliceutil.Map(paths, func(path string) string {
		return dirReplacer.Replace(path)
	})
}

func (mgr *manager) Add(partialPath string, r Resource) {
	mgr.resources[partialPath] = r
}

func (mgr *manager) Get(partialPath string) (Resource, error) {
	if resource, found := mgr.resources[partialPath]; found {
		return resource, nil
	}

	path := filepath.Join(mgr.fileIndex.RootDir(), partialPath)
	ext := filepath.Ext(partialPath)
	typeName := mgr.dotReplacer.Replace(ext)

	t, found := mgr.typeRegistry.Get(typeName)
	if !found {
		err := fmt.Errorf("type '%s' not found", typeName)

		return nil, err
	}

	resource, err := t.Load(path)
	if err != nil {
		err = fmt.Errorf("failed to load %s: %w", typeName, err)

		return nil, err
	}

	if err = resource.Initialize(mgr); err != nil {
		return nil, fmt.Errorf("failed to initialize: %w", err)
	}

	mgr.resources[partialPath] = resource

	return resource, nil
}

func (mgr *manager) Clear() {
	maps.Clear(mgr.resources)
}
