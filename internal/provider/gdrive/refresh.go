package gdrive

import (
	"context"
	"fmt"
	"os"

	"github.com/mdanialr/cron-upload/internal/config"
	"github.com/mdanialr/cron-upload/internal/provider/gdrive/token"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v2"
)

// Refresh exchange authorization code for new refresh token and save it to 'config.Provider.Auth'.
func Refresh(conf *config.Model) error {
	// 1. Read Provider.Cred file
	b, err := os.ReadFile(conf.Provider.Cred)
	if err != nil {
		return fmt.Errorf("failed to read Provider.Cred file: %s\n", err)
	}
	oConfig, err := google.ConfigFromJSON(b, drive.DriveScope)
	if err != nil {
		return fmt.Errorf("failed to create oauth2.Config instance from Provider.Cred file: %s", err)
	}

	// 2. Get new token
	tok, err := refreshToken(oConfig)
	if err != nil {
		return fmt.Errorf("failed to get new token from web: %s", err)
	}

	// 3. Delete old token.json file
	os.Remove(conf.Provider.Token)
	// 4. Save newly retrieved token to token.json file
	if err := token.SaveToken(conf.Provider.Token, tok); err != nil {
		return fmt.Errorf("failed to save new oauth2.Token instance to token.json file in: %s with error: %s\n", conf.Provider.Token, err)
	}

	// 5. Prepare NewToken instance then assign new refresh token to instance
	newRefreshTokenI := token.NewToken{
		RefreshToken: tok.RefreshToken,
		ClientID:     oConfig.ClientID,
		ClientSecret: oConfig.ClientSecret,
		TokenUrl:     oConfig.Endpoint.TokenURL,
	}

	// 6. Delete old file Provider.Auth file
	os.Remove(conf.Provider.Auth)
	// 7. Save new Provider.Auth file
	if err := token.SaveRefreshToken(conf.Provider.Auth, newRefreshTokenI); err != nil {
		return fmt.Errorf("failed to save new Provider.Auth file: %s", err)
	}

	return nil
}

// refreshToken exchange authorization code for new refresh.
func refreshToken(config *oauth2.Config) (*oauth2.Token, error) {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline, oauth2.ApprovalForce)
	fmt.Printf("Go to the following link in your browser then type the authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		return nil, fmt.Errorf("failed to read authorization code: %s", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve token from web: %v", err)
	}

	return tok, nil
}
