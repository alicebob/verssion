package w

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
func GetPage(page string) (Page, error) {
	p := Page{
		Page: page,
		T:    time.Now().UTC(),
	}

	// no redirects
	r, err := client.Get(WikiURL(page))
	if err != nil {
		return p, err
	}
	defer r.Body.Close()

	switch code := r.StatusCode; code {
	case 200:
		p.StableVersion, p.Homepage = StableVersion(r.Body)
		if p.StableVersion == "" {
			fmt.Errorf("%q: no version found", page)
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
		for _, r := range t.Rows {
			if len(r) < 2 {
				continue
			}
			k, v := r[0], r[1]
			switch k {
			case "Stable release", "Latest release":
				stable = v
			case "Official website", "Website":
				if homepage == "" && v != "" {
					homepage = v
				}
			}
		}
	}
	return stable, homepage
}

func WikiURL(page string) string {
	return "https://en.wikipedia.org/wiki/" + page
}

// title version of a wikipage path
func Title(page string) string {
	return strings.Replace(page, "_", " ", -1)
}

func Titles(pages []string) []string {
	var titles []string
	for _, p := range pages {
		titles = append(titles, Title(p))
	}
	return titles
}
