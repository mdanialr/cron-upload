package token

import (
	"encoding/json"
	"fmt"
	"os"
)

// SaveRefreshToken Saves a Provider.Auth to a file path.
func SaveRefreshToken(path string, auth NewToken) error {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("failed to write Provider.Auth to path: %s\n", err)
	}
	defer f.Close()

	return json.NewEncoder(f).Encode(auth)
}
