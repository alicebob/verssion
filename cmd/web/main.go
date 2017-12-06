package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/alicebob/verssion/core"
	"github.com/alicebob/verssion/web"
)

var (
	baseURL = flag.String("base", "http://localhost:3141", "base URL")
	dbURL   = flag.String("db", "postgresql:///w", "postgres URL")
	listen  = flag.String("listen", ":3141", "http listen")
)

func main() {
	flag.Parse()
	pages := flag.Args()
	if len(pages) != 0 {
		fmt.Fprintf(os.Stderr, "no args accepted\n")
		os.Exit(2)
	}

	db, err := core.NewPostgres(*dbURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "pg: %s\n", err)
		os.Exit(2)
	}

	mux := web.Mux(*baseURL, db, web.WikiFetcher())

	fmt.Printf("listening on %s...\n", *listen)
	log.Fatal(http.ListenAndServe(*listen, mux))
}
