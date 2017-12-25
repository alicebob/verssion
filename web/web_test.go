package web_test

import (
	"net/http/httptest"
	"testing"

	"github.com/alicebob/verssion/core"
	"github.com/alicebob/verssion/web"
)

func TestWeb(t *testing.T) {
	var (
		db = core.NewMemory()
		m  = web.Mux("/", db, nil, "")
	)
	s := httptest.NewServer(m)
	defer s.Close()
	db.Store(core.Page{Page: "Debian", StableVersion: "my version"})

	{
		status, _ := get(t, s, "/robots.txt")
		if have, want := status, 200; have != want {
			t.Fatalf("have %v, want %v", have, want)
		}
	}

	{
		status, _ := get(t, s, "/nosuch.txt")
		if have, want := status, 404; have != want {
			t.Fatalf("have %v, want %v", have, want)
		}
	}
}
