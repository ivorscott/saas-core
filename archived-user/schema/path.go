package schema

import (
	"path/filepath"
	"runtime"
)

var (
	_, b, _, _ = runtime.Caller(0)
	basepath   = filepath.Dir(b)
)

// PWD provided an absolute path to the caller
func PWD() string {
	return basepath
}
