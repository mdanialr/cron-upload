package gdrive

const (
	MIMEShortcut     = "application/vnd.google-apps.shortcut" // MIMEShortcut MIMEType of the file for shortcut which unique only to Google Drive
	MIMEFile         = "application/vnd.google-apps.file"     // MIMEFile MIMEType of the file for file which unique only to Google Drive
	MIMEFolder       = "application/vnd.google-apps.folder"   // MIMEFolder MIMEType of the file for folder which unique only to Google Drive
	MIMEDocs         = "application/vnd.google-apps.document" // MIMEDocs MIMEType of the file for Google Docs which unique only to Google Drive
	MIMEZip          = "application/zip"                      // MIMEZip MIMEType of the file for zip archive format (.zip)
	MIMEPlainText    = "text/plain"                           // MIMEPlainText MIMEType of the file for plain text format (.txt)
	MIMERichText     = "application/rtf"                      // MIMERichText MIMEType of the file for rich text format (.rtf)
	MIMEJpeg         = "image/jpeg"                           // MIMEJpeg MIMEType of the file for image with jpeg format (.jpeg)
	MIMEPng          = "image/png"                            // MIMEPng MIMEType of the file for image with png format (.png)
	MIMESvg          = "image/svg+xml"                        // MIMESvg MIMEType of the file for image with svg format (.svg)
	FieldId          = "id"                                   // FieldId id of the file
	FieldName        = "name"                                 // FieldName name of the file
	FieldMIME        = "mimeType"                             // FieldMIME define what mime type is used
	FieldDescription = "description"                          // FieldDescription description of the file
	FieldParents     = "parents"                              // FieldParents parents of the file, which is the folder's id where this file is located
	FieldCreatedAt   = "createdTime"                          // FieldCreatedAt the time at which the file was created
)
