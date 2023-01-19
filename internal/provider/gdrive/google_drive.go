package gdrive

import (
	"fmt"

	"github.com/mdanialr/cron-upload/internal/provider"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/googleapi"
)

// NewGoogleDriveProvider return provider that use Google Drive as the cloud provider.
func NewGoogleDriveProvider(svc *drive.FilesService) provider.Cloud {
	return &googleDrive{
		svc: svc,
	}
}

type googleDrive struct {
	svc *drive.FilesService
}

func (g *googleDrive) GetFolders(parent ...string) ([]*provider.Payload, error) {
	query := fmt.Sprintf("%s = '%s'", FieldMIME, MIMEFolder)
	// add first parent as query if provided
	if len(parent) > 0 {
		query = fmt.Sprintf("%s and '%s' in parents", query, parent[0])
	}

	return g.queryList(query)
}

func (g *googleDrive) CreateFolder(name string, parent ...string) (string, error) {
	folder := &drive.File{
		Name:     name,
		MimeType: MIMEFolder,
	}
	// append the first parent if provided
	if len(parent) > 0 {
		folder.Parents = append(folder.Parents, parent[0])
	}
	cr, err := g.svc.Create(folder).Do()
	if err != nil {
		return "", fmt.Errorf("failed to create a folder with name '%s': %s", name, err)
	}

	return cr.Id, nil
}

func (g *googleDrive) GetFiles(folderId string) ([]*provider.Payload, error) {
	query := fmt.Sprintf("%s != '%s' and '%s' in parents", FieldMIME, MIMEFolder, folderId)
	return g.queryList(query)
}

func (g *googleDrive) UploadFile(payload *provider.Payload, chunkSize ...int) (*provider.Payload, error) {
	defer payload.File.Close()
	var uploadChunkSize int
	var resPayload provider.Payload

	// set chunk size if provided
	if len(chunkSize) > 0 {
		uploadChunkSize = chunkSize[0]
	}

	fl := drive.File{
		Name:    payload.Name,
		Parents: payload.Parent,
	}
	newFl, err := g.svc.Create(&fl).Media(payload.File, googleapi.ChunkSize(uploadChunkSize)).Do()
	if err != nil {
		return nil, fmt.Errorf("failed to upload file with name '%s': %s", payload.Name, err)
	}

	resPayload.Id = newFl.Id
	resPayload.Name = newFl.Name
	resPayload.Parent = payload.Parent

	return &resPayload, nil
}

func (g *googleDrive) Delete(id string) error {
	if err := g.svc.Delete(id).Do(); err != nil {
		return fmt.Errorf("failed to delete a file/folder with id '%s': %s", id, err)
	}
	return nil
}

// queryList helper to do List API call using the provided string as the query.
func (g *googleDrive) queryList(query string) ([]*provider.Payload, error) {
	ls, err := g.svc.List().Fields(
		toField("files/"+FieldId),
		toField("files/"+FieldName),
		toField("files/"+FieldParents),
		toField("files/"+FieldCreatedAt),
	).Q(query).Do()
	if err != nil {
		return nil, fmt.Errorf("failed to query for folder: %s", err)
	}

	if ls != nil {
		if len(ls.Files) > 0 {
			var res []*provider.Payload
			for _, fl := range ls.Files {
				res = append(res, &provider.Payload{
					Id:        fl.Id,
					Name:      fl.Name,
					Parent:    fl.Parents,
					CreatedAt: fl.CreatedTime,
				})
			}
			return res, nil
		}
	}

	return nil, nil
}
