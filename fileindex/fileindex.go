package fileindex

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"time"

	collect "github.com/sxyazi/go-collection"
	"golang.org/x/exp/maps"
)

//go:generate mockgen -destination=mock_fileindex/mockfileindex.go . FileIndex

type FileFilter struct {
	Extension string
	Directory string
}

type FileFilterOpt func(*FileFilter)

type FileIndex interface {
	RootDir() string
	LastScanComplete() time.Time
	Paths(filterOpts ...FileFilterOpt) []string
	PathDirs() []string
	PathExts() []string

	Scan() error
}

type index struct {
	rootDir          string
	lastScanComplete time.Time
	pathsByExt       map[string][]string
	pathsByDir       map[string][]string
	allPaths         []string
}

func New(rootDir string) FileIndex {
	return &index{
		rootDir:          rootDir,
		lastScanComplete: time.Time{},
		pathsByExt:       map[string][]string{},
		pathsByDir:       map[string][]string{},
		allPaths:         []string{},
	}
}

func FilterByDir(dir string) FileFilterOpt {
	return func(f *FileFilter) {
		f.Directory = dir
	}
}

func FilterByExt(ext string) FileFilterOpt {
	return func(f *FileFilter) {
		f.Extension = ext
	}
}

func (idx *index) RootDir() string {
	return idx.rootDir
}

func (idx *index) LastScanComplete() time.Time {
	return idx.lastScanComplete
}

func (idx *index) PathDirs() []string {
	return maps.Keys(idx.pathsByDir)
}

func (idx *index) PathExts() []string {
	return maps.Keys(idx.pathsByExt)
}

func (idx *index) Paths(filterOpts ...FileFilterOpt) []string {
	filter := &FileFilter{}
	for _, opt := range filterOpts {
		opt(filter)
	}

	switch {
	case filter.Directory != "" && filter.Extension != "":
		pathsByDir := idx.pathsByDir[filter.Directory]
		pathsByExt := idx.pathsByExt[filter.Extension]

		return collect.Merge(pathsByDir, pathsByExt)
	case filter.Directory != "":
		return idx.pathsByDir[filter.Directory]
	case filter.Extension != "":
		return idx.pathsByExt[filter.Extension]
	}

	return idx.allPaths
}

func (idx *index) Scan() error {
	paths := []string{}
	info, err := os.Stat(idx.rootDir)
	if err != nil {
		return fmt.Errorf("failed to stat root dir '%s': %w", idx.rootDir, err)
	}

	if !info.IsDir() {
		return fmt.Errorf("root '%s' is not a dir", idx.rootDir)
	}

	walkErr := filepath.WalkDir(idx.rootDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() && d.Type().IsRegular() {
			paths = append(paths, path)
		}

		return nil
	})

	if walkErr != nil {
		return fmt.Errorf("failed to walk root dir '%s': %w", idx.rootDir, walkErr)
	}

	pathsByExt := map[string][]string{}
	pathsByDir := map[string][]string{}

	for _, path := range paths {
		dir := filepath.Dir(path)
		ext := filepath.Ext(path)

		pathsByDir[dir] = append(pathsByDir[dir], path)
		pathsByExt[ext] = append(pathsByExt[ext], path)
	}

	idx.allPaths = paths
	idx.pathsByExt = pathsByExt
	idx.pathsByDir = pathsByDir
	idx.lastScanComplete = time.Now()

	return nil
}
