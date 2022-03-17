package pcloud

import "net/url"

// LogoutResponse json response from pCloud API after sending request to log out and invalidate token.
type LogoutResponse struct {
	IsDeleted bool `json:"auth_deleted"` // Determine whether token invalidation was successful or not.
}

// GetLogoutUrl generate url that could be used to log out meaning remove the given token.
func GetLogoutUrl(a string) string {
	val := url.Values{"auth": {a}}
	uri := url.URL{Scheme: Scheme, Host: EndPoint, Path: Logout, RawQuery: val.Encode()}

	return uri.String()
}
