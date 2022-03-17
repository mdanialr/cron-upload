package pcloud

import "net/url"

// DigestResponse holds standard response when requesting digest from pCloud API.
type DigestResponse struct {
	Digest string `json:"digest"` // the digest that only valid for 30s.
}

// GetDigestUrl generate url that could be used to request digest from API.
func GetDigestUrl() string {
	val := url.URL{Scheme: Scheme, Host: EndPoint, Path: GetDigest}
	return val.String()
}
