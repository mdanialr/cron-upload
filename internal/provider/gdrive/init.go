package gdrive

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/mdanialr/cron-upload/internal/config"
	"github.com/mdanialr/cron-upload/internal/provider/gdrive/token"
)

// InitToken initialize token using auth.json file to retrieve token.json for authentication.
func InitToken(conf *config.Model, client *http.Client) error {
	// 1. Prepare NewToken instance
	newTokenI := token.NewToken{}

	// 2. Read auth.json and inject their values to NewToken instance
	b, err := os.ReadFile(conf.Provider.Auth)
	if err != nil {
		return fmt.Errorf("failed to read auth.json file in: %s with error: %s\n", conf.Provider.Auth, err)
	}
	if err := json.Unmarshal(b, &newTokenI); err != nil {
		return fmt.Errorf("failed to binding auth.json to NewToken model: %s\n", err)
	}

	newToken, err := newTokenI.RenewToken(client)
	if err != nil {
		return fmt.Errorf("failed to renew token: %s\n", err)
	}

	// 3. Delete old token.json file
	os.Remove(conf.Provider.Token)
	// 4. Save newly retrieved token to token.json file
	if err := token.SaveToken(conf.Provider.Token, newToken); err != nil {
		return fmt.Errorf("failed to save new oauth2.Token instance to token.json file in: %s with error: %s\n", conf.Provider.Token, err)
	}

	return nil
}
