package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	flag.Usage = Usage
	flag.Parse()
	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(127)
	}
	fmt.Println("package to watch:", flag.Arg(0))
}

func Usage() {
	fmt.Println("roadrunner path/to/pkg")
	flag.PrintDefaults()
}
