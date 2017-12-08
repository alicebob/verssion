// test helpers

package web_test

import (
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"strings"
	"testing"
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

func contains(t *testing.T, body string, needles ...string) {
	t.Helper()
	for _, need := range needles {
		if in, want := body, need; !strings.Contains(in, want) {
			t.Fatalf("no %q found", want)
		}
	}
}
