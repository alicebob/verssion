// get latest versions from wikipedia

package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/alicebob/w/w"
)

var dburl = flag.String("db", "postgresql:///w", "postgres URL")

func main() {
	flag.Parse()
	pages := flag.Args()
	if len(pages) == 0 {
		fmt.Fprintf(os.Stderr, "no pages given\n")
		os.Exit(2)
	}

	fmt.Printf("there we go...\n")

	db, err := w.NewPostgres(*dburl)
	if err != nil {
		fmt.Fprintf(os.Stderr, "pg: %s\n", err)
		os.Exit(2)
	}

	revs, err := w.Revisions(pages)
	if err != nil {
		fmt.Fprintf(os.Stderr, "wikipedia revisions: %s\n", err)
		os.Exit(2)
	}

	for page, rev := range revs {
		old, err := db.Load(page)
		if err != nil {
			fmt.Fprintf(os.Stderr, "db load %q: %s\n", page, err)
			continue
		}
		if old != nil && rev.RevID == old.Revision {
			fmt.Printf("%q: no update (%d/%s)\n", page, rev.RevID, rev.T)
			continue
		}
		new, err := w.GetPage(rev.RevID)
		if err != nil {
			fmt.Fprintf(os.Stderr, "wikipedia %q (%d): %s\n", page, rev.RevID, err)
			continue
		}
		e := w.Entry{
			Page:          page,
			Revision:      rev.RevID,
			T:             rev.T,
			StableVersion: new.StableVersion,
		}
		if err := db.Store(e); err != nil {
			fmt.Fprintf(os.Stderr, "db store %q: %s\n", page, err)
			continue
		}
		// Only a change in wikipedia change, might be no version update
		sv := ""
		if old != nil {
			sv = old.StableVersion
		}
		fmt.Printf("stored %q: %s->%s\n", page, sv, new.StableVersion)
	}
}
