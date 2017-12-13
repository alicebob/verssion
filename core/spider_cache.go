package core

import (
	"fmt"
	"sync"
	"time"
)

const (
	cacheErr = 30 * time.Second
	CacheOK  = 6 * time.Hour
)

// SpiderCache wraps a Spider implementation
type SpiderCache struct {
	mu     sync.Mutex
	spider Spider
	pages  map[string]*last
}

var _ Spider = &SpiderCache{}

type last struct {
	mu        sync.Mutex
	cacheTill time.Time
	page      *Page
	err       error
}

func NewSpiderCache(spider Spider) *SpiderCache {
	return &SpiderCache{
		spider: spider,
		pages:  map[string]*last{},
	}
}

func (u *SpiderCache) cachedFetch(page string) (*Page, error) {
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
	l.page, l.err = u.spider.Spider(page)
	c := CacheOK
	if l.err != nil {
		c = cacheErr
	}
	l.cacheTill = now.Add(c)
	return l.page, l.err
}

// Fetch the most recent version (or a cache).
// Follows redirects.
func (u *SpiderCache) fetch(page string, redir int) (*Page, error) {
	if redir < 0 {
		return nil, fmt.Errorf("%q: too many redirects", page)
	}

	p, err := u.cachedFetch(page)
	if err == nil {
		return p, nil
	}
	if red, ok := err.(ErrRedirect); ok {
		return u.fetch(red.To, redir-1)
	}
	return nil, err
}

func (u *SpiderCache) Spider(page string) (*Page, error) {
	return u.fetch(page, 10)
}
