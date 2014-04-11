/*
Package directory implements functionality for navigating
and listing directories, including size calculations.

Directory paths are always returned without a trailing slash.
*/
package directory

import (
	"io/ioutil"
	"os"
)

// Structure representing a directory entry.
type Entry struct {
	Name        string
	Size        int64
	IsDirectory bool
}

// Calculates and returns the size (in
// bytes) of the directory for the given path.
func Size(path string, result chan int64) {
	var size int64

	// Read the directory entries.
	entries, _ := ioutil.ReadDir(path)

	// Sum the entry sizes, recursing if necessary.
	for _, entry := range entries {
		if os.FileMode.IsDir(entry.Mode()) {
			recursiveResult := make(chan int64)
			go Size(path+"/"+entry.Name(), recursiveResult)
			size += <-recursiveResult
		} else {
			size += entry.Size()
		}
	}
	result <- size
}

// Returns a list of entries (and their sizes) for the given
// path. The current and parent (./..) entries are not included.
func Entries(path string) (entries []*Entry) {
	var size int64

	// Read the directory entries.
	dirEntries, _ := ioutil.ReadDir(path)
	entries = make([]*Entry, len(dirEntries))

	for index, entry := range dirEntries {
		entryInfo, _ := os.Stat(path + "/" + entry.Name())

		// Figure out the entry's size differently
		// depending on whether or not it's a directory.
		if entryInfo.IsDir() {
			result := make(chan int64)
			go Size(path+"/"+entry.Name(), result)
			size = <-result
		} else {
			size = entryInfo.Size()
		}

		entries[index] = &Entry{Name: entry.Name(), Size: size, IsDirectory: entryInfo.IsDir()}
	}
	return
}
