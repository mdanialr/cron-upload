package pcloud

import (
	"crypto/sha1"
	"encoding/hex"
	"net/url"
	"strings"
)

// User holds credential for authentication.
type User struct {
	Username string // An email.
	Password string // Plain password.
	Token    string // Would be populated later.
}

// TokenResponse json response from pCloud API after sending request to generate token.
type TokenResponse struct {
	Result int    `json:"result"` // non 0 result is errors.
	Auth   string `json:"auth"`   // Generated token from pCloud API.
}

// GenerateTokenUrl Generate url that could be used to get a token for authentication.
func (u *User) GenerateTokenUrl(digest string) string {
	// https://docs.pcloud.com/methods/intro/authentication.html
	// sha1( password + sha1( lowercase of username ) + digest)
	userHash := sha1.Sum([]byte(strings.ToLower(u.Username)))
	userDig := hex.EncodeToString(userHash[:])
	digHash := sha1.Sum([]byte(u.Password + userDig + digest))
	dig := hex.EncodeToString(digHash[:])

	val := url.Values{
		"getauth":        {"1"},
		"username":       {u.Username},
		"digest":         {digest},
		"passworddigest": {dig},
	}
	uri := url.URL{Scheme: Scheme, Host: EndPoint, Path: Login, RawQuery: val.Encode()}

	return uri.String()
}
