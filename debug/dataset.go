package debug

import "golang.org/x/exp/slices"

type Dataset struct {
	sortedKeys []string
	data       map[string]string
}

func NewDataset() *Dataset {
	return &Dataset{
		sortedKeys: []string{},
		data:       map[string]string{},
	}
}
func (im *Dataset) Set(key, val string) {
	if idx, found := slices.BinarySearch(im.sortedKeys, key); !found {
		im.sortedKeys = slices.Insert(im.sortedKeys, idx, key)
	}

	im.data[key] = val
}

func (im *Dataset) Get(key string) (string, bool) {
	if val, found := im.data[key]; found {
		return val, true
	}

	return "", false
}

func (im *Dataset) SortedKeys() []string {
	return im.sortedKeys
}
