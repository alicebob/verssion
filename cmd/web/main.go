package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/alicebob/w/w"
	"github.com/julienschmidt/httprouter"
)

var (
	dburl  = flag.String("db", "postgresql:///w", "postgres URL")
	listen = flag.String("listen", ":3141", "http listen")
)

func main() {
	flag.Parse()
	pages := flag.Args()
	if len(pages) != 0 {
		fmt.Fprintf(os.Stderr, "no args accepted\n")
		os.Exit(2)
	}

	db, err := w.NewPostgres(*dburl)
	if err != nil {
		fmt.Fprintf(os.Stderr, "pg: %s\n", err)
		os.Exit(2)
	}

	r := httprouter.New()
	r.GET("/", indexHandler(db))
	fmt.Printf("listening on %s...\n", *listen)
	log.Fatal(http.ListenAndServe(*listen, r))
}

func indexHandler(db w.DB) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		fmt.Fprintf(w, "hello world<br />\n")
		es, _ := db.Recent()
		for _, e := range es {
			fmt.Fprintf(w, "%s: %s<br />\n", e.Page, e.StableVersion)
		}
	}
}
