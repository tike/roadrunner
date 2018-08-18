package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/tike/roadrunner/internal/scan"
)

func main() {
	flag.Usage = Usage
	flag.Parse()
	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(127)
	}
	fmt.Println("package to watch:", flag.Arg(0))
	packs, err := scan.GetPkgs(flag.Arg(0))
	if err != nil {
		fmt.Println("couldn't find package", err)
		os.Exit(1)
	}
	fmt.Println("watching dependecies:")
	f := "%-20s %-60s %s\n"
	fmt.Printf(f, "name", "import", "src")
	for _, imp := range packs.SortedNames() {
		pkg := packs[imp]
		fmt.Printf(f, pkg.Name, imp, pkg.ImportPath)
	}
}

func Usage() {
	fmt.Println("roadrunner path/to/pkg")
	flag.PrintDefaults()
}
