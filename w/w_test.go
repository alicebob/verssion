package w

import (
	"os"
	"testing"
)

func TestStableVersion(t *testing.T) {
	type cas struct {
		Filename string
		Version  string
	}
	cases := []cas{
		{
			Filename: "git.html",
			Version:  "2.14.2 / 22 September 2017",
		},
		{
			Filename: "debian.html",
			Version:  "9.2 (Stretch)",
		},
	}

	for _, c := range cases {
		r, err := os.Open("./data/" + c.Filename)
		if err != nil {
			t.Fatal(err)
		}
		defer r.Close()
		version := StableVersion(r)
		if have, want := version, c.Version; have != want {
			t.Errorf("have %v, want %v", have, want)
		}
	}
}
