package internal

import (
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

// Watch one or more files, but instead of watching the file directly it watches
// the parent directory. This solves various issues where files are frequently
// renamed, such as editors saving them.
func WatchFiles(patterns ...string) {
	if len(patterns) < 1 {
		panic("must specify at least one file to watch")
	}

	// Create a new watcher.
	w, err := fsnotify.NewWatcher()
	if err != nil {
		panic(fmt.Sprintf("creating a new watcher: %s", err))
	}
	defer w.Close()

	// Start listening for events.
	go fileLoop(w, patterns)

	files := make([]string, 0, 10)
	filepath.WalkDir(".", func(path string, d fs.DirEntry, err error) error {

		for _, pattern := range patterns {
			if !d.IsDir() && matches(path, pattern) {
				files = append(files, path)
			}
		}

		return nil
	})

	fmt.Printf("Files to watch: %+v", files)

	// Add all files from the commandline.
	for _, p := range files {
		st, err := os.Lstat(p)
		if err != nil {
			panic(fmt.Sprintf("%s", err))
		}

		if st.IsDir() {
			panic(fmt.Sprintf("%q is a directory, not a file", p))
		}

		// Watch the directory, not the file itself.
		err = w.Add(filepath.Dir(p))
		if err != nil {
			panic(fmt.Sprintf("%q: %s", p, err))
		}
	}

	fmt.Println("ready; press ^C to panic")
	<-make(chan struct{}) // Block forever
}

func matches(changedFileName string, pattern string) bool {
	match, err := path.Match(pattern, changedFileName)
	fmt.Printf("%s matches pattern %s: %t\n", changedFileName, pattern, match)
	if err != nil {
		return false
	}

	return match
}

func fileLoop(w *fsnotify.Watcher, files []string) {
	i := 0
	for {
		select {
		// Read from Errors.
		case err, ok := <-w.Errors:
			if !ok { // Channel was closed (i.e. Watcher.Close() was called).
				return
			}
			fmt.Printf("ERROR: %s", err)
		// Read from Events.
		case e, ok := <-w.Events:
			if !ok { // Channel was closed (i.e. Watcher.Close() was called).
				return
			}

			// Ignore files we're not interested in. Can use a
			// map[string]struct{} if you have a lot of files, but for just a
			// few files simply looping over a slice is faster.
			var found bool
			for _, f := range files {
				if matches(f, e.Name) {
					found = true
					break
				}
			}
			if !found {
				continue
			}

			// Just print the event nicely aligned, and keep track how many
			// events we've seen.
			i++
			fmt.Printf("%d %v", i, e)
		}
	}
}
