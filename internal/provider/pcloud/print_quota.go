package pcloud

import "net/url"

// QuotaResponse holds standard response from pCloud API to check available & used
// storage quota.
type QuotaResponse struct {
	Result    int   `json:"result"`    // non 0 result is errors.
	Quota     int64 `json:"quota"`     // Available quota. In bytes.
	UsedQuota int64 `json:"usedquota"` // Used quota. In bytes.
}

// GetQuotaUrl generate url that could be used to get info regarding available
// and used storage quota.
func GetQuotaUrl(a string) string {
	val := url.Values{"auth": {a}}
	uri := url.URL{Scheme: Scheme, Host: EndPoint, Path: Login, RawQuery: val.Encode()}

	return uri.String()
}
