package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alicebob/verssion/core"
	"github.com/alicebob/verssion/web"
)

var (
	baseURL = flag.String("base", "http://localhost:3141", "base URL")
	dbURL   = flag.String("db", "postgresql:///verssion", "postgres URL")
	listen  = flag.String("listen", ":3141", "http listen")
	static  = flag.String("static", "", "subdir with static files")
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

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

	go web.CacheLoop(db)
	go func() {
		for {
			if err := db.UpdateViews(ctx); err != nil {
				fmt.Printf("update views failed: %s\n", err)
			}
			select {
			case <-ctx.Done():
				break
			case <-time.After(7 * time.Minute):
			}
		}
	}()

	mux := web.Mux(*baseURL, db, web.WikiSpider(), *static)

	fmt.Printf("listening on %s...\n", *listen)
	log.Fatal(http.ListenAndServe(*listen, mux))
}
