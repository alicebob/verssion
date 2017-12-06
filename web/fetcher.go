package web

import (
	"log"
	"time"

	"github.com/alicebob/verssion/core"
)

type Fetcher func(page string) (*core.Page, error)

// NotFetcher doesn't fetch a page. Use in tests.
func NotFetcher() Fetcher {
	return func(string) (*core.Page, error) {
		return nil, nil
	}
}

var _ Fetcher = NotFetcher()

// WikiFetcher loads from wikipedia
func WikiFetcher() Fetcher {
	up := NewUpdate()
	return func(page string) (*core.Page, error) {
		return up.Fetch(page, 10)
	}
}

var _ Fetcher = WikiFetcher()

// loadPage returns a the lastest from the DB if that's recent enough, or uses
// the fetcher to spider the page
func loadPage(page string, db core.DB, fetch Fetcher) (*core.Page, error) {
	{
		p, err := db.Last(page)
		if err != nil {
			return nil, err
		}
		// Recent enough version found in the db
		if p != nil && p.T.After(time.Now().Add(-cacheOK)) {
			return p, nil
		}
	}
	log.Printf("go fetch %q", page)
	p, err := fetch(page)
	if err != nil {
		return nil, err
	}
	// can happen with the NotFetcher
	if p == nil {
		return nil, core.ErrNotFound{Page: page}
	}

	if err := db.Store(*p); err != nil {
		return nil, err
	}

	return p, nil
}
