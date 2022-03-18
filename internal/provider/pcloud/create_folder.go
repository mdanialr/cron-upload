package pcloud

import (
	"fmt"
	"net/url"
)

// CreateFolderResponse standard response from pCloud API call for creating folder.
type CreateFolderResponse struct {
	Result    int                      `json:"result"`   // non 0 result is errors.
	IsCreated bool                     `json:"created"`  // Whether this folder already created or not before this API call requested.
	Metadata  CreateFolderMetaResponse `json:"metadata"` // Contains metadata about this folder.
}

// CreateFolderMetaResponse additional metadata for this folder.
type CreateFolderMetaResponse struct {
	Path string `json:"path"`     // Full path for this folder.
	Name string `json:"name"`     // Name of this folder.
	Id   int    `json:"folderid"` // Id of this folder. This is so important for uploading files.
}

// GetCreateFolderUrl generate url that could be used to create folder from the
// given path value.
func GetCreateFolderUrl(a string, p string) string {
	val := url.Values{"auth": {a}}
	uri := url.URL{Scheme: Scheme, Host: EndPoint, Path: CreateFolder, RawQuery: val.Encode()}

	res := fmt.Sprintf("%s&path=%s", uri.String(), p)
	return res
}
