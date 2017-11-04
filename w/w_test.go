package w

import (
	"os"
	"testing"
)

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
			Homepage: "git-scm.com",
		},
		{
			Filename: "debian.html",
			Version:  "9.2 (Stretch)",
			Homepage: "www.debian.org",
		},
		{
			Filename: "postgresql.html",
			Version:  "10.0 / 5 October 2017",
			Homepage: "postgresql.org",
		},
	}

	for _, c := range cases {
		r, err := os.Open("./data/" + c.Filename)
		if err != nil {
			t.Fatal(err)
		}
		defer r.Close()
		stable, homepage := StableVersion(r)
		if have, want := stable, c.Version; have != want {
			t.Errorf("have %v, want %v", have, want)
		}
		if have, want := homepage, c.Homepage; have != want {
			t.Errorf("have %v, want %v", have, want)
		}
	}
}
