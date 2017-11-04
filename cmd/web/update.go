package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	libw "github.com/alicebob/w/w"
)

const (
	cacheErr = 30 * time.Second
	cacheOK  = 6 * time.Hour
)

type update struct {
	db    libw.DB
	mu    sync.Mutex
	pages map[string]*last
}

type last struct {
	mu        sync.Mutex
	cacheTill time.Time
	page      libw.Page
	err       error
}

func newUpdate(db libw.DB) *update {
	return &update{
		db:    db,
		pages: map[string]*last{},
	}
}

func (u *update) fetch(page string) (libw.Page, error) {
	p, err := libw.GetPage(page)
	if err != nil {
		return p, err
	}
	return p, u.db.Store(p)

}

func (u *update) cachedFetch(page string) (libw.Page, error) {
	u.mu.Lock()
	l, ok := u.pages[page]
	if !ok {
		l = &last{}
		// load last spider from the DB
		p, err := u.db.Last(page)
		if err != nil {
			log.Printf("last %q: %s", page, err)
		} else {
			l.page = *p
			l.cacheTill = p.T.Add(cacheOK)
		}
		u.pages[page] = l
	}
	u.mu.Unlock()

	now := time.Now().UTC()
	if now.Before(l.cacheTill) {
		log.Printf("cached %q...", page)
		return l.page, l.err
	}

	l.page, l.err = u.fetch(page)
	c := cacheOK
	if l.err != nil {
		if _, ok := l.err.(libw.ErrRedirect); !ok {
			c = cacheErr
		}
	}
	l.cacheTill = now.Add(c)
	return l.page, l.err
}

// Fetch the most recent version (or a cache).
// Follows redirects.
func (u *update) Fetch(page string, redir int) (*libw.Page, error) {
	if redir < 0 {
		return nil, fmt.Errorf("%q: too many redirects", page)
	}

	p, err := u.cachedFetch(page)
	if err == nil {
		return &p, nil
	}
	if red, ok := err.(libw.ErrRedirect); ok {
		return u.Fetch(red.To, redir-1)
	}
	return nil, err
}
