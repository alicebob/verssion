// get latest versions from wikipedia, and store when changed

package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/alicebob/w/w"
)

var (
	dburl = flag.String("db", "postgresql:///w", "postgres URL")
	sleep = flag.Duration("sleep", 4*time.Second, "sleep between pages")
)

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

	for _, page := range pages {
		if err := handlePage(db, page); err != nil {
			fmt.Fprintf(os.Stderr, "%q: %s", page, err)
			// TODO: retry?
			time.Sleep(1 * time.Second)
			continue
		}
		time.Sleep(*sleep)
	}
}

func handlePage(db w.DB, page string) error {
	old, err := db.Load(page)
	if err != nil {
		return err
	}
	new, err := w.GetPage(page)
	if err != nil {
		return err
	}

	if old != nil && new.StableVersion == old.StableVersion {
		fmt.Printf("no update: %q (%s)\n", page, new.StableVersion)
		return nil
	}
	e := w.Entry{
		Page:          page,
		T:             time.Now().UTC(),
		StableVersion: new.StableVersion,
	}
	if err := db.Store(e); err != nil {
		return err
	}
	// Only a change in wikipedia change, might be no version update
	sv := ""
	if old != nil {
		sv = old.StableVersion
	}
	fmt.Printf("stored %q: %s-> %s\n", page, sv, e.StableVersion)
	return nil
}
