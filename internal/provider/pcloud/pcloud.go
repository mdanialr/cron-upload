package pcloud

const (
	Scheme                   = "https"
	EndPoint                 = "eapi.pcloud.com"         // Because registered user in Europe region. If in America use api.pcloud.com instead.
	GetDigest                = "getdigest"               // Returns a digest for digest authentication. Digests are valid for 30 seconds.
	Login                    = "userinfo"                // Generate token.
	Logout                   = "logout"                  // Gets a token and invalidates it (Delete a token).
	TOKENS                   = "listtokens"              // Get a list with the currently active tokens associated with the current user.
	CreateFolder             = "createfolderifnotexists" // Creates a folder if the folder doesn't exist or returns the existing folder's metadata.
	LIST_FOLDERS             = "listfolder"              // Receive data for a folder.
	DELETE_FOLDER            = "deletefolder"            // Deletes a folder. The folder must be empty.
	DELETE_FOLDERS_RECURSIVE = "deletefolderrecursive"   // This function deletes files, directories, and removes sharing. Use with extreme care.
)

// StdResponse standard response from pCloud API that always return 'return' that determine
// whether API call is success or failure. Response with non '0' return value are errors.
type StdResponse struct {
	Result               int `json:"result"` // non 0 result is errors.
	DigestResponse           // Used by digest.
	TokenResponse            // Used by login API calls to get token for authentication.
	LogoutResponse           // Used by logout API calls to invalidate token.
	QuotaResponse            // Used by print storage quota.
	CreateFolderResponse     // Used by create folder.
}
