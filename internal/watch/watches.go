package watch

import (
	"log"
	"path/filepath"
	"time"

	"github.com/fsnotify/fsnotify"
)

func Watch(files []string, cmd []string, threshold time.Duration, verbose bool) (chan<- struct{}, chan<- struct{}, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, nil, err
	}

	for _, file := range files {
		err = watcher.Add(file)
		if err != nil {
			log.Printf("\033[31m adding %s: %s\033[0m", file, err)
			return nil, nil, err
		}
	}
	done, restart := startWatch(watcher, cmd, threshold, verbose)
	return done, restart, nil
}

func startWatch(watcher *fsnotify.Watcher, cmd []string, threshold time.Duration, verbose bool) (chan<- struct{}, chan<- struct{}) {
	done := make(chan struct{})
	restart := make(chan struct{})

	pkg, args := cmd[0], cmd[1:]
	go doWatch(watcher, pkg, args, done, restart, threshold, verbose)

	return done, restart
}

func doWatch(watcher *fsnotify.Watcher, pkg string, args []string, done, restart <-chan struct{}, threshold time.Duration, verbose bool) {
	defer watcher.Close()
	filteredEvents := dropRedundant(threshold, watcher.Events, verbose)

	cancel := initiateRun(pkg, args, verbose)
	for {
		select {
		case <-restart:
			cancel()
			cancel = initiateRun(pkg, args, verbose)
		case event := <-filteredEvents:
			if verbose {
				log.Printf("event: %s", event)
			}
			if filepath.Ext(event.Name) != ".go" {
				break
			}
			switch event.Op {
			case fsnotify.Create, fsnotify.Write, fsnotify.Rename, fsnotify.Remove:
				cancel()
				log.Printf("event: \033[32m%s\033[0m", event)
				cancel = initiateRun(pkg, args, verbose)
			}
		case err := <-watcher.Errors:
			log.Printf("\033[31mwatch error: %s\033[0m", err)
		case <-done:
			cancel()
			return
		}
	}
}

func dropRedundant(threshold time.Duration, events <-chan fsnotify.Event, verbose bool) <-chan fsnotify.Event {
	filtered := make(chan fsnotify.Event)
	go func() {
		thres := time.NewTicker(threshold)
		defer thres.Stop()

		for event := range events {
			select {
			case <-thres.C:
				filtered <- event
				if verbose {
					log.Printf("passed: %s", event)
				}
			default:
				if verbose {
					log.Printf("filtered: %s", event)
				}
				continue
			}
		}
		close(filtered)
	}()

	return filtered
}

func initiateRun(pkg string, args []string, verbose bool) func() {
	cancel, err := fullCycle(pkg, args, verbose)
	if err != nil {
		log.Printf("\033[31mstarting %s %s: %s\033[0m", pkg, args, err)
		return func() {}
	}
	log.Println("\033[33mhit CTRL+D to restart the inner process\033[0m")
	log.Println("\033[31mhit CTRL+C to exit\033[0m")
	return cancel
}
