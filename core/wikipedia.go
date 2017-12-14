package core

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"golang.org/x/net/html"
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
		doc, err := html.Parse(r.Body)
		if err != nil {
			return nil, err
		}

		canonical, err := CanonicalPage(doc)
		if err == nil && canonical != "" && canonical != page {
			return nil, ErrRedirect{Page: page, To: canonical}
		}

		p.StableVersion, p.Homepage = StableVersion(doc)
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

// given a full wikipedia url returns the page name
// also works with parts of the wikipedia url
func WikiBasePage(href string) string {
	u, err := url.Parse(href)
	if err != nil {
		return href
	}
	p := u.Path
	p = strings.TrimPrefix(p, "/wiki/")
	p = strings.Replace(p, " ", "_", -1)
	return p
}
