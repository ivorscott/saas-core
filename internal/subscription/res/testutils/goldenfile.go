package testutils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

const goldenDir = "golden"

// GoldenConfig manages golden file updates.
type GoldenConfig struct {
	Update bool
}

// NewGoldenConfig creates a new GoldenConfig manager.
func NewGoldenConfig(shouldUpdate bool) *GoldenConfig {
	return &GoldenConfig{
		Update: shouldUpdate,
	}
}

// ShouldUpdate returns true if golden files should update.
func (u *GoldenConfig) ShouldUpdate() bool {
	return u.Update
}

// LoadGoldenFile loads json from golden files directory and returns the json as string. Returns early with
// empty string when the file doesn't exist.
func LoadGoldenFile(target interface{}, name string) string {
	file := filepath.Clean(resFile(goldenDir, name))
	if _, err := os.ReadFile(file); err != nil {
		// abort -- file doesn't exist
		return ""
	}
	return LoadJSON(target, goldenDir, name)
}

// LoadJSON loads a json file from the golden files directory and returns the json as string. LoadJSON panics on errors.
func LoadJSON(target interface{}, pathElem ...string) string {
	file := filepath.Clean(resFile(pathElem...))

	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		panic(fmt.Sprintf("unable to read test data file %v: %s", file, err))
	}

	if target != nil {
		if err = json.Unmarshal(bytes, target); err != nil {
			panic(fmt.Sprintf("unable to unmarshal json data: %s", err))
		}
	}

	return string(bytes)
}

// SaveGoldenFile saves a data structure as json to a file in the golden files directory. SaveGoldenFile panics on errors.
func SaveGoldenFile(data interface{}, name string) {
	file := filepath.Clean(resFile(goldenDir, name))
	if _, err := os.ReadFile(file); err != nil {
		if _, err = os.OpenFile(file, os.O_CREATE, 0600); err != nil {
			panic("can't create new golden file")
		}
	}
	SaveJSON(data, goldenDir, name)
}

// SaveJSON saves a data structure as json to a file in the golden files directory. SaveJSON panics on errors.
func SaveJSON(data interface{}, pathElem ...string) {
	encoded, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		panic(fmt.Sprintf("unable to marshal test data: %s", err))
	}

	encoded = append(encoded, byte('\n'))
	file := resFile(pathElem...)

	err = ioutil.WriteFile(file, encoded, os.ModePerm)
	if err != nil {
		panic(fmt.Sprintf("unable to write test data file %v: %s", file, err))
	}
}

// MarshalJSON simple helper that marshals data into a JSON string.
func MarshalJSON(data interface{}) string {
	encoded, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		panic(fmt.Sprintf("unable to marshal test data: %s", err))
	}
	return string(encoded)
}
