package core

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

var client = http.Client{
	CheckRedirect: func(*http.Request, []*http.Request) error {
		return http.ErrUseLastResponse
	},
}

type ErrRedirect struct {
	Page, To string
}

func (e ErrRedirect) Error() string {
	return fmt.Sprintf("%q: see page %q", e.Page, e.To)
}

type ErrNotFound struct {
	Page string
}

func (e ErrNotFound) Error() string {
	return fmt.Sprintf("%q: no such page", e.Page)
}

// GetPage downloads and parses given wikipage
func GetPage(page, url string) (Page, error) {
	p := Page{
		Page: page,
		T:    time.Now().UTC(),
	}

	// no redirects
	r, err := client.Get(url)
	if err != nil {
		return p, err
	}
	defer r.Body.Close()

	switch code := r.StatusCode; code {
	case 200:
		p.StableVersion, p.Homepage = StableVersion(r.Body)
		if p.StableVersion == "" {
			return p, fmt.Errorf("%q: no version found", page)
		}
		return p, nil
	case 301:
		loc, err := r.Location()
		if err != nil {
			return p, err
		}
		to := strings.TrimPrefix(loc.Path, "/wiki/")
		return p, ErrRedirect{Page: page, To: to}
	case 404:
		return p, ErrNotFound{Page: page}
	default:
		return p, fmt.Errorf("%q: wikipedia error (status: %d)", page, code)
	}
}

func StableVersion(n io.Reader) (string, string) {
	var stable, homepage string

	ts, err := FindTables(n)
	if err != nil {
		return "", ""
	}
	for _, t := range ts {
		for i, r := range t.Rows {
			if len(r) == 0 {
				continue
			}
			v := ""
			if len(r) > 1 {
				v = r[1]
			}
			switch k := r[0]; k {
			case "Stable release", "Latest release", "Last release":
				stable = v
			case "Stable release(s) [±]":
				// Firefox, has a table with versions. The version is in the
				// next row.
				if len(t.Rows) > i {
					if nextRow := t.Rows[i+1]; len(nextRow) > 0 {
						stable = nextRow[0]
					}
				}
			case "Official website", "Website":
				if homepage == "" && v != "" {
					homepage = v
				}
			}
		}
	}
	return stable, homepage
}

// title version of a wikipage path
func Title(page string) string {
	s := strings.Replace(page, "_", " ", -1)
	// remove some common disambiguations
	s = strings.Replace(s, " (software)", "", -1)
	s = strings.Replace(s, " (programming language)", "", -1)
	return s
}

func Titles(pages []string) []string {
	var titles []string
	for _, p := range pages {
		titles = append(titles, Title(p))
	}
	return titles
}
