package worker

import (
	"path/filepath"
	"sync"
	"time"

	pv "github.com/mdanialr/cron-upload/internal/provider"
	w "github.com/mdanialr/cron-upload/internal/worker"
	h "github.com/mdanialr/cron-upload/pkg/helper"
	"github.com/mdanialr/cron-upload/pkg/logger"
)

// NewWorker return new first version worker.
func NewWorker(g *sync.WaitGroup, log logger.Writer, cloud pv.Cloud) w.I {
	return &worker{
		g:     g,
		log:   log,
		cloud: cloud,
	}
}

// worker the first version of worker implementation, in case there is a new
// way to implement how this app should work in the future.
type worker struct {
	g     *sync.WaitGroup
	log   logger.Writer
	cloud pv.Cloud
}

func (w *worker) BuildRoutesProvider() map[string]string {
	var result = make(map[string]string)
	payloads, err := w.cloud.GetFolders()
	if err != nil {
		w.log.WriteErr(err)
		return result
	}
	for _, payload := range payloads {
		route := pv.LookupRoute(payloads, payload.Name)
		result[route] = payload.Id
	}
	return result
}

func (w *worker) CheckCreateRoute(parentId, name string) (string, error) {
	// check whether the given folder name is already exist in folder id
	payloads, err := w.cloud.GetFolders(parentId)
	if err != nil {
		return "", err
	}
	// return the id if already exists in the cloud provider
	for _, payload := range payloads {
		if payload.Name == name {
			return payload.Id, nil
		}
	}
	// then create new folder
	id, err := w.cloud.CreateFolder(name, parentId)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (w *worker) ListExpiredFiles(channels w.Channels, folderId string, expiry uint) {
	defer close(channels.In)
	defer close(channels.Out)

	payloads, err := w.cloud.GetFiles(folderId)
	if err != nil {
		w.log.WriteErr("failed to get files in", folderId, "from cloud provider:", err)
		return
	}
	// just skip if there is no files
	if len(payloads) < 1 {
		return
	}
	// check files expiry
	for _, payload := range payloads {
		fmtTime, _ := time.Parse(time.RFC3339, payload.CreatedAt)
		if h.ToWib(time.Now()).After(h.ToWib(fmtTime).Add(time.Minute * time.Duration(expiry))) {
			// feed expired files
			channels.Out <- payload
			continue
		}
		// feed unexpired files
		channels.In <- payload
	}
}

func (w *worker) ListUnmatchedFiles(channels w.Channels, folderId string, localFiles ...string) {
	defer close(channels.Out)
	defer w.g.Done()

	// wait & grab all files in the cloud provider
	var fileInCloud []string
	for payload := range channels.In {
		fileInCloud = append(fileInCloud, payload.Name)
	}
	// match the local files against files in the cloud
	for _, localFile := range localFiles {
		var match bool
		for _, cloudFile := range fileInCloud {
			// don't upload local files that's already exist in the cloud provider
			// based on the file name
			if filepath.Base(localFile) == cloudFile {
				match = true
			}
		}
		if !match {
			// if unmatched then feed it to Out channel after successfully opening the file
			pay, err := pv.NewWithFile(localFile, folderId)
			if err != nil {
				w.log.WriteErr(err)
				continue
			}
			channels.Out <- pay
		}
	}
}

func (w *worker) DeleteFile(channels w.Channels) {
	defer w.g.Done()

	for payload := range channels.In {
		w.log.WriteInf(h.LogStart("DELETE", payload.Name))
		// delete for every incoming file id
		if err := w.cloud.Delete(payload.Id); err != nil {
			w.log.WriteErr(err, "with filename:", payload.Name)
		}
		w.log.WriteInf(h.LogDone("DELETE", payload.Name))
	}
}

func (w *worker) UploadFile(channels w.Channels) {
	defer w.g.Done()

	for payload := range channels.In {
		w.log.WriteInf(h.LogStart("UPLOAD", payload.Name))
		// upload for every incoming file reader
		if _, err := w.cloud.UploadFile(payload); err != nil {
			w.log.WriteErr(err)
		}
		w.log.WriteInf(h.LogDone("UPLOAD", payload.Name))
	}
}
