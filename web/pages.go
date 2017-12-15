package web

import (
	"sort"
	"strings"
	"sync"

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
		retc   = make(chan string, len(pages))
		errorc = make(chan error, len(pages))
		wg     sync.WaitGroup
	)

	for _, p := range pages {
		wg.Add(1)
		go func(p string) {
			defer wg.Done()
			n, err := StoreSpider(db, spider, p)
			if err != nil {
				errorc <- err
				return
			}
			retc <- n.Page
		}(p)
	}
	wg.Wait()
	close(retc)
	close(errorc)

	var (
		res    []string
		errors []error
	)
	for r := range retc {
		res = append(res, r)
	}
	for e := range errorc {
		errors = append(errors, e)
	}
	return res, errors
}
