// get latest versions from wikipedia, and store when changed

package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	libw "github.com/alicebob/verssion/w"
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

	db, err := libw.NewPostgres(*dburl)
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

func handlePage(db libw.DB, page string) error {
	p, err := libw.GetPage(page)
	if err != nil {
		return err
	}

	sv := p.StableVersion
	if sv == "" {
		return fmt.Errorf("no stable version found")
	}

	return db.Store(libw.Page{
		Page:          page,
		T:             time.Now().UTC(),
		StableVersion: sv,
		Homepage:      p.Homepage,
	})
}
