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
	}

	for _, c := range cases {
		r, err := os.Open("./data/" + c.Filename)
		if err != nil {
			t.Fatal(err)
		}
		defer r.Close()
		wp := StableVersion(r)
		if have, want := wp.StableVersion, c.Version; have != want {
			t.Errorf("have %v, want %v", have, want)
		}
		if have, want := wp.Homepage, c.Homepage; have != want {
			t.Errorf("have %v, want %v", have, want)
		}
	}
}
