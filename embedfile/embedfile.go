package embedfile

import (
	"embed"
)

var (
	//go:embed assets
	s embed.FS
)

func ReadFs(path string) ([]byte, error) {
	data, err := s.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return data, nil
}
