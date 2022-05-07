package gdrive

const (
	MIMEShortcut     = "application/vnd.google-apps.shortcut" // MIMEType of the file for shortcut which unique only to Google Drive
	MIMEFile         = "application/vnd.google-apps.file"     // MIMEType of the file for file which unique only to Google Drive
	MIMEFolder       = "application/vnd.google-apps.folder"   // MIMEType of the file for folder which unique only to Google Drive
	MIMEDocs         = "application/vnd.google-apps.document" // MIMEType of the file for Google Docs which unique only to Google Drive
	MIMEZip          = "application/zip"                      // MIMEType of the file for zip archive format (.zip)
	MIMEPlainText    = "text/plain"                           // MIMEType of the file for plain text format (.txt)
	MIMERichText     = "application/rtf"                      // MIMEType of the file for rich text format (.rtf)
	MIMEJpeg         = "image/jpeg"                           // MIMEType of the file for image with jpeg format (.jpeg)
	MIMEPng          = "image/png"                            // MIMEType of the file for image with png format (.png)
	MIMESvg          = "image/svg+xml"                        // MIMEType of the file for image with svg format (.svg)
	FieldId          = "id"                                   // id of the file
	FieldName        = "name"                                 // name of the file
	FieldMIME        = "mimeType"                             // a field to define what mime type would be used
	FieldDescription = "description"                          // description of the file
	FieldParents     = "parents"                              // parents of the file, which is the folder's id where this file is located
)
