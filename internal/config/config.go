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
	LogDir     string   `yaml:"log"` // Where info & error log for this app is written.
	LogFile    *os.File // File instance that would be using by logger to write into.
	Provider   Provider `yaml:"provider"`    // detail about which provider is used.
	Upload     Upload   `yaml:"upload"`      // detail about folders that would be uploaded.
	RootFolder string   `yaml:"root_folder"` // all Upload folders would be created inside this root folder.
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

	if m.Provider.Name == "" {
		return fmt.Errorf("`provider.name` field is required")
	}
	if m.Provider.Name == "drive" {
		if m.Provider.Auth == "" {
			return fmt.Errorf("`provider.auth` field is required")
		}
		if m.Provider.Cred == "" {
			return fmt.Errorf("`provider.cred` is required")
		}

		// if provided then:
		if m.Provider.Token != "" {
			// make sure has leading and trailing slash
			if !strings.HasPrefix(m.Provider.Token, "/") {
				m.Provider.Token = "/" + m.Provider.Token
			}
			if !strings.HasSuffix(m.Provider.Token, "/") {
				m.Provider.Token += "/"
			}
		}
		// if not provided then use this default value
		if m.Provider.Token == "" {
			m.Provider.Token = "/tmp/cron-upload-token.json"
		}
		// append file name if it does not have one yet
		if !strings.HasSuffix(m.Provider.Token, "cron-upload-token.json") {
			m.Provider.Token += "cron-upload-token.json"
		}
	}

	if m.RootFolder == "" {
		m.RootFolder = "Cron-Backups"
	}

	return nil
}
