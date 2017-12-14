package core

import (
	"os"
	"testing"

	"golang.org/x/net/html"
)

func TestCanonicalPage(t *testing.T) {
	type cas struct {
		Filename  string
		Canonical string
	}
	cases := []cas{
		{
			Filename:  "git.html",
			Canonical: "Git",
		},
		{
			Filename:  "debian.html",
			Canonical: "Debian",
		},
		{
			Filename:  "leftpad.html",
			Canonical: "Npm_(software)",
		},
	}

	for _, c := range cases {
		r, err := os.Open("./data/" + c.Filename)
		if err != nil {
			t.Fatal(err)
		}
		defer r.Close()

		doc, err := html.Parse(r)
		if err != nil {
			t.Fatal(err)
		}

		can, err := CanonicalPage(doc)
		if err != nil {
			t.Fatal(err)
		}
		if have, want := can, c.Canonical; have != want {
			t.Errorf("have %q, want %q", have, want)
		}
	}
}

func TestStableVersion(t *testing.T) {
	type cas struct {
		Filename string
		Version  string
		Homepage string
	}
	cases := []cas{
		{
			Filename: "git.html",
			Version:  "2.14.2 / 22 September 2017",
			Homepage: "[git-scm.com](https://git-scm.com)",
		},
		{
			Filename: "debian.html",
			Version:  "9.2 (Stretch) (October 7, 2017)",
			Homepage: "[www.debian.org](https://www.debian.org)",
		},
		{
			Filename: "postgresql.html",
			Version:  "10.0 / 5 October 2017",
			Homepage: "[postgresql.org](https://postgresql.org)",
		},
		{
			Filename: "python.html",
			Version:  "3.6.3 / 3 October 2017\n2.7.14 / 16 September 2017",
			Homepage: "[www.python.org](https://www.python.org/)",
		},
		{
			Filename: "firefox.html",
			Version:  "Standard 56.0.2 / 26 October 2017\nESR 52.4.1 / 9 October 2017",
			Homepage: "[mozilla.org/firefox](https://mozilla.org/firefox)",
		},
		{
			Filename: "pine.html",
			Version:  "4.64 (September 28, 2005)",
			Homepage: "[www.washington.edu/pine](http://www.washington.edu/pine)",
		},
		{
			Filename: "systemd.html",
			Version:  "235 (October 6, 2017)",
			Homepage: "[freedesktop.org/.../systemd/](https://freedesktop.org/wiki/Software/systemd/)",
		},
		{
			Filename: "mariadb.html",
			Version:  "10.2.11 (November 28, 2017)",
			Homepage: "[mariadb.org](https://mariadb.org/), [mariadb.com](https://mariadb.com/)",
		},
	}

	for _, c := range cases {
		r, err := os.Open("./data/" + c.Filename)
		if err != nil {
			t.Fatal(err)
		}
		defer r.Close()

		doc, err := html.Parse(r)
		if err != nil {
			t.Fatal(err)
		}

		stable, homepage := StableVersion(doc)
		if have, want := stable, c.Version; have != want {
			t.Errorf("have %q, want %q", have, want)
		}
		if have, want := homepage, c.Homepage; have != want {
			t.Errorf("have %q, want %q", have, want)
		}
	}
}

func TestTitle(t *testing.T) {
	for title, want := range map[string]string{
		"Foo":                        "Foo",
		"Foo_bar":                    "Foo bar",
		"Foo (not software)":         "Foo (not software)",
		"Foo (software)":             "Foo",
		"Foo (programming language)": "Foo",
	} {
		if have := Title(title); have != want {
			t.Errorf("%q: have %q, want %q", title, have, want)
		}
	}
}
