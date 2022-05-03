package token

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"golang.org/x/oauth2"
)

// NewToken necessary format to request new token.
type NewToken struct {
	ClientID     string `json:"client_id"`     // oauth2 client id.
	ClientSecret string `json:"client_secret"` // oauth2 client secret.
	TokenUrl     string `json:"token_uri"`     // token url where the new token request would be sent. Should exist in credentials.json file.
	RefreshToken string `json:"refresh"`       // refresh token for exchanging new token to google apis.
}

// RenewToken get new token from google apis.
func (n *NewToken) RenewToken(cl *http.Client) (*oauth2.Token, error) {
	urlValue := url.Values{
		"client_id":     {n.ClientID},
		"client_secret": {n.ClientSecret},
		"refresh_token": {n.RefreshToken}, "grant_type": {"refresh_token"},
	}

	resp, err := cl.PostForm(n.TokenUrl, urlValue)
	if err != nil {
		return nil, fmt.Errorf("failed sent POST request to renew token: %s\n", err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body after sent POST request: %s\n", err)
	}
	defer resp.Body.Close()

	var newToken *oauth2.Token
	if err := json.Unmarshal(body, &newToken); err != nil {
		return nil, fmt.Errorf("failed to bind response body to oauth2.Token model: %s\n", err)
	}

	return newToken, nil
}
