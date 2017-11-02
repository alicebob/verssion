package w

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

type WikiPage struct {
	StableVersion string
	Homepage      string
}

var client = &http.Client{
	CheckRedirect: func(*http.Request, []*http.Request) error {
		return http.ErrUseLastResponse
	},
}

// GetPage downloads and parses given wikipage
func GetPage(page string) (WikiPage, error) {
	p, err := client.Get(WikiURL(page))
	if err != nil {
		return WikiPage{}, err
	}
	defer p.Body.Close()
	if p.StatusCode != 200 {
		return WikiPage{}, fmt.Errorf("not ok: %d", p.StatusCode)
	}

	return StableVersion(p.Body), nil
}

func StableVersion(n io.Reader) WikiPage {
	ts, err := FindTables(n)
	if err != nil {
		return WikiPage{}
	}
	p := WikiPage{}
	for _, t := range ts {
		for _, r := range t.Rows {
			if len(r) < 2 {
				continue
			}
			k, v := r[0], r[1]
			switch k {
			case "Stable release", "Latest release":
				p.StableVersion = strings.Split(v, ";")[0]
			case "Official website", "Website":
				if v != "" {
					p.Homepage = v
				}
			}
		}
	}
	return p
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
