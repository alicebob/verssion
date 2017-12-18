// Compare versions as found on https://verssion.one (and hence wikipedia)
// against the version published on the website of each project.
// This is fuzzy.

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

	"golang.org/x/net/html"
)

var (
	baseURL = flag.String("base", "https://verssion.one", "verssion URL")
	verbose = flag.Bool("v", false, "verbose")
)

func main() {
	flag.Parse()

	for _, f := range flag.Args() {
		fmt.Printf("f: %s\n", f)
		lines, err := readFile(f)
		if err != nil {
			log.Fatal(err)
		}
		for _, l := range lines {
			fmt.Printf("l: %s\n", l)
			if len(l) != 2 {
				log.Printf("invalid line: %q", l)
				continue
			}
			if err := lookat(l[0], l[1]); err != nil {
				log.Printf(err.Error())
			}
		}
	}
}

// readFile reads lines and splits them in fields.
// Lines starting with a '#' are skipped.
func readFile(file string) ([][]string, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var ls [][]string
	s := bufio.NewScanner(f)
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
	txt, err := randomSite(url)
	if err != nil {
		return fmt.Errorf("%s: %s", cur.Page, err)
	}
	if *verbose {
		fmt.Printf("%s: txt:%q\n", cur.Page, txt)
	}

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
	res, err := http.Get(*baseURL + "/p/" + p + "/?format=json")
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if code := res.StatusCode; code != 200 {
		return nil, fmt.Errorf("%s: HTTP %d", p, code)
	}

	r := &Page{}
	if err := json.NewDecoder(res.Body).Decode(r); err != nil {
		return nil, err
	}
	return r, nil
}

func randomSite(url string) (string, error) {
	res, err := http.Get(url)
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
