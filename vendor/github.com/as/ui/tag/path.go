package tag

import (
	"path"
	"path/filepath"
)

func IsAbs(name string) bool {
	return filepath.IsAbs(name) || path.IsAbs(name)
}
