package pcloud

import "net/url"

// GetClearTrashUrl generate url that could be used for clearing trash.
func GetClearTrashUrl(a string) string {
	val := url.Values{"auth": {a}, "folderid": {"0"}}
	uri := url.URL{Scheme: Scheme, Host: EndPoint, Path: ClearTrash, RawQuery: val.Encode()}

	return uri.String()
}
