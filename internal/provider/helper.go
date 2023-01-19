package provider

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// LookupRoute recursively lookup parent's name for the given route in the
// payloads stack.
func LookupRoute(payloads []*Payload, route string) string {
	for _, payload := range payloads {
		firstPart := strings.Split(route, "/")[0]
		if payload.Name == firstPart {
			if len(payload.Parent) > 0 {
				r := LookupRouteName(payloads, payload.Parent[0])            // search for the parent's route name
				return LookupRoute(payloads, fmt.Sprintf("%s/%s", r, route)) // recursively search for the parents until it has no parents
			}
		}
	}
	return route
}

// LookupRouteName search for matched id in the payloads stack from the given
// id.
func LookupRouteName(payloads []*Payload, id string) string {
	for _, payload := range payloads {
		if payload.Id == id {
			return payload.Name
		}
	}
	return ""
}

// NewWithFile return new payload instance along with io.ReadCloser for File
// field from the given filepath. Also add the given parents to Parents field.
// return error if failed to read the given filepath.
func NewWithFile(filePath string, parents ...string) (*Payload, error) {
	// try open the given filepath and return error if any
	flInstance, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file path: %s", err)
	}
	// create new payload along with the file reader
	payload := Payload{
		Name:   filepath.Base(flInstance.Name()),
		Parent: parents,
		File:   flInstance,
	}
	return &payload, nil
}
