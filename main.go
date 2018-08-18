package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/tike/roadrunner/internal/scan"
)

// flags
var (
	verbose bool
)

func main() {
	flagFoo()

	fmt.Println("package to watch:", flag.Arg(0))
	pkgs, err := scan.GetPkgs(flag.Arg(0))
	if err != nil {
		fmt.Println("couldn't find package", err)
		os.Exit(1)
	}
	if verbose {
		printPackageList(pkgs)
	}
	_, err = watch(pkgs.FileList(), flag.Args())
	if err != nil {
		fmt.Println("error setting up file watcher:", err)
	}
	time.Sleep(1 * time.Minute)
}

func printPackageList(packs scan.PkgPack) {
	fmt.Println("watching packages:")
	f := "%-20s %-60s %s\n"
	fmt.Printf(f, "name", "import", "src")
	for _, imp := range packs.SortedNames() {
		pkg := packs[imp]
		fmt.Printf(f, pkg.Name, imp, pkg.ImportPath)
	}
}

func flagFoo() {
	flag.Usage = Usage
	flag.BoolVar(&verbose, "v", true, "verbose output")
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
