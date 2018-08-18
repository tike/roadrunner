package watch

import (
	"context"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

func fullCycle(pkg string, args []string, verbose bool) (context.CancelFunc, error) {
	if err := build(pkg); err != nil {
		return nil, err
	}
	return run(pkg, args, verbose), nil
}

func build(pkg string) error {
	cmd := exec.Command("go", "install", "-v", pkg)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	log.Printf("building %s", pkg)
	return cmd.Run()
}

func run(pkg string, args []string, verbose bool) context.CancelFunc {
	name := filepath.Base(pkg)
	ctx, cancel := context.WithCancel(context.Background())

	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	go func() {
		log.Printf("running %s", cmd.Args)
		if err := cmd.Run(); err != nil {
			log.Printf("%s terminated (%s)", cmd.Args, err)
		}
	}()

	return func() {
		if verbose {
			log.Printf("terminating %s", cmd.Args)
		}
		cancel()
	}
}
