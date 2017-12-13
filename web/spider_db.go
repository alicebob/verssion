package web

import (
	"log"
	"time"

	"github.com/alicebob/verssion/core"
)

// StoreSpider tries to load a spider from the db
func StoreSpider(db core.DB, spider core.Spider, page string) (*core.Page, error) {
	{
		p, err := db.Last(page)
		if err != nil {
			return nil, err
		}
		// Recent enough version found in the db
		if p != nil && p.T.After(time.Now().Add(-core.CacheOK)) {
			return p, nil
		}
	}
	log.Printf("go fetch %q", page)
	p, err := spider.Spider(page)
	if err != nil {
		return nil, err
	}

	if p != nil {
		if err := db.Store(*p); err != nil {
			return nil, err
		}
	}

	return p, nil
}
