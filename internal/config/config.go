package config

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

// Model holds all data from config file.
type Model struct {
	LogDir    string   `yaml:"log"`        // Where info & error log for this app is written.
	MaxWorker uint8    `yaml:"max_worker"` // Max number of workers that do the job which is upload file to cloud provider.
	LogFile   *os.File // File instance that would be using by logger to write into.
	Provider  Provider `yaml:"provider"` // detail about which provider is used.
}

// NewConfig read io.Reader then map and load the value to the returned Model.
func NewConfig(fileBuf io.Reader) (mod *Model, err error) {
	buf := new(bytes.Buffer)

	if _, err := buf.ReadFrom(fileBuf); err != nil {
		return mod, fmt.Errorf("failed to read from file buffer: %v", err)
	}

	if err := yaml.Unmarshal(buf.Bytes(), &mod); err != nil {
		return mod, fmt.Errorf("failed to unmarshal: %v", err)
	}

	return
}

// Sanitization check and sanitize config Model's instance.
func (m *Model) Sanitization() error {
	if m.LogDir == "" {
		m.LogDir = "/tmp/"
	}
	if !strings.HasPrefix(m.LogDir, "/") {
		m.LogDir = "/" + m.LogDir
	}
	if !strings.HasSuffix(m.LogDir, "/") {
		m.LogDir += "/"
	}

	if m.MaxWorker < 1 {
		m.MaxWorker = 1
	}

	if m.Provider.Name == "" {
		return fmt.Errorf("`provider.name` field is required")
	}

	return nil
}
