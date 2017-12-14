package web

import (
	"log"
	"sort"
	"strings"

	"github.com/alicebob/verssion/core"
)

// from textarea to pages
func toPages(q string) []string {
	var (
		ps []string
	)
	for _, l := range strings.Split(q, "\n") {
		l = strings.TrimSpace(l)
		if len(l) == 0 {
			continue
		}
		ps = append(ps, core.WikiBasePage(l))
	}
	return ps
}

func unique(ps []string) []string {
	m := map[string]struct{}{}
	for _, p := range ps {
		m[p] = struct{}{}
	}
	res := make([]string, 0, len(m))
	for p := range m {
		res = append(res, p)
	}
	sort.Strings(res)
	return res
}

func runUpdates(db core.DB, spider core.Spider, pages []string) ([]string, []error) {
	var (
		ret    []string
		errors []error
	)

	for _, p := range pages {
		if n, err := StoreSpider(db, spider, p); err != nil {
			log.Printf("update %q: %s", p, err)
			errors = append(errors, err)
		} else {
			ret = append(ret, n.Page)
		}
	}
	return ret, errors
}
