package config

import (
	"fmt"
	"strings"
)

// Upload holds actual data about which path that would be uploaded and what the folder's name in provider.
type Upload []struct {
	Folders Folder `yaml:"folder"` // all folders.
}

// Folder holds detail folders that would be uploaded to provider.
type Folder struct {
	Name string `yaml:"name"` // the name of the folder that would be created (if not exist) in provider.
	Path string `yaml:"path"` // actual full path where target files would be uploaded to provider.
}

// Sanitization sanitize every single folder.
func (f *Folder) Sanitization() error {
	if f.Path == "" {
		return fmt.Errorf("`path` field is required")
	}
	if f.Name == "" {
		f.Name = f.Path
	}

	f.Name = strings.TrimLeft(f.Name, "/")
	f.Name = strings.TrimRight(f.Name, "/")

	return nil
}

// Sanitization sanitize all folders one by one.
func (u *Upload) Sanitization() error {
	for i, uu := range *u {
		if err := uu.Folders.Sanitization(); err != nil {
			return fmt.Errorf("failed to sanitizing #%d folder in config: %s", i, err)
		}
	}
	return nil
}
