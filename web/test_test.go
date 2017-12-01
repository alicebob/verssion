// test helpers
package web_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

// get is a helper for GETs
func get(t *testing.T, s *httptest.Server, url string) (int, string) {
	c := s.Client()
	c.CheckRedirect = func(*http.Request, []*http.Request) error {
		return http.ErrUseLastResponse
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
