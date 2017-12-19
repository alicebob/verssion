// Spiders pages from verssion (to be sure they are updated), and optionally
// compares the version against the version published on the website of each
// project.
// That's a fuzzy process.

package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"golang.org/x/net/html"
)

var (
	baseURL = flag.String("base", "https://verssion.one", "verssion URL")
	sleep   = flag.Duration("sleep", time.Second, "sleep between spiders")

	client = &http.Client{
		CheckRedirect: func(*http.Request, []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
)

func main() {
	flag.Parse()

	files := flag.Args()
	if len(files) == 0 {
		files = []string{"-"}
	}
	for _, f := range files {
		lines, err := readFile(f)
		if err != nil {
			log.Fatal(err)
		}
		for _, l := range lines {
			switch len(l) {
			case 1:
				// only spider verssion
				fmt.Printf("ping %s\n", l[0])
				_, err := verssion(l[0])
				if err != nil {
					log.Print(err.Error())
				}
			case 2:
				// compare verssion against version
				if err := lookat(l[0], l[1]); err != nil {
					log.Print(err.Error())
				}
			default:
				log.Printf("invalid line: %q", strings.Join(l, " "))
			}
			time.Sleep(*sleep)
		}
	}
}

// readFile reads lines and splits them in fields.
// Lines starting with a '#' are skipped.
func readFile(file string) ([][]string, error) {
	var f io.Reader
	if file == "-" {
		f = os.Stdin
	} else {
		fh, err := os.Open(file)
		if err != nil {
			return nil, err
		}
		defer fh.Close()
		f = fh
	}

	var (
		s  = bufio.NewScanner(f)
		ls [][]string
	)
	for s.Scan() {
		line := s.Text()
		if strings.HasPrefix(line, "#") {
			continue
		}
		fs := strings.Fields(line)
		if len(fs) == 0 {
			continue
		}
		ls = append(ls, fs)
	}
	return ls, nil
}

func lookat(p, url string) error {
	cur, err := verssion(p)
	if err != nil {
		return err
	}
	fmt.Printf("%s: verssion %q (=~ %q)\n", cur.Page, cur.StableVersion, guessVersion(cur.StableVersion))
	txt, err := getText(url)
	if err != nil {
		return fmt.Errorf("%s: %s", cur.Page, err)
	}
	// fmt.Printf("%s: txt:%q\n", cur.Page, txt)

	fmt.Printf("%s: site =~ %q\n", cur.Page, guessVersion(txt))

	if have, want := guessVersion(cur.StableVersion), guessVersion(txt); have != want {
		log.Printf("%s: we have %q, website has %q", cur.Page, have, want)
	}
	return nil
}

type Page struct {
	Page          string `json:"page"`
	Title         string `json:"title"`
	StableVersion string `json:"stable_version"`
	Homepage      string `json:"homepage"`
}

func verssion(p string) (*Page, error) {
	res, err := client.Get(*baseURL + "/p/" + p + "/?format=json")
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	switch code := res.StatusCode; code {
	case 200:
		// all fine
	case 301, 302:
		l, _ := res.Location()
		return nil, fmt.Errorf("%s: HTTP %d %s", p, code, l)
	default:
		return nil, fmt.Errorf("%s: HTTP %d", p, code)
	}

	r := &Page{}
	return r, json.NewDecoder(res.Body).Decode(r)
}

// getText gets a text version of a URL
func getText(url string) (string, error) {
	res, err := client.Get(url)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	if code := res.StatusCode; code != 200 {
		return "", fmt.Errorf("HTTP %d", code)
	}

	return asText(res.Body)
}

func asText(s io.Reader) (string, error) {
	doc, err := html.Parse(s)
	if err != nil {
		return "", err
	}
	txt := ""
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode {
			switch n.Data {
			case "head", "noscript", "script", "style":
				return
			}
		}
		if n.Type == html.TextNode {
			txt += n.Data
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
	return txt, nil
}

var versionre = regexp.MustCompile(`v?\d+(\.\d+)+`)

// finds the first thing which looks like a version number. Or an empty string.
func guessVersion(s string) string {
	return versionre.FindString(s)
}
