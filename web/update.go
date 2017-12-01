package web

import (
	"fmt"
	"sync"
	"time"

	libw "github.com/alicebob/verssion/w"
)

const (
	cacheErr = 30 * time.Second
	cacheOK  = 6 * time.Hour
)

type Update struct {
	mu    sync.Mutex
	pages map[string]*last
}

type last struct {
	mu        sync.Mutex
	cacheTill time.Time
	page      libw.Page
	err       error
}

func NewUpdate() *Update {
	return &Update{
		pages: map[string]*last{},
	}
}

func (u *Update) fetch(page string) (libw.Page, error) {
	return libw.GetPage(page, WikiURL(page))
}

func (u *Update) cachedFetch(page string) (libw.Page, error) {
	u.mu.Lock()
	l, ok := u.pages[page]
	if !ok {
		l = &last{}
		u.pages[page] = l
	}
	u.mu.Unlock()

	l.mu.Lock()
	defer l.mu.Unlock()
	now := time.Now()
	if !l.cacheTill.IsZero() && now.Before(l.cacheTill) {
		return l.page, l.err
	}
	l.page, l.err = u.fetch(page)
	c := cacheOK
	if l.err != nil {
		c = cacheErr
	}
	l.cacheTill = now.Add(c)
	return l.page, l.err
}

// Fetch the most recent version (or a cache).
// Follows redirects.
func (u *Update) Fetch(page string, redir int) (*libw.Page, error) {
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

func WikiURL(page string) string {
	return "https://en.wikipedia.org/wiki/" + page
}
