package w

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Page struct {
	StableVersion string
}

var client = &http.Client{
	CheckRedirect: func(*http.Request, []*http.Request) error {
		return http.ErrUseLastResponse
	},
}

// GetPage downloads and parses given wikipage
func GetPage(page string) (Page, error) {
	p, err := client.Get(wikiURL(page))
	if err != nil {
		return Page{}, err
	}
	defer p.Body.Close()
	if p.StatusCode != 200 {
		return Page{}, fmt.Errorf("not ok: %d", p.StatusCode)
	}

	return Page{
		StableVersion: StableVersion(p.Body),
	}, nil
}

func StableVersion(n io.Reader) string {
	ts, err := FindTables(n)
	if err != nil {
		return ""
	}
	for _, t := range ts {
		for _, r := range t.Rows {
			if len(r) < 2 {
				continue
			}
			k, v := r[0], r[1]
			switch k {
			case "Stable release", "Latest release":
				return strings.Split(v, ";")[0]
			}
		}
	}
	return ""
}

func wikiURL(page string) string {
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
