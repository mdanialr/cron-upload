package pcloud

import "net/url"

// GetUploadFileUrl generate url that could be used to upload a file to pCLoud.
func GetUploadFileUrl(a string, folderId string) string {
	val := url.Values{"auth": {a}, "folderid": {folderId}, "nopartial": {"1"}, "renameifexists": {"1"}}
	uri := url.URL{Scheme: Scheme, Host: EndPoint, Path: UploadFile, RawQuery: val.Encode()}

	return uri.String()
}
