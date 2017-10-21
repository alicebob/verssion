package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/alicebob/w/w"
)

func main() {
	flag.Parse()
	for _, v := range flag.Args() {
		p, err := w.GetParseTree(v)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: %s\n", v, err)
			os.Exit(2)
		}
		fmt.Printf(p + "\n")
	}
}
