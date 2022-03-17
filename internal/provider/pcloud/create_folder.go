package pcloud

import (
	"fmt"
	"net/url"
)

type CreateFolderResponse struct {
	IsCreated bool `json:"created"` // Whether this folder already created or not before this API call requested.
}

// GetCreateFolderUrl generate url that could be used to create folder from the
// given path value.
func GetCreateFolderUrl(a string, p string) string {
	val := url.Values{"auth": {a}}
	uri := url.URL{Scheme: Scheme, Host: EndPoint, Path: CreateFolder, RawQuery: val.Encode()}

	res := fmt.Sprintf("%s&path=%s", uri.String(), p)
	return res
}
