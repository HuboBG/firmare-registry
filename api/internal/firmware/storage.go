package firmware

import "path/filepath"

// Storage defines filesystem layout for firmware binaries.
type Storage struct {
	BaseDir string
}

func (s Storage) Dir(typeName, version string) string {
	return filepath.Join(s.BaseDir, typeName, version)
}

func (s Storage) FilePath(typeName, version string) string {
	return filepath.Join(s.Dir(typeName, version), "firmware.bin")
}
