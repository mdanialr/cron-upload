package pcloud

import "net/url"

// DeleteFileResponse standard response from pCloud API after sending delete file request.
type DeleteFileResponse struct {
	Result int            `json:"result"`   // Should be non-0 value if there are any errors.
	Meta   DeleteFileMeta `json:"metadata"` // Metadata contain deleted file.
}

// DeleteFileMeta all metadata contain deleted file.
type DeleteFileMeta struct {
	IsDeleted bool   `json:"isdeleted"` // Status boolean that should be true if the file successfully deleted.
	Name      string `json:"name"`      // The name of the deleted file.
}

// GetDeleteFileUrl generate url that could be used to delete a file from the given fileId.
func GetDeleteFileUrl(a, fileId string) string {
	val := url.Values{"auth": {a}, "fileid": {fileId}}
	uri := url.URL{Scheme: Scheme, Host: EndPoint, Path: DeleteFile, RawQuery: val.Encode()}

	return uri.String()
}
