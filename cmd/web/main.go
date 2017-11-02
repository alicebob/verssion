package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	libw "github.com/alicebob/w/w"
	"github.com/julienschmidt/httprouter"
)

var (
	baseURL = flag.String("base", "http://localhost:3141", "base URL")
	dbURL   = flag.String("db", "postgresql:///w", "postgres URL")
	listen  = flag.String("listen", ":3141", "http listen")
	updates = flag.Bool("update", true, "update pages")
)

func main() {
	flag.Parse()
	pages := flag.Args()
	if len(pages) != 0 {
		fmt.Fprintf(os.Stderr, "no args accepted\n")
		os.Exit(2)
	}

	db, err := libw.NewPostgres(*dbURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "pg: %s\n", err)
		os.Exit(2)
	}

	up := newUpdate(db)
	if !*updates {
		up = nil
	}

	r := httprouter.New()
	r.GET("/", indexHandler(db))
	r.GET("/adhoc/", adhocHandler(db))
	r.GET("/adhoc/atom.xml", adhocAtomHandler(db, up, *baseURL))
	r.GET("/curated/", newCuratedHandler(db, up))
	r.POST("/curated/", newCuratedHandler(db, up))
	r.GET("/curated/:id/", curatedHandler(db, *baseURL))
	r.GET("/curated/:id/edit.html", curatedEditHandler(db, up, *baseURL))
	r.POST("/curated/:id/edit.html", curatedEditHandler(db, up, *baseURL))
	r.GET("/curated/:id/atom.xml", curatedAtomHandler(db, up, *baseURL))
	r.GET("/p/:page/", pageHandler(db, up))
	fmt.Printf("listening on %s...\n", *listen)
	log.Fatal(http.ListenAndServe(*listen, r))
}

func indexHandler(db libw.DB) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		es, _ := db.Recent()
		runTmpl(w, indexTempl, map[string]interface{}{
			"title":   "hello world",
			"entries": es,
		})
	}
}

func adhocURL(pages []string) string {
	u, err := url.Parse(*baseURL)
	if err != nil {
		panic(err)
	}
	u.Path += "/adhoc/atom.xml"
	u.RawQuery = url.Values{
		"p": pages,
	}.Encode()
	return u.String()
}
