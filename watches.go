package main

import (
	"log"

	"github.com/fsnotify/fsnotify"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func watch(files []string, cmd []string) (chan<- bool, error) {
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
	return doWatch(watcher, cmd), nil
}

func doWatch(watcher *fsnotify.Watcher, cmd []string) chan<- bool {
	done := make(chan bool)
	pkg, args := cmd[0], cmd[1:]
	go func() {
		defer watcher.Close()
		for {
			select {
			case event := <-watcher.Events:
				log.Println("event:", event)
				fullCycle(pkg, args...)
			case err := <-watcher.Errors:
				log.Println("error:", err)
			case <-done:
				return
			}
		}
	}()
	return done
}
