package main

import (
	"log"

	"github.com/fsnotify/fsnotify"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func watch(files ...string) (chan<- bool, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		err = watcher.Add(file)
		if err != nil {
			log.Println(file)
			return nil, err
		}
	}
	return doWatch(watcher), nil
}

func doWatch(watcher *fsnotify.Watcher) chan<- bool {
	done := make(chan bool)
	go func() {
		defer watcher.Close()
		for {
			select {
			case event := <-watcher.Events:
				log.Println("event:", event)
				if event.Op&fsnotify.Write == fsnotify.Write {
					log.Println("modified file:", event.Name)
				}
			case err := <-watcher.Errors:
				log.Println("error:", err)
			case <-done:
				return
			}
		}
	}()
	return done
}
