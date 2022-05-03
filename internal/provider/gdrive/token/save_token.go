package token

import (
	"encoding/json"
	"fmt"
	"os"

	"golang.org/x/oauth2"
)

// SaveToken Saves a token to a file path.
func SaveToken(path string, token *oauth2.Token) error {
	fmt.Println("Saving credential file to:", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("failed to write oauth2.Token to path: %s\n", err)
	}
	defer f.Close()

	return json.NewEncoder(f).Encode(token)
}
