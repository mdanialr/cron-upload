package provider

import (
	"fmt"
	"io"

	"golang.org/x/exp/slices"
)

type cloud string

var (
	// GoogleDrive one of the supported cloud provider.
	GoogleDrive cloud = "drive"
	// S3Bucket one of the supported cloud provider.
	S3Bucket cloud = "s3"
	// Support contain all supported cloud providers.
	Support []cloud
)

func init() {
	Support = append(Support, GoogleDrive) // support for Google Drive
	Support = append(Support, S3Bucket)    // support for AWS S3 Bucket
}

// Payload common data structure that's used by provider.
type Payload struct {
	Name      string        // Name can be a file or folder name.
	Id        string        // Id can be a file or folder id.
	Parent    []string      // Parent the parent folder where this file/folder is in.
	CreatedAt string        // CreatedAt the time at which the file was created (RFC-3339 date-time)
	File      io.ReadCloser // File can be a file reader that can be used to upload.
}

// Cloud provider interface. Every supported/implemented cloud storage
// provider should use this interface as the guideline.
type Cloud interface {
	// GetFolders retrieve all folders from provider. Optionally Use the
	// first given parent as the folder's parent. Does not return error
	// if no data found but return error if the API call failed.
	GetFolders(parent ...string) ([]*Payload, error)
	// CreateFolder create new folder by the given name as the folder name
	// and create it inside the provided first parent id. Return the id of
	// the newly created folder if success.
	CreateFolder(name string, parent ...string) (string, error)
	// GetFiles retrieve all non-folder data for the given folder id.
	// Does not return error if no data found but return error if the
	// API call failed.
	GetFiles(folderId string) ([]*Payload, error)
	// UploadFile upload the given payload and optionally set the upload chunk
	// size to reduce memory allocation but may slow down upload duration.
	// Name field as the file name, Parent field as the folder id where this
	// file will be uploaded, File field is the binary data of the file. Should
	// defer close the File field.
	UploadFile(payload *Payload, chunkSize ...int) (*Payload, error)
	// Delete do delete the given id. id can be either the id of a file
	// or folder.
	Delete(id string) error
}

// ValidateSupportedClouds make sure the provided provider name is currently
// supported.
func ValidateSupportedClouds(providerName string) error {
	if !slices.Contains(Support, cloud(providerName)) {
		return fmt.Errorf("the given provider name is not supported at this moment. currently support: %s", Support)
	}
	return nil
}
