package web_test

import (
	"net/http/httptest"
	"testing"

	"github.com/alicebob/verssion/core"
	"github.com/alicebob/verssion/web"
)

func TestIndex(t *testing.T) {
	var (
		db = core.NewMemory()
		m  = web.Mux("/", db, web.NotFetcher(), "")
	)
	s := httptest.NewServer(m)
	defer s.Close()
	db.Store(core.Page{Page: "Debian", StableVersion: "my version"})
	db.Store(core.Page{Page: "Glasgow_Haskell_Compiler", StableVersion: "8.2.1 / July 22, 2017"})
	db.Store(core.Page{Page: "WithLink", StableVersion: "a [link](https://link.me)!"})

	status, body := get(t, s, "")
	if have, want := status, 200; have != want {
		t.Fatalf("have %v, want %v", have, want)
	}

	contains(t, body,
		"<title>VeRSSion",
		"Debian",
		"my version",
		"Glasgow Haskell Compiler", // titleification test
		"Glasgow_Haskell_Compiler",
		"a link!",
	)
}
