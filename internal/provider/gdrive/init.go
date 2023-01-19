package gdrive

import (
	"context"
	"fmt"
	"log"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

// Init return ready to use Google Apis client that use the given service
// account token path as the credential.
func Init(serviceTokenPath string) (*drive.FilesService, error) {
	// init context and read the given token filepath at once
	ctx := context.Background()
	tk, _ := os.ReadFile(serviceTokenPath)
	// create new Google Api credential based on the above token
	token, err := google.CredentialsFromJSON(ctx, tk, drive.DriveScope)
	if err != nil {
		return nil, fmt.Errorf("failed to init Google Drive client: %s", err)
	}
	// create new http client along with the oauth token for Google Api call
	cl := oauth2.NewClient(ctx, token.TokenSource)
	// create new Google Drive client service
	svc, err := drive.NewService(ctx, option.WithHTTPClient(cl))
	if err != nil {
		log.Fatalln("failed to create drive service instance:", err)
	}
	// return just the file service instead of drive service
	return drive.NewFilesService(svc), nil
}
