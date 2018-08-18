package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/tike/roadrunner/internal/scan"
	"github.com/tike/roadrunner/internal/watch"
)

// flags
var (
	verbose   bool
	threshold time.Duration
)

func init() {
	log.SetFlags(log.Lshortfile | log.Ltime)
	log.SetPrefix(" [roadrunner] ")
}

func main() {
	flagFoo()

	log.Printf("\033[32mpackage to watch: %s\033[0m", flag.Arg(0))
	pkgs, err := scan.GetPkgs(flag.Arg(0))
	if err != nil {
		log.Printf("\033[31mcouldn't find package %s\033[0m", err)
		os.Exit(1)
	}
	if verbose {
		printPackageList(pkgs)
	}

	done, restart, err := watch.Watch(pkgs.FileList(), flag.Args(), threshold, verbose)
	if err != nil {
		log.Printf("\033[31merror setting up file watcher: %s\033[0m", err)
	}

	b := make([]byte, 1)
	for {
		if _, err := os.Stdin.Read(b); err != nil {
			if err == io.EOF {
				log.Println("\033[32mrestarting\033[0m")
				restart <- struct{}{}
				continue
			}
			log.Printf("\033[31merror reading STDIN: %s\033m", err)
		}
	}
	done <- struct{}{}
}

func printPackageList(packs scan.PkgPack) {
	log.Println("watching packages:")
	f := "%-20s %-60s %s\n"
	log.Printf(f, "name", "import", "src")
	for _, imp := range packs.SortedNames() {
		pkg := packs[imp]
		log.Printf(f, pkg.Name, imp, pkg.ImportPath)
	}
}

func flagFoo() {
	flag.Usage = Usage
	flag.BoolVar(&verbose, "v", false, "verbose output")
	flag.DurationVar(&threshold, "t", 10*time.Second, "only proceess one event per time interval (others will be dropped)")
	flag.Parse()
	if flag.NArg() < 1 {
		flag.Usage()
		os.Exit(127)
	}
}

func Usage() {
	fmt.Println("roadrunner path/to/pkg")
	flag.PrintDefaults()
}
