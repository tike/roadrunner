package main

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

func fullCycle(pkg string, args ...string) {
	log.Printf("start building %s", pkg)
	if err := build(pkg); err != nil {
		log.Printf("building %s: %v", pkg, err)
	}
	log.Printf("done building %s, running %s %v", pkg, pkg, args)
	if err := run(pkg, args...); err != nil {
		log.Printf("running %s: %v", pkg, err)
	}
}

func build(pkg string) error {
	return do("go", "install", "-v", pkg)
}

func run(pkg string, args ...string) error {
	name := filepath.Base(pkg)
	return do(name, args...)
}

func do(cmd string, args ...string) error {
	e := exec.Command(cmd, args...)
	e.Stdout = os.Stdout
	e.Stderr = os.Stderr
	return e.Run()
}
