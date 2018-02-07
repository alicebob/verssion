// test helpers

package web_test

import (
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/alicebob/verssion/core"
)

var noRedirect = func(*http.Request, []*http.Request) error {
	return http.ErrUseLastResponse
}

// get is a helper for GETs
func get(t *testing.T, s *httptest.Server, url string) (int, string) {
	t.Helper()

	c := s.Client()
	c.CheckRedirect = noRedirect
	if c.Jar == nil {
		jar, err := cookiejar.New(nil)
		if err != nil {
			t.Fatal(err)
		}
		c.Jar = jar
	}
	r, err := s.Client().Get(s.URL + url)
	if err != nil {
		t.Fatal(err)
	}
	defer r.Body.Close()

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		t.Fatal(err)
	}

	return r.StatusCode, string(b)
}

// getRedirect does a GET and expects a redirect. Returns the Location header
func getRedirect(t *testing.T, s *httptest.Server, url string) string {
	t.Helper()

	c := s.Client()
	c.CheckRedirect = noRedirect
	r, err := s.Client().Get(s.URL + url)
	if err != nil {
		t.Fatal(err)
	}
	defer r.Body.Close()
	switch c := r.StatusCode; c {
	case 301, 302:
	default:
		t.Fatalf("not a redirect: %d", c)
	}
	return r.Header.Get("Location")
}

type test func(*testing.T, string)

func mustcontain(need string) test {
	return func(t *testing.T, body string) {
		t.Helper()
		if in, want := body, need; !strings.Contains(in, want) {
			t.Fatalf("no %q found in %q", want, in)
		}
	}
}
func mustnotcontain(need string) test {
	return func(t *testing.T, body string) {
		t.Helper()
		if in, wantnot := body, need; strings.Contains(in, wantnot) {
			t.Fatalf("%q found in %q", wantnot, in)
		}
	}
}

func with(t *testing.T, body string, tests ...test) {
	t.Helper()
	for _, test_ := range tests {
		test_(t, body)
	}
}

type FixedSpider struct {
	pages map[string]core.Page
	errs  map[string]error
}

func NewFixedSpider(pages ...core.Page) *FixedSpider {
	ps := map[string]core.Page{}
	for _, p := range pages {
		ps[p.Page] = p
	}
	return &FixedSpider{
		pages: ps,
		errs:  map[string]error{},
	}
}

var _ core.Spider = NewFixedSpider()

func (s *FixedSpider) Spider(page string) (*core.Page, error) {
	if err, ok := s.errs[page]; ok && err != nil {
		return nil, err
	}

	p, ok := s.pages[page]
	if !ok {
		return nil, core.ErrNotFound{Page: page}
	}
	return &p, nil
}

func (s *FixedSpider) SetError(page string, err error) {
	s.errs[page] = err
}
