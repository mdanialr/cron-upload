package scan

import (
	"fmt"
	"io/ioutil"
	"sort"
	"time"
)

// byDateAsc implements sort.Interface for []time.Time based on the creation time and ascending sorted.
type byDateAsc []sortedFile

func (a byDateAsc) Len() int           { return len(a) }
func (a byDateAsc) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byDateAsc) Less(i, j int) bool { return a[i].date.Before(a[j].date) }

// sortedFile custom struct for sorting purpose only.
type sortedFile struct {
	name string    // the name of the file.
	date time.Time // the latest modified time of the file.
}

// Files scan the given directory and return the list of the filename that already
// sorted by date in ascending direction.
func Files(dir string) (result []string, err error) {
	var sorted []sortedFile

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return []string{}, fmt.Errorf("failed to read the given dir: %s", err)
	}

	for _, fl := range files {
		// append only if the data is NOT a directory
		if !fl.IsDir() {
			sorted = append(sorted, sortedFile{
				fmt.Sprintf("%s/%s", dir, fl.Name()),
				fl.ModTime(),
			})
		}
	}

	sort.Sort(byDateAsc(sorted))
	for _, st := range sorted {
		result = append(result, st.name)
	}

	return result, nil
}
