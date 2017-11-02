package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/alicebob/w/w"
)

const (
	cacheErr = 30 * time.Second
	cacheOK  = 1 * time.Hour
)

type update struct {
	db    w.DB
	mu    sync.Mutex
	pages map[string]*last
}

type last struct {
	mu       sync.Mutex
	waitTill time.Time
}

func newUpdate(db w.DB) *update {
	return &update{
		db:    db,
		pages: map[string]*last{},
	}
}

func (u *update) Update(page string) error {
	u.mu.Lock()
	l, ok := u.pages[page]
	if !ok {
		l = &last{}
		u.pages[page] = l
	}
	u.mu.Unlock()

	now := time.Now()

	l.mu.Lock()
	defer l.mu.Unlock()
	if now.Before(l.waitTill) {
		log.Printf("cached %q...", page)
		return nil
	}

	log.Printf("updating %q...", page)
	l.waitTill = now.Add(cacheErr)
	p, err := w.GetPage(page)
	if err != nil {
		return err
	}

	sv := p.StableVersion
	if sv == "" {
		return fmt.Errorf("no version number found")
	}
	if err := u.db.Store(w.Page{
		Page:          page,
		T:             time.Now().UTC(),
		StableVersion: sv,
		Homepage:      p.Homepage,
	}); err != nil {
		return err
	}

	l.waitTill = now.Add(cacheOK)
	return nil
}
