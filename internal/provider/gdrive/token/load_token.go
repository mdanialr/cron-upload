package token

import (
	"encoding/json"
	"os"

	"golang.org/x/oauth2"
)

// LoadToken retrieves an oauth2 token from a local file.
func LoadToken(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var tok *oauth2.Token
	err = json.NewDecoder(f).Decode(&tok)
	return tok, err
}
