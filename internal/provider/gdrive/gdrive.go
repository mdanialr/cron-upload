package gdrive

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/mdanialr/cron-upload/internal/config"
	"github.com/mdanialr/cron-upload/internal/provider/gdrive/token"
	"github.com/mdanialr/cron-upload/internal/scan"
	"golang.org/x/oauth2"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

// GoogleDrive run the job which is upload all files to Google Drive provider.
func GoogleDrive(conf *config.Model) {
	oAuthConfig := &oauth2.Config{}
	ctx := context.Background()

	tok, err := token.LoadToken(conf.Provider.Token)
	if err != nil {
		log.Println("failed to read token.json:", err)
	}
	client := oAuthConfig.Client(ctx, tok)

	dr, err := drive.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Println("failed to create new drive service instance:", err)
	}

	if err := token.CheckTokenValidity(dr); err != nil {
		// 1. Prepare NewToken instance
		newTokenI := token.NewToken{}

		// 2. Read auth.json and inject their values to NewToken instance
		b, err := os.ReadFile("auth.json")
		if err != nil {
			log.Fatalln("failed to read auth.json file:", err)
		}
		if err := json.Unmarshal(b, &newTokenI); err != nil {
			log.Fatalln("failed to binding auth.json to NewToken model:", err)
		}

		cl := &http.Client{}
		newToken, err := newTokenI.RenewToken(cl)
		if err != nil {
			log.Fatalln("failed to get new token:", err)
		}

		// 3. Delete old token.json file
		os.Remove(conf.Provider.Token)
		// 4. Save new token to token.json file
		if err := token.SaveToken(conf.Provider.Token, newToken); err != nil {
			log.Fatalln("failed to save new oauth2.Token instance to token.json file:", err)
		}
	}

	// START-
	var (
		rootIdFolder          = ""
		currentParentIdFolder = ""
	)

	// Note: make sure to empty trash first. Otherwise, root folder could never be able to created
	if err := dr.Files.EmptyTrash().Do(); err != nil {
		log.Println("failed to empty trash:", err)
	}

	// 1. Search RootFolder in MyDrive
	q := fmt.Sprintf("mimeType = '%s' and name = '%s'", MIMEFolder, conf.RootFolder)
	folder, err := dr.Files.List().Q(q).Do()
	if err != nil {
		log.Fatalf("failed to query for a root folder with a name %s: %s\n", conf.RootFolder, err)
	}
	if len(folder.Files) > 0 {
		rootIdFolder = folder.Files[0].Id
	}
	// 2. If not exist yet, then create it
	if len(folder.Files) <= 0 {
		fl := &drive.File{Name: conf.RootFolder, MimeType: MIMEFolder}
		newFl, err := dr.Files.Create(fl).Do()
		if err != nil {
			log.Fatalf("failed to create root folder: %s with error: %s\n", conf.RootFolder, err)
		}
		rootIdFolder = newFl.Id
	}

	// 3. Loop through all folders in upload section
	for _, up := range conf.Upload {
		up.Folders.Sanitization()
		currentParentIdFolder = rootIdFolder

		// 4. Loop through all folder tree in Folders.Name
		foldersToBeChecked := strings.Split(up.Folders.Name, "/")
		lastFolderTree := foldersToBeChecked[len(foldersToBeChecked)-1]
		for _, folderCheck := range foldersToBeChecked {

			// 5. Search this folder's name starting from RootFolder as parent folder
			q := fmt.Sprintf(
				"mimeType = '%s' and name = '%s' and '%s' in parents",
				MIMEFolder, folderCheck, currentParentIdFolder,
			)
			folder, err := dr.Files.List().Q(q).Do()
			if err != nil {
				log.Fatalf("failed to query for a folder with a name: %s under a parent folder's id: %s\n",
					folderCheck, currentParentIdFolder,
				)
			}
			// If found then change the current parent folder to this folder's id
			if len(folder.Files) > 0 {
				currentParentIdFolder = folder.Files[0].Id
			}

			// 6. If not found or not exist yet, then create it and change the current parent folder to this folder's id
			if len(folder.Files) <= 0 {
				fl := &drive.File{Name: folderCheck, MimeType: MIMEFolder,
					Parents: []string{currentParentIdFolder},
				}
				newFl, err := dr.Files.Create(fl).Do()
				if err != nil {
					log.Fatalln("failed to create a folder:", err)
				}
				currentParentIdFolder = newFl.Id
			}

			// 7. If we reach the last in folder tree then upload all files to this folder
			if folderCheck == lastFolderTree {
				allFiles, err := scan.Files(up.Folders.Path)
				if err != nil {
					log.Fatalf("failed to scan and read folder path: %s with error: %s\n", up.Folders.Path, err)
				}

				if len(allFiles) > 0 {
					var soonToBeDeletedFiles []string
					// 8. Before uploading. Remember to take notes all files (id) that reside in this folder that
					// fulfill the requirements to be deleted. like, maybe already pass the date. BUT do this only if
					// retain_days of this section is NOT 0.
					if up.Folders.Retain > 0 {
						query := fmt.Sprintf("'%s' in parents and %s != '%s'",
							currentParentIdFolder, FieldMIME, MIMEFolder,
						)
						list, _ := dr.Files.List().Q(query).Fields(
							"files(id)",
							"files(createdTime)",
							//"files(name)",
						).Do()
						for _, fl := range list.Files {
							t, err := time.Parse(time.RFC3339, fl.CreatedTime)
							if err != nil {
								log.Println("failed to parse time:", err)
							}
							sinceCreate := math.Round(time.Since(t).Hours())
							if uint(sinceCreate) > ((up.Folders.Retain * 24) - 1) {
								soonToBeDeletedFiles = append(soonToBeDeletedFiles, fl.Id)
							}
						}
					}

					// Then we are ready to upload the files to Google Drive's folder
					for _, fl := range allFiles {
						flInstance, err := os.Open(fl)
						if err != nil {
							log.Fatalf("failed to open file: %s with error: %s\n", fl, err)
						}
						defer flInstance.Close()
						fl := &drive.File{
							Parents: []string{currentParentIdFolder},
							Name:    filepath.Base(flInstance.Name()),
						}
						uploadFl, err := dr.Files.Create(fl).Media(flInstance).Fields(
							FieldId, FieldMIME,
							FieldName, FieldParents,
						).Do()
						if err != nil {
							log.Fatalf("failed to upload file: %s with error: %s\n", uploadFl.Name, err)
						}
					}

					// 9. Lastly, delete all soon to be deleted files using their id
					for _, filesToDelete := range soonToBeDeletedFiles {
						if err := dr.Files.Delete(filesToDelete).Do(); err != nil {
							log.Fatalf("failed to delete a file or a folder with id: %s and error: %s\n", filesToDelete, err)
						}
					}
				}
			}
		}
	}
	// END-
}
