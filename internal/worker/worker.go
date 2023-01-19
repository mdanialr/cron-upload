package worker

import "github.com/mdanialr/cron-upload/internal/provider"

// I every worker implementation should follow the provided guideline/function
// signature here.
type I interface {
	// BuildRoutesProvider build routes from cloud provider's folder structure
	// along with the target folder id for each route.
	BuildRoutesProvider() map[string]string
	// CheckCreateRoute check whether the given folder name is already exist
	// inside the parent id otherwise create new one in the cloud provider
	// and use the given parent id as the parent and name as the folder name.
	// return the newly created folder id.
	CheckCreateRoute(parentId, name string) (string, error)
	// ListExpiredFiles list all files that's inside the given folder id in the
	// cloud provider then return unexpired payload to In channel and expired
	// payload within the given expiry duration (in minutes) to the Out
	// channel. Make sure to close both In & Out channels thereafter.
	ListExpiredFiles(channels Channels, folderId string, expiry uint)
	// ListUnmatchedFiles match the incoming payload from In chanel with the
	// given local file names, then feed the unmatched payload to Out chanel.
	// May use In chanel from ListExpiredFiles as feeder of In chanel. Make
	// sure to close Out chanel thereafter.
	ListUnmatchedFiles(channels Channels, folderId string, localFiles ...string)
	// DeleteFile delete file id that's fed from the in channel. May use Out
	// chanel from ListExpiredFiles as feeder of In chanel.
	DeleteFile(channels Channels)
	// UploadFile upload payload that's fed from the given channel. May use Out
	// chanel from ListUnmatchedFiles as feeder of In chanel.
	UploadFile(channels Channels)
}

// Channels contain the necessary channels that's passed between signature.
// Can be used by ListExpiredFiles, ListUnmatchedFiles, DeleteFile &
// UploadFile.
type Channels struct {
	In, Out chan *provider.Payload
}
