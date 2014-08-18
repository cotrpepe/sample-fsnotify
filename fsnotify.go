package main

import (
	"log"
	"github.com/go-fsnotify/fsnotify"
	"os"
	"path/filepath"
)

func main() {
	if len(os.Args) != 2 {
		log.Println("Usage: go run " + filepath.Base(os.Args[0]) + ".go directory_name")
		return
	}
	dirname := os.Args[1]

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan bool)

	go func() {
		for {
			select {
			case event:= <-watcher.Events:
				log.Println("event: ", event)
				switch {
				case event.Op&fsnotify.Write == fsnotify.Write:
					log.Println("Modified file: ", event.Name)
				case event.Op&fsnotify.Create == fsnotify.Create:
					log.Println("Created file: ", event.Name)
				case event.Op&fsnotify.Remove == fsnotify.Remove:
					log.Println("Removed file: ", event.Name)
				case event.Op&fsnotify.Rename == fsnotify.Rename:
					log.Println("Renamed file: ", event.Name)
				case event.Op&fsnotify.Chmod == fsnotify.Chmod:
					log.Println("File changed permission: ", event.Name)
				}
			case err:= <-watcher.Errors:
				log.Println("error: ", err)
				done <-true
			}
		}
	}()

	err = watcher.Add(dirname)
	if err != nil {
		log.Fatal(err)
	}

	<-done
}
