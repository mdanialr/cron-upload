package scan

import (
	"fmt"
	"os"
	"sort"
	"time"
)

// byDateAsc implements sort.Interface for []time.Time based on the creation time and ascending sorted.
type byDateAsc []sortedFile

func (a byDateAsc) Len() int           { return len(a) }
func (a byDateAsc) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byDateAsc) Less(i, j int) bool { return a[i].Date.Before(a[j].Date) }

// sortedFile custom struct for sorting purpose only.
type sortedFile struct {
	Name string    // Name the name of the file.
	Date time.Time // the latest modified time of the file.
}

// FilesAsc scan the given directory and return the list of the filename that already
// sorted by date in ascending direction.
func FilesAsc(dir string) (result []string, err error) {
	var sorted []sortedFile

	files, err := os.ReadDir(dir)
	if err != nil {
		return []string{}, fmt.Errorf("failed to read the given dir: %s", err)
	}

	for _, fl := range files {
		// append only if the data is NOT a directory
		if !fl.IsDir() {
			info, _ := fl.Info()
			sorted = append(sorted, sortedFile{
				Name: fmt.Sprintf("%s/%s", dir, fl.Name()),
				Date: info.ModTime(),
			})
		}
	}

	sort.Sort(byDateAsc(sorted))
	for _, st := range sorted {
		result = append(result, st.Name)
	}

	return result, nil
}
