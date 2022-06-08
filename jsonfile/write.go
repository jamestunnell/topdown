package jsonfile

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"io/ioutil"
)

func Write(path string, val any) error {
	d, err := json.MarshalIndent(val, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to unmarshal value: %w", err)
	}

	if err = ioutil.WriteFile(path, d, fs.ModePerm); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}
