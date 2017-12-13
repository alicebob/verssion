package core

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

const (
	DefaultUserAgent = "verssion_bot/1.0 (https://verssion.one)"
)

var client = http.Client{
	CheckRedirect: func(*http.Request, []*http.Request) error {
		return http.ErrUseLastResponse
	},
}

type WikipediaSpider struct {
	URL       func(string) string
	UserAgent string
}

var _ Spider = &WikipediaSpider{}

// NewWikipediaSpider makes a spider. url makes a URL from a page.
func NewWikipediaSpider(url func(string) string) *WikipediaSpider {
	return &WikipediaSpider{
		URL:       url,
		UserAgent: DefaultUserAgent,
	}
}

// Spider downloads and parses given wikipage
func (spider *WikipediaSpider) Spider(page string) (*Page, error) {
	p := Page{
		Page: page,
		T:    time.Now().UTC(),
	}

	log.Printf("wiki fetch %q", page)
	req, err := http.NewRequest("GET", spider.URL(page), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", spider.UserAgent)

	// no redirects
	r, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	switch code := r.StatusCode; code {
	case 200:
		p.StableVersion, p.Homepage = StableVersion(r.Body)
		if p.StableVersion == "" {
			return nil, ErrNoVersion{Page: page}
		}
		return &p, nil
	case 301:
		loc, err := r.Location()
		if err != nil {
			return nil, err
		}
		to := strings.TrimPrefix(loc.Path, "/wiki/")
		return nil, ErrRedirect{Page: page, To: to}
	case 404:
		return nil, ErrNotFound{Page: page}
	default:
		return nil, fmt.Errorf("%q: wikipedia error (status: %d)", page, code)
	}
}
