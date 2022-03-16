package pcloud

const (
	Scheme                   = "https"
	EndPoint                 = "eapi.pcloud.com"         // Because registered user in Europe region. If in America use api.pcloud.com instead.
	GetDigest                = "getdigest"               // Returns a digest for digest authentication. Digests are valid for 30 seconds.
	Login                    = "userinfo"                // Generate token.
	LOGOUT                   = "logout"                  // Gets a token and invalidates it (Delete a token).
	TOKENS                   = "listtokens"              // Get a list with the currently active tokens associated with the current user.
	CREATE_FOLDER            = "createfolderifnotexists" // Creates a folder if the folder doesn't exist or returns the existing folder's metadata.
	LIST_FOLDERS             = "listfolder"              // Receive data for a folder.
	DELETE_FOLDER            = "deletefolder"            // Deletes a folder. The folder must be empty.
	DELETE_FOLDERS_RECURSIVE = "deletefolderrecursive"   // This function deletes files, directories, and removes sharing. Use with extreme care.
)
